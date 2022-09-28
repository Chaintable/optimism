package op_batcher

import (
	"time"

	"github.com/urfave/cli"

	"github.com/ethereum-optimism/optimism/op-batcher/flags"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	oppprof "github.com/ethereum-optimism/optimism/op-service/pprof"
	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
)

type Config struct {
	/* Required Params */

	// L1EthRpc is the HTTP provider URL for L1.
	L1EthRpc string

<<<<<<< HEAD
	// L2EthRpc is the HTTP provider URL for L2.
	L2EthRpc string

	// RollupRpc is the HTTP provider URL for the rollup node.
=======
	// L2EthRpc is the HTTP provider URL for the L2 execution engine.
	L2EthRpc string

	// RollupRpc is the HTTP provider URL for the L2 rollup node.
>>>>>>> v0.5.23
	RollupRpc string

	// MinL1TxSize is the minimum size of a batch tx submitted to L1.
	MinL1TxSize uint64

	// MaxL1TxSize is the maximum size of a batch tx submitted to L1.
	MaxL1TxSize uint64

<<<<<<< HEAD
=======
	// ChannelTimeout is the maximum amount of time to attempt completing an opened channel,
	// as opposed to submitting missing blocks in new channels
	ChannelTimeout uint64

>>>>>>> v0.5.23
	// PollInterval is the delay between querying L2 for more transaction
	// and creating a new batch.
	PollInterval time.Duration

	// NumConfirmations is the number of confirmations which we will wait after
	// appending new batches.
	NumConfirmations uint64

	// SafeAbortNonceTooLowCount is the number of ErrNonceTooLowObservations
	// required to give up on a tx at a particular nonce without receiving
	// confirmation.
	SafeAbortNonceTooLowCount uint64

	// ResubmissionTimeout is time we will wait before resubmitting a
	// transaction.
	ResubmissionTimeout time.Duration

	// Mnemonic is the HD seed used to derive the wallet private keys for both
	// the sequence and proposer. Must be used in conjunction with
	// SequencerHDPath and ProposerHDPath.
	Mnemonic string

	// SequencerHDPath is the derivation path used to obtain the private key for
	// batched submission of sequencer transactions.
	SequencerHDPath string

<<<<<<< HEAD
<<<<<<< HEAD
	// SequencerHistoryDBFilename is the filename of the database used to track
	// the latest L2 sequencer batches that were published.
	SequencerHistoryDBFilename string

	// SequencerGenesisHash is the genesis hash of the L2 chain.
	SequencerGenesisHash string

=======
>>>>>>> v0.5.23
=======
	// PrivateKey is the private key used to submit sequencer transactions.
	PrivateKey string

>>>>>>> v0.5.24
	// SequencerBatchInboxAddress is the address in which to send batch
	// transactions.
	SequencerBatchInboxAddress string

	RPCConfig oprpc.CLIConfig

	/* Optional Params */

	LogConfig oplog.CLIConfig

	MetricsConfig opmetrics.CLIConfig

	PprofConfig oppprof.CLIConfig
}

func (c Config) Check() error {
	if err := c.RPCConfig.Check(); err != nil {
		return err
	}
	if err := c.LogConfig.Check(); err != nil {
		return err
	}
	if err := c.MetricsConfig.Check(); err != nil {
		return err
	}
	if err := c.PprofConfig.Check(); err != nil {
		return err
	}
	return nil
}

// NewConfig parses the Config from the provided flags or environment variables.
func NewConfig(ctx *cli.Context) Config {
	return Config{
		/* Required Flags */
		L1EthRpc:                   ctx.GlobalString(flags.L1EthRpcFlag.Name),
		L2EthRpc:                   ctx.GlobalString(flags.L2EthRpcFlag.Name),
		RollupRpc:                  ctx.GlobalString(flags.RollupRpcFlag.Name),
		MinL1TxSize:                ctx.GlobalUint64(flags.MinL1TxSizeBytesFlag.Name),
		MaxL1TxSize:                ctx.GlobalUint64(flags.MaxL1TxSizeBytesFlag.Name),
<<<<<<< HEAD
=======
		ChannelTimeout:             ctx.GlobalUint64(flags.ChannelTimeoutFlag.Name),
>>>>>>> v0.5.23
		PollInterval:               ctx.GlobalDuration(flags.PollIntervalFlag.Name),
		NumConfirmations:           ctx.GlobalUint64(flags.NumConfirmationsFlag.Name),
		SafeAbortNonceTooLowCount:  ctx.GlobalUint64(flags.SafeAbortNonceTooLowCountFlag.Name),
		ResubmissionTimeout:        ctx.GlobalDuration(flags.ResubmissionTimeoutFlag.Name),
		Mnemonic:                   ctx.GlobalString(flags.MnemonicFlag.Name),
		SequencerHDPath:            ctx.GlobalString(flags.SequencerHDPathFlag.Name),
<<<<<<< HEAD
<<<<<<< HEAD
		SequencerHistoryDBFilename: ctx.GlobalString(flags.SequencerHistoryDBFilenameFlag.Name),
		SequencerGenesisHash:       ctx.GlobalString(flags.SequencerGenesisHashFlag.Name),
=======
>>>>>>> v0.5.23
=======
		PrivateKey:                 ctx.GlobalString(flags.PrivateKeyFlag.Name),
>>>>>>> v0.5.24
		SequencerBatchInboxAddress: ctx.GlobalString(flags.SequencerBatchInboxAddressFlag.Name),
		RPCConfig:                  oprpc.ReadCLIConfig(ctx),
		LogConfig:                  oplog.ReadCLIConfig(ctx),
		MetricsConfig:              opmetrics.ReadCLIConfig(ctx),
		PprofConfig:                oppprof.ReadCLIConfig(ctx),
	}
}
