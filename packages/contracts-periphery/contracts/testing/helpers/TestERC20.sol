// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestERC20 is ERC20 {
    constructor() ERC20("TEST", "TST") {}
=======
import { ERC20 } from "@rari-capital/solmate/src/tokens/ERC20.sol";

contract TestERC20 is ERC20 {
    constructor() ERC20("TEST", "TST", 18) {}
>>>>>>> v0.5.23

    function mint(address to, uint256 value) public {
        _mint(to, value);
    }
}
