// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { ERC20 } from "@rari-capital/solmate/src/tokens/ERC20.sol";
import { ERC721 } from "@rari-capital/solmate/src/tokens/ERC721.sol";
import { Transactor } from "./Transactor.sol";

/**
 * @title AssetReceiver
 * @notice AssetReceiver is a minimal contract for receiving funds assets in the form of either
 * ETH, ERC20 tokens, or ERC721 tokens. Only the contract owner may withdraw the assets.
 */
contract AssetReceiver is Transactor {
    /**
<<<<<<< HEAD
     * Emitted when ETH is received by this address.
=======
     * @notice Emitted when ETH is received by this address.
     *
<<<<<<< HEAD
     * @param from Address that sent ETH to this contract.
>>>>>>> v0.5.23
=======
     * @param from   Address that sent ETH to this contract.
     * @param amount Amount of ETH received.
>>>>>>> @eth-optimism/l2geth@0.5.27
     */
    event ReceivedETH(address indexed from, uint256 amount);

    /**
<<<<<<< HEAD
     * Emitted when ETH is withdrawn from this address.
=======
     * @notice Emitted when ETH is withdrawn from this address.
     *
     * @param withdrawer Address that triggered the withdrawal.
     * @param recipient  Address that received the withdrawal.
     * @param amount     ETH amount withdrawn.
>>>>>>> v0.5.23
     */
    event WithdrewETH(address indexed withdrawer, address indexed recipient, uint256 amount);

    /**
<<<<<<< HEAD
     * Emitted when ERC20 tokens are withdrawn from this address.
=======
     * @notice Emitted when ERC20 tokens are withdrawn from this address.
     *
     * @param withdrawer Address that triggered the withdrawal.
     * @param recipient  Address that received the withdrawal.
     * @param asset      Address of the token being withdrawn.
     * @param amount     ERC20 amount withdrawn.
>>>>>>> v0.5.23
     */
    event WithdrewERC20(
        address indexed withdrawer,
        address indexed recipient,
        address indexed asset,
        uint256 amount
    );

    /**
<<<<<<< HEAD
     * Emitted when ERC721 tokens are withdrawn from this address.
=======
     * @notice Emitted when ERC20 tokens are withdrawn from this address.
     *
     * @param withdrawer Address that triggered the withdrawal.
     * @param recipient  Address that received the withdrawal.
     * @param asset      Address of the token being withdrawn.
     * @param id         Token ID being withdrawn.
>>>>>>> v0.5.23
     */
    event WithdrewERC721(
        address indexed withdrawer,
        address indexed recipient,
        address indexed asset,
        uint256 id
    );

    /**
     * @param _owner Initial contract owner.
     */
    constructor(address _owner) Transactor(_owner) {}

    /**
<<<<<<< HEAD
     * Make sure we can receive ETH.
=======
     * @notice Make sure we can receive ETH.
>>>>>>> v0.5.23
     */
    receive() external payable {
        emit ReceivedETH(msg.sender, msg.value);
    }

    /**
<<<<<<< HEAD
     * Withdraws full ETH balance to the recipient.
=======
     * @notice Withdraws full ETH balance to the recipient.
>>>>>>> v0.5.23
     *
     * @param _to Address to receive the ETH balance.
     */
    function withdrawETH(address payable _to) external onlyOwner {
        withdrawETH(_to, address(this).balance);
    }

    /**
<<<<<<< HEAD
     * Withdraws partial ETH balance to the recipient.
     *
     * @param _to Address to receive the ETH balance.
=======
     * @notice Withdraws partial ETH balance to the recipient.
     *
     * @param _to     Address to receive the ETH balance.
>>>>>>> v0.5.23
     * @param _amount Amount of ETH to withdraw.
     */
    function withdrawETH(address payable _to, uint256 _amount) public onlyOwner {
        // slither-disable-next-line reentrancy-unlimited-gas
        (bool success, ) = _to.call{ value: _amount }("");
        emit WithdrewETH(msg.sender, _to, _amount);
    }

    /**
<<<<<<< HEAD
     * Withdraws full ERC20 balance to the recipient.
     *
     * @param _asset ERC20 token to withdraw.
     * @param _to Address to receive the ERC20 balance.
=======
     * @notice Withdraws full ERC20 balance to the recipient.
     *
     * @param _asset ERC20 token to withdraw.
     * @param _to    Address to receive the ERC20 balance.
>>>>>>> v0.5.23
     */
    function withdrawERC20(ERC20 _asset, address _to) external onlyOwner {
        withdrawERC20(_asset, _to, _asset.balanceOf(address(this)));
    }

    /**
<<<<<<< HEAD
     * Withdraws partial ERC20 balance to the recipient.
     *
     * @param _asset ERC20 token to withdraw.
     * @param _to Address to receive the ERC20 balance.
=======
     * @notice Withdraws partial ERC20 balance to the recipient.
     *
     * @param _asset  ERC20 token to withdraw.
     * @param _to     Address to receive the ERC20 balance.
>>>>>>> v0.5.23
     * @param _amount Amount of ERC20 to withdraw.
     */
    function withdrawERC20(
        ERC20 _asset,
        address _to,
        uint256 _amount
    ) public onlyOwner {
        // slither-disable-next-line unchecked-transfer
        _asset.transfer(_to, _amount);
        // slither-disable-next-line reentrancy-events
        emit WithdrewERC20(msg.sender, _to, address(_asset), _amount);
    }

    /**
<<<<<<< HEAD
     * Withdraws ERC721 token to the recipient.
     *
     * @param _asset ERC721 token to withdraw.
     * @param _to Address to receive the ERC721 token.
     * @param _id Token ID of the ERC721 token to withdraw.
=======
     * @notice Withdraws ERC721 token to the recipient.
     *
     * @param _asset ERC721 token to withdraw.
     * @param _to    Address to receive the ERC721 token.
     * @param _id    Token ID of the ERC721 token to withdraw.
>>>>>>> v0.5.23
     */
    function withdrawERC721(
        ERC721 _asset,
        address _to,
        uint256 _id
    ) external onlyOwner {
        _asset.transferFrom(address(this), _to, _id);
        // slither-disable-next-line reentrancy-events
        emit WithdrewERC721(msg.sender, _to, address(_asset), _id);
    }
}
