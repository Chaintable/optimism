/* Imports: External */
import { DeployFunction } from 'hardhat-deploy/dist/types'
import '@eth-optimism/hardhat-deploy-config'
import '@nomiclabs/hardhat-ethers'
import 'hardhat-deploy'

const deployFn: DeployFunction = async (hre) => {
  const { deployer } = await hre.getNamedAccounts()
  const { deploy } = hre.deployments

<<<<<<< HEAD
<<<<<<< HEAD
  const { deploy } = await hre.deployments.deterministic('Proxy', {
    salt: hre.ethers.utils.solidityKeccak256(
      ['string'],
      ['OptimismMintableERC721FactoryProxy']
    ),
    from: deployer,
    args: [hre.deployConfig.ddd],
    log: true,
  })
=======
  const { deploy } = await hre.deployments.deterministic(
    'OptimismMintableERC721FactoryProxy',
    {
      contract: 'Proxy',
      salt: hre.ethers.utils.solidityKeccak256(
        ['string'],
        ['OptimismMintableERC721FactoryProxy']
      ),
      from: deployer,
      args: [hre.deployConfig.ddd],
      log: true,
    }
  )
>>>>>>> v0.5.23
=======
  console.log(
    `Deploying OptimismMintableERC721FactoryProxy to ${hre.network.name}`
  )
  console.log(`Using deployer ${deployer}`)
>>>>>>> @eth-optimism/l2geth@0.5.27

  // Deploy the OptimismMintableERC721FactoryProxy with
  // the deployer as the admin. The admin and implementation
  // will be updated with the deployment of the implementation
  await deploy('OptimismMintableERC721FactoryProxy', {
    contract: 'Proxy',
    from: deployer,
    args: [deployer],
    log: true,
    waitConfirmations: 1,
  })

  const Deployment__OptimismMintableERC721FactoryProxy =
    await hre.deployments.get('OptimismMintableERC721FactoryProxy')
  console.log(
    `OptimismMintableERC721FactoryProxy deployed to ${Deployment__OptimismMintableERC721FactoryProxy.address}`
  )
}

deployFn.tags = ['OptimismMintableERC721FactoryProxy']

export default deployFn
