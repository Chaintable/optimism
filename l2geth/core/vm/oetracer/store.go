// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package oetracer

import (
	"context"

	"github.com/ethereum-optimism/optimism/l2geth/common"
	"github.com/ethereum-optimism/optimism/l2geth/ethdb/leveldb"
)

var txTraceStore Store = (*TxTraceDB)(nil)

var txEnable bool = false

type TxTraceDB struct {
	kv *leveldb.Database
}

// ReadTxTrace retrieve tx trace data from underlying kv store.
func (t *TxTraceDB) ReadTxTrace(ctx context.Context, txHash common.Hash) ([]byte, error) {
	return t.kv.Get(txHash.Bytes())
}

// WriteTxTrace save tx trace data to underlying kv store.
func (t *TxTraceDB) WriteTxTrace(ctx context.Context, txHash common.Hash, trace []byte) error {
	return t.kv.Put(txHash.Bytes(), trace)
}

// Store contains all the methods for tx-trace to interact with the underlying database.
type Store interface {
	// ReadTxTrace retrieve tracing result from underlying database.
	ReadTxTrace(ctx context.Context, txHash common.Hash) ([]byte, error)
	// WriteTxTrace write tracing result to underlying database.
	WriteTxTrace(ctx context.Context, txHash common.Hash, trace []byte) error
}

// OpenTxTraceDB ...
func OpenTxTraceDB(path string) error {
	db, err := leveldb.New(path, 256, 0, "")
	if err != nil {
		return err
	}
	txTraceStore = &TxTraceDB{kv: db}
	txEnable = true
	return nil
}

// CloseTxTraceDB ...
func CloseTxTraceDB() error {
	if !txEnable {
		return nil
	}
	return txTraceStore.(*TxTraceDB).kv.Close()
}

// GetTxTraceStore ...
func GetTxTraceStore() Store {
	if !txEnable {
		return nil
	}
	return txTraceStore
}
