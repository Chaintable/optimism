// SPDX-License-Identifier: MIT
<<<<<<< HEAD
pragma solidity 0.8.10;

/**
 * @title Burner
 * @notice Burner self-destructs on creation and sends all ETH to itself, removing all ETH given to
<<<<<<< HEAD
 * the contract from the circulating supply.
=======
 *         the contract from the circulating supply. Self-destructing is the only way to remove ETH
 *         from the circulating supply.
>>>>>>> v0.5.23
 */
contract Burner {
    constructor() payable {
        selfdestruct(payable(address(this)));
    }
}
=======
pragma solidity 0.8.15;
>>>>>>> v0.5.24

/**
 * @title Burn
 * @notice Utilities for burning stuff.
 */
library Burn {
    /**
     * Burns a given amount of ETH.
     *
     * @param _amount Amount of ETH to burn.
     */
    function eth(uint256 _amount) internal {
        new Burner{ value: _amount }();
    }

    /**
     * Burns a given amount of gas.
     *
     * @param _amount Amount of gas to burn.
     */
<<<<<<< HEAD
    function gas(uint256 _amount) internal {
=======
    function gas(uint256 _amount) internal view {
>>>>>>> v0.5.23
        uint256 i = 0;
        uint256 initialGas = gasleft();
        while (initialGas - gasleft() < _amount) {
            ++i;
        }
    }
}

/**
 * @title Burner
 * @notice Burner self-destructs on creation and sends all ETH to itself, removing all ETH given to
 *         the contract from the circulating supply. Self-destructing is the only way to remove ETH
 *         from the circulating supply.
 */
contract Burner {
    constructor() payable {
        selfdestruct(payable(address(this)));
    }
}
