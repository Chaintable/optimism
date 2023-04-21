// Copyright 2021 The go-ethereum Authors
// // This file is part of the go-ethereum library.
// //
// // The go-ethereum library is free software: you can redistribute it and/or modify
// // it under the terms of the GNU Lesser General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.
// //
// // The go-ethereum library is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// // GNU Lesser General Public License for more details.
// //
// // You should have received a copy of the GNU Lesser General Public License
// // along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// // Package ethapi implements the general Ethereum API functions.

package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum-optimism/optimism/l2geth/common"
	"github.com/ethereum-optimism/optimism/l2geth/core"
	"github.com/ethereum-optimism/optimism/l2geth/core/state"
	"github.com/ethereum-optimism/optimism/l2geth/core/types"
	"github.com/ethereum-optimism/optimism/l2geth/core/vm"
	"github.com/ethereum-optimism/optimism/l2geth/core/vm/oetracer"
	txtrace "github.com/ethereum-optimism/optimism/l2geth/core/vm/oetracer"

	"github.com/ethereum/go-ethereum/log"
)

// PublicTxTraceAPI provides an API to tracing transaction or block information.
// // It offers only methods that operate on public data that is freely available to anyone.
type PublicTxTraceAPI struct {
	e *Ethereum
}

// NewPublicTxTraceAPI creates a new trace API.
func NewPublicTxTraceAPI(e *Ethereum) *PublicTxTraceAPI {
	return &PublicTxTraceAPI{e: e}
}

// Transaction trace_transaction function returns transaction traces.
func (api *PublicTxTraceAPI) Transaction(ctx context.Context, txHash common.Hash) (interface{}, error) {
	if oetracer.GetTxTraceStore() != nil {
		raw, err := oetracer.GetTxTraceStore().ReadTxTrace(ctx, txHash)
		if err != nil {
			goto replay
		}
		if bytes.Equal(raw, []byte{}) { // empty response
			goto replay
		}

		var res interface{}
		if err := json.Unmarshal(raw, &res); err != nil {
			return []byte{}, err
		}
		return res, nil
	}

replay:
	log.Warn("tx trace store is nil, fallback to default trace method", "txHash", txHash)

	if api.e.blockchain == nil {
		return []byte{}, fmt.Errorf("blockchain corruput")
	}

	tx, blockHash, blockNumber, index, err := api.e.APIBackend.GetTransaction(ctx, txHash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("transaction %#v not found", txHash)
	}
	// It shouldn't happen in practice.
	if blockNumber == 0 {
		return nil, fmt.Errorf("genesis is not traceable")
	}

	txctx := &txTraceContext{
		tx:    tx,
		index: int(index),
		block: blockHash,
	}

	msg, vmctx, statedb, err := computeTxEnv(api.e, blockHash, int(index), defaultTraceReexec)
	if err != nil {
		return nil, err
	}
	// Trace the transaction and return
	return api.traceTx(ctx, msg, vmctx, txctx, statedb)
}

// traceTx configures a new tracer according to the provided configuration, and
// executes the given message in the provided environment. The return value will
// be as parity's one.
func (api *PublicTxTraceAPI) traceTx(ctx context.Context, message core.Message, vmctx vm.Context, txctx *txTraceContext, statedb *state.StateDB) (interface{}, error) {
	var (
		tracer *txtrace.StructLogger
		err    error
	)

	// Construct trace logger to record result as parity's one
	tracer = txtrace.NewTraceStructLogger(nil)

	// Fill essential info into logger
	tracer.SetFrom(message.From())
	tracer.SetTo(message.To())
	tracer.SetValue(*message.Value())
	tracer.SetGasUsed(message.Gas())
	tracer.SetBlockHash(txctx.block)
	tracer.SetBlockNumber(vmctx.BlockNumber)
	tracer.SetTx(txctx.tx.Hash())
	tracer.SetTxIndex(uint(txctx.index))

	// Run the transaction with tracing enabled.
	vmenv := vm.NewEVM(vmctx, statedb, api.e.blockchain.Config(), vm.Config{Debug: true, Tracer: tracer})
	_, _, failed, err := core.ApplyMessage(vmenv, message, new(core.GasPool).AddGas(message.Gas()))
	if err != nil {
		return nil, fmt.Errorf("tracing failed: %v", err)
	}
	if failed {
		log.Warn("apply message with transaction tracing failed", "txHash", txctx.tx.Hash())
	}

	tracer.Finalize()
	return tracer.GetResult(), nil
}

// txTraceContext is the contextual infos about a transaction before it gets run.
type txTraceContext struct {
	tx    *types.Transaction
	index int         // Index of the transaction within the block
	block common.Hash // Hash of the block containing the transaction
}
