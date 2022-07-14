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

  await deploy('L1CrossDomainMessenger', {
    from: deployer,
    args: [],
    log: true,
    waitConfirmations: 1,
  })

  const provider = hre.ethers.provider.getSigner(deployer)
  const oracle = await hre.deployments.get('L2OutputOracle')
  const messenger = await hre.deployments.get('L1CrossDomainMessenger')

  const L1CrossDomainMessenger = new Contract(
    messenger.address,
    messenger.abi,
    provider
  )

  const tx = await L1CrossDomainMessenger.initialize(oracle.address)
  const receipt = await tx.wait()
  console.log(`${receipt.transactionHash}: initialize(${oracle.address})`)
=======
  const { deployConfig } = hre

  await deploy('L1CrossDomainMessengerProxy', {
    contract: 'Proxy',
    from: deployer,
    args: [deployer],
    log: true,
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

  const portal = await hre.deployments.get('OptimismPortalProxy')

  await deploy('L1CrossDomainMessenger', {
    from: deployer,
    args: [portal.address],
    log: true,
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

  const proxy = await hre.deployments.get('L1CrossDomainMessengerProxy')
  const Proxy = await hre.ethers.getContractAt('Proxy', proxy.address)
  const messenger = await hre.deployments.get('L1CrossDomainMessenger')

  const L1CrossDomainMessenger = await hre.ethers.getContractAt(
    'L1CrossDomainMessenger',
    proxy.address
  )

  const upgradeTx = await Proxy.upgradeToAndCall(
    messenger.address,
    L1CrossDomainMessenger.interface.encodeFunctionData('initialize')
  )
  await upgradeTx.wait()

  const portalAddress = await L1CrossDomainMessenger.portal()
  if (portalAddress !== portal.address) {
    throw new Error('portal misconfigured')
  }
>>>>>>> v0.5.23
}

deployFn.tags = ['L1CrossDomainMessenger']

export default deployFn
