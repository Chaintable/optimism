package op_e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
<<<<<<< HEAD
	"io/ioutil"
	"math/big"
	"os"
<<<<<<< HEAD
=======
	"math/big"
>>>>>>> v0.5.23
=======
	"path"
>>>>>>> @eth-optimism/l2geth@0.5.27
	"strings"
	"testing"
	"time"

	bss "github.com/ethereum-optimism/optimism/op-batcher"
	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-node/eth"
>>>>>>> @eth-optimism/l2geth@0.5.27
	"github.com/ethereum-optimism/optimism/op-node/metrics"
>>>>>>> v0.5.23
	rollupNode "github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/driver"
	"github.com/ethereum-optimism/optimism/op-node/testlog"
	l2os "github.com/ethereum-optimism/optimism/op-proposer"
<<<<<<< HEAD
<<<<<<< HEAD

=======
>>>>>>> v0.5.23
=======
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
<<<<<<< HEAD
>>>>>>> v0.5.24
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	geth_eth "github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"
)

var (
	testingJWTSecret = [32]byte{123}
)

func DefaultSystemConfig(t *testing.T) SystemConfig {
	secrets, err := e2eutils.DefaultMnemonicConfig.Secrets()
	require.NoError(t, err)
	addresses := secrets.Addresses()

	return SystemConfig{
		Secrets: secrets,

		Premine: make(map[common.Address]*big.Int),

		DeployConfig: &genesis.DeployConfig{
			L1ChainID:   900,
			L2ChainID:   901,
			L2BlockTime: 2,

			FinalizationPeriodSeconds: 60 * 60 * 24,
			MaxSequencerDrift:         10,
			SequencerWindowSize:       30,
			ChannelTimeout:            10,
			P2PSequencerAddress:       addresses.SequencerP2P,
			BatchInboxAddress:         common.Address{0: 0x52, 19: 0xff}, // tbd
			BatchSenderAddress:        addresses.Batcher,

			L2OutputOracleSubmissionInterval: 4,
			L2OutputOracleStartingTimestamp:  -1,
			L2OutputOracleProposer:           addresses.Proposer,
			L2OutputOracleOwner:              common.Address{}, // tbd

			L1BlockTime:                 2,
			L1GenesisBlockNonce:         4660,
			CliqueSignerAddress:         addresses.CliqueSigner,
			L1GenesisBlockTimestamp:     hexutil.Uint64(time.Now().Unix()),
			L1GenesisBlockGasLimit:      5_000_000,
			L1GenesisBlockDifficulty:    uint642big(1),
			L1GenesisBlockMixHash:       common.Hash{},
			L1GenesisBlockCoinbase:      common.Address{},
			L1GenesisBlockNumber:        0,
			L1GenesisBlockGasUsed:       0,
			L1GenesisBlockParentHash:    common.Hash{},
			L1GenesisBlockBaseFeePerGas: uint642big(7),

			L2GenesisBlockNonce:         0,
			L2GenesisBlockExtraData:     []byte{},
			L2GenesisBlockGasLimit:      5_000_000,
			L2GenesisBlockDifficulty:    uint642big(1),
			L2GenesisBlockMixHash:       common.Hash{},
			L2GenesisBlockCoinbase:      common.Address{0: 0x12},
			L2GenesisBlockNumber:        0,
			L2GenesisBlockGasUsed:       0,
			L2GenesisBlockParentHash:    common.Hash{},
			L2GenesisBlockBaseFeePerGas: uint642big(7),

			OptimismBaseFeeRecipient:    common.Address{0: 0x52, 19: 0xf0}, // tbd
			OptimismL1FeeRecipient:      common.Address{0: 0x52, 19: 0xf1},
			OptimismL2FeeRecipient:      common.Address{0: 0x52, 19: 0xf2}, // tbd
			L2CrossDomainMessengerOwner: common.Address{0: 0x52, 19: 0xf3}, // tbd
			GasPriceOracleOwner:         addresses.Alice,                   // tbd
			GasPriceOracleOverhead:      0,
			GasPriceOracleScalar:        0,
			GasPriceOracleDecimals:      0,
			DeploymentWaitConfirmations: 1,

			EIP1559Elasticity:  2,
			EIP1559Denominator: 8,

			FundDevAccounts: true,
		},
		L1InfoPredeployAddress: predeploys.L1BlockAddr,
		JWTFilePath:            writeDefaultJWT(t),
		JWTSecret:              testingJWTSecret,
		Nodes: map[string]*rollupNode.Config{
			"verifier": {
				Driver: driver.Config{
					VerifierConfDepth:  0,
					SequencerConfDepth: 0,
					SequencerEnabled:   false,
				},
				L1EpochPollInterval: time.Second * 4,
			},
			"sequencer": {
				Driver: driver.Config{
					VerifierConfDepth:  0,
					SequencerConfDepth: 0,
					SequencerEnabled:   true,
				},
				// Submitter PrivKey is set in system start for rollup nodes where sequencer = true
				RPC: rollupNode.RPCConfig{
					ListenAddr:  "127.0.0.1",
					ListenPort:  0,
					EnableAdmin: true,
				},
				L1EpochPollInterval: time.Second * 4,
			},
		},
		Loggers: map[string]log.Logger{
			"verifier":  testlog.Logger(t, log.LvlInfo).New("role", "verifier"),
			"sequencer": testlog.Logger(t, log.LvlInfo).New("role", "sequencer"),
			"batcher":   testlog.Logger(t, log.LvlInfo).New("role", "batcher"),
			"proposer":  testlog.Logger(t, log.LvlCrit).New("role", "proposer"),
		},
		P2PTopology: nil, // no P2P connectivity by default
	}
}

func writeDefaultJWT(t *testing.T) string {
	// Sadly the geth node config cannot load JWT secret from memory, it has to be a file
	jwtPath := path.Join(t.TempDir(), "jwt_secret")
	if err := os.WriteFile(jwtPath, []byte(hexutil.Encode(testingJWTSecret[:])), 0600); err != nil {
		t.Fatalf("failed to prepare jwt file for geth: %v", err)
	}
	return jwtPath
}

type L2OOContractConfig struct {
	SubmissionFrequency   *big.Int
<<<<<<< HEAD
<<<<<<< HEAD
	L2StartTime           *big.Int
	L2BlockTime           *big.Int
	GenesisL2Output       [32]byte
	HistoricalTotalBlocks *big.Int
=======
	L2StartingBlock       *big.Int
	GenesisL2Output       [32]byte
	HistoricalTotalBlocks *big.Int
	L2StartingTimeStamp   *big.Int
	L2BlockTime           *big.Int
>>>>>>> v0.5.23
=======
	HistoricalTotalBlocks *big.Int
>>>>>>> v0.5.24
}

type DepositContractConfig struct {
	L2Oracle           common.Address
	FinalizationPeriod *big.Int
}

type SystemConfig struct {
	Secrets                *e2eutils.Secrets
	L1InfoPredeployAddress common.Address

	DeployConfig *genesis.DeployConfig

	JWTFilePath string
	JWTSecret   [32]byte

<<<<<<< HEAD
<<<<<<< HEAD
	Nodes        map[string]*rollupNode.Config // Per node config. Don't use populate rollup.Config
	Loggers      map[string]log.Logger
	RollupConfig rollup.Config // Shared rollup configs
=======
=======
	Premine        map[common.Address]*big.Int
>>>>>>> @eth-optimism/l2geth@0.5.27
	Nodes          map[string]*rollupNode.Config // Per node config. Don't use populate rollup.Config
	Loggers        map[string]log.Logger
	ProposerLogger log.Logger
	BatcherLogger  log.Logger
<<<<<<< HEAD
	RollupConfig   rollup.Config // Shared rollup configs
>>>>>>> v0.5.23

	L1BlockTime uint64
=======
>>>>>>> @eth-optimism/l2geth@0.5.27

	// map of outbound connections to other nodes. Node names prefixed with "~" are unconnected but linked.
	// A nil map disables P2P completely.
	// Any node name not in the topology will not have p2p enabled.
	P2PTopology map[string][]string
<<<<<<< HEAD
<<<<<<< HEAD
=======

	BaseFeeRecipient common.Address
	L1FeeRecipient   common.Address
>>>>>>> v0.5.23
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
}

type System struct {
	cfg SystemConfig

	RollupConfig *rollup.Config

	// Connections to running nodes
<<<<<<< HEAD
<<<<<<< HEAD
	nodes                      map[string]*node.Node
	backends                   map[string]*eth.Ethereum
	Clients                    map[string]*ethclient.Client
	RolupGenesis               rollup.Genesis
	rollupNodes                map[string]*rollupNode.OpNode
	l2OutputSubmitter          *l2os.L2OutputSubmitter
	sequencerHistoryDBFileName string
	batchSubmitter             *bss.BatchSubmitter
	L2OOContractAddr           common.Address
	DepositContractAddr        common.Address
	Mocknet                    mocknet.Mocknet
=======
	nodes               map[string]*node.Node
	backends            map[string]*eth.Ethereum
	Clients             map[string]*ethclient.Client
	RolupGenesis        rollup.Genesis
	rollupNodes         map[string]*rollupNode.OpNode
	l2OutputSubmitter   *l2os.L2OutputSubmitter
	batchSubmitter      *bss.BatchSubmitter
	L2OOContractAddr    common.Address
	DepositContractAddr common.Address
	Mocknet             mocknet.Mocknet
>>>>>>> v0.5.23
}

func precompileAlloc() core.GenesisAlloc {
	alloc := make(map[common.Address]core.GenesisAccount)
	var addr [common.AddressLength]byte
	for i := 0; i < 256; i++ {
		addr[common.AddressLength-1] = byte(i)
		alloc[addr] = core.GenesisAccount{Balance: common.Big1}
	}
	return alloc
}

func cliqueExtraData(w accounts.Wallet, signers []string) []byte {
	// 32 Empty bytes
	ret := make([]byte, 32)
	// Signer addresses
	for _, signer := range signers {
		address := deriveAddress(w, signer)
		// Was not able to automatically do this
		for i := 0; i < len(address); i++ {
			ret = append(ret, address[i])
		}
	}
	// 65 Empty bytes
	t := make([]byte, 65)
	return append(ret, t...)
=======
	Nodes             map[string]*node.Node
	Backends          map[string]*geth_eth.Ethereum
	Clients           map[string]*ethclient.Client
	RollupNodes       map[string]*rollupNode.OpNode
	L2OutputSubmitter *l2os.L2OutputSubmitter
	BatchSubmitter    *bss.BatchSubmitter
	Mocknet           mocknet.Mocknet
>>>>>>> @eth-optimism/l2geth@0.5.27
}

func (sys *System) Close() {
	if sys.L2OutputSubmitter != nil {
		sys.L2OutputSubmitter.Stop()
	}
	if sys.BatchSubmitter != nil {
		sys.BatchSubmitter.Stop()
	}
<<<<<<< HEAD
	if sys.sequencerHistoryDBFileName != "" {
		_ = os.Remove(sys.sequencerHistoryDBFileName)
	}
=======
>>>>>>> v0.5.23

	for _, node := range sys.RollupNodes {
		node.Close()
	}
	for _, node := range sys.Nodes {
		node.Close()
	}
	sys.Mocknet.Close()
}

func (cfg SystemConfig) Start() (*System, error) {
	sys := &System{
		cfg:         cfg,
		Nodes:       make(map[string]*node.Node),
		Backends:    make(map[string]*geth_eth.Ethereum),
		Clients:     make(map[string]*ethclient.Client),
		RollupNodes: make(map[string]*rollupNode.OpNode),
	}
	didErrAfterStart := false
	defer func() {
		if didErrAfterStart {
			for _, node := range sys.RollupNodes {
				node.Close()
			}
			for _, node := range sys.Nodes {
				node.Close()
			}
		}
	}()

	l1Genesis, err := genesis.BuildL1DeveloperGenesis(cfg.DeployConfig)
	if err != nil {
		return nil, err
	}

	for addr, amount := range cfg.Premine {
		if existing, ok := l1Genesis.Alloc[addr]; ok {
			l1Genesis.Alloc[addr] = core.GenesisAccount{
				Code:    existing.Code,
				Storage: existing.Storage,
				Balance: amount,
				Nonce:   existing.Nonce,
			}
		} else {
			l1Genesis.Alloc[addr] = core.GenesisAccount{
				Balance: amount,
				Nonce:   0,
			}
		}
	}

	l1Block := l1Genesis.ToBlock()
	l2Genesis, err := genesis.BuildL2DeveloperGenesis(cfg.DeployConfig, l1Block, nil)
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	// Create the L2 Outputsubmitter Address and set it here because it needs to be derived from the accounts
	l2OOSubmitter, err := wallet.PrivateKey(accounts.Account{
		URL: accounts.URL{
			Path: cfg.L2OutputHDPath,
		},
	})
	if err != nil {
		return nil, err
	}
	l2OutputSubmitterAddr := crypto.PubkeyToAddress(l2OOSubmitter.PublicKey)

	// Genesis
	eth := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

	l1Alloc := precompileAlloc()
	l2Alloc := precompileAlloc()

	for path, amt := range cfg.Premine {
		balance := big.NewInt(int64(amt))
		balance.Mul(balance, eth)
		addr := deriveAddress(wallet, path)
		l1Alloc[addr] = core.GenesisAccount{Balance: balance}
		l2Alloc[addr] = core.GenesisAccount{Balance: balance}
	}

	l2Alloc[cfg.L1InfoPredeployAddress] = core.GenesisAccount{Code: common.FromHex(bindings.L1BlockDeployedBin), Balance: common.Big0}
<<<<<<< HEAD
	l2Alloc[common.HexToAddress(predeploys.L2ToL1MessagePasser)] = core.GenesisAccount{Code: common.FromHex(bindings.L2ToL1MessagePasserDeployedBin), Balance: common.Big0}
=======
	l2Alloc[predeploys.L2ToL1MessagePasserAddr] = core.GenesisAccount{Code: common.FromHex(bindings.L2ToL1MessagePasserDeployedBin), Balance: common.Big0}
	l2Alloc[predeploys.GasPriceOracleAddr] = core.GenesisAccount{Code: common.FromHex(bindings.GasPriceOracleDeployedBin), Balance: common.Big0, Storage: map[common.Hash]common.Hash{
		// storage for GasPriceOracle to have transctorPath wallet as owner
		common.BigToHash(big.NewInt(0)): common.HexToHash("0x8A0A996b22B103B500Cd0F20d62dF2Ba3364D295"),
	}}
>>>>>>> v0.5.23

	genesisTimestamp := uint64(time.Now().Unix())

	l1Genesis := &core.Genesis{
		Config: &params.ChainConfig{
			ChainID:             cfg.L1ChainID,
			HomesteadBlock:      common.Big0,
			EIP150Block:         common.Big0,
			EIP155Block:         common.Big0,
			EIP158Block:         common.Big0,
			ByzantiumBlock:      common.Big0,
			ConstantinopleBlock: common.Big0,
			PetersburgBlock:     common.Big0,
			IstanbulBlock:       common.Big0,
			BerlinBlock:         common.Big0,
			LondonBlock:         common.Big0,
			Clique: &params.CliqueConfig{
				Period: cfg.L1BlockTime,
				Epoch:  30000,
			},
		},
		Alloc:      l1Alloc,
		Difficulty: common.Big1,
		ExtraData:  cliqueExtraData(wallet, []string{cfg.CliqueSignerDerivationPath}),
		GasLimit:   5000000,
		Nonce:      4660,
		Timestamp:  genesisTimestamp,
		BaseFee:    big.NewInt(7),
	}
	l2Genesis := &core.Genesis{
		Config: &params.ChainConfig{
			ChainID:                 cfg.L2ChainID,
			HomesteadBlock:          common.Big0,
			EIP150Block:             common.Big0,
			EIP155Block:             common.Big0,
			EIP158Block:             common.Big0,
			ByzantiumBlock:          common.Big0,
			ConstantinopleBlock:     common.Big0,
			PetersburgBlock:         common.Big0,
			IstanbulBlock:           common.Big0,
			BerlinBlock:             common.Big0,
			LondonBlock:             common.Big0,
			MergeNetsplitBlock:      common.Big0,
			TerminalTotalDifficulty: common.Big0,
<<<<<<< HEAD
=======
			Optimism: &params.OptimismConfig{
				BaseFeeRecipient:   cfg.BaseFeeRecipient,
				L1FeeRecipient:     cfg.L1FeeRecipient,
				EIP1559Elasticity:  2,
				EIP1559Denominator: 8,
			},
>>>>>>> v0.5.23
		},
		Alloc:      l2Alloc,
		Difficulty: common.Big1,
		GasLimit:   5000000,
		Nonce:      0,
		// must be equal (or higher, while within bounds) as the L1 anchor point of the rollup
		Timestamp: genesisTimestamp,
		BaseFee:   big.NewInt(7),
=======
	makeRollupConfig := func() rollup.Config {
		return rollup.Config{
			Genesis: rollup.Genesis{
				L1: eth.BlockID{
					Hash:   l1Block.Hash(),
					Number: 0,
				},
				L2: eth.BlockID{
					Hash:   l2Genesis.ToBlock().Hash(),
					Number: 0,
				},
				L2Time: uint64(cfg.DeployConfig.L1GenesisBlockTimestamp),
			},
			BlockTime:              cfg.DeployConfig.L2BlockTime,
			MaxSequencerDrift:      cfg.DeployConfig.MaxSequencerDrift,
			SeqWindowSize:          cfg.DeployConfig.SequencerWindowSize,
			ChannelTimeout:         cfg.DeployConfig.ChannelTimeout,
			L1ChainID:              cfg.L1ChainIDBig(),
			L2ChainID:              cfg.L2ChainIDBig(),
			P2PSequencerAddress:    cfg.DeployConfig.P2PSequencerAddress,
			FeeRecipientAddress:    l2Genesis.Coinbase,
			BatchInboxAddress:      cfg.DeployConfig.BatchInboxAddress,
			BatchSenderAddress:     cfg.DeployConfig.BatchSenderAddress,
			DepositContractAddress: predeploys.DevOptimismPortalAddr,
		}
>>>>>>> @eth-optimism/l2geth@0.5.27
	}
	defaultConfig := makeRollupConfig()
	sys.RollupConfig = &defaultConfig

	// Initialize nodes
	l1Node, l1Backend, err := initL1Geth(&cfg, l1Genesis)
	if err != nil {
		return nil, err
	}
	sys.Nodes["l1"] = l1Node
	sys.Backends["l1"] = l1Backend

	for name := range cfg.Nodes {
		node, backend, err := initL2Geth(name, big.NewInt(int64(cfg.DeployConfig.L2ChainID)), l2Genesis, cfg.JWTFilePath)
		if err != nil {
			return nil, err
		}
		sys.Nodes[name] = node
		sys.Backends[name] = backend
	}

	// Start
	err = l1Node.Start()
	if err != nil {
		didErrAfterStart = true
		return nil, err
	}
	err = l1Backend.StartMining(1)
	if err != nil {
		didErrAfterStart = true
		return nil, err
	}
	for name, node := range sys.Nodes {
		if name == "l1" {
			continue
		}
		err = node.Start()
		if err != nil {
			didErrAfterStart = true
			return nil, err
		}
	}

	// Configure connections to L1 and L2 for rollup nodes.
	// TODO: refactor testing to use in-process rpc connections instead of websockets.

	l1EndpointConfig := l1Node.WSEndpoint()
	useHTTP := os.Getenv("OP_E2E_USE_HTTP") == "true"
	if useHTTP {
		log.Info("using HTTP client")
		l1EndpointConfig = l1Node.HTTPEndpoint()
	}

	for name, rollupCfg := range cfg.Nodes {
		l2EndpointConfig := sys.Nodes[name].WSAuthEndpoint()
		if useHTTP {
			l2EndpointConfig = sys.Nodes[name].HTTPAuthEndpoint()
		}
		rollupCfg.L1 = &rollupNode.L1EndpointConfig{
			L1NodeAddr: l1EndpointConfig,
			L1TrustRPC: false,
		}
		rollupCfg.L2 = &rollupNode.L2EndpointConfig{
			L2EngineAddr:      l2EndpointConfig,
			L2EngineJWTSecret: cfg.JWTSecret,
		}
	}

	// Geth Clients
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	l1Srv, err := l1Node.RPCHandler()
	if err != nil {
		didErrAfterStart = true
		return nil, err
	}
	l1Client := ethclient.NewClient(rpc.DialInProc(l1Srv))
	sys.Clients["l1"] = l1Client
	for name, node := range sys.Nodes {
		client, err := ethclient.DialContext(ctx, node.WSEndpoint())
		if err != nil {
			didErrAfterStart = true
			return nil, err
		}
		sys.Clients[name] = client
	}

<<<<<<< HEAD
	// Rollup Genesis
	l1GenesisID, _ := getGenesisInfo(l1Client)
	var l2Client *ethclient.Client
	for name, client := range sys.Clients {
		if name != "l1" {
			l2Client = client
			break
		}
	}
	l2GenesisID, l2GenesisTime := getGenesisInfo(l2Client)

	sys.RolupGenesis = rollup.Genesis{
		L1:     l1GenesisID,
		L2:     l2GenesisID,
		L2Time: l2GenesisTime,
	}

	sys.cfg.RollupConfig.Genesis = sys.RolupGenesis
	sys.cfg.RollupConfig.BatchSenderAddress = batchSubmitterAddr
	sys.cfg.RollupConfig.P2PSequencerAddress = p2pSignerAddr
<<<<<<< HEAD
<<<<<<< HEAD
	sys.cfg.L2OOCfg.L2StartTime = new(big.Int).SetUint64(l2GenesisTime)
=======
	sys.cfg.L2OOCfg.L2StartingBlock = new(big.Int).SetUint64(l2GenesisID.Number)
	sys.cfg.L2OOCfg.L2StartingTimeStamp = new(big.Int).SetUint64(l2Genesis.Timestamp)
	sys.cfg.L2OOCfg.L2BlockTime = new(big.Int).SetUint64(2)
>>>>>>> v0.5.23
=======
>>>>>>> v0.5.24

	// Deploy Deposit Contract
	deployerPrivKey, err := sys.wallet.PrivateKey(accounts.Account{
		URL: accounts.URL{
<<<<<<< HEAD
			Path: sys.cfg.DeployerHDPath,
=======
			Path: cfg.DeployerHDPath,
>>>>>>> v0.5.23
		},
	})
	if err != nil {
		return nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(deployerPrivKey, cfg.L1ChainID)
	if err != nil {
		return nil, err
	}

	// empty genesis L2 output.
	// Technically this may need to be computed with l2.ComputeL2OutputRoot(...),
	// but there are no fraud proofs active in the test.
	genesisL2Output := [32]byte{}

	// Deploy contracts
	sys.L2OOContractAddr, _, _, err = bindings.DeployL2OutputOracle(
		opts,
		l1Client,
		sys.cfg.L2OOCfg.SubmissionFrequency,
<<<<<<< HEAD
<<<<<<< HEAD
		sys.cfg.L2OOCfg.L2BlockTime,
		sys.cfg.L2OOCfg.GenesisL2Output,
		sys.cfg.L2OOCfg.HistoricalTotalBlocks,
		sys.cfg.L2OOCfg.L2StartTime,
		l2OutputSubmitterAddr,
=======
		sys.cfg.L2OOCfg.GenesisL2Output,
=======
		genesisL2Output,
>>>>>>> v0.5.24
		sys.cfg.L2OOCfg.HistoricalTotalBlocks,
		new(big.Int).SetUint64(l2GenesisID.Number),
		new(big.Int).SetUint64(l2Genesis.Timestamp),
		new(big.Int).SetUint64(sys.cfg.RollupConfig.BlockTime),
		l2OutputSubmitterAddr,
		crypto.PubkeyToAddress(deployerPrivKey.PublicKey),
>>>>>>> v0.5.23
	)
	sys.cfg.DepositCFG.L2Oracle = sys.L2OOContractAddr
	if err != nil {
		return nil, err
	}
	var tx *types.Transaction
	sys.DepositContractAddr, tx, _, err = bindings.DeployOptimismPortal(
		opts,
		l1Client,
		sys.cfg.DepositCFG.L2Oracle,
		sys.cfg.DepositCFG.FinalizationPeriod,
	)
	if err != nil {
		return nil, err
	}

	// Wait up to 6 blocks to deploy the Optimism portal
	_, err = waitForTransaction(tx.Hash(), l1Client, 6*time.Second*time.Duration(cfg.L1BlockTime))
	if err != nil {
		return nil, fmt.Errorf("waiting for OptimismPortal: %w", err)
=======
	_, err = waitForBlock(big.NewInt(2), l1Client, 6*time.Second*time.Duration(cfg.DeployConfig.L1BlockTime))
	if err != nil {
		return nil, fmt.Errorf("waiting for blocks: %w", err)
>>>>>>> @eth-optimism/l2geth@0.5.27
	}

	sys.Mocknet = mocknet.New()

	p2pNodes := make(map[string]*p2p.Prepared)
	if cfg.P2PTopology != nil {
		// create the peer if it doesn't exist yet.
		initHostMaybe := func(name string) (*p2p.Prepared, error) {
			if p, ok := p2pNodes[name]; ok {
				return p, nil
			}
			h, err := sys.Mocknet.GenPeer()
			if err != nil {
				return nil, fmt.Errorf("failed to init p2p host for node %s", name)
			}
			h.Network()
			_, ok := cfg.Nodes[name]
			if !ok {
				return nil, fmt.Errorf("node %s from p2p topology not found in actual nodes map", name)
			}
			// TODO we can enable discv5 in the testnodes to test discovery of new peers.
			// Would need to mock though, and the discv5 implementation does not provide nice mocks here.
			p := &p2p.Prepared{
				HostP2P:   h,
				LocalNode: nil,
				UDPv5:     nil,
			}
			p2pNodes[name] = p
			return p, nil
		}
		for k, vs := range cfg.P2PTopology {
			peerA, err := initHostMaybe(k)
			if err != nil {
				return nil, fmt.Errorf("failed to setup mocknet peer %s", k)
			}
			for _, v := range vs {
				v = strings.TrimPrefix(v, "~")
				peerB, err := initHostMaybe(v)
				if err != nil {
					return nil, fmt.Errorf("failed to setup mocknet peer %s (peer of %s)", v, k)
				}
				if _, err := sys.Mocknet.LinkPeers(peerA.HostP2P.ID(), peerB.HostP2P.ID()); err != nil {
					return nil, fmt.Errorf("failed to setup mocknet link between %s and %s", k, v)
				}
				// connect the peers after starting the full rollup node
			}
		}
	}
<<<<<<< HEAD
=======

	// Don't log state snapshots in test output
	snapLog := log.New()
	snapLog.SetHandler(log.DiscardHandler())

>>>>>>> v0.5.23
	// Rollup nodes
	for name, nodeConfig := range cfg.Nodes {
		c := *nodeConfig // copy
		c.Rollup = makeRollupConfig()

		if p, ok := p2pNodes[name]; ok {
			c.P2P = p

<<<<<<< HEAD
			if c.Sequencer {
=======
			if c.Driver.SequencerEnabled {
<<<<<<< HEAD
>>>>>>> v0.5.23
				c.P2PSigner = &p2p.PreparedSigner{Signer: p2p.NewLocalSigner(p2pSignerPrivKey)}
=======
				c.P2PSigner = &p2p.PreparedSigner{Signer: p2p.NewLocalSigner(cfg.Secrets.SequencerP2P)}
>>>>>>> @eth-optimism/l2geth@0.5.27
			}
		}

<<<<<<< HEAD
		node, err := rollupNode.New(context.Background(), &c, cfg.Loggers[name], cfg.Loggers[name], "")
=======
		node, err := rollupNode.New(context.Background(), &c, cfg.Loggers[name], snapLog, "", metrics.NewMetrics(""))
>>>>>>> v0.5.23
		if err != nil {
			didErrAfterStart = true
			return nil, err
		}
		err = node.Start(context.Background())
		if err != nil {
			didErrAfterStart = true
			return nil, err
		}
		sys.RollupNodes[name] = node
	}

	if cfg.P2PTopology != nil {
		// We only set up the connections after starting the actual nodes,
		// so GossipSub and other p2p protocols can be started before the connections go live.
		// This way protocol negotiation happens correctly.
		for k, vs := range cfg.P2PTopology {
			peerA := p2pNodes[k]
			for _, v := range vs {
				unconnected := strings.HasPrefix(v, "~")
				if unconnected {
					v = v[1:]
				}
				if !unconnected {
					peerB := p2pNodes[v]
					if _, err := sys.Mocknet.ConnectPeers(peerA.HostP2P.ID(), peerB.HostP2P.ID()); err != nil {
						return nil, fmt.Errorf("failed to setup mocknet connection between %s and %s", k, v)
					}
				}
			}
		}
	}

	// L2Output Submitter
	sys.L2OutputSubmitter, err = l2os.NewL2OutputSubmitter(l2os.Config{
		L1EthRpc:                  sys.Nodes["l1"].WSEndpoint(),
		L2EthRpc:                  sys.Nodes["sequencer"].WSEndpoint(),
		RollupRpc:                 sys.RollupNodes["sequencer"].HTTPEndpoint(),
		L2OOAddress:               predeploys.DevL2OutputOracleAddr.String(),
		PollInterval:              50 * time.Millisecond,
		NumConfirmations:          1,
		ResubmissionTimeout:       3 * time.Second,
		SafeAbortNonceTooLowCount: 3,
<<<<<<< HEAD
		LogLevel:                  "info",
		LogTerminal:               true,
		Mnemonic:                  sys.cfg.Mnemonic,
		L2OutputHDPath:            sys.cfg.L2OutputHDPath,
<<<<<<< HEAD
	}, "", log.New())
=======
=======
		LogConfig: oplog.CLIConfig{
			Level:  "info",
			Format: "text",
		},
<<<<<<< HEAD
		Mnemonic:       sys.cfg.Mnemonic,
		L2OutputHDPath: sys.cfg.L2OutputHDPath,
>>>>>>> v0.5.24
=======
		PrivateKey: hexPriv(cfg.Secrets.Proposer),
>>>>>>> @eth-optimism/l2geth@0.5.27
	}, "", sys.cfg.Loggers["proposer"])
>>>>>>> v0.5.23
	if err != nil {
		return nil, fmt.Errorf("unable to setup l2 output submitter: %w", err)
	}

	if err := sys.L2OutputSubmitter.Start(); err != nil {
		return nil, fmt.Errorf("unable to start l2 output submitter: %w", err)
	}

<<<<<<< HEAD
	sequencerHistoryDBFile, err := ioutil.TempFile("", "bss.*.json")
	if err != nil {
		return nil, fmt.Errorf("unable to create sequencer history db file: %w", err)
	}
	sys.sequencerHistoryDBFileName = sequencerHistoryDBFile.Name()
	if err = sequencerHistoryDBFile.Close(); err != nil {
		return nil, fmt.Errorf("unable to close sequencer history db file: %w", err)
	}

=======
>>>>>>> v0.5.23
	// Batch Submitter
<<<<<<< HEAD
	sys.batchSubmitter, err = bss.NewBatchSubmitter(bss.Config{
<<<<<<< HEAD
		L1EthRpc:                   sys.nodes["l1"].WSEndpoint(),
		L2EthRpc:                   sys.nodes["sequencer"].WSEndpoint(),
		RollupRpc:                  rollupEndpoint,
		MinL1TxSize:                1,
		MaxL1TxSize:                120000,
<<<<<<< HEAD
=======
		ChannelTimeout:             sys.cfg.RollupConfig.ChannelTimeout,
>>>>>>> v0.5.23
		PollInterval:               50 * time.Millisecond,
		NumConfirmations:           1,
		ResubmissionTimeout:        5 * time.Second,
		SafeAbortNonceTooLowCount:  3,
<<<<<<< HEAD
		LogLevel:                   "info",
		LogTerminal:                true,
		Mnemonic:                   sys.cfg.Mnemonic,
		SequencerHDPath:            sys.cfg.BatchSubmitterHDPath,
		SequencerHistoryDBFilename: sys.sequencerHistoryDBFileName,
		SequencerGenesisHash:       sys.RolupGenesis.L2.Hash.String(),
		SequencerBatchInboxAddress: sys.cfg.RollupConfig.BatchInboxAddress.String(),
	}, "", log.New())
=======
		LogLevel:                   "info", // ignored if started in-process this way
		LogTerminal:                true,   // ignored
=======
		L1EthRpc:                  sys.nodes["l1"].WSEndpoint(),
		L2EthRpc:                  sys.nodes["sequencer"].WSEndpoint(),
		RollupRpc:                 rollupEndpoint,
=======
	sys.BatchSubmitter, err = bss.NewBatchSubmitter(bss.Config{
		L1EthRpc:                  sys.Nodes["l1"].WSEndpoint(),
		L2EthRpc:                  sys.Nodes["sequencer"].WSEndpoint(),
		RollupRpc:                 sys.RollupNodes["sequencer"].HTTPEndpoint(),
>>>>>>> @eth-optimism/l2geth@0.5.27
		MinL1TxSize:               1,
		MaxL1TxSize:               120000,
		ChannelTimeout:            cfg.DeployConfig.ChannelTimeout,
		PollInterval:              50 * time.Millisecond,
		NumConfirmations:          1,
		ResubmissionTimeout:       5 * time.Second,
		SafeAbortNonceTooLowCount: 3,
		LogConfig: oplog.CLIConfig{
			Level:  "info",
			Format: "text",
		},
<<<<<<< HEAD
>>>>>>> v0.5.24
		Mnemonic:                   sys.cfg.Mnemonic,
		SequencerHDPath:            sys.cfg.BatchSubmitterHDPath,
		SequencerBatchInboxAddress: sys.cfg.RollupConfig.BatchInboxAddress.String(),
=======
		PrivateKey:                 hexPriv(cfg.Secrets.Batcher),
		SequencerBatchInboxAddress: cfg.DeployConfig.BatchInboxAddress.String(),
>>>>>>> @eth-optimism/l2geth@0.5.27
	}, sys.cfg.Loggers["batcher"])
>>>>>>> v0.5.23
	if err != nil {
		return nil, fmt.Errorf("failed to setup batch submitter: %w", err)
	}

	if err := sys.BatchSubmitter.Start(); err != nil {
		return nil, fmt.Errorf("unable to start batch submitter: %w", err)
	}

	return sys, nil
}

func (cfg SystemConfig) L1ChainIDBig() *big.Int {
	return new(big.Int).SetUint64(cfg.DeployConfig.L1ChainID)
}

func (cfg SystemConfig) L2ChainIDBig() *big.Int {
	return new(big.Int).SetUint64(cfg.DeployConfig.L2ChainID)
}

func uint642big(in uint64) *hexutil.Big {
	b := new(big.Int).SetUint64(in)
	hu := hexutil.Big(*b)
	return &hu
}

func hexPriv(in *ecdsa.PrivateKey) string {
	b := e2eutils.EncodePrivKey(in)
	return hexutil.Encode(b)
}
