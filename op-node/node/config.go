package node

import (
<<<<<<< HEAD
	"fmt"

	"github.com/ethereum-optimism/optimism/op-node/p2p"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
=======
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/driver"
>>>>>>> v0.5.23
)

type Config struct {
	L1 L1EndpointSetup
	L2 L2EndpointSetup

<<<<<<< HEAD
	Rollup rollup.Config

	// Sequencer flag, enables sequencing
	Sequencer bool
=======
	Driver driver.Config

	Rollup rollup.Config
>>>>>>> v0.5.23

	// P2PSigner will be used for signing off on published content
	// if the node is sequencing and if the p2p stack is enabled
	P2PSigner p2p.SignerSetup

	RPC RPCConfig

	P2P p2p.SetupP2P

<<<<<<< HEAD
=======
	Metrics MetricsConfig

<<<<<<< HEAD
>>>>>>> v0.5.23
=======
	Pprof PprofConfig

	// Used to poll the L1 for new finalized or safe blocks
	L1EpochPollInterval time.Duration

>>>>>>> v0.5.24
	// Optional
	Tracer    Tracer
	Heartbeat HeartbeatConfig
}

type RPCConfig struct {
	ListenAddr  string
	ListenPort  int
	EnableAdmin bool
}

func (cfg *RPCConfig) HttpEndpoint() string {
	return fmt.Sprintf("http://%s:%d", cfg.ListenAddr, cfg.ListenPort)
}

<<<<<<< HEAD
// Check verifies that the given configuration makes sense
func (cfg *Config) Check() error {
	if err := cfg.L2.Check(); err != nil {
		return fmt.Errorf("l2 endpoint config error: %v", err)
	}
	if err := cfg.Rollup.Check(); err != nil {
		return fmt.Errorf("rollup config error: %v", err)
	}
	if cfg.P2P != nil {
		if err := cfg.P2P.Check(); err != nil {
			return fmt.Errorf("p2p config error: %v", err)
=======
type MetricsConfig struct {
	Enabled    bool
	ListenAddr string
	ListenPort int
}

func (m MetricsConfig) Check() error {
	if !m.Enabled {
		return nil
	}

	if m.ListenPort < 0 || m.ListenPort > math.MaxUint16 {
		return errors.New("invalid metrics port")
	}

	return nil
}

type PprofConfig struct {
	Enabled    bool
	ListenAddr string
	ListenPort string
}

func (p PprofConfig) Check() error {
	return nil
}

type HeartbeatConfig struct {
	Enabled bool
	Moniker string
	URL     string
}

// Check verifies that the given configuration makes sense
func (cfg *Config) Check() error {
	if err := cfg.L2.Check(); err != nil {
		return fmt.Errorf("l2 endpoint config error: %w", err)
	}
	if err := cfg.Rollup.Check(); err != nil {
		return fmt.Errorf("rollup config error: %w", err)
	}
	if err := cfg.Metrics.Check(); err != nil {
		return fmt.Errorf("metrics config error: %w", err)
	}
	if err := cfg.Pprof.Check(); err != nil {
		return fmt.Errorf("pprof config error: %w", err)
	}
	if cfg.P2P != nil {
		if err := cfg.P2P.Check(); err != nil {
			return fmt.Errorf("p2p config error: %w", err)
>>>>>>> v0.5.23
		}
	}
	return nil
}
