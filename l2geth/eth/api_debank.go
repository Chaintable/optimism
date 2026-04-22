package eth

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"

	ptracer "github.com/Chaintable/pipeline/tracer"
	ptypes "github.com/Chaintable/pipeline/types"
	"github.com/Chaintable/pipeline/util"
	"github.com/ethereum-optimism/optimism/l2geth/common"
	"github.com/ethereum-optimism/optimism/l2geth/common/hexutil"
	"github.com/ethereum-optimism/optimism/l2geth/core"
	"github.com/ethereum-optimism/optimism/l2geth/core/types"
	"github.com/ethereum-optimism/optimism/l2geth/core/vm"
	"github.com/ethereum-optimism/optimism/l2geth/log"
	"github.com/ethereum-optimism/optimism/l2geth/rpc"
)

// DebankAPI provides the trace_debankBlock RPC method.
type DebankAPI struct {
	eth *Ethereum

	genesisOnce   sync.Once
	genesisOutput *ptypes.DebankOutPut
	genesisErr    error
}

// NewDebankAPI creates a new DebankAPI instance.
func NewDebankAPI(eth *Ethereum) *DebankAPI {
	return &DebankAPI{eth: eth}
}

// DebankBlock replays a block and returns the DebankOutPut.
func (api *DebankAPI) DebankBlock(ctx context.Context, blockNr rpc.BlockNumber) (*ptypes.DebankOutPut, error) {
	blockchain := api.eth.BlockChain()
	chainConfig := blockchain.Config()

	// === Genesis Block ===
	if blockNr == 0 {
		return api.debankGenesisBlock(blockchain)
	}

	// === Regular Block ===
	block := blockchain.GetBlockByNumber(uint64(blockNr))
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", blockNr)
	}
	parentBlock := blockchain.GetBlockByNumber(uint64(blockNr - 1))
	if parentBlock == nil {
		return nil, fmt.Errorf("parent block #%d not found", blockNr-1)
	}

	statedb, err := blockchain.StateAt(parentBlock.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to get state at %s: %v", parentBlock.Root().Hex(), err)
	}
	receipts := blockchain.GetReceiptsByHash(block.Hash())
	if len(receipts) == 0 && len(block.Transactions()) > 0 {
		return nil, fmt.Errorf("receipts not found for block #%d", blockNr)
	}

	// 使用 Pipeline 的 RPCTracer
	tracer := ptracer.NewRPCTracer()
	tracer.OnBlockStart(block)
	statedb.OnLog = tracer.OnLog

	signer := types.MakeSigner(chainConfig, block.Number())
	gasPool := new(core.GasPool).AddGas(block.GasLimit())

	for i, tx := range block.Transactions() {
		msg, err := tx.AsMessage(signer)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tx %d to message: %v", i, err)
		}
		statedb.Prepare(tx.Hash(), block.Hash(), i)

		tracer.OnTxStart(tx, msg.From())

		vmConfig := vm.Config{Debug: true, Tracer: tracer.GetCallTracer()}
		evmContext := core.NewEVMContext(msg, block.Header(), blockchain, nil)
		vmenv := vm.NewEVM(evmContext, statedb, chainConfig, vmConfig)

		_, _, _, err = core.ApplyMessage(vmenv, msg, gasPool)
		tracer.OnTxEnd(receipts[i], err)
		statedb.Finalise(chainConfig.IsEIP158(block.Number()))
	}

	// Verify state root
	root := statedb.IntermediateRoot(chainConfig.IsEIP158(block.Number()))
	if root != block.Root() {
		return nil, fmt.Errorf("state root mismatch: got %s, want %s", root.Hex(), block.Root().Hex())
	}

	snapDestructs, snapAccounts, snapStorage, codes := statedb.Output()
	return tracer.GetOutPut(parentBlock.Root(), root, snapDestructs, snapAccounts, snapStorage, codes), nil
}

// debankGenesisBlock handles genesis block (block 0).
// The first call triggers an async background export; subsequent calls return
// the cached result or a "still initializing" error.
func (api *DebankAPI) debankGenesisBlock(blockchain *core.BlockChain) (*ptypes.DebankOutPut, error) {
	api.genesisOnce.Do(func() {
		go func() {
			log.Info("Genesis block export started")
			output, err := api.buildGenesisOutput(blockchain)
			if err != nil {
				log.Error("Genesis block export failed", "err", err)
			} else {
				log.Info("Genesis block export completed")
			}
			api.genesisOutput = output
			api.genesisErr = err
		}()
	})

	if api.genesisOutput != nil {
		return api.genesisOutput, nil
	}
	if api.genesisErr != nil {
		return nil, api.genesisErr
	}
	return nil, fmt.Errorf("genesis block is still exporting, please wait")
}

func (api *DebankAPI) buildGenesisOutput(blockchain *core.BlockChain) (*ptypes.DebankOutPut, error) {
	block := blockchain.GetBlockByNumber(0)
	if block == nil {
		return nil, fmt.Errorf("genesis block not found")
	}

	statedb, err := blockchain.StateAt(block.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to get genesis state at %s: %v", block.Root().Hex(), err)
	}

	log.Info("Starting genesis block dump")
	dump := statedb.RawDump(false, false, false)
	log.Info("Genesis block dump completed", "accounts", len(dump.Accounts))

	header := util.BuildPilelineBlockHeader(block)
	blockDiff := ptracer.GenesisDumpToStateDiff(dump)
	blockDiff.Hash = header.StateRoot

	blockFile := &ptypes.BlockFile{
		Block:            util.BuildPipelineBlock(block),
		Txs:              make([]ptypes.Transaction, 0),
		Events:           make([]ptypes.Event, 0),
		Traces:           make([]ptypes.Trace, 0),
		ErrorEvents:      make([]ptypes.Event, 0),
		ErrorTraces:      make([]ptypes.Trace, 0),
		StorageContracts: make([]string, 0),
	}

	zeroAddr := "0x0000000000000000000000000000000000000000"
	txIdx := int64(0)

	sortedAddrs := make([]common.Address, 0, len(dump.Accounts))
	for addr := range dump.Accounts {
		sortedAddrs = append(sortedAddrs, addr)
	}
	sort.Slice(sortedAddrs, func(i, j int) bool {
		return sortedAddrs[i].Hex() < sortedAddrs[j].Hex()
	})

	for _, addr := range sortedAddrs {
		acc := dump.Accounts[addr]
		addrLower := strings.ToLower(addr.Hex())

		if len(acc.Storage) > 0 {
			blockFile.StorageContracts = append(blockFile.StorageContracts, addrLower)
		}

		balance, ok := new(big.Int).SetString(acc.Balance, 10)
		if !ok {
			balance = big.NewInt(0)
		}

		if balance.Sign() > 0 {
			txID := fmt.Sprintf("0xgenesis01%013d%s", 0, addrLower)
			blockFile.Txs = append(blockFile.Txs, ptypes.Transaction{
				ID:               txID,
				From:             zeroAddr,
				To:               addrLower,
				Gas:              big.NewInt(0),
				GasPrice:         big.NewInt(0),
				GasUsed:          big.NewInt(0),
				Status:           true,
				GasFeeCap:        big.NewInt(0),
				GasTipCap:        big.NewInt(0),
				Input:            []byte{},
				Nonce:            big.NewInt(0),
				TransactionIndex: txIdx,
				Value:            (*hexutil.Big)(balance),
			})
			blockFile.Traces = append(blockFile.Traces, ptypes.Trace{
				ID:                util.ToHash([]string{txID, "", "0"}),
				From:              zeroAddr,
				To:                addrLower,
				Gas:               big.NewInt(0),
				GasUsed:           big.NewInt(0),
				Input:             []byte{},
				Output:            []byte{},
				Value:             (*hexutil.Big)(balance),
				CallCreateType:    "call",
				CallType:          "call",
				TxID:              txID,
				ParentTraceID:     "",
				PosInParentTrace:  0,
				SelfStorageChange: false,
				StorageChange:     false,
				Subtraces:         0,
				TraceAddress:      []int64{},
			})
			txIdx++
		}

		code := common.Hex2Bytes(acc.Code)
		if len(code) > 0 {
			txID := fmt.Sprintf("0xgenesis02%013d%s", 0, addrLower)
			blockFile.Txs = append(blockFile.Txs, ptypes.Transaction{
				ID:               txID,
				From:             zeroAddr,
				To:               addrLower,
				Gas:              big.NewInt(0),
				GasPrice:         big.NewInt(0),
				GasUsed:          big.NewInt(0),
				Status:           true,
				GasFeeCap:        big.NewInt(0),
				GasTipCap:        big.NewInt(0),
				Input:            code,
				Nonce:            big.NewInt(0),
				TransactionIndex: txIdx,
				Value:            (*hexutil.Big)(big.NewInt(0)),
			})
			blockFile.Traces = append(blockFile.Traces, ptypes.Trace{
				ID:                util.ToHash([]string{txID, "", "0"}),
				From:              zeroAddr,
				To:                addrLower,
				Gas:               big.NewInt(0),
				GasUsed:           big.NewInt(0),
				Input:             code,
				Output:            code,
				Value:             (*hexutil.Big)(big.NewInt(0)),
				CallCreateType:    "create",
				CallType:          "",
				TxID:              txID,
				ParentTraceID:     "",
				PosInParentTrace:  0,
				SelfStorageChange: false,
				StorageChange:     false,
				Subtraces:         0,
				TraceAddress:      []int64{},
			})
			txIdx++
		}
	}

	nativeTokenAddr := "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	nativeTokenTxID := fmt.Sprintf("0xgenesis03%013d%s", 0, nativeTokenAddr)
	blockFile.Txs = append(blockFile.Txs, ptypes.Transaction{
		ID:               nativeTokenTxID,
		From:             zeroAddr,
		To:               nativeTokenAddr,
		Gas:              big.NewInt(0),
		GasPrice:         big.NewInt(0),
		GasUsed:          big.NewInt(0),
		Status:           true,
		GasFeeCap:        big.NewInt(0),
		GasTipCap:        big.NewInt(0),
		Input:            []byte{},
		Nonce:            big.NewInt(0),
		TransactionIndex: txIdx,
		Value:            (*hexutil.Big)(big.NewInt(0)),
	})
	blockFile.Traces = append(blockFile.Traces, ptypes.Trace{
		ID:                util.ToHash([]string{nativeTokenTxID, "", "0"}),
		From:              zeroAddr,
		To:                nativeTokenAddr,
		Gas:               big.NewInt(0),
		GasUsed:           big.NewInt(0),
		Input:             []byte{},
		Output:            []byte{},
		Value:             (*hexutil.Big)(big.NewInt(0)),
		CallCreateType:    "create",
		CallType:          "",
		TxID:              nativeTokenTxID,
		ParentTraceID:     "",
		PosInParentTrace:  0,
		SelfStorageChange: false,
		StorageChange:     false,
		Subtraces:         0,
		TraceAddress:      []int64{},
	})

	stateDiffBytes, err := util.EncodeToRlp(blockDiff)
	if err != nil {
		log.Error("Failed to encode genesis state diff", "err", err)
		stateDiffBytes = []byte{}
	}

	return &ptypes.DebankOutPut{
		BlockFile:      blockFile,
		Header:         header,
		StateDiff:      hexutil.Bytes(stateDiffBytes),
		ValidationHash: blockFile.Validation().ValidationHash,
	}, nil
}
