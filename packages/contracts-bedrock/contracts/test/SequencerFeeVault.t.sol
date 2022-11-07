// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Bridge_Initializer } from "./CommonTest.t.sol";

import { SequencerFeeVault } from "../L2/SequencerFeeVault.sol";
<<<<<<< HEAD
import { L2StandardBridge } from "../L2/L2StandardBridge.sol";
<<<<<<< HEAD
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
=======
=======
import { StandardBridge } from "../universal/StandardBridge.sol";
>>>>>>> @eth-optimism/l2geth@0.5.29
import { Predeploys } from "../libraries/Predeploys.sol";

contract SequencerFeeVault_Test is Bridge_Initializer {
<<<<<<< HEAD
    SequencerFeeVault vault =
        SequencerFeeVault(payable(Predeploys.SEQUENCER_FEE_WALLET));
>>>>>>> v0.5.24
=======
    SequencerFeeVault vault = SequencerFeeVault(payable(Predeploys.SEQUENCER_FEE_WALLET));
>>>>>>> @eth-optimism/l2geth@0.5.27
    address constant recipient = address(256);

    event Withdrawal(uint256 value, address to, address from);

    function setUp() public override {
        super.setUp();
<<<<<<< HEAD

<<<<<<< HEAD
        vm.etch(
<<<<<<< HEAD
<<<<<<< HEAD
            Lib_PredeployAddresses.SEQUENCER_FEE_WALLET,
=======
            PredeployAddresses.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.23
=======
            Predeploys.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.24
            address(new SequencerFeeVault()).code
        );
=======
        vm.etch(Predeploys.SEQUENCER_FEE_WALLET, address(new SequencerFeeVault()).code);
>>>>>>> @eth-optimism/l2geth@0.5.27

        vm.store(
<<<<<<< HEAD
<<<<<<< HEAD
            Lib_PredeployAddresses.SEQUENCER_FEE_WALLET,
=======
            PredeployAddresses.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.23
=======
            Predeploys.SEQUENCER_FEE_WALLET,
>>>>>>> v0.5.24
            bytes32(uint256(0)),
            bytes32(uint256(uint160(recipient)))
        );
=======
        vm.etch(Predeploys.SEQUENCER_FEE_WALLET, address(new SequencerFeeVault(recipient)).code);
>>>>>>> @eth-optimism/l2geth@0.5.29
    }

    function test_minWithdrawalAmount() external {
        assertEq(vault.MIN_WITHDRAWAL_AMOUNT(), 10 ether);
    }

    function test_constructor() external {
        assertEq(vault.l1FeeWallet(), recipient);
    }

    function test_receive() external {
        assertEq(address(vault).balance, 0);

        vm.prank(alice);
        (bool success, ) = address(vault).call{ value: 100 }(hex"");

        assertEq(success, true);
        assertEq(address(vault).balance, 100);
    }

    function test_revertWithdraw() external {
        assert(address(vault).balance < vault.MIN_WITHDRAWAL_AMOUNT());

        vm.expectRevert(
            "FeeVault: withdrawal amount must be greater than minimum withdrawal amount"
        );
        vault.withdraw();
    }

    function test_withdraw() external {
        vm.deal(address(vault), vault.MIN_WITHDRAWAL_AMOUNT() + 1);

        vm.expectEmit(true, true, true, true);
        emit Withdrawal(address(vault).balance, vault.RECIPIENT(), address(this));

        vm.expectCall(
<<<<<<< HEAD
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
=======
            Predeploys.L2_STANDARD_BRIDGE,
            address(vault).balance,
            abi.encodeWithSelector(
<<<<<<< HEAD
                L2StandardBridge.withdrawTo.selector,
                Predeploys.LEGACY_ERC20_ETH,
>>>>>>> v0.5.24
=======
                StandardBridge.bridgeETHTo.selector,
>>>>>>> @eth-optimism/l2geth@0.5.29
                vault.l1FeeWallet(),
                20000,
                bytes("")
            )
        );

        vault.withdraw();
    }
}
