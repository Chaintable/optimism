// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
<<<<<<< HEAD
import { Lib_PredeployAddresses } from "../libraries/Lib_PredeployAddresses.sol";
import { StandardBridge } from "../universal/StandardBridge.sol";
=======
import { PredeployAddresses } from "../libraries/PredeployAddresses.sol";
=======
import { Predeploys } from "../libraries/Predeploys.sol";
>>>>>>> v0.5.24
import { StandardBridge } from "../universal/StandardBridge.sol";
import { Semver } from "../universal/Semver.sol";
>>>>>>> v0.5.23
import { OptimismMintableERC20 } from "../universal/OptimismMintableERC20.sol";

/**
 * @custom:proxied
 * @custom:predeploy 0x4200000000000000000000000000000000000010
 * @title L2StandardBridge
 * @notice The L2StandardBridge is responsible for transfering ETH and ERC20 tokens between L1 and
<<<<<<< HEAD
 *         L2. ERC20 tokens sent to L1 are escrowed within this contract.
<<<<<<< HEAD
 */
contract L2StandardBridge is StandardBridge {
=======
=======
 *         L2. In the case that an ERC20 token is native to L2, it will be escrowed within this
 *         contract. If the ERC20 token is native to L1, it will be burnt.
>>>>>>> @eth-optimism/l2geth@0.5.27
 *         Note that this contract is not intended to support all variations of ERC20 tokens.
 *         Examples of some token types that may not be properly supported by this contract include,
 *         but are not limited to: tokens with transfer fees, rebasing tokens, and
 *         tokens with blocklists.
 */
contract L2StandardBridge is StandardBridge, Semver {
>>>>>>> v0.5.23
    /**
     * @custom:legacy
     * @notice Emitted whenever a withdrawal from L2 to L1 is initiated.
     *
<<<<<<< HEAD
     * @param _l1Token Address of the token on L1.
     * @param _l2Token Address of the corresponding token on L2.
     * @param _from    Address of the withdrawer.
     * @param _to      Address of the recipient on L1.
     * @param _amount  Amount of the ERC20 withdrawn.
     * @param _data    Extra data attached to the withdrawal.
     */
    event WithdrawalInitiated(
        address indexed _l1Token,
        address indexed _l2Token,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
=======
     * @param l1Token   Address of the token on L1.
     * @param l2Token   Address of the corresponding token on L2.
     * @param from      Address of the withdrawer.
     * @param to        Address of the recipient on L1.
     * @param amount    Amount of the ERC20 withdrawn.
     * @param extraData Extra data attached to the withdrawal.
     */
    event WithdrawalInitiated(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
>>>>>>> v0.5.23
    );

    /**
     * @custom:legacy
     * @notice Emitted whenever an ERC20 deposit is finalized.
     *
<<<<<<< HEAD
     * @param _l1Token Address of the token on L1.
     * @param _l2Token Address of the corresponding token on L2.
     * @param _from    Address of the depositor.
     * @param _to      Address of the recipient on L2.
     * @param _amount  Amount of the ERC20 deposited.
     * @param _data    Extra data attached to the deposit.
     */
    event DepositFinalized(
        address indexed _l1Token,
        address indexed _l2Token,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
=======
     * @param l1Token   Address of the token on L1.
     * @param l2Token   Address of the corresponding token on L2.
     * @param from      Address of the depositor.
     * @param to        Address of the recipient on L2.
     * @param amount    Amount of the ERC20 deposited.
     * @param extraData Extra data attached to the deposit.
     */
    event DepositFinalized(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
>>>>>>> v0.5.23
    );

    /**
<<<<<<< HEAD
     * @custom:legacy
     * @notice Emitted whenever a deposit fails.
     *
<<<<<<< HEAD
     * @param _l1Token Address of the token on L1.
     * @param _l2Token Address of the corresponding token on L2.
     * @param _from    Address of the depositor.
     * @param _to      Address of the recipient on L2.
     * @param _amount  Amount of the ERC20 deposited.
     * @param _data    Extra data attached to the deposit.
     */
    event DepositFailed(
        address indexed _l1Token,
        address indexed _l2Token,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
    );

    /**
     * @notice Initializes the L2StandardBridge.
     *
     * @param _otherBridge Address of the L1StandardBridge.
     */
    function initialize(address payable _otherBridge) public {
        _initialize(payable(Lib_PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER), _otherBridge);
=======
     * @param l1Token   Address of the token on L1.
     * @param l2Token   Address of the corresponding token on L2.
     * @param from      Address of the depositor.
     * @param to        Address of the recipient on L2.
     * @param amount    Amount of the ERC20 deposited.
     * @param extraData Extra data attached to the deposit.
     */
    event DepositFailed(
        address indexed l1Token,
        address indexed l2Token,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
     * @custom:semver 0.0.2
     *
     * @param _otherBridge Address of the L1StandardBridge.
     */
<<<<<<< HEAD
    constructor(address payable _otherBridge) Semver(0, 0, 1) {
        initialize(_otherBridge);
    }

    /**
     * @notice Initializer.
     *
     * @param _otherBridge Address of the L1StandardBridge.
     */
    function initialize(address payable _otherBridge) public initializer {
        __StandardBridge_init(payable(PredeployAddresses.L2_CROSS_DOMAIN_MESSENGER), _otherBridge);
>>>>>>> v0.5.23
    }
=======
    constructor(address payable _otherBridge)
        Semver(0, 0, 2)
        StandardBridge(payable(Predeploys.L2_CROSS_DOMAIN_MESSENGER), _otherBridge)
    {}
>>>>>>> v0.5.24

    /**
     * @custom:legacy
     * @notice Initiates a withdrawal from L2 to L1.
     *
     * @param _l2Token     Address of the L2 token to withdraw.
     * @param _amount      Amount of the L2 token to withdraw.
     * @param _minGasLimit Minimum gas limit to use for the transaction.
<<<<<<< HEAD
     * @param _data        Extra data attached to the withdrawal.
=======
     * @param _extraData   Extra data attached to the withdrawal.
>>>>>>> v0.5.23
     */
    function withdraw(
        address _l2Token,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
    ) external payable virtual {
        _initiateWithdrawal(_l2Token, msg.sender, msg.sender, _amount, _minGasLimit, _data);
=======
        bytes calldata _extraData
    ) external payable virtual onlyEOA {
        _initiateWithdrawal(_l2Token, msg.sender, msg.sender, _amount, _minGasLimit, _extraData);
>>>>>>> v0.5.23
    }

    /**
     * @custom:legacy
     * @notice Initiates a withdrawal from L2 to L1 to a target account on L1.
<<<<<<< HEAD
=======
     *         Note that if ETH is sent to a contract on L1 and the call fails, then that ETH will
     *         be locked in the L1StandardBridge. ETH may be recoverable if the call can be
     *         successfully replayed by increasing the amount of gas supplied to the call. If the
     *         call will fail for any amount of gas, then the ETH will be locked permanently.
>>>>>>> v0.5.23
     *
     * @param _l2Token     Address of the L2 token to withdraw.
     * @param _to          Recipient account on L1.
     * @param _amount      Amount of the L2 token to withdraw.
     * @param _minGasLimit Minimum gas limit to use for the transaction.
<<<<<<< HEAD
     * @param _data        Extra data attached to the withdrawal.
=======
     * @param _extraData   Extra data attached to the withdrawal.
>>>>>>> v0.5.23
     */
    function withdrawTo(
        address _l2Token,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
    ) external payable virtual {
        _initiateWithdrawal(_l2Token, msg.sender, _to, _amount, _minGasLimit, _data);
=======
        bytes calldata _extraData
    ) external payable virtual {
        _initiateWithdrawal(_l2Token, msg.sender, _to, _amount, _minGasLimit, _extraData);
>>>>>>> v0.5.23
    }

    /**
     * @custom:legacy
     * @notice Finalizes a deposit from L1 to L2.
     *
<<<<<<< HEAD
     * @param _l1Token Address of the L1 token to deposit.
     * @param _l2Token Address of the corresponding L2 token.
     * @param _from    Address of the depositor.
     * @param _to      Address of the recipient.
     * @param _amount  Amount of the tokens being deposited.
     * @param _data    Extra data attached to the deposit.
=======
     * @param _l1Token   Address of the L1 token to deposit.
     * @param _l2Token   Address of the corresponding L2 token.
     * @param _from      Address of the depositor.
     * @param _to        Address of the recipient.
     * @param _amount    Amount of the tokens being deposited.
     * @param _extraData Extra data attached to the deposit.
>>>>>>> v0.5.23
     */
    function finalizeDeposit(
        address _l1Token,
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
<<<<<<< HEAD
        bytes calldata _data
    ) external payable virtual {
        if (_l1Token == address(0) && _l2Token == Lib_PredeployAddresses.OVM_ETH) {
            finalizeBridgeETH(_from, _to, _amount, _data);
        } else {
            finalizeBridgeERC20(_l2Token, _l1Token, _from, _to, _amount, _data);
        }
        emit DepositFinalized(_l1Token, _l2Token, _from, _to, _amount, _data);
=======
        bytes calldata _extraData
    ) external payable virtual {
        if (_l1Token == address(0) && _l2Token == Predeploys.LEGACY_ERC20_ETH) {
            finalizeBridgeETH(_from, _to, _amount, _extraData);
        } else {
            finalizeBridgeERC20(_l2Token, _l1Token, _from, _to, _amount, _extraData);
        }

        emit DepositFinalized(_l1Token, _l2Token, _from, _to, _amount, _extraData);
>>>>>>> v0.5.23
    }

    /**
     * @custom:legacy
     * @notice Internal function to a withdrawal from L2 to L1 to a target account on L1.
     *
     * @param _l2Token     Address of the L2 token to withdraw.
     * @param _from        Address of the withdrawer.
     * @param _to          Recipient account on L1.
     * @param _amount      Amount of the L2 token to withdraw.
     * @param _minGasLimit Minimum gas limit to use for the transaction.
<<<<<<< HEAD
     * @param _data        Extra data attached to the withdrawal.
=======
     * @param _extraData   Extra data attached to the withdrawal.
>>>>>>> v0.5.23
     */
    function _initiateWithdrawal(
        address _l2Token,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
    ) internal {
        address l1Token = OptimismMintableERC20(_l2Token).l1Token();
        if (_l2Token == Lib_PredeployAddresses.OVM_ETH) {
            _initiateBridgeETH(_from, _to, _amount, _minGasLimit, _data);
        } else {
            _initiateBridgeERC20(_l2Token, l1Token, _from, _to, _amount, _minGasLimit, _data);
        }
        emit WithdrawalInitiated(l1Token, _l2Token, msg.sender, _to, _amount, _data);
=======
        bytes calldata _extraData
    ) internal {
        address l1Token = OptimismMintableERC20(_l2Token).l1Token();
        if (_l2Token == Predeploys.LEGACY_ERC20_ETH) {
            _initiateBridgeETH(_from, _to, _amount, _minGasLimit, _extraData);
        } else {
            _initiateBridgeERC20(_l2Token, l1Token, _from, _to, _amount, _minGasLimit, _extraData);
        }

        emit WithdrawalInitiated(l1Token, _l2Token, _from, _to, _amount, _extraData);
>>>>>>> v0.5.23
    }
}
