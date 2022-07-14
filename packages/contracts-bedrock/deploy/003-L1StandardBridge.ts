/* Imports: Internal */
import { DeployFunction } from 'hardhat-deploy/dist/types'
<<<<<<< HEAD
import { Contract } from 'ethers'
import 'hardhat-deploy'
=======
import 'hardhat-deploy'
import '@nomiclabs/hardhat-ethers'
import '@eth-optimism/hardhat-deploy-config'
>>>>>>> v0.5.23

const deployFn: DeployFunction = async (hre) => {
  const { deploy } = hre.deployments
  const { deployer } = await hre.getNamedAccounts()
<<<<<<< HEAD

  await deploy('L1StandardBridge', {
    from: deployer,
    args: [],
    log: true,
    waitConfirmations: 1,
  })

  const provider = hre.ethers.provider.getSigner(deployer)

  const messenger = await hre.deployments.get('L1CrossDomainMessenger')
  const bridge = await hre.deployments.get('L1StandardBridge')

  const L1StandardBridge = new Contract(bridge.address, bridge.abi, provider)

  const tx = await L1StandardBridge.initialize(messenger.address)
  const receipt = await tx.wait()
  console.log(`${receipt.transactionHash}: initialize(${messenger.address})`)
=======
  const { deployConfig } = hre

  await deploy('L1StandardBridgeProxy', {
    contract: 'Proxy',
    from: deployer,
    args: [deployer],
    log: true,
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

  const messenger = await hre.deployments.get('L1CrossDomainMessengerProxy')

  await deploy('L1StandardBridge', {
    from: deployer,
    args: [messenger.address],
    log: true,
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

  const proxy = await hre.deployments.get('L1StandardBridgeProxy')
  const Proxy = await hre.ethers.getContractAt('Proxy', proxy.address)
  const bridge = await hre.deployments.get('L1StandardBridge')

  const L1StandardBridge = await hre.ethers.getContractAt(
    'L1StandardBridge',
    proxy.address
  )

  const upgradeTx = await Proxy.upgradeToAndCall(
    bridge.address,
    L1StandardBridge.interface.encodeFunctionData('initialize(address)', [
      messenger.address,
    ])
  )
  await upgradeTx.wait()

  if (messenger.address !== (await L1StandardBridge.messenger())) {
    throw new Error('misconfigured messenger')
  }
>>>>>>> v0.5.23
}

deployFn.tags = ['L1StandardBridge']

export default deployFn
