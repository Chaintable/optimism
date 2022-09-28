// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { L1Block } from "../L2/L1Block.sol";
<<<<<<< HEAD
import { L1BlockNumber } from "../L2/L1BlockNumber.sol";
<<<<<<< HEAD
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";
=======
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";
>>>>>>> v0.5.23
=======
import { L1BlockNumber } from "../legacy/L1BlockNumber.sol";
import { Predeploys } from "../libraries/Predeploys.sol";
>>>>>>> v0.5.24

contract L1BlockNumberTest is Test {
    L1Block lb;
    L1BlockNumber bn;

    function setUp() external {
<<<<<<< HEAD
<<<<<<< HEAD
        vm.etch(Lib_PredeployAddresses.L1_BLOCK_ATTRIBUTES, address(new L1Block()).code);
        lb = L1Block(Lib_PredeployAddresses.L1_BLOCK_ATTRIBUTES);
=======
        vm.etch(PredeployAddresses.L1_BLOCK_ATTRIBUTES, address(new L1Block()).code);
        lb = L1Block(PredeployAddresses.L1_BLOCK_ATTRIBUTES);
>>>>>>> v0.5.23
=======
        vm.etch(Predeploys.L1_BLOCK_ATTRIBUTES, address(new L1Block()).code);
        lb = L1Block(Predeploys.L1_BLOCK_ATTRIBUTES);
>>>>>>> v0.5.24
        bn = new L1BlockNumber();
        vm.prank(lb.DEPOSITOR_ACCOUNT());
        lb.setL1BlockValues(uint64(999), uint64(2), 3, keccak256(abi.encode(1)), uint64(4));
    }

    function test_getL1BlockNumber() external {
        assertEq(bn.getL1BlockNumber(), 999);
    }

    function test_fallback() external {
        (bool success, bytes memory ret) = address(bn).call(hex"");
        assertEq(success, true);
        assertEq(ret, abi.encode(999));
    }

    function test_receive() external {
        (bool success, bytes memory ret) = address(bn).call{ value: 1 }(hex"");
        assertEq(success, true);
        assertEq(ret, abi.encode(999));
    }
}
