package bridge

import (
	"context"
	"time"

<<<<<<< HEAD
	"github.com/ethereum-optimism/optimism/indexer/bindings/l2bridge"
<<<<<<< HEAD:indexer/services/l2/bridge/filter.go
	"github.com/ethereum-optimism/optimism/l2geth/accounts/abi/bind"
=======
=======
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
>>>>>>> v0.5.24
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
>>>>>>> v0.5.23:go/indexer/services/l2/bridge/filter.go
)

// clientRetryInterval is the interval to wait between retrying client API
// calls.
var clientRetryInterval = 5 * time.Second

// FilterWithdrawalInitiatedWithRetry retries the given func until it succeeds,
// waiting for clientRetryInterval duration after every call.
func FilterWithdrawalInitiatedWithRetry(ctx context.Context, filterer *bindings.L2StandardBridgeFilterer, opts *bind.FilterOpts) (*bindings.L2StandardBridgeWithdrawalInitiatedIterator, error) {
	for {
		ctxt, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
		opts.Context = ctxt
		res, err := filterer.FilterWithdrawalInitiated(opts, nil, nil, nil)
		cancel()
		if err == nil {
			return res, nil
		}
		logger.Error("Error fetching filter", "err", err)
		time.Sleep(clientRetryInterval)
	}
}

// FilterDepositFinalizedWithRetry retries the given func until it succeeds,
// waiting for clientRetryInterval duration after every call.
func FilterDepositFinalizedWithRetry(ctx context.Context, filterer *bindings.L2StandardBridgeFilterer, opts *bind.FilterOpts) (*bindings.L2StandardBridgeDepositFinalizedIterator, error) {
	for {
		ctxt, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
		opts.Context = ctxt
		res, err := filterer.FilterDepositFinalized(opts, nil, nil, nil)
		cancel()
		if err == nil {
			return res, nil
		}
		logger.Error("Error fetching filter", "err", err)
		time.Sleep(clientRetryInterval)
	}
}

// FilterDepositFailedWithRetry retries the given func until it succeeds,
// waiting for clientRetryInterval duration after every call.
func FilterDepositFailedWithRetry(ctx context.Context, filterer *bindings.L2StandardBridgeFilterer, opts *bind.FilterOpts) (*bindings.L2StandardBridgeDepositFailedIterator, error) {
	for {
		ctxt, cancel := context.WithTimeout(ctx, DefaultConnectionTimeout)
		opts.Context = ctxt
		res, err := filterer.FilterDepositFailed(opts, nil, nil, nil)
		cancel()
		if err == nil {
			return res, nil
		}
		logger.Error("Error fetching filter", "err", err)
		time.Sleep(clientRetryInterval)
	}
}
