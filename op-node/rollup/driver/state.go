package driver

import (
	"context"
	"encoding/json"
	"fmt"
<<<<<<< HEAD
=======
	"io"
>>>>>>> v0.5.23
	gosync "sync"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
<<<<<<< HEAD
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
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
	"github.com/ethereum/go-ethereum/log"
)

// SyncStatus is a snapshot of the driver
type SyncStatus struct {
	// CurrentL1 is the block that the derivation process is currently at,
	// this may not be fully derived into L2 data yet.
	// If the node is synced, this matches the HeadL1, minus the verifier confirmation distance.
	CurrentL1 eth.L1BlockRef `json:"current_l1"`
	// HeadL1 is the perceived head of the L1 chain, no confirmation distance.
	// The head is not guaranteed to build on the other L1 sync status fields,
	// as the node may be in progress of resetting to adapt to a L1 reorg.
	HeadL1 eth.L1BlockRef `json:"head_l1"`
	// UnsafeL2 is the absolute tip of the L2 chain,
	// pointing to block data that has not been submitted to L1 yet.
	// The sequencer is building this, and verifiers may also be ahead of the
	// SafeL2 block if they sync blocks via p2p or other offchain sources.
	UnsafeL2 eth.L2BlockRef `json:"unsafe_l2"`
	// SafeL2 points to the L2 block that was derived from the L1 chain.
	// This point may still reorg if the L1 chain reorgs.
	SafeL2 eth.L2BlockRef `json:"safe_l2"`
	// FinalizedL2 points to the L2 block that was derived fully from
	// finalized L1 information, thus irreversible.
	FinalizedL2 eth.L2BlockRef `json:"finalized_l2"`
}

type state struct {
	// Chain State
	l1Head      eth.L1BlockRef // Latest recorded head of the L1 Chain, independent of derivation work
	l2Head      eth.L2BlockRef // L2 Unsafe Head
	l2SafeHead  eth.L2BlockRef // L2 Safe Head - this is the head of the L2 chain as derived from L1
	l2Finalized eth.L2BlockRef // L2 Block that will never be reversed

	// The derivation pipeline is reset whenever we reorg.
	// The derivation pipeline determines the new l2SafeHead.
	derivation DerivationPipeline

	// When the derivation pipeline is waiting for new data to do anything
	idleDerivation bool

	// Requests for sync status. Synchronized with event loop to avoid reading an inconsistent sync status.
	syncStatusReq chan chan SyncStatus

	// Rollup config: rollup chain configuration
	Config *rollup.Config

	// Driver config: verifier and sequencer settings
	DriverConfig *Config
>>>>>>> v0.5.23

	// Connections (in/out)
	l1Heads          chan eth.L1BlockRef
	unsafeL2Payloads chan *eth.ExecutionPayload
	l1               L1Chain
	l2               L2Chain
	output           outputInterface
	network          Network // may be nil, network for is optional

	log         log.Logger
	snapshotLog log.Logger
	done        chan struct{}

	wg gosync.WaitGroup
}

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
	output outputInterface, derivationPipeline DerivationPipeline, network Network) *state {
	return &state{
		derivation:       derivationPipeline,
		idleDerivation:   false,
		syncStatusReq:    make(chan chan SyncStatus, 10),
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
		sequencer:        sequencer,
=======
>>>>>>> v0.5.23
		l1Heads:          make(chan eth.L1BlockRef, 10),
		unsafeL2Payloads: make(chan *eth.ExecutionPayload, 10),
	}
}

// Start starts up the state loop. The context is only for initialization.
// The loop will have been started iff err is not nil.
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

	s.derivation.Reset()

	s.wg.Add(1)
	go s.eventLoop()

>>>>>>> v0.5.23
	return nil
}

func (s *state) Close() error {
<<<<<<< HEAD
	close(s.done)
=======
	s.done <- struct{}{}
>>>>>>> v0.5.23
	s.wg.Wait()
	return nil
}

func (s *state) OnL1Head(ctx context.Context, head eth.L1BlockRef) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.l1Heads <- head:
		return nil
	}
}

func (s *state) OnUnsafeL2Payload(ctx context.Context, payload *eth.ExecutionPayload) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.unsafeL2Payloads <- payload:
		return nil
	}
}

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
	// We don't need to do anything if the head hasn't changed.
	if s.l1Head.Hash == newL1Head.Hash {
		s.log.Trace("Received L1 head signal that is the same as the current head", "l1Head", newL1Head)
	} else if s.l1Head.Hash == newL1Head.ParentHash {
		// We got a new L1 block whose parent hash is the same as the current L1 head. Means we're
		// dealing with a linear extension (new block is the immediate child of the old one).
		s.log.Debug("L1 head moved forward", "l1Head", newL1Head)
	} else {
		// New L1 block is not the same as the current head or a single step linear extension.
		// This could either be a long L1 extension, or a reorg. Both can be handled the same way.
		s.log.Warn("L1 Head signal indicates an L1 re-org", "old_l1_head", s.l1Head, "new_l1_head_parent", newL1Head.ParentHash, "new_l1_head", newL1Head)
	}
	s.l1Head = newL1Head
>>>>>>> v0.5.23
}

// findL1Origin determines what the next L1 Origin should be.
// The L1 Origin is either the L2 Head's Origin, or the following L1 block
// if the next L2 block's time is greater than or equal to the L2 Head's Origin.
func (s *state) findL1Origin(ctx context.Context) (eth.L1BlockRef, error) {
	// If we are at the head block, don't do a lookup.
	if s.l2Head.L1Origin.Hash == s.l1Head.Hash {
		return s.l1Head, nil
	}

	// Grab a reference to the current L1 origin block.
	currentOrigin, err := s.l1.L1BlockRefByHash(ctx, s.l2Head.L1Origin.Hash)
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
			"l2_head", s.l2Head, "l2_head_time", s.l2Head.Time,
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
	// TODO: Can add confirmation depth here if we want.
=======
>>>>>>> v0.5.23
	if s.l2Head.Time+s.Config.BlockTime >= nextOrigin.Time {
		return nextOrigin, nil
	}

	return currentOrigin, nil
}

// createNewL2Block builds a L2 block on top of the L2 Head (unsafe). Used by Sequencer nodes to
// construct new L2 blocks. Verifier nodes will use handleEpoch instead.
func (s *state) createNewL2Block(ctx context.Context) error {
	// Figure out which L1 origin block we're going to be building on top of.
	l1Origin, err := s.findL1Origin(ctx)
	if err != nil {
		s.log.Error("Error finding next L1 Origin", "err", err)
		return err
	}

	// Rollup is configured to not start producing blocks until a specific L1 block has been
	// reached. Don't produce any blocks until we're at that genesis block.
	if l1Origin.Number < s.Config.Genesis.L1.Number {
		s.log.Info("Skipping block production because the next L1 Origin is behind the L1 genesis", "next", l1Origin.ID(), "genesis", s.Config.Genesis.L1)
		return nil
	}

	// Should never happen. Sequencer will halt if we get into this situation somehow.
	nextL2Time := s.l2Head.Time + s.Config.BlockTime
	if nextL2Time < l1Origin.Time {
		s.log.Error("Cannot build L2 block for time before L1 origin",
			"l2Head", s.l2Head, "nextL2Time", nextL2Time, "l1Origin", l1Origin, "l1OriginTime", l1Origin.Time)
		return fmt.Errorf("cannot build L2 block on top %s for time %d before L1 origin %s at time %d",
			s.l2Head, nextL2Time, l1Origin, l1Origin.Time)
	}

	// Actually create the new block.
<<<<<<< HEAD
	newUnsafeL2Head, payload, err := s.output.createNewBlock(ctx, s.l2Head, s.l2SafeHead.ID(), s.l2Finalized, l1Origin)
=======
	newUnsafeL2Head, payload, err := s.output.createNewBlock(ctx, s.l2Head, s.l2SafeHead.ID(), s.l2Finalized.ID(), l1Origin)
>>>>>>> v0.5.23
	if err != nil {
		s.log.Error("Could not extend chain as sequencer", "err", err, "l2UnsafeHead", s.l2Head, "l1Origin", l1Origin)
		return err
	}

	// Update our L2 head block based on the new unsafe block we just generated.
<<<<<<< HEAD
	s.l2Head = newUnsafeL2Head
=======
	s.derivation.SetUnsafeHead(newUnsafeL2Head)
	s.l2Head = newUnsafeL2Head

>>>>>>> v0.5.23
	s.log.Info("Sequenced new l2 block", "l2Head", s.l2Head, "l1Origin", s.l2Head.L1Origin, "txs", len(payload.Transactions), "time", s.l2Head.Time)

	if s.network != nil {
		if err := s.network.PublishL2Payload(ctx, payload); err != nil {
			s.log.Warn("failed to publish newly created block", "id", payload.ID(), "err", err)
			return err
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
func (s *state) eventLoop() {
>>>>>>> v0.5.23
	defer s.wg.Done()
	s.log.Info("State loop started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start a ticker to produce L2 blocks at a constant rate. Ticker will only run if we're
	// running in Sequencer mode.
	var l2BlockCreationTickerCh <-chan time.Time
<<<<<<< HEAD
	if s.sequencer {
=======
	if s.DriverConfig.SequencerEnabled {
>>>>>>> v0.5.23
		l2BlockCreationTicker := time.NewTicker(time.Duration(s.Config.BlockTime) * time.Second)
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
	// reqStep requests that a driver stpe be taken. Won't deadlock if the channel is full.
	// TODO: Rename step request
=======
	// reqStep requests a derivation step to be taken. Won't deadlock if the channel is full.
>>>>>>> v0.5.23
	reqStep := func() {
		select {
		case stepReqCh <- struct{}{}:
		// Don't deadlock if the channel is already full
		default:
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
=======
			if !s.idleDerivation {
				s.log.Warn("not creating block, node is deriving new l2 data", "head_l1", s.l1Head)
				break
			}
>>>>>>> v0.5.23
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := s.createNewL2Block(ctx)
			cancel()
			if err != nil {
				s.log.Error("Error creating new L2 block", "err", err)
			}

			// We need to catch up to the next origin as quickly as possible. We can do this by
			// requesting a new block ASAP instead of waiting for the next tick.
<<<<<<< HEAD
			// TODO: If we want to consider confirmations, need to consider here too.
			if s.l1Head.Number > s.l2Head.L1Origin.Number {
=======
			// We don't request a block if the confirmation depth is not met.
			if s.l1Head.Number > s.l2Head.L1Origin.Number+s.DriverConfig.SequencerConfDepth {
>>>>>>> v0.5.23
				s.log.Trace("Asking for a second L2 block asap", "l2Head", s.l2Head)
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
			reqStep()

		case newL1Head := <-s.l1Heads:
			s.log.Info("new l1 Head")
			s.snapshot("New L1 Head")
			s.handleNewL1Block(newL1Head)
			reqStep() // a new L1 head may mean we have the data to not get an EOF again.
		case <-stepReqCh:
			s.idleDerivation = false
			s.log.Debug("Derivation process step", "onto_origin", s.derivation.Progress().Origin, "onto_closed", s.derivation.Progress().Closed)
			stepCtx, cancel := context.WithTimeout(ctx, time.Second*10) // TODO pick a timeout for executing a single step
			err := s.derivation.Step(stepCtx)
			cancel()
			if err == io.EOF {
				s.log.Debug("Derivation process went idle", "progress", s.derivation.Progress().Origin)
				s.idleDerivation = true
				continue
			} else if err != nil {
				// If the pipeline corrupts, e.g. due to a reorg, simply reset it
				s.log.Warn("Derivation pipeline is reset", "err", err)
				s.derivation.Reset()
			} else {
				finalized, safe, unsafe := s.derivation.Finalized(), s.derivation.SafeL2Head(), s.derivation.UnsafeL2Head()
				// log sync progress when it changes
				if s.l2Finalized != finalized || s.l2SafeHead != safe || s.l2Head != unsafe {
					s.log.Info("Sync progress", "finalized", finalized, "safe", safe, "unsafe", unsafe)
				}
				// update the heads
				s.l2Finalized = finalized
				s.l2SafeHead = safe
				s.l2Head = unsafe
				reqStep() // continue with the next step if we can
			}
		case respCh := <-s.syncStatusReq:
			respCh <- SyncStatus{
				CurrentL1:   s.derivation.Progress().Origin,
				HeadL1:      s.l1Head,
				UnsafeL2:    s.l2Head,
				SafeL2:      s.l2SafeHead,
				FinalizedL2: s.l2Finalized,
			}
>>>>>>> v0.5.23
		case <-s.done:
			return
		}
	}
}

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

func (s *state) snapshot(event string) {
	s.snapshotLog.Info("Rollup State Snapshot",
		"event", event,
		"l1Head", deferJSONString{s.l1Head},
		"l1Current", deferJSONString{s.derivation.Progress().Origin},
		"l2Head", deferJSONString{s.l2Head},
		"l2SafeHead", deferJSONString{s.l2SafeHead},
		"l2FinalizedHead", deferJSONString{s.l2Finalized})
>>>>>>> v0.5.23
}
