// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
/**
 * @title OVM_DeployerWhitelist
 * @dev The Deployer Whitelist is a temporary predeploy used to provide additional safety during the
 * initial phases of our mainnet roll out. It is owned by the Optimism team, and defines accounts
 * which are allowed to deploy contracts on Layer2. The Execution Manager will only allow an
 * ovmCREATE or ovmCREATE2 operation to proceed if the deployer's address whitelisted.
 */
contract DeployerWhitelist {
    /**********
     * Events *
     **********/

    event OwnerChanged(address oldOwner, address newOwner);
    event WhitelistStatusChanged(address deployer, bool whitelisted);
    event WhitelistDisabled(address oldOwner);

    /**********************
     * Contract Constants *
     **********************/

    // WARNING: When owner is set to address(0), the whitelist is disabled.
    address public owner;
    mapping(address => bool) public whitelist;

    /**********************
     * Function Modifiers *
     **********************/

    /**
     * Blocks functions to anyone except the contract owner.
     */
    modifier onlyOwner() {
        require(msg.sender == owner, "Function can only be called by the owner of this contract.");
        _;
    }

    /********************
     * Public Functions *
     ********************/

    /**
     * Adds or removes an address from the deployment whitelist.
     * @param _deployer Address to update permissions for.
=======
import { Semver } from "../universal/Semver.sol";

/**
 * @custom:legacy
 * @custom:proxied
 * @custom:predeployed 0x4200000000000000000000000000000000000002
 * @title DeployerWhitelist
 * @notice DeployerWhitelist is a legacy contract that was originally used to act as a whitelist of
 *         addresses allowed to the Optimism network. The DeployerWhitelist has since been
 *         disabled, but the code is kept in state for the sake of full backwards compatibility.
 *         As of the Bedrock upgrade, the DeployerWhitelist is completely unused by the Optimism
 *         system and could, in theory, be removed entirely.
 */
contract DeployerWhitelist is Semver {
    /**
     * @notice Address of the owner of this contract. Note that when this address is set to
     *         address(0), the whitelist is disabled.
     */
    address public owner;

    /**
     * @notice Mapping of deployer addresses to boolean whitelist status.
     */
    mapping(address => bool) public whitelist;

    /**
     * @notice Emitted when the owner of this contract changes.
     *
     * @param oldOwner Address of the previous owner.
     * @param newOwner Address of the new owner.
     */
    event OwnerChanged(address oldOwner, address newOwner);

    /**
     * @notice Emitted when the whitelist status of a deployer changes.
     *
     * @param deployer    Address of the deployer.
     * @param whitelisted Boolean indicating whether the deployer is whitelisted.
     */
    event WhitelistStatusChanged(address deployer, bool whitelisted);

    /**
     * @notice Emitted when the whitelist is disabled.
     *
     * @param oldOwner Address of the final owner of the whitelist.
     */
    event WhitelistDisabled(address oldOwner);

    /**
     * @notice Blocks functions to anyone except the contract owner.
     */
    modifier onlyOwner() {
        require(
            msg.sender == owner,
            "DeployerWhitelist: function can only be called by the owner of this contract"
        );
        _;
    }

    /**
     * @custom:semver 0.0.1
     */
    constructor() Semver(0, 0, 1) {}

    /**
     * @notice Adds or removes an address from the deployment whitelist.
     *
     * @param _deployer      Address to update permissions for.
>>>>>>> v0.5.23
     * @param _isWhitelisted Whether or not the address is whitelisted.
     */
    function setWhitelistedDeployer(address _deployer, bool _isWhitelisted) external onlyOwner {
        whitelist[_deployer] = _isWhitelisted;
        emit WhitelistStatusChanged(_deployer, _isWhitelisted);
    }

    /**
<<<<<<< HEAD
     * Updates the owner of this contract.
     * @param _owner Address of the new owner.
     */
    // slither-disable-next-line external-function
    function setOwner(address _owner) public onlyOwner {
=======
     * @notice Updates the owner of this contract.
     *
     * @param _owner Address of the new owner.
     */
    function setOwner(address _owner) external onlyOwner {
>>>>>>> v0.5.23
        // Prevent users from setting the whitelist owner to address(0) except via
        // enableArbitraryContractDeployment. If you want to burn the whitelist owner, send it to
        // any other address that doesn't have a corresponding knowable private key.
        require(
            _owner != address(0),
<<<<<<< HEAD
            "OVM_DeployerWhitelist: can only be disabled via enableArbitraryContractDeployment"
=======
            "DeployerWhitelist: can only be disabled via enableArbitraryContractDeployment"
>>>>>>> v0.5.23
        );

        emit OwnerChanged(owner, _owner);
        owner = _owner;
    }

    /**
<<<<<<< HEAD
     * Permanently enables arbitrary contract deployment and deletes the owner.
=======
     * @notice Permanently enables arbitrary contract deployment and deletes the owner.
>>>>>>> v0.5.23
     */
    function enableArbitraryContractDeployment() external onlyOwner {
        emit WhitelistDisabled(owner);
        owner = address(0);
    }

    /**
<<<<<<< HEAD
     * Checks whether an address is allowed to deploy contracts.
     * @param _deployer Address to check.
     * @return _allowed Whether or not the address can deploy contracts.
=======
     * @notice Checks whether an address is allowed to deploy contracts.
     *
     * @param _deployer Address to check.
     *
     * @return Whether or not the address can deploy contracts.
>>>>>>> v0.5.23
     */
    function isDeployerAllowed(address _deployer) external view returns (bool) {
        return (owner == address(0) || whitelist[_deployer]);
    }
}
