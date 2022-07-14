/* Imports: Internal */
import { DeployFunction } from 'hardhat-deploy/dist/types'
import 'hardhat-deploy'
<<<<<<< HEAD
=======
import '@nomiclabs/hardhat-ethers'
import '@eth-optimism/hardhat-deploy-config'
>>>>>>> v0.5.23

const deployFn: DeployFunction = async (hre) => {
  const { deploy, get } = hre.deployments
  const { deployer } = await hre.getNamedAccounts()
<<<<<<< HEAD
=======
  const { deployConfig } = hre

  await deploy('OptimismPortalProxy', {
    contract: 'Proxy',
    from: deployer,
    args: [deployer],
    log: true,
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

>>>>>>> v0.5.23
  const oracle = await get('L2OutputOracle')

  await deploy('OptimismPortal', {
    from: deployer,
    args: [oracle.address, 2],
    log: true,
<<<<<<< HEAD
    waitConfirmations: 1,
  })
=======
    waitConfirmations: deployConfig.deploymentWaitConfirmations,
  })

  const proxy = await hre.deployments.get('OptimismPortalProxy')
  const Proxy = await hre.ethers.getContractAt('Proxy', proxy.address)

  const OptimismPortal = await hre.ethers.getContractAt(
    'OptimismPortal',
    proxy.address
  )

  const portal = await hre.deployments.get('OptimismPortal')
  const tx = await Proxy.upgradeToAndCall(
    portal.address,
    OptimismPortal.interface.encodeFunctionData('initialize()')
  )
  await tx.wait()

  const l2Oracle = await OptimismPortal.L2_ORACLE()
  if (l2Oracle !== oracle.address) {
    throw new Error('L2 Oracle mismatch')
  }
>>>>>>> v0.5.23
}

deployFn.tags = ['OptimismPortal']

export default deployFn
