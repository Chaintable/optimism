// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {
    OwnableUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { OptimismMintableERC721 } from "./OptimismMintableERC721.sol";
<<<<<<< HEAD
=======
import { Semver } from "@eth-optimism/contracts-bedrock/contracts/universal/Semver.sol";
>>>>>>> v0.5.23

/**
 * @title OptimismMintableERC721Factory
 * @notice Factory contract for creating OptimismMintableERC721 contracts.
 */
<<<<<<< HEAD
contract OptimismMintableERC721Factory is OwnableUpgradeable {
    /**
     * @notice Contract version number.
     */
    uint8 public constant VERSION = 1;

=======
contract OptimismMintableERC721Factory is Semver, OwnableUpgradeable {
>>>>>>> v0.5.23
    /**
     * @notice Emitted whenever a new OptimismMintableERC721 contract is created.
     *
     * @param remoteToken Address of the token on the remote domain.
     * @param localToken  Address of the token on the this domain.
     */
    event OptimismMintableERC721Created(address indexed remoteToken, address indexed localToken);

    /**
     * @notice Address of the ERC721 bridge on this network.
     */
    address public bridge;

    /**
<<<<<<< HEAD
=======
     * @notice Chain ID for the remote network.
     */
    uint256 public remoteChainId;

    /**
>>>>>>> v0.5.23
     * @notice Tracks addresses created by this factory.
     */
    mapping(address => bool) public isStandardOptimismMintableERC721;

    /**
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
    function createStandardOptimismMintableERC721(
        address _remoteToken,
        string memory _name,
        string memory _symbol
    ) external {
        require(
            _remoteToken != address(0),
            "OptimismMintableERC721Factory: L1 token address cannot be address(0)"
        );
<<<<<<< HEAD
=======

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

        isStandardOptimismMintableERC721[address(localToken)] = true;
        emit OptimismMintableERC721Created(_remoteToken, address(localToken));
    }
}
