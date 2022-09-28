<<<<<<< HEAD
//SPDX-License-Identifier: MIT
pragma solidity 0.8.10;

import { L1Block } from "./L1Block.sol";
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";
=======
// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { L1Block } from "../L2/L1Block.sol";
import { Predeploys } from "../libraries/Predeploys.sol";
import { Semver } from "../universal/Semver.sol";
>>>>>>> v0.5.23

/**
 * @custom:legacy
 * @custom:proxied
 * @custom:predeploy 0x4200000000000000000000000000000000000013
 * @title L1BlockNumber
 * @notice L1BlockNumber is a legacy contract that fills the roll of the OVM_L1BlockNumber contract
 *         in the old version of the Optimism system. Only necessary for backwards compatibility.
 *         If you want to access the L1 block number going forward, you should use the L1Block
 *         contract instead.
 */
<<<<<<< HEAD
contract L1BlockNumber {
=======
contract L1BlockNumber is Semver {
    /**
     * @custom:semver 0.0.1
     */
    constructor() Semver(0, 0, 1) {}

>>>>>>> v0.5.23
    /**
     * @notice Returns the L1 block number.
     */
    receive() external payable {
        uint256 l1BlockNumber = getL1BlockNumber();
        assembly {
            mstore(0, l1BlockNumber)
            return(0, 32)
        }
    }

    /**
     * @notice Returns the L1 block number.
     */
    // solhint-disable-next-line no-complex-fallback
    fallback() external payable {
        uint256 l1BlockNumber = getL1BlockNumber();
        assembly {
            mstore(0, l1BlockNumber)
            return(0, 32)
        }
    }

    /**
     * @notice Retrieves the latest L1 block number.
     *
     * @return Latest L1 block number.
     */
    function getL1BlockNumber() public view returns (uint256) {
<<<<<<< HEAD:packages/contracts-bedrock/contracts/L2/L1BlockNumber.sol
<<<<<<< HEAD
        return L1Block(Lib_PredeployAddresses.L1_BLOCK_ATTRIBUTES).number();
=======
        return L1Block(PredeployAddresses.L1_BLOCK_ATTRIBUTES).number();
>>>>>>> v0.5.23
=======
        return L1Block(Predeploys.L1_BLOCK_ATTRIBUTES).number();
>>>>>>> v0.5.24:packages/contracts-bedrock/contracts/legacy/L1BlockNumber.sol
    }
}
