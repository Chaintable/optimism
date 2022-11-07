// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
<<<<<<< HEAD
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";
import { OptimismPortal } from "./OptimismPortal.sol";
import { CrossDomainMessenger } from "../universal/CrossDomainMessenger.sol";
=======
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";
=======
import { Predeploys } from "../libraries/Predeploys.sol";
>>>>>>> v0.5.24
import { OptimismPortal } from "./OptimismPortal.sol";
import { CrossDomainMessenger } from "../universal/CrossDomainMessenger.sol";
import { Semver } from "../universal/Semver.sol";
>>>>>>> v0.5.23

/**
 * @custom:proxied
 * @title L1CrossDomainMessenger
 * @notice The L1CrossDomainMessenger is a message passing interface between L1 and L2 responsible
 *         for sending and receiving data on the L1 side. Users are encouraged to use this
 *         interface instead of interacting with lower-level contracts directly.
 */
<<<<<<< HEAD
contract L1CrossDomainMessenger is CrossDomainMessenger {
    /**
     * @notice Address of the OptimismPortal.
     */
    OptimismPortal public portal;

    /**
     * @notice Initializes the L1CrossDomainMessenger.
     *
     * @param _portal Address of the OptimismPortal to send and receive messages through.
     */
    function initialize(OptimismPortal _portal) external {
        portal = _portal;

        address[] memory blockedSystemAddresses = new address[](1);
        blockedSystemAddresses[0] = address(this);

        _initialize(Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER, blockedSystemAddresses);
=======
contract L1CrossDomainMessenger is CrossDomainMessenger, Semver {
    /**
     * @notice Address of the OptimismPortal.
     */
    OptimismPortal public immutable PORTAL;

    /**
     * @custom:semver 0.0.1
     *
     * @param _portal Address of the OptimismPortal contract on this network.
     */
    constructor(OptimismPortal _portal)
        Semver(0, 0, 1)
        CrossDomainMessenger(Predeploys.L2_CROSS_DOMAIN_MESSENGER)
    {
        PORTAL = _portal;
        initialize(address(0));
    }

    /**
     * @notice Initializer.
     *
     * @param _owner Address of the initial owner of this contract.
     */
<<<<<<< HEAD
    function initialize() public initializer {
<<<<<<< HEAD
        address[] memory blockedSystemAddresses = new address[](1);
        blockedSystemAddresses[0] = address(this);
<<<<<<< HEAD
        __CrossDomainMessenger_init(
            PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER,
            blockedSystemAddresses
        );
>>>>>>> v0.5.23
    }

    /**
     * @notice Checks whether the message being sent from the other messenger.
     *
     * @return True if the message was sent from the messenger, false otherwise.
     */
    function _isSystemMessageSender() internal view override returns (bool) {
        return msg.sender == address(portal) && portal.l2Sender() == otherMessenger;
=======
        __CrossDomainMessenger_init(Predeploys.L2_CROSS_DOMAIN_MESSENGER, blockedSystemAddresses);
>>>>>>> v0.5.24
=======
        __CrossDomainMessenger_init();
>>>>>>> @eth-optimism/l2geth@0.5.27
=======
    function initialize(address _owner) public initializer {
        __CrossDomainMessenger_init();
        _transferOwnership(_owner);
>>>>>>> @eth-optimism/l2geth@0.5.29
    }

    /**
     * @inheritdoc CrossDomainMessenger
     */
    function _sendMessage(
        address _to,
        uint64 _gasLimit,
        uint256 _value,
        bytes memory _data
    ) internal override {
        PORTAL.depositTransaction{ value: _value }(_to, _value, _gasLimit, false, _data);
    }

    /**
     * @inheritdoc CrossDomainMessenger
     */
    function _isOtherMessenger() internal view override returns (bool) {
        return msg.sender == address(PORTAL) && PORTAL.l2Sender() == OTHER_MESSENGER;
    }

    /**
     * @inheritdoc CrossDomainMessenger
     */
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(PORTAL);
    }
}
