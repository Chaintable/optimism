/* Imports: External */
import { DeployFunction } from 'hardhat-deploy/dist/types'
import { ethers } from 'hardhat'

const deployFn: DeployFunction = async (hre) => {
  const { deployer } = await hre.getNamedAccounts()

  await hre.deployments.deploy('OptimismMintableERC721Factory', {
    from: deployer,
<<<<<<< HEAD
    args: [ethers.constants.AddressZero],
=======
    args: [ethers.constants.AddressZero, 0],
>>>>>>> v0.5.23
    log: true,
  })
}

deployFn.tags = ['OptimismMintableERC721FactoryImplementation']
deployFn.dependencies = ['OptimismMintableERC721FactoryProxy']

export default deployFn
