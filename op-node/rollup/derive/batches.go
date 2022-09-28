package derive

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"bytes"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum/go-ethereum/core/types"
)

func BatchesFromEVMTransactions(config *rollup.Config, txLists []types.Transactions) ([]*BatchData, []error) {
	var out []*BatchData
	var errs []error
	l1Signer := config.L1Signer()
	for i, txs := range txLists {
		for j, tx := range txs {
			if to := tx.To(); to != nil && *to == config.BatchInboxAddress {
				seqDataSubmitter, err := l1Signer.Sender(tx) // optimization: only derive sender if To is correct
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid signature: tx list: %d, tx: %d, err: %w", i, j, err))
					continue // bad signature, ignore
				}
				// some random L1 user might have sent a transaction to our batch inbox, ignore them
				if seqDataSubmitter != config.BatchSenderAddress {
					errs = append(errs, fmt.Errorf("unauthorized batch submitter: tx list: %d, tx: %d", i, j))
					continue // not an authorized batch submitter, ignore
				}
				batches, err := DecodeBatches(config, bytes.NewReader(tx.Data()))
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid batch: tx list: %d, tx: %d, err: %w", i, j, err))
					continue
				}
				out = append(out, batches...)
			}
		}
	}
	return out, errs
}

func FilterBatches(config *rollup.Config, epoch rollup.Epoch, minL2Time uint64, maxL2Time uint64, batches []*BatchData) (out []*BatchData) {
	uniqueTime := make(map[uint64]struct{})
	for _, batch := range batches {
		if !ValidBatch(batch, config, epoch, minL2Time, maxL2Time) {
=======
	"errors"
	"fmt"

=======
>>>>>>> v0.5.24
	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

<<<<<<< HEAD
var DifferentEpoch = errors.New("batch is of different epoch")

func FilterBatches(log log.Logger, config *rollup.Config, epoch eth.BlockID, minL2Time uint64, maxL2Time uint64, batches []*BatchData) (out []*BatchData) {
	uniqueTime := make(map[uint64]struct{})
	for _, batch := range batches {
		if err := ValidBatch(batch, config, epoch, minL2Time, maxL2Time); err != nil {
			if err == DifferentEpoch {
				log.Trace("ignoring batch of different epoch", "expected_epoch", epoch,
					"epoch", batch.Epoch(), "timestamp", batch.Timestamp, "txs", len(batch.Transactions))
			} else {
				log.Warn("filtered batch", "expected_epoch", epoch, "min", minL2Time, "max", maxL2Time,
					"epoch", batch.Epoch(), "timestamp", batch.Timestamp, "txs", len(batch.Transactions), "err", err)
			}
>>>>>>> v0.5.23
			continue
		}
		// Check if we have already seen a batch for this L2 block
		if _, ok := uniqueTime[batch.Timestamp]; ok {
<<<<<<< HEAD
=======
			log.Warn("duplicate batch", "epoch", batch.Epoch(), "timestamp", batch.Timestamp, "txs", len(batch.Transactions))
>>>>>>> v0.5.23
			// block already exists, batch is duplicate (first batch persists, others are ignored)
			continue
		}
		uniqueTime[batch.Timestamp] = struct{}{}
		out = append(out, batch)
	}
	return
}

<<<<<<< HEAD
func ValidBatch(batch *BatchData, config *rollup.Config, epoch rollup.Epoch, minL2Time uint64, maxL2Time uint64) bool {
	if batch.Epoch != epoch {
		// Batch was tagged for past or future epoch,
		// i.e. it was included too late or depends on the given L1 block to be processed first.
		return false
	}
	if (batch.Timestamp-config.Genesis.L2Time)%config.BlockTime != 0 {
		return false // bad timestamp, not a multiple of the block time
	}
	if batch.Timestamp < minL2Time {
		return false // old batch
	}
	// limit timestamp upper bound to avoid huge amount of empty blocks
	if batch.Timestamp >= maxL2Time {
		return false // too far in future
	}
	for _, txBytes := range batch.Transactions {
		if len(txBytes) == 0 {
			return false // transaction data must not be empty
		}
		if txBytes[0] == types.DepositTxType {
			return false // sequencers may not embed any deposits into batch data
		}
	}
	return true
}

// FillMissingBatches turns a collection of batches to the input batches for a series of blocks
func FillMissingBatches(batches []*BatchData, epoch, blockTime, minL2Time, nextL1Time uint64) []*BatchData {
=======
func ValidBatch(batch *BatchData, config *rollup.Config, epoch eth.BlockID, minL2Time uint64, maxL2Time uint64) error {
	if batch.EpochNum != rollup.Epoch(epoch.Number) {
		// Batch was tagged for past or future epoch,
		// i.e. it was included too late or depends on the given L1 block to be processed first.
		// This is a very common error, batches may just be buffered for a later epoch.
		return DifferentEpoch
=======
type BatchWithL1InclusionBlock struct {
	L1InclusionBlock eth.L1BlockRef
	Batch            *BatchData
}

type BatchValidity uint8

const (
	// BatchDrop indicates that the batch is invalid, and will always be in the future, unless we reorg
	BatchDrop = iota
	// BatchAccept indicates that the batch is valid and should be processed
	BatchAccept
	// BatchUndecided indicates we are lacking L1 information until we can proceed batch filtering
	BatchUndecided
	// BatchFuture indicates that the batch may be valid, but cannot be processed yet and should be checked again later
	BatchFuture
)

// CheckBatch checks if the given batch can be applied on top of the given l2SafeHead, given the contextual L1 blocks the batch was included in.
// The first entry of the l1Blocks should match the origin of the l2SafeHead. One or more consecutive l1Blocks should be provided.
// In case of only a single L1 block, the decision whether a batch is valid may have to stay undecided.
func CheckBatch(cfg *rollup.Config, log log.Logger, l1Blocks []eth.L1BlockRef, l2SafeHead eth.L2BlockRef, batch *BatchWithL1InclusionBlock) BatchValidity {
	// add details to the log
	log = log.New(
		"batch_timestamp", batch.Batch.Timestamp,
		"parent_hash", batch.Batch.ParentHash,
		"batch_epoch", batch.Batch.Epoch(),
		"txs", len(batch.Batch.Transactions),
	)

	// sanity check we have consistent inputs
	if len(l1Blocks) == 0 {
		log.Warn("missing L1 block input, cannot proceed with batch checking")
		return BatchUndecided
>>>>>>> v0.5.24
	}
	epoch := l1Blocks[0]
	if epoch.Hash != l2SafeHead.L1Origin.Hash {
		log.Warn("safe L2 head L1 origin does not match batch first l1 block (current epoch)",
			"safe_l2", l2SafeHead, "safe_origin", l2SafeHead.L1Origin, "epoch", epoch)
		return BatchUndecided
	}

	nextTimestamp := l2SafeHead.Time + cfg.BlockTime
	if batch.Batch.Timestamp > nextTimestamp {
		log.Trace("received out-of-order batch for future processing after next batch", "next_timestamp", nextTimestamp)
		return BatchFuture
	}
	if batch.Batch.Timestamp < nextTimestamp {
		log.Warn("dropping batch with old timestamp", "min_timestamp", nextTimestamp)
		return BatchDrop
	}

	// dependent on above timestamp check. If the timestamp is correct, then it must build on top of the safe head.
	if batch.Batch.ParentHash != l2SafeHead.Hash {
		log.Warn("ignoring batch with mismatching parent hash", "current_safe_head", l2SafeHead.Hash)
		return BatchDrop
	}

	// Filter out batches that were included too late.
	if uint64(batch.Batch.EpochNum)+cfg.SeqWindowSize < batch.L1InclusionBlock.Number {
		log.Warn("batch was included too late, sequence window expired")
		return BatchDrop
	}

	// Check the L1 origin of the batch
	batchOrigin := epoch
	if uint64(batch.Batch.EpochNum) < epoch.Number {
		log.Warn("dropped batch, epoch is too old", "minimum", epoch.ID())
		// batch epoch too old
		return BatchDrop
	} else if uint64(batch.Batch.EpochNum) == epoch.Number {
		// Batch is sticking to the current epoch, continue.
	} else if uint64(batch.Batch.EpochNum) == epoch.Number+1 {
		// With only 1 l1Block we cannot look at the next L1 Origin.
		// Note: This means that we are unable to determine validity of a batch
		// without more information. In this case we should bail out until we have
		// more information otherwise the eager algorithm may diverge from a non-eager
		// algorithm.
		if len(l1Blocks) < 2 {
			log.Info("eager batch wants to advance epoch, but could not without more L1 blocks", "current_epoch", epoch.ID())
			return BatchUndecided
		}
		batchOrigin = l1Blocks[1]
	} else {
		log.Warn("batch is for future epoch too far ahead, while it has the next timestamp, so it must be invalid", "current_epoch", epoch.ID())
		return BatchDrop
	}

<<<<<<< HEAD
// FillMissingBatches turns a collection of batches to the input batches for a series of blocks
func FillMissingBatches(batches []*BatchData, epoch eth.BlockID, blockTime, minL2Time, nextL1Time uint64) []*BatchData {
>>>>>>> v0.5.23
	m := make(map[uint64]*BatchData)
	// The number of L2 blocks per sequencing window is variable, we do not immediately fill to maxL2Time:
	// - ensure at least 1 block
	// - fill up to the next L1 block timestamp, if higher, to keep up with L1 time
	// - fill up to the last valid batch, to keep up with L2 time
	newHeadL2Timestamp := minL2Time
<<<<<<< HEAD
	if nextL1Time > newHeadL2Timestamp+blockTime {
		newHeadL2Timestamp = nextL1Time - blockTime
	}
	for _, b := range batches {
		m[b.BatchV1.Timestamp] = b
=======
	if nextL1Time > newHeadL2Timestamp+1 {
		newHeadL2Timestamp = nextL1Time - 1
	}
	for _, b := range batches {
		m[b.Timestamp] = b
>>>>>>> v0.5.23
		if b.Timestamp > newHeadL2Timestamp {
			newHeadL2Timestamp = b.Timestamp
		}
	}
	var out []*BatchData
	for t := minL2Time; t <= newHeadL2Timestamp; t += blockTime {
		b, ok := m[t]
		if ok {
			out = append(out, b)
		} else {
<<<<<<< HEAD
			out = append(out, &BatchData{
				BatchV1{
					Epoch:     rollup.Epoch(epoch),
					Timestamp: t,
				},
			})
		}
=======
			out = append(out,
				&BatchData{
					BatchV1{
						EpochNum:  rollup.Epoch(epoch.Number),
						EpochHash: epoch.Hash,
						Timestamp: t,
					},
				})
		}

>>>>>>> v0.5.23
=======
	if batch.Batch.EpochHash != batchOrigin.Hash {
		log.Warn("batch is for different L1 chain, epoch hash does not match", "expected", batchOrigin.ID())
		return BatchDrop
	}

	// If we ran out of sequencer time drift, then we drop the batch and produce an empty batch instead,
	// as the sequencer is not allowed to include anything past this point without moving to the next epoch.
	if max := batchOrigin.Time + cfg.MaxSequencerDrift; batch.Batch.Timestamp > max {
		log.Warn("batch exceeded sequencer time drift, sequencer must adopt new L1 origin to include transactions again", "max_time", max)
		return BatchDrop
	}

	// We can do this check earlier, but it's a more intensive one, so we do this last.
	for i, txBytes := range batch.Batch.Transactions {
		if len(txBytes) == 0 {
			log.Warn("transaction data must not be empty, but found empty tx", "tx_index", i)
			return BatchDrop
		}
		if txBytes[0] == types.DepositTxType {
			log.Warn("sequencers may not embed any deposits into batch data, but found tx that has one", "tx_index", i)
			return BatchDrop
		}
>>>>>>> v0.5.24
	}

	return BatchAccept
}
