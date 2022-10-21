// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { OptimismMintableERC721 } from "./OptimismMintableERC721.sol";
import { Semver } from "@eth-optimism/contracts-bedrock/contracts/universal/Semver.sol";

/**
 * @title OptimismMintableERC721Factory
 * @notice Factory contract for creating OptimismMintableERC721 contracts.
 */
<<<<<<< HEAD
<<<<<<< HEAD
contract OptimismMintableERC721Factory is OwnableUpgradeable {
    /**
     * @notice Contract version number.
     */
    uint8 public constant VERSION = 1;

=======
contract OptimismMintableERC721Factory is Semver, OwnableUpgradeable {
>>>>>>> v0.5.23
=======
contract OptimismMintableERC721Factory is Semver {
>>>>>>> @eth-optimism/l2geth@0.5.27
    /**
     * @notice Emitted whenever a new OptimismMintableERC721 contract is created.
     *
     * @param localToken  Address of the token on the this domain.
     * @param remoteToken Address of the token on the remote domain.
     * @param deployer    Address of the initiator of the deployment
     */
    event OptimismMintableERC721Created(
        address indexed localToken,
        address indexed remoteToken,
        address deployer
    );

    /**
     * @notice Address of the ERC721 bridge on this network.
     */
    address public immutable bridge;

    /**
<<<<<<< HEAD
=======
     * @notice Chain ID for the remote network.
     */
    uint256 public immutable remoteChainId;

    /**
>>>>>>> v0.5.23
     * @notice Tracks addresses created by this factory.
     */
    mapping(address => bool) public isOptimismMintableERC721;

    /**
<<<<<<< HEAD
<<<<<<< HEAD
     * @param _bridge Address of the ERC721 bridge on this network.
     */
    constructor(address _bridge) {
        intialize(_bridge);
=======
     * @custom:semver 0.0.1
     *
     * @param _bridge Address of the ERC721 bridge on this network.
     */
    constructor(address _bridge, uint256 _remoteChainId) Semver(0, 0, 1) {
        initialize(_bridge, _remoteChainId);
>>>>>>> v0.5.23
    }

    /**
     * @notice Initializes the factory.
     *
     * @param _bridge Address of the ERC721 bridge on this network.
     */
<<<<<<< HEAD
    function intialize(address _bridge) public reinitializer(VERSION) {
        bridge = _bridge;
=======
    function initialize(address _bridge, uint256 _remoteChainId) public initializer {
        bridge = _bridge;
        remoteChainId = _remoteChainId;
>>>>>>> v0.5.23

        // Initialize upgradable OZ contracts
        __Ownable_init();
=======
     * @custom:semver 1.0.0
     *
     * @param _bridge Address of the ERC721 bridge on this network.
     */
    constructor(address _bridge, uint256 _remoteChainId) Semver(1, 0, 0) {
        require(
            _bridge != address(0),
            "OptimismMintableERC721Factory: bridge cannot be address(0)"
        );
        require(
            _remoteChainId != 0,
            "OptimismMintableERC721Factory: remote chain id cannot be zero"
        );

        bridge = _bridge;
        remoteChainId = _remoteChainId;
>>>>>>> @eth-optimism/l2geth@0.5.27
    }

    /**
     * @notice Creates an instance of the standard ERC721.
     *
     * @param _remoteToken Address of the corresponding token on the other domain.
<<<<<<< HEAD
     * @param _name ERC721 name.
     * @param _symbol ERC721 symbol.
=======
     * @param _name        ERC721 name.
     * @param _symbol      ERC721 symbol.
>>>>>>> v0.5.23
     */
    function createOptimismMintableERC721(
        address _remoteToken,
        string memory _name,
        string memory _symbol
    ) external returns (address) {
        require(
            _remoteToken != address(0),
            "OptimismMintableERC721Factory: L1 token address cannot be address(0)"
        );
<<<<<<< HEAD
=======

<<<<<<< HEAD
>>>>>>> v0.5.23
        require(
            bridge != address(0),
            "OptimismMintableERC721Factory: bridge address must be initialized"
        );

        OptimismMintableERC721 localToken = new OptimismMintableERC721(
            bridge,
<<<<<<< HEAD
=======
            remoteChainId,
>>>>>>> v0.5.23
            _remoteToken,
            _name,
            _symbol
        );
=======
        address localToken = address(
            new OptimismMintableERC721(bridge, remoteChainId, _remoteToken, _name, _symbol)
        );

        isOptimismMintableERC721[localToken] = true;
        emit OptimismMintableERC721Created(localToken, _remoteToken, msg.sender);
>>>>>>> @eth-optimism/l2geth@0.5.27

        return localToken;
    }
}
