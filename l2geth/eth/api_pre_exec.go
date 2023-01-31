package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum-optimism/optimism/l2geth/accounts/abi"
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
	ChainId  *big.Int        `json:"chainId,omitempty"`
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	Data     *hexutil.Bytes  `json:"data"`
	Input    *hexutil.Bytes  `json:"input"`
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
	vmCfg := *api.e.APIBackend.eth.blockchain.GetVMConfig()
	vmCfg.Debug = true
	vmCfg.Tracer = tracer
	evm, vmError, err := api.e.APIBackend.GetEVM(ctx, msg, state, header, &vmCfg)
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

const (
	UnKnown            = 1000
	InsufficientBalane = 1001
	Reverted           = 1002
)

type PreError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PreResult struct {
	Trace     *[]txtrace.ActionTrace `json:"trace"`
	Logs      []*types.Log           `json:"logs"`
	StateDiff txtrace.StateDiff      `json:"stateDiff"`
	Error     PreError               `json:"error"`
	GasUsed   uint64                 `json:"gasUsed"`
}

func toPreError(err error, failed bool, result []byte) PreError {
	preErr := PreError{
		Code: UnKnown,
	}
	if err != nil {
		preErr.Msg = err.Error()
	}
	if failed {
		preErr.Code = Reverted
		if result != nil {
			preErr.Msg, _ = abi.UnpackRevert(result)
		}
		return preErr
	}
	if strings.HasPrefix(preErr.Msg, "out of gas") {
		preErr.Code = Reverted
	}
	if strings.HasPrefix(preErr.Msg, "insufficient funds") {
		preErr.Code = InsufficientBalane
	}
	if strings.HasPrefix(preErr.Msg, "insufficient balance") {
		preErr.Code = InsufficientBalane
	}
	return preErr
}

func (api *PreExecAPI) TraceMany(ctx context.Context, origins []PreExecTx) ([]PreResult, error) {
	preResList := make([]PreResult, 0)
	state, header, err := api.e.APIBackend.StateAndHeaderByNumberOrHash(ctx, rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if state == nil || err != nil {
		return nil, err
	}
	for i := 0; i < len(origins); i++ {
		origin := origins[i]
		if origin.Nonce == nil {
			preResList = append(preResList, PreResult{
				Error: PreError{
					Code: UnKnown,
					Msg:  "nonce is nil",
				},
			})
			continue
		}
		if i > 0 && (uint64)(*origin.Nonce) <= (uint64)(*origins[i-1].Nonce) {
			preResList = append(preResList, PreResult{
				Error: PreError{
					Code: UnKnown,
					Msg:  fmt.Sprintf("nonce decreases, tx index %d has nonce %d, tx index %d has nonce %d", i-1, (uint64)(*origins[i-1].Nonce), i, (uint64)(*origin.Nonce)),
				},
			})
			continue
		}

		// Set sender address or use a default if none specified
		var addr common.Address
		if origin.From == nil {
			if !rcfg.UsingOVM {
				if wallets := api.e.APIBackend.AccountManager().Wallets(); len(wallets) > 0 {
					if accounts := wallets[0].Accounts(); len(accounts) > 0 {
						addr = accounts[0].Address
					}
				}
			}
		} else {
			addr = *origin.From
		}
		// Set default gas & gas price if none were set
		gas := uint64(math.MaxUint64 / 2)
		if origin.Gas != nil {
			gas = uint64(*origin.Gas)
		}
		gasPrice := new(big.Int)
		if origin.GasPrice != nil {
			gasPrice = origin.GasPrice.ToInt()
		}

		value := new(big.Int)
		if origin.Value != nil {
			value = origin.Value.ToInt()
		}

		var data []byte
		if origin.Data != nil {
			data = []byte(*origin.Data)
		}

		blockNumber := header.Number
		timestamp := header.Time
		if rcfg.UsingOVM {
			block, err := api.e.APIBackend.BlockByNumber(ctx, rpc.BlockNumber(header.Number.Uint64()))
			if err != nil {
				preResList = append(preResList, PreResult{
					Error: PreError{
						Code: UnKnown,
						Msg:  err.Error(),
					},
				})
				continue
			}
			if block != nil {
				txs := block.Transactions()
				if header.Number.Uint64() != 0 {
					if len(txs) != 1 {
						preResList = append(preResList, PreResult{
							Error: PreError{
								Code: UnKnown,
								Msg:  fmt.Sprintf("block %d has more than 1 transaction", header.Number.Uint64()),
							},
						})
						continue
					}
					tx := txs[0]
					blockNumber = tx.L1BlockNumber()
					timestamp = tx.L1Timestamp()
				}
			}
		}
		msg := types.NewMessage(addr, origin.To, 0, value, gas, gasPrice, data, false, blockNumber, timestamp, types.QueueOriginSequencer)
		txHash := common.BigToHash(big.NewInt(int64(i)))
		tracer := txtrace.NewTraceStructLogger(nil)
		tracer.SetFrom(msg.From())
		tracer.SetTo(msg.To())
		tracer.SetValue(*msg.Value())
		tracer.SetGasUsed(msg.Gas())
		tracer.SetBlockHash(header.Hash())
		tracer.SetBlockNumber(header.Number)
		tracer.SetTxIndex(0)
		// fix panic: use value copy replace pointer reference
		vmCfg := *api.e.APIBackend.eth.blockchain.GetVMConfig()
		vmCfg.Debug = true
		vmCfg.Tracer = tracer
		evm, vmError, err := api.e.APIBackend.GetEVM(ctx, msg, state, header, &vmCfg)
		if err != nil {
			preResList = append(preResList, PreResult{
				Error: PreError{
					Code: UnKnown,
					Msg:  err.Error(),
				},
			})
			continue
		}
		state.Prepare(txHash, header.Hash(), i)
		gp := new(core.GasPool).AddGas(math.MaxUint64)
		result, usedGas, failed, err := core.ApplyMessage(evm, msg, gp)
		if failed {
			preRes := PreResult{
				Error:   toPreError(err, failed, result),
				GasUsed: usedGas,
			}
			preResList = append(preResList, preRes)
			continue
		}
		if err := vmError(); err != nil {
			preRes := PreResult{
				Error:   toPreError(err, failed, result),
				GasUsed: usedGas,
			}
			preResList = append(preResList, preRes)
			continue
		}
		if err != nil {
			preRes := PreResult{
				Error:   toPreError(err, failed, result),
				GasUsed: usedGas,
			}
			preResList = append(preResList, preRes)
			continue
		}
		tracer.Finalize()
		preRes := PreResult{
			Trace:     tracer.GetResult(),
			Logs:      state.GetLogs(txHash),
			StateDiff: tracer.GetStateDiff(),
			GasUsed:   usedGas,
		}
		if preRes.Error.Msg == "" && preRes.Trace != nil && len(*preRes.Trace) > 0 && (*preRes.Trace)[0].Error != "" {
			preRes.Error = PreError{
				Code: Reverted,
				Msg:  (*preRes.Trace)[0].Error,
			}
		}
		preResList = append(preResList, preRes)
	}
	return preResList, nil
}
