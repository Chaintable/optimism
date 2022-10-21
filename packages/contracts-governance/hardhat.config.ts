import dotenv from 'dotenv'
<<<<<<< HEAD
=======
import { HardhatUserConfig } from 'hardhat/config'
<<<<<<< HEAD
>>>>>>> v0.5.23
=======
import { getenv } from '@eth-optimism/core-utils'
>>>>>>> @eth-optimism/l2geth@0.5.27
import '@nomiclabs/hardhat-ethers'
import '@nomiclabs/hardhat-etherscan'
import '@nomiclabs/hardhat-waffle'
import 'hardhat-gas-reporter'
import 'solidity-coverage'
<<<<<<< HEAD
<<<<<<< HEAD
import { task, types } from 'hardhat/config'
import { providers, utils, Wallet } from 'ethers'
import { CrossChainMessenger } from '@eth-optimism/sdk'
=======
>>>>>>> v0.5.23
=======
import '@eth-optimism/hardhat-deploy-config'
import 'hardhat-deploy'
>>>>>>> @eth-optimism/l2geth@0.5.27

import './scripts/deploy-token'
import './scripts/multi-send'
import './scripts/mint-initial-supply'
import './scripts/generate-merkle-root'
import './scripts/create-airdrop-json'
import './scripts/deploy-distributor'
import './scripts/test-claims'
import './scripts/create-distributor-json'
<<<<<<< HEAD

dotenv.config()

task('accounts', 'Prints the list of accounts').setAction(async (args, hre) => {
  const accounts = await hre.ethers.getSigners()

  for (const account of accounts) {
    console.log(account.address)
  }
})

task('deposit', 'Deposits funds onto Optimism.')
  .addParam('to', 'Recipient address.', null, types.string)
  .addParam('amountEth', 'Amount in ETH to send.', null, types.string)
  .addParam('l1ProviderUrl', '', process.env.L1_PROVIDER_URL, types.string)
  .addParam('l2ProviderUrl', '', process.env.L2_PROVIDER_URL, types.string)
  .addParam('privateKey', '', process.env.PRIVATE_KEY, types.string)
  .setAction(async (args) => {
    const { to, amountEth, l1ProviderUrl, l2ProviderUrl, privateKey } = args
    if (!l1ProviderUrl || !l2ProviderUrl || !privateKey) {
      throw new Error(
        'You must define --l1-provider-url, --l2-provider-url, --private-key in your environment.'
      )
    }

    const l1Provider = new providers.JsonRpcProvider(l1ProviderUrl)
    const l1Wallet = new Wallet(privateKey, l1Provider)
    const messenger = new CrossChainMessenger({
      l1SignerOrProvider: l1Wallet,
      l2SignerOrProvider: l2ProviderUrl,
      l1ChainId: (await l1Provider.getNetwork()).chainId,
    })

    const amountWei = utils.parseEther(amountEth)
    console.log(`Depositing ${amountEth} ETH to ${to}...`)
    const tx = await messenger.depositETH(amountWei, {
      recipient: to,
    })
    console.log(`Got TX hash ${tx.hash}. Waiting...`)
    await tx.wait()

    const l2Provider = new providers.JsonRpcProvider(l2ProviderUrl)

    const l1WalletOnL2 = new Wallet(privateKey, l2Provider)
    await l1WalletOnL2.sendTransaction({
      to,
      value: utils.parseEther(amountEth),
    })

    const balance = await l2Provider.getBalance(to)
    console.log('Funded account balance', balance.toString())
    console.log('Done.')
  })

const privKey = process.env.PRIVATE_KEY || '0x' + '11'.repeat(32)

/**
 * @type import("hardhat/config").HardhatUserConfig
 */
module.exports = {
  solidity: '0.8.12',
=======
import './scripts/deposit'

dotenv.config()

const privKey = process.env.PRIVATE_KEY || '0x' + '11'.repeat(32)

const config: HardhatUserConfig = {
  solidity: {
    version: '0.8.12',
    settings: {
      outputSelection: {
        '*': {
          '*': ['metadata', 'storageLayout'],
        },
      },
    },
  },
>>>>>>> v0.5.23
  networks: {
    optimism: {
      chainId: 17,
      url: 'http://localhost:8545',
<<<<<<< HEAD
      saveDeployments: false,
=======
>>>>>>> v0.5.23
    },
    'optimism-kovan': {
      chainId: 69,
      url: 'https://kovan.optimism.io',
      accounts: [privKey],
    },
<<<<<<< HEAD
    'optimism-nightly': {
      chainId: 421,
      url: 'https://goerli-nightly-us-central1-a-sequencer.optimism.io',
      saveDeployments: true,
=======
    'optimism-goerli': {
      chainId: 420,
      url: 'https://goerli.optimism.io',
      accounts: [privKey],
    },
    'optimism-nightly': {
      chainId: 421,
      url: 'https://goerli-nightly-us-central1-a-sequencer.optimism.io',
>>>>>>> v0.5.23
      accounts: [privKey],
    },
    'optimism-mainnet': {
      chainId: 10,
      url: 'https://mainnet.optimism.io',
      accounts: [privKey],
    },
    'hardhat-node': {
<<<<<<< HEAD
      url: 'http://localhost:9545',
      saveDeployments: false,
=======
      url: 'http://localhost:8545',
>>>>>>> v0.5.23
    },
  },
  paths: {
    deployConfig: 'deploy-config',
  },
  deployConfigSpec: {
    upgrader: {
      type: 'address',
    },
  },
  gasReporter: {
    enabled: process.env.REPORT_GAS !== undefined,
    currency: 'USD',
  },
  etherscan: {
    apiKey: process.env.ETHERSCAN_API_KEY,
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
<<<<<<< HEAD
=======

export default config
>>>>>>> v0.5.23
