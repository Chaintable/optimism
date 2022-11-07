package driver

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
)

<<<<<<< HEAD
type Driver struct {
	s *state
}

<<<<<<< HEAD
<<<<<<< HEAD
type BatchSubmitter interface {
	Submit(config *rollup.Config, batches []*derive.BatchData) (common.Hash, error)
}

type Downloader interface {
	InfoByHash(ctx context.Context, hash common.Hash) (derive.L1Info, error)
	Fetch(ctx context.Context, blockHash common.Hash) (derive.L1Info, types.Transactions, types.Receipts, error)
	FetchAllTransactions(ctx context.Context, window []eth.BlockID) ([]types.Transactions, error)
}

type Engine interface {
	GetPayload(ctx context.Context, payloadId eth.PayloadID) (*eth.ExecutionPayload, error)
	ForkchoiceUpdate(ctx context.Context, state *eth.ForkchoiceState, attr *eth.PayloadAttributes) (*eth.ForkchoiceUpdatedResult, error)
	NewPayload(ctx context.Context, payload *eth.ExecutionPayload) error
	PayloadByHash(context.Context, common.Hash) (*eth.ExecutionPayload, error)
	PayloadByNumber(context.Context, *big.Int) (*eth.ExecutionPayload, error)
}

type L1Chain interface {
	L1BlockRefByNumber(context.Context, uint64) (eth.L1BlockRef, error)
	L1BlockRefByHash(context.Context, common.Hash) (eth.L1BlockRef, error)
	L1HeadBlockRef(context.Context) (eth.L1BlockRef, error)
	L1Range(ctx context.Context, base eth.BlockID, max uint64) ([]eth.BlockID, error)
}

type L2Chain interface {
	ForkchoiceUpdate(ctx context.Context, state *eth.ForkchoiceState, attr *eth.PayloadAttributes) (*eth.ForkchoiceUpdatedResult, error)
=======
=======
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
type Metrics interface {
	RecordPipelineReset()
	RecordSequencingError()
	RecordPublishingError()
	RecordDerivationError()

	RecordReceivedUnsafePayload(payload *eth.ExecutionPayload)

	RecordL1Ref(name string, ref eth.L1BlockRef)
	RecordL2Ref(name string, ref eth.L2BlockRef)

	RecordUnsafePayloadsBuffer(length uint64, memSize uint64, next eth.BlockID)

	SetDerivationIdle(idle bool)

	RecordL1ReorgDepth(d uint64)
	CountSequencedTxs(count int)
}

<<<<<<< HEAD
>>>>>>> v0.5.24
type Downloader interface {
	InfoByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, error)
	Fetch(ctx context.Context, blockHash common.Hash) (eth.BlockInfo, types.Transactions, eth.ReceiptsFetcher, error)
}

=======
>>>>>>> @eth-optimism/l2geth@0.5.27
type L1Chain interface {
	derive.L1Fetcher
	L1BlockRefByLabel(context.Context, eth.BlockLabel) (eth.L1BlockRef, error)
}

type L2Chain interface {
	derive.Engine
<<<<<<< HEAD
	L2BlockRefHead(ctx context.Context) (eth.L2BlockRef, error)
>>>>>>> v0.5.23
	L2BlockRefByNumber(ctx context.Context, l2Num *big.Int) (eth.L2BlockRef, error)
=======
	L2BlockRefByLabel(ctx context.Context, label eth.BlockLabel) (eth.L2BlockRef, error)
>>>>>>> v0.5.24
	L2BlockRefByHash(ctx context.Context, l2Hash common.Hash) (eth.L2BlockRef, error)
	L2BlockRefByNumber(ctx context.Context, num uint64) (eth.L2BlockRef, error)
}

<<<<<<< HEAD
type outputInterface interface {
	// insertEpoch creates and inserts one epoch on top of the safe head. It prefers blocks it creates to what is recorded in the unsafe chain.
	// It returns the new L2 head and L2 Safe head and if there was a reorg. This function must return if there was a reorg otherwise the L2 chain must be traversed.
	insertEpoch(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.L2BlockRef, l2Finalized eth.BlockID, l1Input []eth.BlockID) (eth.L2BlockRef, eth.L2BlockRef, bool, error)

	// createNewBlock builds a new block based on the L2 Head, L1 Origin, and the current mempool.
	createNewBlock(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.BlockID, l2Finalized eth.BlockID, l1Origin eth.L1BlockRef) (eth.L2BlockRef, *eth.ExecutionPayload, error)

	// processBlock simply tries to add the block to the chain, reorging if necessary, and updates the forkchoice of the engine.
	processBlock(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.BlockID, l2Finalized eth.BlockID, payload *eth.ExecutionPayload) error
=======
type DerivationPipeline interface {
	Reset()
	Step(ctx context.Context) error
	SetUnsafeHead(head eth.L2BlockRef)
	AddUnsafePayload(payload *eth.ExecutionPayload)
	Finalize(ref eth.L1BlockRef)
	FinalizedL1() eth.L1BlockRef
	Finalized() eth.L2BlockRef
	SafeL2Head() eth.L2BlockRef
	UnsafeL2Head() eth.L2BlockRef
	Origin() eth.L1BlockRef
}

type L1StateIface interface {
	HandleNewL1HeadBlock(head eth.L1BlockRef)
	HandleNewL1SafeBlock(safe eth.L1BlockRef)
	HandleNewL1FinalizedBlock(finalized eth.L1BlockRef)

	L1Head() eth.L1BlockRef
	L1Safe() eth.L1BlockRef
	L1Finalized() eth.L1BlockRef
}

type L1OriginSelectorIface interface {
	FindL1Origin(ctx context.Context, l1Head eth.L1BlockRef, l2Head eth.L2BlockRef) (eth.L1BlockRef, error)
}

type SequencerIface interface {
	StartBuildingBlock(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.BlockID, l2Finalized eth.BlockID, l1Origin eth.L1BlockRef) error
	CompleteBuildingBlock(ctx context.Context) (*eth.ExecutionPayload, error)

	// createNewBlock builds a new block based on the L2 Head, L1 Origin, and the current mempool.
<<<<<<< HEAD
	createNewBlock(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.BlockID, l2Finalized eth.BlockID, l1Origin eth.L1BlockRef) (eth.L2BlockRef, *eth.ExecutionPayload, error)
>>>>>>> v0.5.23
=======
	CreateNewBlock(ctx context.Context, l2Head eth.L2BlockRef, l2SafeHead eth.BlockID, l2Finalized eth.BlockID, l1Origin eth.L1BlockRef) (eth.L2BlockRef, *eth.ExecutionPayload, error)
>>>>>>> @eth-optimism/l2geth@0.5.27
}

type Network interface {
	// PublishL2Payload is called by the driver whenever there is a new payload to publish, synchronously with the driver main loop.
	PublishL2Payload(ctx context.Context, payload *eth.ExecutionPayload) error
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func NewDriver(cfg rollup.Config, l2 *l2.Source, l1 *l1.Source, network Network, log log.Logger, snapshotLog log.Logger, sequencer bool) *Driver {
=======
func NewDriver(driverCfg *Config, cfg *rollup.Config, l2 *l2.Source, l1 *l1.Source, network Network, log log.Logger, snapshotLog log.Logger) *Driver {
>>>>>>> v0.5.23
=======
func NewDriver(driverCfg *Config, cfg *rollup.Config, l2 L2Chain, l1 L1Chain, network Network, log log.Logger, snapshotLog log.Logger, metrics Metrics) *Driver {
>>>>>>> v0.5.24
	output := &outputImpl{
		Config: cfg,
		dl:     l1,
		l2:     l2,
		log:    log,
	}
<<<<<<< HEAD
	return &Driver{
		s: NewState(log, snapshotLog, cfg, l1, l2, output, network, sequencer),
	}
=======

	var state *state
	verifConfDepth := NewConfDepth(driverCfg.VerifierConfDepth, func() eth.L1BlockRef { return state.l1Head }, l1)
	derivationPipeline := derive.NewDerivationPipeline(log, cfg, verifConfDepth, l2, metrics)
	state = NewState(driverCfg, log, snapshotLog, cfg, l1, l2, output, derivationPipeline, network, metrics)
	return &Driver{s: state}
>>>>>>> v0.5.23
}

func (d *Driver) OnL1Head(ctx context.Context, head eth.L1BlockRef) error {
	return d.s.OnL1Head(ctx, head)
}

func (d *Driver) OnL1Safe(ctx context.Context, safe eth.L1BlockRef) error {
	return d.s.OnL1Safe(ctx, safe)
}

func (d *Driver) OnL1Finalized(ctx context.Context, finalized eth.L1BlockRef) error {
	return d.s.OnL1Finalized(ctx, finalized)
}

func (d *Driver) OnUnsafeL2Payload(ctx context.Context, payload *eth.ExecutionPayload) error {
	return d.s.OnUnsafeL2Payload(ctx, payload)
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
func (d *Driver) SyncStatus(ctx context.Context) (*SyncStatus, error) {
=======
func (d *Driver) ResetDerivationPipeline(ctx context.Context) error {
	return d.s.ResetDerivationPipeline(ctx)
}

func (d *Driver) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
>>>>>>> v0.5.24
	return d.s.SyncStatus(ctx)
}

>>>>>>> v0.5.23
func (d *Driver) Start(ctx context.Context) error {
	return d.s.Start(ctx)
}
func (d *Driver) Close() error {
	return d.s.Close()
=======
// NewDriver composes an events handler that tracks L1 state, triggers L2 derivation, and optionally sequences new L2 blocks.
func NewDriver(driverCfg *Config, cfg *rollup.Config, l2 L2Chain, l1 L1Chain, network Network, log log.Logger, snapshotLog log.Logger, metrics Metrics) *Driver {
	sequencer := NewSequencer(log, cfg, l1, l2)
	l1State := NewL1State(log, metrics)
	findL1Origin := NewL1OriginSelector(log, cfg, l1, driverCfg.SequencerConfDepth)
	verifConfDepth := NewConfDepth(driverCfg.VerifierConfDepth, l1State.L1Head, l1)
	derivationPipeline := derive.NewDerivationPipeline(log, cfg, verifConfDepth, l2, metrics)

	return &Driver{
		l1State:          l1State,
		derivation:       derivationPipeline,
		idleDerivation:   false,
		stateReq:         make(chan chan struct{}),
		forceReset:       make(chan chan struct{}, 10),
		config:           cfg,
		driverConfig:     driverCfg,
		done:             make(chan struct{}),
		log:              log,
		snapshotLog:      snapshotLog,
		l1:               l1,
		l2:               l2,
		l1OriginSelector: findL1Origin,
		sequencer:        sequencer,
		network:          network,
		metrics:          metrics,
		l1HeadSig:        make(chan eth.L1BlockRef, 10),
		l1SafeSig:        make(chan eth.L1BlockRef, 10),
		l1FinalizedSig:   make(chan eth.L1BlockRef, 10),
		unsafeL2Payloads: make(chan *eth.ExecutionPayload, 10),
	}
>>>>>>> @eth-optimism/l2geth@0.5.27
}
