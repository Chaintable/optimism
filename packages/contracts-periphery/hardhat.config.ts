<<<<<<< HEAD
<<<<<<< HEAD
import { HardhatUserConfig } from 'hardhat/types'
=======
=======
import assert from 'assert'

>>>>>>> v0.5.24
import { HardhatUserConfig, subtask } from 'hardhat/config'
import { TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS } from 'hardhat/builtin-tasks/task-names'
>>>>>>> v0.5.23
import { getenv } from '@eth-optimism/core-utils'
import * as dotenv from 'dotenv'

import { configSpec } from './src/config/deploy'

// Hardhat plugins
import '@nomiclabs/hardhat-ethers'
import '@nomiclabs/hardhat-waffle'
import '@nomiclabs/hardhat-etherscan'
import '@eth-optimism/hardhat-deploy-config'
<<<<<<< HEAD
=======
import '@typechain/hardhat'
>>>>>>> v0.5.23
import 'solidity-coverage'
import 'hardhat-gas-reporter'
import 'hardhat-deploy'

// Hardhat tasks
import './tasks'

// Load environment variables from .env
dotenv.config()

<<<<<<< HEAD
=======
subtask(TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS).setAction(
  async (_, __, runSuper) => {
    const paths = await runSuper()

    return paths.filter((p: string) => !p.endsWith('.t.sol'))
  }
)

<<<<<<< HEAD
>>>>>>> v0.5.23
=======
assert(
  !(getenv('PRIVATE_KEY') && getenv('LEDGER_ADDRESS')),
  'use only one of PRIVATE_KEY or LEDGER_ADDRESS'
)

const accounts = getenv('PRIVATE_KEY')
  ? [getenv('PRIVATE_KEY')]
  : (undefined as any)

>>>>>>> v0.5.24
const config: HardhatUserConfig = {
  networks: {
    optimism: {
      chainId: 10,
      url: 'https://mainnet.optimism.io',
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('OPTIMISTIC_ETHERSCAN_API_KEY'),
        },
      },
    },
    'optimism-kovan': {
      chainId: 69,
      url: 'https://kovan.optimism.io',
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('OPTIMISTIC_ETHERSCAN_API_KEY'),
        },
      },
    },
    'optimism-goerli': {
      chainId: 420,
      url: 'https://goerli.optimism.io',
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('OPTIMISTIC_ETHERSCAN_API_KEY'),
        },
      },
    },
    ethereum: {
      chainId: 1,
      url: `https://mainnet.infura.io/v3/${getenv('INFURA_PROJECT_ID')}`,
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('ETHEREUM_ETHERSCAN_API_KEY'),
        },
      },
    },
    goerli: {
      chainId: 5,
      url: `https://goerli.infura.io/v3/${getenv('INFURA_PROJECT_ID')}`,
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('ETHEREUM_ETHERSCAN_API_KEY'),
        },
      },
    },
    ropsten: {
      chainId: 3,
      url: `https://ropsten.infura.io/v3/${getenv('INFURA_PROJECT_ID')}`,
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('ETHEREUM_ETHERSCAN_API_KEY'),
        },
      },
    },
    kovan: {
      chainId: 42,
      url: `https://kovan.infura.io/v3/${getenv('INFURA_PROJECT_ID')}`,
      accounts,
      verify: {
        etherscan: {
          apiKey: getenv('ETHEREUM_ETHERSCAN_API_KEY'),
        },
      },
    },
  },
  paths: {
    deployConfig: './config/deploy',
  },
  deployConfigSpec: configSpec,
<<<<<<< HEAD
  external: {
    contracts: [
      {
        artifacts: '../contracts-bedrock/artifacts',
      },
    ],
  },
  mocha: {
    timeout: 50000,
  },
=======
  mocha: {
    timeout: 50000,
  },
  typechain: {
    outDir: 'dist/types',
    target: 'ethers-v5',
  },
>>>>>>> v0.5.23
  solidity: {
    compilers: [
      {
        version: '0.8.15',
        settings: {
          optimizer: { enabled: true, runs: 10_000 },
        },
      },
    ],
    settings: {
      metadata: {
        bytecodeHash: 'none',
      },
      outputSelection: {
        '*': {
          '*': ['metadata', 'storageLayout'],
        },
      },
    },
  },
  namedAccounts: {
    deployer: {
      default: getenv('LEDGER_ADDRESS')
        ? `ledger://${getenv('LEDGER_ADDRESS')}`
        : 0,
      hardhat: 0,
    },
  },
}

export default config
