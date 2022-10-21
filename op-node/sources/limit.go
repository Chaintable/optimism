package sources

import (
	"context"
	"sync"

<<<<<<< HEAD
=======
	"github.com/ethereum-optimism/optimism/op-node/client"
<<<<<<< HEAD
<<<<<<< HEAD:op-node/l1/request_sema.go

>>>>>>> v0.5.23
=======
>>>>>>> v0.5.24:op-node/sources/limit.go
=======
	"github.com/ethereum/go-ethereum"
>>>>>>> @eth-optimism/l2geth@0.5.27
	"github.com/ethereum/go-ethereum/rpc"
)

type limitClient struct {
<<<<<<< HEAD
	c    RPCClient
=======
	c    client.RPC
>>>>>>> v0.5.23
	sema chan struct{}
	wg   sync.WaitGroup
}

// LimitRPC limits concurrent RPC requests (excluding subscriptions) to a given number by wrapping the client with a semaphore.
<<<<<<< HEAD
func LimitRPC(c RPCClient, concurrentRequests int) RPCClient {
=======
func LimitRPC(c client.RPC, concurrentRequests int) client.RPC {
>>>>>>> v0.5.23
	return &limitClient{
		c: c,
		// the capacity of the channel determines how many go-routines can concurrently execute requests with the wrapped client.
		sema: make(chan struct{}, concurrentRequests),
	}
}

func (lc *limitClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	lc.wg.Add(1)
	defer lc.wg.Done()
	lc.sema <- struct{}{}
	defer func() { <-lc.sema }()
	return lc.c.BatchCallContext(ctx, b)
}

func (lc *limitClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	lc.wg.Add(1)
	defer lc.wg.Done()
	lc.sema <- struct{}{}
	defer func() { <-lc.sema }()
	return lc.c.CallContext(ctx, result, method, args...)
}

func (lc *limitClient) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (ethereum.Subscription, error) {
	// subscription doesn't count towards request limit
	return lc.c.EthSubscribe(ctx, channel, args...)
}

func (lc *limitClient) Close() {
	lc.wg.Wait()
	close(lc.sema)
	lc.c.Close()
}
