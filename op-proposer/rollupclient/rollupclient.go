package rollupclient

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/eth"
<<<<<<< HEAD
	"github.com/ethereum-optimism/optimism/op-node/node"
=======
	"github.com/ethereum-optimism/optimism/op-node/rollup/driver"
>>>>>>> v0.5.23
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type RollupClient struct {
	rpc *rpc.Client
}

func NewRollupClient(rpc *rpc.Client) *RollupClient {
	return &RollupClient{rpc}
}

<<<<<<< HEAD
func (r *RollupClient) GetBatchBundle(
	ctx context.Context,
	req *node.BatchBundleRequest,
) (*node.BatchBundleResponse, error) {

	var batchResponse = new(node.BatchBundleResponse)
	err := r.rpc.CallContext(ctx, &batchResponse, "optimism_getBatchBundle", req)
	return batchResponse, err
}

=======
>>>>>>> v0.5.23
func (r *RollupClient) OutputAtBlock(ctx context.Context, blockNum *big.Int) ([]eth.Bytes32, error) {
	var output []eth.Bytes32
	err := r.rpc.CallContext(ctx, &output, "optimism_outputAtBlock", hexutil.EncodeBig(blockNum))
	return output, err
}
<<<<<<< HEAD
=======

func (r *RollupClient) SyncStatus(ctx context.Context) (*driver.SyncStatus, error) {
	var output *driver.SyncStatus
	err := r.rpc.CallContext(ctx, &output, "optimism_syncStatus")
	return output, err
}

func (r *RollupClient) Version(ctx context.Context) (string, error) {
	var output string
	err := r.rpc.CallContext(ctx, &output, "optimism_version")
	return output, err
}
>>>>>>> v0.5.23
