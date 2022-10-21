// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { AddressAliasHelper } from "../vendor/AddressAliasHelper.sol";
import { Predeploys } from "../libraries/Predeploys.sol";
import { CrossDomainMessenger } from "../universal/CrossDomainMessenger.sol";
import { Semver } from "../universal/Semver.sol";
import { L2ToL1MessagePasser } from "./L2ToL1MessagePasser.sol";

/**
 * @custom:proxied
 * @custom:predeploy 0x4200000000000000000000000000000000000007
 * @title L2CrossDomainMessenger
 * @notice The L2CrossDomainMessenger is a high-level interface for message passing between L1 and
 *         L2 on the L2 side. Users are generally encouraged to use this contract instead of lower
 *         level message passing contracts.
 */
<<<<<<< HEAD
contract L2CrossDomainMessenger is CrossDomainMessenger {
    /**
     * @notice Initializes the L2CrossDomainMessenger.
     *
     * @param _l1CrossDomainMessenger Address of the L1CrossDomainMessenger contract.
     */
    function initialize(address _l1CrossDomainMessenger) external {
        address[] memory blockedSystemAddresses = new address[](2);
        blockedSystemAddresses[0] = address(this);
        blockedSystemAddresses[1] = Lib_PredeployAddresses.L2_TO_L1_MESSAGE_PASSER;

        _initialize(_l1CrossDomainMessenger, blockedSystemAddresses);
=======
contract L2CrossDomainMessenger is CrossDomainMessenger, Semver {
    /**
     * @custom:semver 0.0.1
     *
     * @param _l1CrossDomainMessenger Address of the L1CrossDomainMessenger contract.
     */
    constructor(address _l1CrossDomainMessenger)
        Semver(0, 0, 1)
        CrossDomainMessenger(_l1CrossDomainMessenger)
    {
        initialize();
    }

    /**
     * @notice Initializer.
     */
<<<<<<< HEAD
    function initialize(address _l1CrossDomainMessenger) public initializer {
        address[] memory blockedSystemAddresses = new address[](2);
        blockedSystemAddresses[0] = Predeploys.L2_CROSS_DOMAIN_MESSENGER;
        blockedSystemAddresses[1] = Predeploys.L2_TO_L1_MESSAGE_PASSER;
        __CrossDomainMessenger_init(_l1CrossDomainMessenger, blockedSystemAddresses);
>>>>>>> v0.5.23
=======
    function initialize() public initializer {
        __CrossDomainMessenger_init();
>>>>>>> @eth-optimism/l2geth@0.5.27
    }

    /**
     * @custom:legacy
     * @notice Legacy getter for the remote messenger. Use otherMessenger going forward.
     *
     * @return Address of the L1CrossDomainMessenger contract.
     */
<<<<<<< HEAD
    function l1CrossDomainMessenger() public returns (address) {
=======
    function l1CrossDomainMessenger() public view returns (address) {
>>>>>>> v0.5.23
        return otherMessenger;
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
<<<<<<< HEAD
<<<<<<< HEAD
        L2ToL1MessagePasser(payable(Lib_PredeployAddresses.L2_TO_L1_MESSAGE_PASSER))
            .initiateWithdrawal{ value: _value }(_to, _gasLimit, _data);
=======
        L2ToL1MessagePasser(payable(PredeployAddresses.L2_TO_L1_MESSAGE_PASSER)).initiateWithdrawal{
=======
        L2ToL1MessagePasser(payable(Predeploys.L2_TO_L1_MESSAGE_PASSER)).initiateWithdrawal{
>>>>>>> v0.5.24
            value: _value
        }(_to, _gasLimit, _data);
>>>>>>> v0.5.23
    }

    /**
     * @inheritdoc CrossDomainMessenger
     */
    function _isOtherMessenger() internal view override returns (bool) {
        return AddressAliasHelper.undoL1ToL2Alias(msg.sender) == otherMessenger;
    }

    /**
     * @inheritdoc CrossDomainMessenger
     */
    function _isUnsafeTarget(address _target) internal view override returns (bool) {
        return _target == address(this) || _target == address(Predeploys.L2_TO_L1_MESSAGE_PASSER);
    }
}
