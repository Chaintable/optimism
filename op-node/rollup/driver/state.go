package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	gosync "sync"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
<<<<<<< HEAD
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	"github.com/ethereum/go-ethereum/log"
)

type state struct {
	// Chain State
	l1Head      eth.L1BlockRef // Latest recorded head of the L1 Chain
	l2Head      eth.L2BlockRef // L2 Unsafe Head
	l2SafeHead  eth.L2BlockRef // L2 Safe Head - this is the head of the L2 chain as derived from L1 (thus it is Sequencer window blocks behind)
	l2Finalized eth.BlockID    // L2 Block that will never be reversed
	l1WindowBuf []eth.BlockID  // l1WindowBuf buffers the next L1 block IDs to derive new L2 blocks from, with increasing block height.

	// Rollup config
	Config    rollup.Config
	sequencer bool
=======
=======
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
>>>>>>> v0.5.24
	"github.com/ethereum/go-ethereum/log"
=======
	"github.com/ethereum-optimism/optimism/op-service/backoff"
>>>>>>> @eth-optimism/l2geth@0.5.27
)

// Deprecated: use eth.SyncStatus instead.
type SyncStatus = eth.SyncStatus

type Driver struct {
	l1State L1StateIface

	// The derivation pipeline is reset whenever we reorg.
	// The derivation pipeline determines the new l2Safe.
	derivation DerivationPipeline

	// When the derivation pipeline is waiting for new data to do anything
	idleDerivation bool

	// Requests for sync status. Synchronized with event loop to avoid reading an inconsistent sync status.
	syncStatusReq chan chan eth.SyncStatus

	// Upon receiving a channel in this channel, the derivation pipeline is forced to be reset.
	// It tells the caller that the reset occurred by closing the passed in channel.
	forceReset chan chan struct{}

	// Rollup config: rollup chain configuration
	config *rollup.Config

	// Driver config: verifier and sequencer settings
<<<<<<< HEAD
	DriverConfig *Config
>>>>>>> v0.5.23
=======
	driverConfig *Config
>>>>>>> @eth-optimism/l2geth@0.5.27

	// L1 Signals:
	//
	// Not all L1 blocks, or all changes, have to be signalled:
	// the derivation process traverses the chain and handles reorgs as necessary,
	// the driver just needs to be aware of the *latest* signals enough so to not
	// lag behind actionable data.
	l1HeadSig      chan eth.L1BlockRef
	l1SafeSig      chan eth.L1BlockRef
	l1FinalizedSig chan eth.L1BlockRef

	// L2 Signals:
	unsafeL2Payloads chan *eth.ExecutionPayload

	l1               L1Chain
	l2               L2Chain
	l1OriginSelector L1OriginSelectorIface
	sequencer        SequencerIface
	network          Network // may be nil, network for is optional

	metrics     Metrics
	log         log.Logger
	snapshotLog log.Logger
	done        chan struct{}

	wg gosync.WaitGroup
}

<<<<<<< HEAD
<<<<<<< HEAD
// NewState creates a new driver state. State changes take effect though the given output.
// Optionally a network can be provided to publish things to other nodes than the engine of the driver.
func NewState(log log.Logger, snapshotLog log.Logger, config rollup.Config, l1Chain L1Chain, l2Chain L2Chain, output outputInterface, network Network, sequencer bool) *state {
	return &state{
		Config:           config,
=======
// NewState creates a new driver state. State changes take effect though
// the given output, derivation pipeline and network interfaces.
func NewState(driverCfg *Config, log log.Logger, snapshotLog log.Logger, config *rollup.Config, l1Chain L1Chain, l2Chain L2Chain,
	output outputInterface, derivationPipeline DerivationPipeline, network Network, metrics Metrics) *state {
	return &state{
		derivation:       derivationPipeline,
		idleDerivation:   false,
		syncStatusReq:    make(chan chan eth.SyncStatus, 10),
		forceReset:       make(chan chan struct{}, 10),
		Config:           config,
		DriverConfig:     driverCfg,
>>>>>>> v0.5.23
		done:             make(chan struct{}),
		log:              log,
		snapshotLog:      snapshotLog,
		l1:               l1Chain,
		l2:               l2Chain,
		output:           output,
		network:          network,
<<<<<<< HEAD
<<<<<<< HEAD
		sequencer:        sequencer,
=======
>>>>>>> v0.5.23
		l1Heads:          make(chan eth.L1BlockRef, 10),
=======
		metrics:          metrics,
		l1HeadSig:        make(chan eth.L1BlockRef, 10),
		l1SafeSig:        make(chan eth.L1BlockRef, 10),
		l1FinalizedSig:   make(chan eth.L1BlockRef, 10),
>>>>>>> v0.5.24
		unsafeL2Payloads: make(chan *eth.ExecutionPayload, 10),
	}
}

// Start starts up the state loop.
// The loop will have been started iff err is not nil.
<<<<<<< HEAD
func (s *state) Start(ctx context.Context) error {
	l1Head, err := s.l1.L1HeadBlockRef(ctx)
	if err != nil {
		return err
	}
<<<<<<< HEAD

	// Check that we are past the genesis
	if l1Head.Number > s.Config.Genesis.L1.Number {
		l2Head, err := s.l2.L2BlockRefByNumber(ctx, nil)
		if err != nil {
			return err
		}
		// Ensure that we are on the correct chain. Note that we cannot rely on rely on the UnsafeHead being more than
		// a sequence window behind the L1 Head and must walk back 1 sequence window as we do not track the end L1 block
		// hash of the sequence window when we derive an L2 block.
		unsafeHead, safeHead, err := sync.FindL2Heads(ctx, l2Head, s.Config.SeqWindowSize, s.l1, s.l2, &s.Config.Genesis)
		if err != nil {
			return err
		}
		s.l2Head = unsafeHead
		s.l2SafeHead = safeHead

	} else {
		// Not yet reached genesis block
		// TODO: Test this codepath. That requires setting up L1, letting it run, and then creating the L2 genesis from there.
		// Note: This will not work for setting the the genesis normally, but if the L1 node is not yet synced we could get this case.
		l2genesis := eth.L2BlockRef{
			Hash:           s.Config.Genesis.L2.Hash,
			Number:         s.Config.Genesis.L2.Number,
			Time:           s.Config.Genesis.L2Time,
			L1Origin:       s.Config.Genesis.L1,
			SequenceNumber: 0,
		}
		s.l2Head = l2genesis
		s.l2SafeHead = l2genesis
	}

	s.l1Head = l1Head

	s.wg.Add(1)
	go s.loop()
=======
	s.l1Head = l1Head
	s.l2Head, _ = s.l2.L2BlockRefByNumber(ctx, nil)

=======
func (s *state) Start(_ context.Context) error {
>>>>>>> v0.5.24
=======
// Start starts up the state loop.
// The loop will have been started iff err is not nil.
func (s *Driver) Start() error {
>>>>>>> @eth-optimism/l2geth@0.5.27
	s.derivation.Reset()

	s.wg.Add(1)
	go s.eventLoop()

>>>>>>> v0.5.23
	return nil
}

<<<<<<< HEAD
func (s *state) Close() error {
<<<<<<< HEAD
	close(s.done)
=======
=======
func (s *Driver) Close() error {
>>>>>>> @eth-optimism/l2geth@0.5.27
	s.done <- struct{}{}
>>>>>>> v0.5.23
	s.wg.Wait()
	return nil
}

// OnL1Head signals the driver that the L1 chain changed the "unsafe" block,
// also known as head of the chain, or "latest".
func (s *Driver) OnL1Head(ctx context.Context, unsafe eth.L1BlockRef) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.l1HeadSig <- unsafe:
		return nil
	}
}

// OnL1Safe signals the driver that the L1 chain changed the "safe",
// also known as the justified checkpoint (as seen on L1 beacon-chain).
func (s *Driver) OnL1Safe(ctx context.Context, safe eth.L1BlockRef) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.l1SafeSig <- safe:
		return nil
	}
}

func (s *Driver) OnL1Finalized(ctx context.Context, finalized eth.L1BlockRef) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.l1FinalizedSig <- finalized:
		return nil
	}
}

func (s *Driver) OnUnsafeL2Payload(ctx context.Context, payload *eth.ExecutionPayload) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.unsafeL2Payloads <- payload:
		return nil
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// l1WindowBufEnd returns the last block that should be used as `base` to L1ChainWindow.
// This is either the last block of the window, or the L1 base block if the window is not populated.
func (s *state) l1WindowBufEnd() eth.BlockID {
	if len(s.l1WindowBuf) == 0 {
		return s.l2SafeHead.L1Origin
	}
	return s.l1WindowBuf[len(s.l1WindowBuf)-1]
}

func (s *state) handleNewL1Block(ctx context.Context, newL1Head eth.L1BlockRef) error {
	// We don't need to do anything if the head hasn't changed.
	if s.l1Head.Hash == newL1Head.Hash {
		s.log.Trace("Received L1 head signal that is the same as the current head", "l1Head", newL1Head)
		return nil
	}

	// We got a new L1 block whose parent hash is the same as the current L1 head. Means we're
	// dealing with a linear extension (new block is the immediate child of the old one). We
	// handle this by simply adding the new block to the window of blocks that we're considering
	// when extending the L2 chain.
	if s.l1Head.Hash == newL1Head.ParentHash {
		s.log.Trace("Linear extension", "l1Head", newL1Head)
		s.l1Head = newL1Head
		if s.l1WindowBufEnd().Hash == newL1Head.ParentHash {
			s.l1WindowBuf = append(s.l1WindowBuf, newL1Head.ID())
		}
		return nil
	}

	// New L1 block is not the same as the current head or a single step linear extension.
	// This could either be a long L1 extension, or a reorg. Both can be handled the same way.
	s.log.Warn("L1 Head signal indicates an L1 re-org", "old_l1_head", s.l1Head, "new_l1_head_parent", newL1Head.ParentHash, "new_l1_head", newL1Head)
	unsafeL2Head, safeL2Head, err := sync.FindL2Heads(ctx, s.l2Head, s.Config.SeqWindowSize, s.l1, s.l2, &s.Config.Genesis)
	if err != nil {
		s.log.Error("Could not get new unsafe L2 head when trying to handle a re-org", "err", err)
		return err
	}
	// Update forkchoice
	fc := eth.ForkchoiceState{
		HeadBlockHash:      unsafeL2Head.Hash,
		SafeBlockHash:      safeL2Head.Hash,
		FinalizedBlockHash: s.l2Finalized.Hash,
	}
	_, err = s.l2.ForkchoiceUpdate(ctx, &fc, nil)
	if err != nil {
		s.log.Error("Could not set new forkchoice when trying to handle a re-org", "err", err)
		return err
	}
	// State Update
	s.l1Head = newL1Head
	s.l1WindowBuf = nil
	s.l2Head = unsafeL2Head
	// Don't advance l2SafeHead past it's current value
	if s.l2SafeHead.Number >= safeL2Head.Number {
		s.l2SafeHead = safeL2Head
	}

	return nil
=======
func (s *state) handleNewL1Block(newL1Head eth.L1BlockRef) {
=======
func (s *state) handleNewL1HeadBlock(head eth.L1BlockRef) {
>>>>>>> v0.5.24
	// We don't need to do anything if the head hasn't changed.
	if s.l1Head == (eth.L1BlockRef{}) {
		s.log.Info("Received first L1 head signal", "l1_head", head)
	} else if s.l1Head.Hash == head.Hash {
		s.log.Trace("Received L1 head signal that is the same as the current head", "l1_head", head)
	} else if s.l1Head.Hash == head.ParentHash {
		// We got a new L1 block whose parent hash is the same as the current L1 head. Means we're
		// dealing with a linear extension (new block is the immediate child of the old one).
		s.log.Debug("L1 head moved forward", "l1_head", head)
	} else {
		if s.l1Head.Number >= head.Number {
			s.metrics.RecordL1ReorgDepth(s.l1Head.Number - head.Number)
		}
		// New L1 block is not the same as the current head or a single step linear extension.
		// This could either be a long L1 extension, or a reorg, or we simply missed a head update.
		s.log.Warn("L1 head signal indicates a possible L1 re-org", "old_l1_head", s.l1Head, "new_l1_head_parent", head.ParentHash, "new_l1_head", head)
	}
<<<<<<< HEAD
	s.l1Head = newL1Head
>>>>>>> v0.5.23
=======
	s.snapshot("New L1 Head")
	s.metrics.RecordL1Ref("l1_head", head)
	s.l1Head = head
}

func (s *state) handleNewL1SafeBlock(safe eth.L1BlockRef) {
	s.log.Info("New L1 safe block", "l1_safe", safe)
	s.metrics.RecordL1Ref("l1_safe", safe)
	s.l1Safe = safe
}

func (s *state) handleNewL1FinalizedBlock(finalized eth.L1BlockRef) {
	s.log.Info("New L1 finalized block", "l1_finalized", finalized)
	s.metrics.RecordL1Ref("l1_finalized", finalized)
	s.l1Finalized = finalized
	s.derivation.Finalize(finalized.ID())
>>>>>>> v0.5.24
}

// findL1Origin determines what the next L1 Origin should be.
// The L1 Origin is either the L2 Head's Origin, or the following L1 block
// if the next L2 block's time is greater than or equal to the L2 Head's Origin.
func (s *state) findL1Origin(ctx context.Context) (eth.L1BlockRef, error) {
	l2Head := s.derivation.UnsafeL2Head()
	// If we are at the head block, don't do a lookup.
	if l2Head.L1Origin.Hash == s.l1Head.Hash {
		return s.l1Head, nil
	}

	// Grab a reference to the current L1 origin block.
	currentOrigin, err := s.l1.L1BlockRefByHash(ctx, l2Head.L1Origin.Hash)
	if err != nil {
		return eth.L1BlockRef{}, err
	}

<<<<<<< HEAD
=======
	if currentOrigin.Number+1+s.DriverConfig.SequencerConfDepth > s.l1Head.Number {
		// TODO: we can decide to ignore confirmation depth if we would be forced
		//  to make an empty block (only deposits) by staying on the current origin.
		s.log.Info("sequencing with old origin to preserve conf depth",
			"current", currentOrigin, "current_time", currentOrigin.Time,
			"l1_head", s.l1Head, "l1_head_time", s.l1Head.Time,
			"l2_head", l2Head, "l2_head_time", l2Head.Time,
			"depth", s.DriverConfig.SequencerConfDepth)
		return currentOrigin, nil
	}

>>>>>>> v0.5.23
	// Attempt to find the next L1 origin block, where the next origin is the immediate child of
	// the current origin block.
	nextOrigin, err := s.l1.L1BlockRefByNumber(ctx, currentOrigin.Number+1)
	if err != nil {
		s.log.Error("Failed to get next origin. Falling back to current origin", "err", err)
		return currentOrigin, nil
	}

	// If the next L2 block time is greater than the next origin block's time, we can choose to
	// start building on top of the next origin. Sequencer implementation has some leeway here and
	// could decide to continue to build on top of the previous origin until the Sequencer runs out
	// of slack. For simplicity, we implement our Sequencer to always start building on the latest
	// L1 block when we can.
<<<<<<< HEAD
<<<<<<< HEAD
	// TODO: Can add confirmation depth here if we want.
=======
>>>>>>> v0.5.23
	if s.l2Head.Time+s.Config.BlockTime >= nextOrigin.Time {
=======
	if l2Head.Time+s.Config.BlockTime >= nextOrigin.Time {
>>>>>>> v0.5.24
		return nextOrigin, nil
	}

	return currentOrigin, nil
}

=======
>>>>>>> @eth-optimism/l2geth@0.5.27
// createNewL2Block builds a L2 block on top of the L2 Head (unsafe). Used by Sequencer nodes to
// construct new L2 blocks. Verifier nodes will use handleEpoch instead.
func (s *Driver) createNewL2Block(ctx context.Context) error {
	l2Head := s.derivation.UnsafeL2Head()
	l2Safe := s.derivation.SafeL2Head()
	l2Finalized := s.derivation.Finalized()

	// Figure out which L1 origin block we're going to be building on top of.
	l1Origin, err := s.l1OriginSelector.FindL1Origin(ctx, s.l1State.L1Head(), l2Head)
	if err != nil {
		s.log.Error("Error finding next L1 Origin", "err", err)
		return err
	}

	// Rollup is configured to not start producing blocks until a specific L1 block has been
	// reached. Don't produce any blocks until we're at that genesis block.
	if l1Origin.Number < s.config.Genesis.L1.Number {
		s.log.Info("Skipping block production because the next L1 Origin is behind the L1 genesis", "next", l1Origin.ID(), "genesis", s.config.Genesis.L1)
		return nil
	}

	// Should never happen. Sequencer will halt if we get into this situation somehow.
	nextL2Time := l2Head.Time + s.config.BlockTime
	if nextL2Time < l1Origin.Time {
		s.log.Error("Cannot build L2 block for time before L1 origin",
			"l2Unsafe", l2Head, "nextL2Time", nextL2Time, "l1Origin", l1Origin, "l1OriginTime", l1Origin.Time)
		return fmt.Errorf("cannot build L2 block on top %s for time %d before L1 origin %s at time %d",
			l2Head, nextL2Time, l1Origin, l1Origin.Time)
	}

	// Actually create the new block.
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	newUnsafeL2Head, payload, err := s.output.createNewBlock(ctx, s.l2Head, s.l2SafeHead.ID(), s.l2Finalized, l1Origin)
=======
	newUnsafeL2Head, payload, err := s.output.createNewBlock(ctx, s.l2Head, s.l2SafeHead.ID(), s.l2Finalized.ID(), l1Origin)
>>>>>>> v0.5.23
=======
	newUnsafeL2Head, payload, err := s.output.createNewBlock(ctx, l2Head, l2Safe.ID(), l2Finalized.ID(), l1Origin)
>>>>>>> v0.5.24
=======
	newUnsafeL2Head, payload, err := s.sequencer.CreateNewBlock(ctx, l2Head, l2Safe.ID(), l2Finalized.ID(), l1Origin)
>>>>>>> @eth-optimism/l2geth@0.5.27
	if err != nil {
		s.log.Error("Could not extend chain as sequencer", "err", err, "l2_parent", l2Head, "l1_origin", l1Origin)
		return err
	}

	// Update our L2 head block based on the new unsafe block we just generated.
<<<<<<< HEAD
	s.l2Head = newUnsafeL2Head
=======
	s.derivation.SetUnsafeHead(newUnsafeL2Head)

<<<<<<< HEAD
>>>>>>> v0.5.23
	s.log.Info("Sequenced new l2 block", "l2Head", s.l2Head, "l1Origin", s.l2Head.L1Origin, "txs", len(payload.Transactions), "time", s.l2Head.Time)
=======
	s.log.Info("Sequenced new l2 block", "l2_unsafe", newUnsafeL2Head, "l1_origin", newUnsafeL2Head.L1Origin, "txs", len(payload.Transactions), "time", newUnsafeL2Head.Time)
	s.metrics.CountSequencedTxs(len(payload.Transactions))
>>>>>>> v0.5.24

	if s.network != nil {
		if err := s.network.PublishL2Payload(ctx, payload); err != nil {
			s.log.Warn("failed to publish newly created block", "id", payload.ID(), "err", err)
			s.metrics.RecordPublishingError()
			// publishing of unsafe data via p2p is optional. Errors are not severe enough to change/halt sequencing but should be logged and metered.
		}
	}

	return nil
}

<<<<<<< HEAD
// handleEpoch attempts to insert a full L2 epoch on top of the L2 Safe Head.
// It ensures that a full sequencing window is available and updates the state as needed.
func (s *state) handleEpoch(ctx context.Context) (bool, error) {
	s.log.Trace("Handling epoch", "l2Head", s.l2Head, "l2SafeHead", s.l2SafeHead)
	// Extend cached window if we do not have enough saved blocks
	// attempt to buffer up to 2x the size of a sequence window of L1 blocks, to speed up later handleEpoch calls
	if len(s.l1WindowBuf) < int(s.Config.SeqWindowSize) {
		nexts, err := s.l1.L1Range(ctx, s.l1WindowBufEnd(), 2*s.Config.SeqWindowSize)
		if err != nil {
			s.log.Error("Could not extend the cached L1 window", "err", err, "l2Head", s.l2Head, "l2SafeHead", s.l2SafeHead, "l1Head", s.l1Head, "window_end", s.l1WindowBufEnd())
			return false, err
		}
		s.l1WindowBuf = append(s.l1WindowBuf, nexts...)

	}
	// Ensure that there are enough blocks in the cached window
	if len(s.l1WindowBuf) < int(s.Config.SeqWindowSize) {
		s.log.Debug("Not enough cached blocks to run step", "cached_window_len", len(s.l1WindowBuf))
		return false, nil
	}

	// Insert the epoch
	window := s.l1WindowBuf[:s.Config.SeqWindowSize]
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	newL2Head, newL2SafeHead, reorg, err := s.output.insertEpoch(ctx, s.l2Head, s.l2SafeHead, s.l2Finalized, window)
	cancel()
	if err != nil {
		// Cannot easily check that s.l1WindowBuf[0].ParentHash == s.l2Safehead.L1Origin.Hash in this function, so if insertEpoch
		// may have found a problem with that, clear the buffer and try again later.
		s.l1WindowBuf = nil
		s.log.Error("Error in running the output step.", "err", err, "l2Head", s.l2Head, "l2SafeHead", s.l2SafeHead)
		return false, err
	}

	// State update
	s.l2Head = newL2Head
	s.l2SafeHead = newL2SafeHead
	s.l1WindowBuf = s.l1WindowBuf[1:]
	s.log.Info("Inserted a new epoch", "l2Head", s.l2Head, "l2SafeHead", s.l2SafeHead, "reorg", reorg)
	// TODO: l2Finalized
	return reorg, nil

}

func (s *state) handleUnsafeL2Payload(ctx context.Context, payload *eth.ExecutionPayload) error {
	if s.l2SafeHead.Number > uint64(payload.BlockNumber) {
		s.log.Info("ignoring unsafe L2 execution payload, already have safe payload", "id", payload.ID())
		return nil
	}

	// Note that the payload may cause reorgs. The l2SafeHead may get out of sync because of this.
	// The engine should never reorg past the finalized block hash however.
	// The engine may attempt syncing via p2p if there is a larger gap in the L2 chain.

	l2Ref, err := derive.PayloadToBlockRef(payload, &s.Config.Genesis)
	if err != nil {
		return fmt.Errorf("failed to derive L2 block ref from payload: %v", err)
	}

	if err := s.output.processBlock(ctx, s.l2Head, s.l2SafeHead.ID(), s.l2Finalized, payload); err != nil {
		return fmt.Errorf("failed to process unsafe L2 payload: %v", err)
	}

	// We successfully processed the block, so update the safe head, while leaving the safe head etc. the same.
	s.l2Head = l2Ref

	return nil
}

// loop is the event loop that responds to L1 changes and internal timers to produce L2 blocks.
func (s *state) loop() {
=======
// the eventLoop responds to L1 changes and internal timers to produce L2 blocks.
<<<<<<< HEAD
func (s *state) eventLoop() {
>>>>>>> v0.5.23
=======
func (s *Driver) eventLoop() {
>>>>>>> @eth-optimism/l2geth@0.5.27
	defer s.wg.Done()
	s.log.Info("State loop started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a ticker to produce L2 blocks at a constant rate. Ticker will only run if we're
	// running in Sequencer mode.
	var l2BlockCreationTickerCh <-chan time.Time
<<<<<<< HEAD
<<<<<<< HEAD
	if s.sequencer {
=======
	if s.DriverConfig.SequencerEnabled {
>>>>>>> v0.5.23
		l2BlockCreationTicker := time.NewTicker(time.Duration(s.Config.BlockTime) * time.Second)
=======
	if s.driverConfig.SequencerEnabled {
		l2BlockCreationTicker := time.NewTicker(time.Duration(s.config.BlockTime) * time.Second)
>>>>>>> @eth-optimism/l2geth@0.5.27
		defer l2BlockCreationTicker.Stop()
		l2BlockCreationTickerCh = l2BlockCreationTicker.C
	}

	// stepReqCh is used to request that the driver attempts to step forward by one L1 block.
	stepReqCh := make(chan struct{}, 1)

	// l2BlockCreationReqCh is used to request that the driver create a new L2 block. Only used if
	// we're running in Sequencer mode, because otherwise we'll be deriving our blocks via the
	// stepping process.
	l2BlockCreationReqCh := make(chan struct{}, 1)

	// reqL2BlockCreation requests that a block be created. Won't deadlock if the channel is full.
	reqL2BlockCreation := func() {
		select {
		case l2BlockCreationReqCh <- struct{}{}:
		// Don't deadlock if the channel is already full
		default:
		}
	}

<<<<<<< HEAD
<<<<<<< HEAD
	// reqStep requests that a driver stpe be taken. Won't deadlock if the channel is full.
	// TODO: Rename step request
=======
	// reqStep requests a derivation step to be taken. Won't deadlock if the channel is full.
>>>>>>> v0.5.23
	reqStep := func() {
=======
	// channel, nil by default (not firing), but used to schedule re-attempts with delay
	var delayedStepReq <-chan time.Time

	// keep track of consecutive failed attempts, to adjust the backoff time accordingly
	bOffStrategy := backoff.Exponential()
	stepAttempts := 0

	// step requests a derivation step to be taken. Won't deadlock if the channel is full.
	step := func() {
>>>>>>> v0.5.24
		select {
		case stepReqCh <- struct{}{}:
		// Don't deadlock if the channel is already full
		default:
		}
	}

	// reqStep requests a derivation step nicely, with a delay if this is a reattempt, or not at all if we already scheduled a reattempt.
	reqStep := func() {
		if stepAttempts > 0 {
			// if this is not the first attempt, we re-schedule with a backoff, *without blocking other events*
			if delayedStepReq == nil {
				delay := bOffStrategy.Duration(stepAttempts)
				s.log.Debug("scheduling re-attempt with delay", "attempts", stepAttempts, "delay", delay)
				delayedStepReq = time.After(delay)
			} else {
				s.log.Debug("ignoring step request, already scheduled re-attempt after previous failure", "attempts", stepAttempts)
			}
		} else {
			step()
		}
	}

	// We call reqStep right away to finish syncing to the tip of the chain if we're behind.
	// reqStep will also be triggered when the L1 head moves forward or if there was a reorg on the
	// L1 chain that we need to handle.
	reqStep()

	for {
		select {
		case <-l2BlockCreationTickerCh:
			s.log.Trace("L2 Creation Ticker")
			s.snapshot("L2 Creation Ticker")
			reqL2BlockCreation()

		case <-l2BlockCreationReqCh:
			s.snapshot("L2 Block Creation Request")
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
			l1Head := s.l1State.L1Head()
>>>>>>> @eth-optimism/l2geth@0.5.27
			if !s.idleDerivation {
				s.log.Warn("not creating block, node is deriving new l2 data", "head_l1", l1Head)
				break
			}
<<<<<<< HEAD
>>>>>>> v0.5.23
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
=======
			ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
>>>>>>> @eth-optimism/l2geth@0.5.27
			err := s.createNewL2Block(ctx)
			cancel()
			if err != nil {
				s.log.Error("Error creating new L2 block", "err", err)
				s.metrics.RecordSequencingError()
				break // if we fail, we wait for the next block creation trigger.
			}

			// We need to catch up to the next origin as quickly as possible. We can do this by
			// requesting a new block ASAP instead of waiting for the next tick.
<<<<<<< HEAD
			// TODO: If we want to consider confirmations, need to consider here too.
			if s.l1Head.Number > s.l2Head.L1Origin.Number {
=======
			// We don't request a block if the confirmation depth is not met.
<<<<<<< HEAD
			if s.l1Head.Number > s.l2Head.L1Origin.Number+s.DriverConfig.SequencerConfDepth {
>>>>>>> v0.5.23
				s.log.Trace("Asking for a second L2 block asap", "l2Head", s.l2Head)
=======
			l2Head := s.derivation.UnsafeL2Head()
<<<<<<< HEAD
			if s.l1Head.Number > l2Head.L1Origin.Number+s.DriverConfig.SequencerConfDepth {
				s.log.Trace("Building another L2 block asap to catch up with L1 head", "l2_unsafe", l2Head, "l2_unsafe_l1_origin", l2Head.L1Origin, "l1_head", s.l1Head)
>>>>>>> v0.5.24
=======
			if l1Head.Number > l2Head.L1Origin.Number+s.driverConfig.SequencerConfDepth {
				s.log.Trace("Building another L2 block asap to catch up with L1 head", "l2_unsafe", l2Head, "l2_unsafe_l1_origin", l2Head.L1Origin, "l1_head", l1Head)
>>>>>>> @eth-optimism/l2geth@0.5.27
				// But not too quickly to minimize busy-waiting for new blocks
				time.AfterFunc(time.Millisecond*10, reqL2BlockCreation)
			}

		case payload := <-s.unsafeL2Payloads:
<<<<<<< HEAD
			s.log.Info("Optimistically processing unsafe L2 execution payload", "id", payload.ID())
			err := s.handleUnsafeL2Payload(ctx, payload)
			if err != nil {
				s.log.Warn("Failed to process L2 execution payload received from p2p", "err", err)
			}

		case newL1Head := <-s.l1Heads:
			s.snapshot("New L1 Head")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := s.handleNewL1Block(ctx, newL1Head)
			cancel()
			if err != nil {
				s.log.Error("Error in handling new L1 Head", "err", err)
			}

			// The block number of the L1 origin for the L2 safe head is at least SeqWindowSize
			// behind the L1 head. We can therefore attempt to shift the safe head forward by at
			// least one L1 block. If the node is holding on to unsafe blocks, this may trigger a
			// reorg on L2 in the case that safe (published) data conflicts with local unsafe
			// block data.
			if s.l1Head.Number-s.l2SafeHead.L1Origin.Number >= s.Config.SeqWindowSize {
				s.log.Trace("Requesting next step", "l1Head", s.l1Head, "l2Head", s.l2Head, "l1Origin", s.l2SafeHead.L1Origin)
				reqStep()
			}

		case <-stepReqCh:
			s.snapshot("Step Request")
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			reorg, err := s.handleEpoch(ctx)
			cancel()
			if err != nil {
				s.log.Error("Error in handling epoch", "err", err)
			}

			if reorg {
				s.log.Warn("Got reorg")

				// If we're in sequencer mode and experiencing a reorg, we should request a new
				// block ASAP. Not strictly necessary but means we'll recover from the reorg much
				// faster than if we waited for the next tick.
				if s.sequencer {
					reqL2BlockCreation()
				}
			}

			// The block number of the L1 origin for the L2 safe head is at least SeqWindowSize
			// behind the L1 head. We can therefore attempt to shift the safe head forward by at
			// least one L1 block. If the node is holding on to unsafe blocks, this may trigger a
			// reorg on L2 in the case that safe (published) data conflicts with local unsafe
			// block data.
			if s.l1Head.Number-s.l2SafeHead.L1Origin.Number >= s.Config.SeqWindowSize {
				s.log.Trace("Requesting next step", "l1Head", s.l1Head, "l2Head", s.l2Head, "l1Origin", s.l2SafeHead.L1Origin)
				reqStep()
			}

=======
			s.snapshot("New unsafe payload")
			s.log.Info("Optimistically queueing unsafe L2 execution payload", "id", payload.ID())
			s.derivation.AddUnsafePayload(payload)
			s.metrics.RecordReceivedUnsafePayload(payload)
			reqStep()

		case newL1Head := <-s.l1HeadSig:
			s.l1State.HandleNewL1HeadBlock(newL1Head)
			reqStep() // a new L1 head may mean we have the data to not get an EOF again.
		case newL1Safe := <-s.l1SafeSig:
			s.l1State.HandleNewL1SafeBlock(newL1Safe)
			// no step, justified L1 information does not do anything for L2 derivation or status
		case newL1Finalized := <-s.l1FinalizedSig:
			s.l1State.HandleNewL1FinalizedBlock(newL1Finalized)
			s.derivation.Finalize(newL1Finalized.ID())
			reqStep() // we may be able to mark more L2 data as finalized now
		case <-delayedStepReq:
			delayedStepReq = nil
			step()
		case <-stepReqCh:
			s.metrics.SetDerivationIdle(false)
			s.idleDerivation = false
			s.log.Debug("Derivation process step", "onto_origin", s.derivation.Origin(), "attempts", stepAttempts)
			err := s.derivation.Step(context.Background())
			stepAttempts += 1 // count as attempt by default. We reset to 0 if we are making healthy progress.
			if err == io.EOF {
				s.log.Debug("Derivation process went idle", "progress", s.derivation.Origin())
				s.idleDerivation = true
				stepAttempts = 0
				s.metrics.SetDerivationIdle(true)
				continue
			} else if err != nil && errors.Is(err, derive.ErrReset) {
				// If the pipeline corrupts, e.g. due to a reorg, simply reset it
				s.log.Warn("Derivation pipeline is reset", "err", err)
				s.derivation.Reset()
				s.metrics.RecordPipelineReset()
				continue
			} else if err != nil && errors.Is(err, derive.ErrTemporary) {
				s.log.Warn("Derivation process temporary error", "attempts", stepAttempts, "err", err)
				reqStep()
				continue
			} else if err != nil && errors.Is(err, derive.ErrCritical) {
				s.log.Error("Derivation process critical error", "err", err)
				return
			} else if err != nil && errors.Is(err, derive.NotEnoughData) {
				stepAttempts = 0 // don't do a backoff for this error
				reqStep()
				continue
			} else if err != nil {
				s.log.Error("Derivation process error", "attempts", stepAttempts, "err", err)
				reqStep()
				continue
			} else {
				stepAttempts = 0
				reqStep() // continue with the next step if we can
			}
		case respCh := <-s.syncStatusReq:
			respCh <- eth.SyncStatus{
				CurrentL1:   s.derivation.Origin(),
				HeadL1:      s.l1State.L1Head(),
				SafeL1:      s.l1State.L1Safe(),
				FinalizedL1: s.l1State.L1Finalized(),
				UnsafeL2:    s.derivation.UnsafeL2Head(),
				SafeL2:      s.derivation.SafeL2Head(),
				FinalizedL2: s.derivation.Finalized(),
			}
<<<<<<< HEAD
>>>>>>> v0.5.23
=======
		case respCh := <-s.forceReset:
			s.log.Warn("Derivation pipeline is manually reset")
			s.derivation.Reset()
			s.metrics.RecordPipelineReset()
			close(respCh)
>>>>>>> v0.5.24
		case <-s.done:
			return
		}
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
func (s *state) snapshot(event string) {
	l1HeadJSON, _ := json.Marshal(s.l1Head)
	l2HeadJSON, _ := json.Marshal(s.l2Head)
	l2SafeHeadJSON, _ := json.Marshal(s.l2SafeHead)
	l2FinalizedHeadJSON, _ := json.Marshal(s.l2Finalized)
	l1WindowBufJSON, _ := json.Marshal(s.l1WindowBuf)

	s.snapshotLog.Info("Rollup State Snapshot",
		"event", event,
		"l1Head", string(l1HeadJSON),
		"l2Head", string(l2HeadJSON),
		"l2SafeHead", string(l2SafeHeadJSON),
		"l2FinalizedHead", string(l2FinalizedHeadJSON),
		"l1WindowBuf", string(l1WindowBufJSON))
=======
func (s *state) SyncStatus(ctx context.Context) (*SyncStatus, error) {
	respCh := make(chan SyncStatus)
=======
// ResetDerivationPipeline forces a reset of the derivation pipeline.
// It waits for the reset to occur. It simply unblocks the caller rather
// than fully cancelling the reset request upon a context cancellation.
func (s *Driver) ResetDerivationPipeline(ctx context.Context) error {
	respCh := make(chan struct{}, 1)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.forceReset <- respCh:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-respCh:
			return nil
		}
	}
}

<<<<<<< HEAD
func (s *state) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	respCh := make(chan eth.SyncStatus)
>>>>>>> v0.5.24
=======
func (s *Driver) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	respCh := make(chan eth.SyncStatus, 1)
>>>>>>> @eth-optimism/l2geth@0.5.27
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.syncStatusReq <- respCh:
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case resp := <-respCh:
			return &resp, nil
		}
	}
}

// deferJSONString helps avoid a JSON-encoding performance hit if the snapshot logger does not run
type deferJSONString struct {
	x any
}

func (v deferJSONString) String() string {
	out, _ := json.Marshal(v.x)
	return string(out)
}

func (s *Driver) snapshot(event string) {
	s.snapshotLog.Info("Rollup State Snapshot",
		"event", event,
<<<<<<< HEAD
		"l1Head", deferJSONString{s.l1Head},
		"l1Current", deferJSONString{s.derivation.Progress().Origin},
<<<<<<< HEAD
		"l2Head", deferJSONString{s.l2Head},
		"l2SafeHead", deferJSONString{s.l2SafeHead},
		"l2FinalizedHead", deferJSONString{s.l2Finalized})
>>>>>>> v0.5.23
=======
=======
		"l1Head", deferJSONString{s.l1State.L1Head()},
		"l1Current", deferJSONString{s.derivation.Origin()},
>>>>>>> @eth-optimism/l2geth@0.5.27
		"l2Head", deferJSONString{s.derivation.UnsafeL2Head()},
		"l2Safe", deferJSONString{s.derivation.SafeL2Head()},
		"l2FinalizedHead", deferJSONString{s.derivation.Finalized()})
>>>>>>> v0.5.24
}
