// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import { Bridge_Initializer } from "./CommonTest.t.sol";

import { SequencerFeeVault } from "../L2/SequencerFeeVault.sol";
import { L2StandardBridge } from "../L2/L2StandardBridge.sol";
<<<<<<< HEAD
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";

contract SequencerFeeVault_Test is Bridge_Initializer {
    SequencerFeeVault vault =
        SequencerFeeVault(payable(Lib_PredeployAddresses.SEQUENCER_FEE_WALLET));
=======
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";

contract SequencerFeeVault_Test is Bridge_Initializer {
    SequencerFeeVault vault =
        SequencerFeeVault(payable(PredeployAddresses.SEQUENCER_FEE_WALLET));
>>>>>>> v0.5.23
    address constant recipient = address(256);

    function setUp() public override {
        super.setUp();

        vm.etch(
<<<<<<< HEAD
            Lib_PredeployAddresses.SEQUENCER_FEE_WALLET,
=======
            PredeployAddresses.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.23
            address(new SequencerFeeVault()).code
        );

        vm.store(
<<<<<<< HEAD
            Lib_PredeployAddresses.SEQUENCER_FEE_WALLET,
=======
            PredeployAddresses.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.23
            bytes32(uint256(0)),
            bytes32(uint256(uint160(recipient)))
        );
    }

    function test_minWithdrawalAmount() external {
        assertEq(
            vault.MIN_WITHDRAWAL_AMOUNT(),
            15 ether
        );
    }

    function test_constructor() external {
        assertEq(
            vault.l1FeeWallet(),
            recipient
        );
    }

    function test_receive() external {
        assertEq(
            address(vault).balance,
            0
        );

        vm.prank(alice);
        address(vault).call{ value: 100 }(hex"");

        assertEq(
            address(vault).balance,
            100
        );
    }

    function test_revertWithdraw() external {
        assert(address(vault).balance < vault.MIN_WITHDRAWAL_AMOUNT());

        vm.expectRevert(
            "OVM_SequencerFeeVault: withdrawal amount must be greater than minimum withdrawal amount"
        );
        vault.withdraw();
    }

    function test_withdraw() external {
        vm.deal(address(vault), vault.MIN_WITHDRAWAL_AMOUNT() + 1);

        vm.expectCall(
<<<<<<< HEAD
            Lib_PredeployAddresses.L2_STANDARD_BRIDGE,
            abi.encodeWithSelector(
                L2StandardBridge.withdrawTo.selector,
                Lib_PredeployAddresses.OVM_ETH,
=======
            PredeployAddresses.L2_STANDARD_BRIDGE,
            abi.encodeWithSelector(
                L2StandardBridge.withdrawTo.selector,
                PredeployAddresses.LEGACY_ERC20_ETH,
>>>>>>> v0.5.23
                vault.l1FeeWallet(),
                address(vault).balance,
                0,
                bytes("")
            )
        );

        vault.withdraw();
    }
}
