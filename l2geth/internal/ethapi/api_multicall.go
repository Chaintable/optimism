package ethapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum-optimism/optimism/l2geth/common"
	"github.com/ethereum-optimism/optimism/l2geth/common/hexutil"
	"github.com/ethereum-optimism/optimism/l2geth/common/math"
	"github.com/ethereum-optimism/optimism/l2geth/core"
	"github.com/ethereum-optimism/optimism/l2geth/core/state"
	"github.com/ethereum-optimism/optimism/l2geth/core/types"
	"github.com/ethereum-optimism/optimism/l2geth/core/vm"
	"github.com/ethereum-optimism/optimism/l2geth/log"
	"github.com/ethereum-optimism/optimism/l2geth/metrics"
	"github.com/ethereum-optimism/optimism/l2geth/rollup/rcfg"
	"github.com/ethereum-optimism/optimism/l2geth/rpc"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type multiCallResp struct {
	Results []*callResult   `json:"results"`
	Stats   *multiCallStats `json:"stats"`
}

type callResult struct {
	Code      int           `json:"code"`
	Err       string        `json:"err"`
	FromCache bool          `json:"fromCache"`
	Result    hexutil.Bytes `json:"result"`
	GasUsed   int64         `json:"gasUsed"`
	TimeCost  float64       `json:"timeCost"`
}

type multiCallStats struct {
	BlockNum     int64       `json:"blockNum"`
	BlockHash    common.Hash `json:"blockHash"`
	BlockTime    int64       `json:"blockTime"`
	Success      bool        `json:"success"`
	CacheEnabled bool        `json:"cacheEnabled"`
}

const (
	singleCallTimeout = 1 * time.Second
	multiCallLimit    = 50

	// client param error
	errCodeTxArgs               = -40000
	errNativeMethodNotFound     = -40001
	errNativeMethodInput        = -40002
	errNativeMethodInputAddress = -40003
	errBlockNotFound            = -40004

	// evm processing error
	errNativeMethodOutput     = -40010
	errNativeMethodStateError = -40011
	errMessageExecuting       = -40012
	errEVMCancelled           = -40013
	errEVMReverted            = -40014
)

var (
	ethMultiCallCacheHit   = metrics.GetOrRegisterMeter("rpc/ethmulticall/cache/hit", nil)
	ethMultiCallCacheCount = metrics.GetOrRegisterMeter("rpc/ethmulticall/cache/count", nil)
)

const (
	nativeAddr = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

var (
	// copied from: accounts/abi/abi_test.go
	Uint8, _   = abi.NewType("uint8", "", nil)
	Uint256, _ = abi.NewType("uint256", "", nil)
	String, _  = abi.NewType("string", "", nil)
	Address, _ = abi.NewType("address", "", nil)

	erc20ABI = abi.ABI{
		Methods: map[string]abi.Method{
			"name":        funcName,
			"symbol":      funcSymbol,
			"decimals":    funcDecimals,
			"totalSupply": funcTotalSupply,
			"balanceOf":   funcBalanceOf,
		},
	}

	funcName = abi.NewMethod("name", "name", abi.Function, "", false, false,
		[]abi.Argument{},
		[]abi.Argument{
			{Name: "", Type: String, Indexed: false},
		},
	)
	funcSymbol = abi.NewMethod("symbol", "symbol", abi.Function, "", false, false,
		[]abi.Argument{},
		[]abi.Argument{
			{Name: "", Type: String, Indexed: false},
		},
	)
	funcDecimals = abi.NewMethod("decimals", "decimals", abi.Function, "", false, false,
		[]abi.Argument{},
		[]abi.Argument{
			{Name: "", Type: Uint8, Indexed: false},
		},
	)
	funcTotalSupply = abi.NewMethod("totalSupply", "totalSupply", abi.Function, "", false, false,
		[]abi.Argument{},
		[]abi.Argument{
			{Name: "", Type: Uint256, Indexed: false},
		},
	)
	funcBalanceOf = abi.NewMethod("balanceOf", "balanceOf", abi.Function, "", false, false,
		[]abi.Argument{
			{Name: "", Type: Address, Indexed: false},
		},
		[]abi.Argument{
			{Name: "", Type: Uint256, Indexed: false},
		},
	)
)

func ethCallCacheKey(b Backend, blockHash common.Hash, to *common.Address, input []byte) string {
	var sb strings.Builder

	h := sha256.New()
	h.Write(input)
	bs := h.Sum(nil)

	sb.Grow(len(bs) + len(to.Bytes()) + len(blockHash))
	sb.Write(blockHash[:])
	sb.Write(bytes.ToLower(to.Bytes()))
	sb.Write(bs)

	return sb.String()
}

func handleNative(ctx context.Context, state *state.StateDB, msg types.Message) ([]byte, int, error) {
	data := msg.Data()
	method, err := erc20ABI.MethodById(data)
	if err != nil {
		return nil, errNativeMethodNotFound, err
	}
	switch method.Name {
	case "name", "symbol":
		res, err := method.Outputs.Pack("ETH")
		if err != nil {
			return nil, errNativeMethodOutput, err
		}
		return res, 0, nil
	case "decimals":
		res, err := method.Outputs.Pack(uint8(18))
		if err != nil {
			return nil, errNativeMethodOutput, err
		}
		return res, 0, nil
	case "totalSupply":
		res, err := method.Outputs.Pack(big.NewInt(1_000_000_000_000_000_000)) // 1 ETH
		if err != nil {
			return nil, errNativeMethodOutput, err
		}
		return res, 0, nil
	case "balanceOf":
		inputs, err := method.Inputs.Unpack(data[4:])
		if err != nil || len(inputs) == 0 {
			return nil, errNativeMethodInput, err
		}
		address, ok := inputs[0].(common.Address)
		if !ok {
			return nil, errNativeMethodInputAddress, fmt.Errorf("input address error")
		}
		balance, err := method.Outputs.Pack(state.GetBalance(address))
		if err != nil {
			return nil, errNativeMethodOutput, err
		}
		if state.Error() != nil {
			return nil, errNativeMethodStateError, state.Error()
		}
		return balance, 0, nil
	default:
		return nil, errNativeMethodNotFound, fmt.Errorf("method not found")
	}
}

func doOneCall(ctx context.Context, b Backend, state *state.StateDB, header *types.Header, arg CallArgs, disableCache bool, vmCfg *vm.Config, globalGasCap *big.Int) (*callResult, error) {
	var err error
	var result = &callResult{}

	start := time.Now()

	if !disableCache {
		// try load result from cache
		cacheKey := ethCallCacheKey(b, header.Hash(), arg.To, *arg.Data)
		ethMultiCallCacheCount.Mark(1)
		if r, ok := b.GetCallCache(cacheKey); ok {
			ethMultiCallCacheHit.Mark(1)
			if res, ok := r.(*callResult); ok {
				res.FromCache = true
				return res, nil
			}
		}
		defer func() {
			// `err` here specifics to non-evm error. Evm internal error won't prevent
			// caching the result
			if err == nil {
				b.SetCallCache(cacheKey, result, int64(len(result.Result)))
			}
		}()
	}
	// make sure this will be called prior to the SetCallCache defer func on returning
	defer func() {
		result.TimeCost = time.Since(start).Seconds()
	}()

	// Set sender address or use a default if none specified
	var addr common.Address
	if arg.From == nil {
		if !rcfg.UsingOVM {
			if wallets := b.AccountManager().Wallets(); len(wallets) > 0 {
				if accounts := wallets[0].Accounts(); len(accounts) > 0 {
					addr = accounts[0].Address
				}
			}
		}
	} else {
		addr = *arg.From
	}

	// Set default gas & gas price if none were set
	gas := uint64(math.MaxUint64 / 2)
	if arg.Gas != nil {
		gas = uint64(*arg.Gas)
	}
	if globalGasCap != nil && globalGasCap.Uint64() < gas {
		log.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap.Uint64()
	}
	gasPrice := new(big.Int)
	if arg.GasPrice != nil {
		gasPrice = arg.GasPrice.ToInt()
	}

	value := new(big.Int)
	if arg.Value != nil {
		value = arg.Value.ToInt()
	}

	var data []byte
	if arg.Data != nil {
		data = []byte(*arg.Data)
	}

	// Currently, the blocknumber and timestamp actually refer to the L1BlockNumber and L1Timestamp
	// attached to each transaction. We need to modify the blocknumber and timestamp to reflect this,
	// or else the result of `eth_call` will not be correct.
	blockNumber := header.Number
	timestamp := header.Time
	if rcfg.UsingOVM {
		block, err := b.BlockByNumber(ctx, rpc.BlockNumber(header.Number.Uint64()))
		if err != nil {
			result.Code = errBlockNotFound
			result.Err = err.Error()
			return result, err
		}
		if block != nil {
			txs := block.Transactions()
			if header.Number.Uint64() != 0 {
				if len(txs) != 1 {
					err := fmt.Errorf("block %d has more than 1 transaction", header.Number.Uint64())
					result.Code = errCodeTxArgs
					result.Err = err.Error()
				}
				tx := txs[0]
				blockNumber = tx.L1BlockNumber()
				timestamp = tx.L1Timestamp()
			}
		}
	}

	// Create new call message
	msg := types.NewMessage(addr, arg.To, 0, value, gas, gasPrice, data, false, blockNumber, timestamp, types.QueueOriginSequencer)

	// skip EVM if requests for native token
	if strings.ToLower(msg.To().Hex()) == nativeAddr {
		res, code, err := handleNative(ctx, state, msg)
		if err != nil {
			result.Code = code
			result.Err = err.Error()
		}
		result.Result = res
		return result, err
	}

	// Get a new instance of the EVM.
	//evm, _, _ := b.GetEVM(ctx, msg, state, header, &vm.Config{NoBaseFee: true})
	// Get a new instance of the EVM.
	evm, _, _ := b.GetEVM(ctx, msg, state, header, vmCfg) // never return error

	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()
	// Execute the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	res, gas, failed, err := core.ApplyMessage(evm, msg, gp)
	if err != nil {
		result.Code = errMessageExecuting
		result.Err = err.Error()
		return result, err
	}

	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		result.Code = errEVMCancelled
		result.Err = fmt.Sprintf("execution cancelled, either fast failed or exceeding timeout(%v)", singleCallTimeout)
		return result, err
	}

	if failed {
		reason, errUnpack := abi.UnpackRevert(res)
		err := errors.New("execution reverted")
		if errUnpack == nil {
			err = fmt.Errorf("execution reverted: %v", reason)
		}
		result.Code = errEVMReverted
		result.Err = err.Error()
		return result, err
	}

	result.Result = (hexutil.Bytes)(res)
	result.GasUsed = int64(gas)

	return result, nil
}

func (s *PublicBlockChainAPI) MultiCall(ctx context.Context, args []CallArgs, blockNrOrHash rpc.BlockNumberOrHash, pfastFail, puseParallel, pdisableCache *bool, overrides *map[common.Address]account) (resp *multiCallResp, err error) {

	// maximum calls check
	if len(args) > multiCallLimit {
		return nil, fmt.Errorf("calls exceed limit, expected: <%v, actual: %v", multiCallLimit, len(args))
	}

	setb := func(p *bool, d bool) bool {
		if p == nil {
			return d
		}
		return *p
	}

	fastFail := setb(pfastFail, true)
	useParallel := setb(puseParallel, true)
	disableCache := setb(pdisableCache, false)

	// check block & state
	state, header, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	var accounts map[common.Address]account
	if overrides != nil {
		accounts = *overrides
	}
	// Override the fields of specified contracts before execution.
	for addr, account := range accounts {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return nil, fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}

	blockTime := header.Time

	ret := make([]*callResult, len(args))
	stats := &multiCallStats{
		BlockNum:     header.Number.Int64(),
		BlockHash:    header.Hash(),
		BlockTime:    int64(blockTime),
		Success:      true,
		CacheEnabled: !disableCache,
	}

	ctx, cancel := context.WithTimeout(ctx, singleCallTimeout)
	defer cancel()

	if useParallel {
		// run in parallel
		var wg sync.WaitGroup
		for i, arg := range args {
			wg.Add(1)
			go func(i int, arg CallArgs) {
				defer wg.Done()

				// state is not reentrancy in concurrent scenarios, so use a copy
				r, _ := doOneCall(ctx, s.b, state.Copy(), header, arg, disableCache, &vm.Config{}, s.b.RPCGasCap())
				ret[i] = r
				if r.Err != "" {
					stats.Success = false
					if fastFail {
						cancel()
					}
					return
				}
			}(i, arg)
		}
		wg.Wait()

		return &multiCallResp{Results: ret, Stats: stats}, nil
	}

	// run in sequence
	failedOnce := false
	for i, arg := range args {
		if failedOnce {
			ret[i] = nil
			continue
		}

		r, _ := doOneCall(ctx, s.b, state, header, arg, disableCache, &vm.Config{}, s.b.RPCGasCap())
		ret[i] = r
		if r.Err != "" {
			stats.Success = false
			if fastFail {
				failedOnce = true
			}
			continue
		}
	}

	return &multiCallResp{Results: ret, Stats: stats}, nil
}
