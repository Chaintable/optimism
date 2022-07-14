/* Imports: External */
import { DeployFunction } from 'hardhat-deploy/dist/types'

const deployFn: DeployFunction = async (hre) => {
  const { deployer } = await hre.getNamedAccounts()

<<<<<<< HEAD
  const { deploy } = await hre.deployments.deterministic('ProxyAdmin', {
    salt: hre.ethers.utils.solidityKeccak256(
      ['string'],
      ['PeripheryProxyAdmin']
    ),
    from: deployer,
    args: [hre.deployConfig.ddd],
    log: true,
  })
=======
  const { deploy } = await hre.deployments.deterministic(
    'PeripheryProxyAdmin',
    {
      contract: 'ProxyAdmin',
      salt: hre.ethers.utils.solidityKeccak256(
        ['string'],
        ['PeripheryProxyAdmin']
      ),
      from: deployer,
      args: [hre.deployConfig.ddd],
      log: true,
    }
  )
>>>>>>> v0.5.23

  await deploy()
}

deployFn.tags = ['PeripheryProxyAdmin']

export default deployFn
