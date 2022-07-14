//SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

/* Testing utilities */
import { Messenger_Initializer } from "./CommonTest.t.sol";
import { L2OutputOracle_Initializer } from "./L2OutputOracle.t.sol";

/* Libraries */
<<<<<<< HEAD
import { AddressAliasHelper } from "../libraries/AddressAliasHelper.sol";
import { Lib_DefaultValues } from "../libraries/Lib_DefaultValues.sol";
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";
import { Lib_CrossDomainUtils } from "../libraries/Lib_CrossDomainUtils.sol";
import { WithdrawalVerifier } from "../libraries/Lib_WithdrawalVerifier.sol";
=======
import { AddressAliasHelper } from "../vendor/AddressAliasHelper.sol";
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";
import { Hashing } from "../libraries/Hashing.sol";
import { Encoding } from "../libraries/Encoding.sol";
>>>>>>> v0.5.23

/* Target contract dependencies */
import { L2OutputOracle } from "../L1/L2OutputOracle.sol";
import { OptimismPortal } from "../L1/OptimismPortal.sol";

<<<<<<< HEAD
import { CrossDomainHashing } from "../libraries/Lib_CrossDomainHashing.sol";

=======
>>>>>>> v0.5.23
/* Target contract */
import { L1CrossDomainMessenger } from "../L1/L1CrossDomainMessenger.sol";

contract L1CrossDomainMessenger_Test is Messenger_Initializer {
    // Receiver address for testing
    address recipient = address(0xabbaacdc);

    function setUp() public override {
        super.setUp();
    }

    // pause: should pause the contract when called by the current owner
    function test_L1MessengerPause() external {
        L1Messenger.pause();
        assert(L1Messenger.paused());
    }

    // pause: should not pause the contract when called by account other than the owner
    function testCannot_L1MessengerPause() external {
        vm.expectRevert("Ownable: caller is not the owner");
        vm.prank(address(0xABBA));
        L1Messenger.pause();
    }

<<<<<<< HEAD
    // the version is encoded in the nonce
    function test_L1MessengerMessageVersion() external {
        assertEq(
            CrossDomainHashing.getVersionFromNonce(L1Messenger.messageNonce()),
=======
    // unpause: should unpause the contract when called by the current owner
    function test_L1MessengerUnpause() external {
        L1Messenger.pause();
        assert(L1Messenger.paused());
        L1Messenger.unpause();
        assert(!L1Messenger.paused());
    }

    // unpause: should not unpause the contract when called by account other than the owner
    function testCannot_L1MessengerUnpause() external {
        vm.expectRevert("Ownable: caller is not the owner");
        vm.prank(address(0xABBA));
        L1Messenger.unpause();
    }

    // the version is encoded in the nonce
    function test_L1MessengerMessageVersion() external {
        (,uint16 version) = Encoding.decodeVersionedNonce(L1Messenger.messageNonce());
        assertEq(
            version,
>>>>>>> v0.5.23
            L1Messenger.MESSAGE_VERSION()
        );
    }

    // sendMessage: should be able to send a single message
    // TODO: this same test needs to be done with the legacy message type
    // by setting the message version to 0
    function test_L1MessengerSendMessage() external {
        // deposit transaction on the optimism portal should be called
        vm.expectCall(
            address(op),
            abi.encodeWithSelector(
                OptimismPortal.depositTransaction.selector,
<<<<<<< HEAD
                Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER,
                0,
                100 + L1Messenger.baseGas(hex"ff"),
                false,
                CrossDomainHashing.getVersionedEncoding(
=======
                PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER,
                0,
                100 + L1Messenger.baseGas(hex"ff"),
                false,
                Encoding.encodeCrossDomainMessage(
>>>>>>> v0.5.23
                    L1Messenger.messageNonce(),
                    alice,
                    recipient,
                    0,
                    100,
                    hex"ff"
                )
            )
        );

        // TransactionDeposited event
        vm.expectEmit(true, true, true, true);
        emit TransactionDeposited(
            AddressAliasHelper.applyL1ToL2Alias(address(L1Messenger)),
<<<<<<< HEAD
            Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER,
=======
            PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER,
>>>>>>> v0.5.23
            0,
            0,
            100 + L1Messenger.baseGas(hex"ff"),
            false,
<<<<<<< HEAD
            CrossDomainHashing.getVersionedEncoding(
=======
            Encoding.encodeCrossDomainMessage(
>>>>>>> v0.5.23
                L1Messenger.messageNonce(),
                alice,
                recipient,
                0,
                100,
                hex"ff"
            )
        );

        // SentMessage event
        vm.expectEmit(true, true, true, true);
<<<<<<< HEAD
        emit SentMessage(
           recipient,
           alice,
           hex"ff",
           L1Messenger.messageNonce(),
           100
        );
=======
        emit SentMessage(recipient, alice, hex"ff", L1Messenger.messageNonce(), 100);
>>>>>>> v0.5.23

        vm.prank(alice);
        L1Messenger.sendMessage(recipient, hex"ff", uint32(100));
    }

    // sendMessage: should be able to send the same message twice
    function test_L1MessengerTwiceSendMessage() external {
        uint256 nonce = L1Messenger.messageNonce();
        L1Messenger.sendMessage(recipient, hex"aa", uint32(500_000));
        L1Messenger.sendMessage(recipient, hex"aa", uint32(500_000));
        // the nonce increments for each message sent
<<<<<<< HEAD
        assertEq(
            nonce + 2,
            L1Messenger.messageNonce()
        );
=======
        assertEq(nonce + 2, L1Messenger.messageNonce());
>>>>>>> v0.5.23
    }

    function test_L1MessengerXDomainSenderReverts() external {
        vm.expectRevert("xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();
    }

    // xDomainMessageSender: should return the xDomainMsgSender address
    // TODO: might need a test contract
    // function test_xDomainSenderSetCorrectly() external {}

    // relayMessage: should send a successful call to the target contract
    function test_L1MessengerRelayMessageSucceeds() external {
        address target = address(0xabcd);
<<<<<<< HEAD
        address sender = Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;
=======
        address sender = PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;
>>>>>>> v0.5.23

        vm.expectCall(target, hex"1111");

        // set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
<<<<<<< HEAD
        uint256 senderSlotIndex = 1;
=======
        uint256 senderSlotIndex = 51;
>>>>>>> v0.5.23
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));
        vm.prank(address(op));

        vm.expectEmit(true, true, true, true);

<<<<<<< HEAD
        bytes32 hash = CrossDomainHashing.getVersionedHash(
            0,
            sender,
            target,
            0,
            0,
            hex"1111"
        );
=======
        bytes32 hash = Hashing.hashCrossDomainMessage(0, sender, target, 0, 0, hex"1111");
>>>>>>> v0.5.23

        emit RelayedMessage(hash);

        L1Messenger.relayMessage(
            0, // nonce
            sender,
            target,
            0, // value
            0,
            hex"1111"
        );

        // the message hash is in the successfulMessages mapping
        assert(L1Messenger.successfulMessages(hash));
        // it is not in the received messages mapping
        assertEq(L1Messenger.receivedMessages(hash), false);
    }

    // relayMessage: should revert if attempting to relay a message sent to an L1 system contract
    function test_L1MessengerRelayMessageToSystemContract() external {
        // set the target to be the OptimismPortal
        address target = address(op);
<<<<<<< HEAD
        address sender = Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;
=======
        address sender = PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;
>>>>>>> v0.5.23
        bytes memory message = hex"1111";

        // set the value of op.l2Sender() to be the L2 Cross Domain Messenger.
        vm.prank(address(op));
        vm.expectRevert("Message cannot be replayed.");
        L1Messenger.relayMessage(0, sender, target, 0, 0, message);

        vm.store(address(op), 0, bytes32(abi.encode(sender)));
        vm.expectRevert("Message cannot be replayed.");
        L1Messenger.relayMessage(0, sender, target, 0, 0, message);
    }

    // relayMessage: the xDomainMessageSender is reset to the original value
    function test_L1MessengerxDomainMessageSenderResets() external {
        vm.expectRevert("xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();

<<<<<<< HEAD
        address sender = Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;
        uint256 senderSlotIndex = 1;
=======
        address sender = PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER;

        uint256 senderSlotIndex = 51;
        bytes32 slotValue = vm.load(address(op), bytes32(senderSlotIndex));

>>>>>>> v0.5.23
        vm.store(address(op), bytes32(senderSlotIndex), bytes32(abi.encode(sender)));
        vm.prank(address(op));
        L1Messenger.relayMessage(0, address(0), address(0), 0, 0, hex"");

        vm.expectRevert("xDomainMessageSender is not set");
        L1Messenger.xDomainMessageSender();
    }

    // relayMessage: should revert if paused
    function test_L1MessengerRelayShouldRevertIfPaused() external {
        vm.prank(L1Messenger.owner());
        L1Messenger.pause();

        vm.expectRevert("Pausable: paused");
        L1Messenger.relayMessage(0, address(0), address(0), 0, 0, hex"");
    }
}
