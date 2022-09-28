// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
/* Library Imports */
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";

/* Contract Imports */
import { L2StandardBridge } from "./L2StandardBridge.sol";
=======
import { Semver } from "../universal/Semver.sol";
import { L2StandardBridge } from "./L2StandardBridge.sol";
<<<<<<< HEAD
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";
>>>>>>> v0.5.23
=======
import { Predeploys } from "../libraries/Predeploys.sol";
>>>>>>> v0.5.24

/**
 * @custom:proxied
 * @custom:predeploy 0x4200000000000000000000000000000000000011
 * @title SequencerFeeVault
 * @notice The SequencerFeeVault is the contract that holds any fees paid to the Sequencer during
 *         transaction processing and block production.
 */
<<<<<<< HEAD
contract SequencerFeeVault {
=======
contract SequencerFeeVault is Semver {
>>>>>>> v0.5.23
    /**
     * @notice Minimum balance before a withdrawal can be triggered.
     */
    uint256 public constant MIN_WITHDRAWAL_AMOUNT = 15 ether;

    /**
     * @notice Wallet that will receive the fees on L1.
     */
    address public l1FeeWallet;

    /**
<<<<<<< HEAD
=======
     * @custom:semver 0.0.1
     */
    constructor() Semver(0, 0, 1) {}

    /**
>>>>>>> v0.5.23
     * @notice Allow the contract to receive ETH.
     */
    receive() external payable {}

    /**
     * @notice Triggers a withdrawal of funds to the L1 fee wallet.
     */
    function withdraw() external {
        require(
            address(this).balance >= MIN_WITHDRAWAL_AMOUNT,
            "SequencerFeeVault: withdrawal amount must be greater than minimum withdrawal amount"
        );

<<<<<<< HEAD
        uint256 balance = address(this).balance;

<<<<<<< HEAD
        L2StandardBridge(payable(Lib_PredeployAddresses.L2_STANDARD_BRIDGE)).withdrawTo{
            value: balance
        }(Lib_PredeployAddresses.OVM_ETH, l1FeeWallet, balance, 0, bytes(""));
=======
        L2StandardBridge(payable(PredeployAddresses.L2_STANDARD_BRIDGE)).withdrawTo{
            value: balance
        }(PredeployAddresses.LEGACY_ERC20_ETH, l1FeeWallet, balance, 0, bytes(""));
>>>>>>> v0.5.23
=======
        L2StandardBridge(payable(Predeploys.L2_STANDARD_BRIDGE)).withdrawTo{
            value: address(this).balance
        }(Predeploys.LEGACY_ERC20_ETH, l1FeeWallet, address(this).balance, 0, bytes(""));
>>>>>>> v0.5.24
    }
}
