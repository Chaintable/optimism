package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/l2geth/common"
	"github.com/ethereum-optimism/optimism/l2geth/common/hexutil"
	"github.com/ethereum-optimism/optimism/l2geth/common/math"
	"github.com/ethereum-optimism/optimism/l2geth/core"
	"github.com/ethereum-optimism/optimism/l2geth/core/types"
	txtrace "github.com/ethereum-optimism/optimism/l2geth/core/vm/oetracer"
	"github.com/ethereum-optimism/optimism/l2geth/log"
	"github.com/ethereum-optimism/optimism/l2geth/rollup/rcfg"
	"github.com/ethereum-optimism/optimism/l2geth/rpc"
)

type PreExecTx struct {
	ChainId              *big.Int        `json:"chainId,omitempty"`
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`
	Data                 *hexutil.Bytes  `json:"data"`
	Input                *hexutil.Bytes  `json:"input"`
}

type PreExecAPI struct {
	e *Ethereum
}

func NewPreExecAPI(e *Ethereum) *PreExecAPI {
	return &PreExecAPI{e: e}
}

func (api *PreExecAPI) GetLogs(ctx context.Context, args *PreExecTx) (*types.Receipt, error) {
	state, header, err := api.e.APIBackend.StateAndHeaderByNumberOrHash(ctx, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if state == nil || err != nil {
		return nil, err
	}
	// Set sender address or use a default if none specified
	var addr common.Address
	if args.From == nil {
		if !rcfg.UsingOVM {
			if wallets := api.e.APIBackend.AccountManager().Wallets(); len(wallets) > 0 {
				if accounts := wallets[0].Accounts(); len(accounts) > 0 {
					addr = accounts[0].Address
				}
			}
		}
	} else {
		addr = *args.From
	}
	// Set default gas & gas price if none were set
	gas := uint64(math.MaxUint64 / 2)
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	gasPrice := new(big.Int)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	}

	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}

	var data []byte
	if args.Data != nil {
		data = []byte(*args.Data)
	}

	// Currently, the blocknumber and timestamp actually refer to the L1BlockNumber and L1Timestamp
	// attached to each transaction. We need to modify the blocknumber and timestamp to reflect this,
	// or else the result of `eth_call` will not be correct.
	blockNumber := header.Number
	timestamp := header.Time
	if rcfg.UsingOVM {
		block, err := api.e.APIBackend.BlockByNumber(ctx, rpc.BlockNumber(header.Number.Uint64()))
		if err != nil {
			return nil, err
		}
		if block != nil {
			txs := block.Transactions()
			if header.Number.Uint64() != 0 {
				if len(txs) != 1 {
					return nil, fmt.Errorf("block %d has more than 1 transaction", header.Number.Uint64())
				}
				tx := txs[0]
				blockNumber = tx.L1BlockNumber()
				timestamp = tx.L1Timestamp()
			}
		}
	}
	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, value, gas, gasPrice, data, false, blockNumber, timestamp, types.QueueOriginSequencer)

	evm, vmError, err := api.e.APIBackend.GetEVM(ctx, msg, state, header, nil)
	if err != nil {
		return nil, err
	}

	// Setup the gas pool (also for unmetered requests)
	// and apply the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	_, gas, failed, _ := core.ApplyMessage(evm, msg, gp)
	if err := vmError(); err != nil {
		return nil, err
	}
	var root []byte
	receipt := types.NewReceipt(root, failed, gas)
	receipt.Logs = state.Logs()
	return receipt, nil
}

// TraceTransaction tracing pre-exec transaction object.
func (api *PreExecAPI) TraceTransaction(ctx context.Context, args *PreExecTx) (interface{}, error) {
	state, header, err := api.e.APIBackend.StateAndHeaderByNumberOrHash(ctx, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if state == nil || err != nil {
		return nil, err
	}
	// Set sender address or use a default if none specified
	var addr common.Address
	if args.From == nil {
		if !rcfg.UsingOVM {
			if wallets := api.e.APIBackend.AccountManager().Wallets(); len(wallets) > 0 {
				if accounts := wallets[0].Accounts(); len(accounts) > 0 {
					addr = accounts[0].Address
				}
			}
		}
	} else {
		addr = *args.From
	}
	// Set default gas & gas price if none were set
	gas := uint64(math.MaxUint64 / 2)
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	gasPrice := new(big.Int)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	}

	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}

	var data []byte
	if args.Data != nil {
		data = []byte(*args.Data)
	}

	// Currently, the blocknumber and timestamp actually refer to the L1BlockNumber and L1Timestamp
	// attached to each transaction. We need to modify the blocknumber and timestamp to reflect this,
	// or else the result of `eth_call` will not be correct.
	blockNumber := header.Number
	timestamp := header.Time
	if rcfg.UsingOVM {
		block, err := api.e.APIBackend.BlockByNumber(ctx, rpc.BlockNumber(header.Number.Uint64()))
		if err != nil {
			return nil, err
		}
		if block != nil {
			txs := block.Transactions()
			if header.Number.Uint64() != 0 {
				if len(txs) != 1 {
					return nil, fmt.Errorf("block %d has more than 1 transaction", header.Number.Uint64())
				}
				tx := txs[0]
				blockNumber = tx.L1BlockNumber()
				timestamp = tx.L1Timestamp()
			}
		}
	}
	// Create new call message
	msg := types.NewMessage(addr, args.To, 0, value, gas, gasPrice, data, false, blockNumber, timestamp, types.QueueOriginSequencer)

	tracer := txtrace.NewTraceStructLogger(nil)
	tracer.SetFrom(msg.From())
	tracer.SetTo(msg.To())
	tracer.SetValue(*msg.Value())
	tracer.SetGasUsed(msg.Gas())
	tracer.SetBlockHash(header.Hash())
	tracer.SetBlockNumber(header.Number)
	tracer.SetTxIndex(0)
	vmCfg := api.e.APIBackend.eth.blockchain.GetVMConfig()
	vmCfg.Debug = true
	vmCfg.Tracer = tracer
	evm, vmError, err := api.e.APIBackend.GetEVM(ctx, msg, state, header, vmCfg)
	if err != nil {
		return nil, err
	}

	// Setup the gas pool (also for unmetered requests)
	// and apply the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	_, _, failed, _ := core.ApplyMessage(evm, msg, gp)
	if err := vmError(); err != nil {
		return nil, err
	}
	if failed {
		log.Warn("apply message with transaction tracing failed", "err", err)
	}
	tracer.Finalize()
	return tracer.GetResult(), nil
}
