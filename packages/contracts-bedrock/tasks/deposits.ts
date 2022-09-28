<<<<<<< HEAD
import { task, types } from 'hardhat/config'
import { Contract, providers, utils, Wallet, Event } from 'ethers'
import dotenv from 'dotenv'

import { DepositTx } from '../src'
=======
/*
 * Copyright (c) 2022, OP Labs PBC (MIT License)
 * https://github.com/ethereum-optimism/optimism
 */

import { task, types } from 'hardhat/config'
import { providers, utils, Wallet, Event } from 'ethers'
import dotenv from 'dotenv'
import 'hardhat-deploy'
import '@nomiclabs/hardhat-ethers'
import { DepositTx } from '@eth-optimism/core-utils'
>>>>>>> v0.5.23

dotenv.config()

const sleep = async (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

task('deposit', 'Deposits funds onto L2.')
  .addParam(
    'l1ProviderUrl',
    'L1 provider URL.',
    'http://localhost:8545',
    types.string
  )
  .addParam(
    'l2ProviderUrl',
    'L2 provider URL.',
    'http://localhost:9545',
    types.string
  )
  .addParam('to', 'Recipient address.', null, types.string)
  .addParam('amountEth', 'Amount in ETH to send.', null, types.string)
  .addOptionalParam(
    'privateKey',
    'Private key to send transaction',
    process.env.PRIVATE_KEY,
    types.string
  )
<<<<<<< HEAD
  .addOptionalParam(
    'depositContractAddr',
    'Address of deposit contract.',
    'deaddeaddeaddeaddeaddeaddeaddeaddead0001',
    types.string
  )
  .setAction(async (args, hre) => {
    const {
      l1ProviderUrl,
      l2ProviderUrl,
      to,
      amountEth,
      depositContractAddr,
      privateKey,
    } = args
    const depositFeedArtifact = await hre.deployments.get('OptimismPortal')
=======
  .setAction(async (args, hre) => {
    const { l1ProviderUrl, l2ProviderUrl, to, amountEth, privateKey } = args
    const proxy = await hre.deployments.get('OptimismPortalProxy')

    const OptimismPortal = await hre.ethers.getContractAt(
      'OptimismPortal',
      proxy.address
    )
>>>>>>> v0.5.23

    const l1Provider = new providers.JsonRpcProvider(l1ProviderUrl)
    const l2Provider = new providers.JsonRpcProvider(l2ProviderUrl)

    let l1Wallet: Wallet | providers.JsonRpcSigner
    if (privateKey) {
      l1Wallet = new Wallet(privateKey, l1Provider)
    } else {
      l1Wallet = l1Provider.getSigner()
    }

    const from = await l1Wallet.getAddress()
    console.log(`Sending from ${from}`)
    const balance = await l1Wallet.getBalance()
    if (balance.eq(0)) {
      throw new Error(`${from} has no balance`)
    }

<<<<<<< HEAD
    const depositFeed = new Contract(
      depositContractAddr,
      depositFeedArtifact.abi
    ).connect(l1Wallet)

=======
>>>>>>> v0.5.23
    const amountWei = utils.parseEther(amountEth)
    const value = amountWei.add(utils.parseEther('0.01'))
    console.log(`Depositing ${amountEth} ETH to ${to}`)

    const preL2Balance = await l2Provider.getBalance(to)
    console.log(`${to} has ${utils.formatEther(preL2Balance)} ETH on L2`)

    // Below adds 0.01 ETH to account for gas.
<<<<<<< HEAD
    const tx = await depositFeed.depositTransaction(
=======
    const tx = await OptimismPortal.depositTransaction(
>>>>>>> v0.5.23
      to,
      amountWei,
      '3000000',
      false,
      [],
      { value }
    )
    console.log(`Got TX hash ${tx.hash}. Waiting...`)
    const receipt = await tx.wait()
<<<<<<< HEAD
<<<<<<< HEAD
=======
    console.log(
      `Included in block ${receipt.blockHash} with index ${receipt.logIndex}`
    )
>>>>>>> v0.5.23
=======
    console.log(`Included in block ${receipt.blockHash}`)
>>>>>>> v0.5.24

    // find the transaction deposited event and derive
    // the deposit transaction from it
    const event = receipt.events.find(
      (e: Event) => e.event === 'TransactionDeposited'
    )
<<<<<<< HEAD
    const l2tx = DepositTx.fromL1Event(event)
<<<<<<< HEAD
=======
    console.log(`Deposit has log index ${event.logIndex}`)
>>>>>>> v0.5.23
=======
    console.log(`Deposit has log index ${event.logIndex}`)
    const l2tx = DepositTx.fromL1Event(event)
>>>>>>> v0.5.24
    const hash = l2tx.hash()
    console.log(`Waiting for L2 TX hash ${hash}`)

    let i = 0
    while (true) {
      const expected = await l2Provider.send('eth_getTransactionByHash', [hash])
      if (expected) {
        console.log('Deposit success')
        console.log(JSON.stringify(expected, null, 2))
        console.log('Receipt:')
        const l2Receipt = await l2Provider.getTransactionReceipt(hash)
        console.log(JSON.stringify(l2Receipt, null, 2))
        break
      }

      if (i % 100 === 0) {
        const postL2Balance = await l2Provider.getBalance(to)
        if (postL2Balance.gt(preL2Balance)) {
          console.log(
            `Unexpected balance increase without detecting deposit transaction`
          )
        }
        const block = await l2Provider.getBlock('latest')
        console.log(`latest block ${block.number}:${block.hash}`)
      }

      await sleep(500)
      i++
    }
  })
