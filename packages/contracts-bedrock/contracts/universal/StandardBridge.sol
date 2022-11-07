// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

<<<<<<< HEAD
/* Interface Imports */
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/* Library Imports */
import { ERC165Checker } from "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";
import { Address } from "@openzeppelin/contracts/utils/Address.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

=======
import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import { ERC165Checker } from "@openzeppelin/contracts/utils/introspection/ERC165Checker.sol";
import { Address } from "@openzeppelin/contracts/utils/Address.sol";
import { SafeERC20 } from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import { SafeCall } from "../libraries/SafeCall.sol";
<<<<<<< HEAD
import { IRemoteToken, IL1Token } from "./SupportedInterfaces.sol";
>>>>>>> v0.5.23
=======
import { IOptimismMintableERC20, ILegacyMintableERC20 } from "./SupportedInterfaces.sol";
>>>>>>> @eth-optimism/l2geth@0.5.27
import { CrossDomainMessenger } from "./CrossDomainMessenger.sol";
import { OptimismMintableERC20 } from "./OptimismMintableERC20.sol";

/**
 * @custom:upgradeable
 * @title StandardBridge
<<<<<<< HEAD
 * This contract can manage a 1:1 bridge between two domains for both
 * ETH (native asset) and ERC20s.
 * This contract should be deployed behind a proxy.
 * TODO: do we want a donateERC20 function as well?
 */
abstract contract StandardBridge {
    using SafeERC20 for IERC20;

    /**********
     * Events *
     **********/

    event ETHBridgeInitiated(
        address indexed _from,
        address indexed _to,
        uint256 _amount,
        bytes _data
    );

    event ETHBridgeFinalized(
        address indexed _from,
        address indexed _to,
        uint256 _amount,
        bytes _data
    );

    event ERC20BridgeInitiated(
        address indexed _localToken,
        address indexed _remoteToken,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
    );

    event ERC20BridgeFinalized(
        address indexed _localToken,
        address indexed _remoteToken,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
    );

    event ERC20BridgeFailed(
        address indexed _localToken,
        address indexed _remoteToken,
        address indexed _from,
        address _to,
        uint256 _amount,
        bytes _data
    );

    /*************
     * Variables *
     *************/

    /**
     * @notice The messenger contract on the same domain
=======
 * @notice StandardBridge is a base contract for the L1 and L2 standard ERC20 bridges.
 */
abstract contract StandardBridge {
    using SafeERC20 for IERC20;

    /**
     * @notice The L2 gas limit set when eth is depoisited using the receive() function.
     */
    uint32 internal constant RECEIVE_DEFAULT_GAS_LIMIT = 200_000;

    /**
     * @notice Messenger contract on this domain.
     */
    CrossDomainMessenger public immutable MESSENGER;

    /**
     * @notice Corresponding bridge on the other domain.
     */
    StandardBridge public immutable OTHER_BRIDGE;

    /**
     * @custom:legacy
     * @custom:spacer messenger
     * @notice Spacer for backwards compatibility.
     */
    address private spacer_0_0_20;

    /**
     * @custom:legacy
     * @custom:spacer l2TokenBridge
     * @notice Spacer for backwards compatibility.
     */
    address private spacer_1_0_20;

    /**
     * @notice Mapping that stores deposits for a given pair of local and remote tokens.
     */
    mapping(address => mapping(address => uint256)) public deposits;

    /**
     * @notice Reserve extra slots (to a total of 50) in the storage layout for future upgrades.
     *         A gap size of 47 was chosen here, so that the first slot used in a child contract
     *         would be a multiple of 50.
     */
    uint256[47] private __gap;

    /**
     * @notice Emitted when an ETH bridge is initiated to the other chain.
     *
     * @param from      Address of the sender.
     * @param to        Address of the receiver.
     * @param amount    Amount of ETH sent.
     * @param extraData Extra data sent with the transaction.
     */
    event ETHBridgeInitiated(
        address indexed from,
        address indexed to,
        uint256 amount,
        bytes extraData
    );

    /**
     * @notice Emitted when an ETH bridge is finalized on this chain.
     *
     * @param from      Address of the sender.
     * @param to        Address of the receiver.
     * @param amount    Amount of ETH sent.
     * @param extraData Extra data sent with the transaction.
     */
    event ETHBridgeFinalized(
        address indexed from,
        address indexed to,
        uint256 amount,
        bytes extraData
    );

    /**
     * @notice Emitted when an ERC20 bridge is initiated to the other chain.
     *
     * @param localToken  Address of the ERC20 on this chain.
     * @param remoteToken Address of the ERC20 on the remote chain.
     * @param from        Address of the sender.
     * @param to          Address of the receiver.
     * @param amount      Amount of the ERC20 sent.
     * @param extraData   Extra data sent with the transaction.
     */
    event ERC20BridgeInitiated(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
     * @notice Emitted when an ERC20 bridge is finalized on this chain.
     *
     * @param localToken  Address of the ERC20 on this chain.
     * @param remoteToken Address of the ERC20 on the remote chain.
     * @param from        Address of the sender.
     * @param to          Address of the receiver.
     * @param amount      Amount of the ERC20 sent.
     * @param extraData   Extra data sent with the transaction.
     */
    event ERC20BridgeFinalized(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
<<<<<<< HEAD
     * @notice Emitted when an ERC20 bridge to this chain fails.
     *
     * @param localToken  Address of the ERC20 on this chain.
     * @param remoteToken Address of the ERC20 on the remote chain.
     * @param from        Address of the sender.
     * @param to          Address of the receiver.
     * @param amount      Amount of ETH sent.
     * @param extraData   Extra data sent with the transaction.
     */
    event ERC20BridgeFailed(
        address indexed localToken,
        address indexed remoteToken,
        address indexed from,
        address to,
        uint256 amount,
        bytes extraData
    );

    /**
<<<<<<< HEAD
     * @notice The L2 gas limit set when eth is depoisited using the receive() function.
     */
    uint32 internal constant RECEIVE_DEFAULT_GAS_LIMIT = 200_000;

    /**
     * @notice Messenger contract on this domain.
>>>>>>> v0.5.23
     */
    CrossDomainMessenger public messenger;

    /**
<<<<<<< HEAD
     * @notice The corresponding bridge on the other domain
     */
    StandardBridge public otherBridge;

    mapping(address => mapping(address => uint256)) public deposits;

    /*************
     * Modifiers *
     *************/

    /**
     * @notice Only allow EOAs to call the functions. Note that this
     * is not safe against contracts calling code during their constructor
=======
     * @notice Corresponding bridge on the other domain.
     */
    StandardBridge public otherBridge;

    /**
     * @notice Mapping that stores deposits for a given pair of local and remote tokens.
     */
    mapping(address => mapping(address => uint256)) public deposits;

    /**
=======
>>>>>>> v0.5.24
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
     * @notice Only allow EOAs to call the functions. Note that this is not safe against contracts
     *         calling code within their constructors, but also doesn't really matter since we're
     *         just trying to prevent users accidentally depositing with smart contract wallets.
>>>>>>> v0.5.23
     */
    modifier onlyEOA() {
        require(
            !Address.isContract(msg.sender),
            "StandardBridge: function can only be called from an EOA"
        );
        _;
    }

    /**
<<<<<<< HEAD
     * @notice Ensures that the caller is the messenger, and that
     * it has the l2Sender value set to the address of the remote Token Bridge.
=======
     * @notice Ensures that the caller is a cross-chain message from the other bridge.
>>>>>>> v0.5.23
     */
    modifier onlyOtherBridge() {
        require(
            msg.sender == address(MESSENGER) &&
                MESSENGER.xDomainMessageSender() == address(OTHER_BRIDGE),
            "StandardBridge: function can only be called from the other bridge"
        );
        _;
    }

<<<<<<< HEAD
=======
    /**
<<<<<<< HEAD
     * @notice Ensures that the caller is this contract.
     */
>>>>>>> v0.5.23
    modifier onlySelf() {
        require(msg.sender == address(this), "StandardBridge: function can only be called by self");
        _;
    }

<<<<<<< HEAD
    /********************
     * Public Functions *
     ********************/

    /**
     * @notice Send ETH to this contract. This is used during upgrades
     */
    function donateETH() external payable {}

    /**
     * @notice EOAs can simply send ETH to this contract to have it be deposited
     * to L2 through the standard bridge.
     */
    receive() external payable onlyEOA {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, 200_000, bytes(""));
    }

    /**
     * @notice Send ETH to the message sender on the remote domain
     */
    function bridgeETH(uint32 _minGasLimit, bytes calldata _data) public payable onlyEOA {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, _minGasLimit, _data);
    }

    /**
     * @notice Send ETH to a specified account on the remote domain
=======
    /**
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
     * @param _messenger   Address of CrossDomainMessenger on this network.
     * @param _otherBridge Address of the other StandardBridge contract.
     */
    constructor(address payable _messenger, address payable _otherBridge) {
        MESSENGER = CrossDomainMessenger(_messenger);
        OTHER_BRIDGE = StandardBridge(_otherBridge);
    }

    /**
     * @notice Allows EOAs to deposit ETH by sending directly to the bridge.
     */
    receive() external payable onlyEOA {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, RECEIVE_DEFAULT_GAS_LIMIT, bytes(""));
    }

    /**
     * @custom:legacy
     * @notice Legacy getter for messenger contract.
     *
     * @return Messenger contract on this domain.
     */
    function messenger() external view returns (CrossDomainMessenger) {
        return MESSENGER;
    }

    /**
     * @notice Sends ETH to the sender's address on the other chain.
     *
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
     */
    function bridgeETH(uint32 _minGasLimit, bytes calldata _extraData) public payable onlyEOA {
        _initiateBridgeETH(msg.sender, msg.sender, msg.value, _minGasLimit, _extraData);
    }

    /**
     * @notice Sends ETH to a receiver's address on the other chain. Note that if ETH is sent to a
     *         smart contract and the call fails, the ETH will be temporarily locked in the
     *         StandardBridge on the other chain until the call is replayed. If the call cannot be
     *         replayed with any amount of gas (call always reverts), then the ETH will be
     *         permanently locked in the StandardBridge on the other chain. ETH will also
     *         be locked if the receiver is the other bridge, because finalizeBridgeETH will revert
     *         in that case.
     *
     * @param _to          Address of the receiver.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function bridgeETHTo(
        address _to,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
    ) public payable {
        _initiateBridgeETH(msg.sender, _to, msg.value, _minGasLimit, _data);
    }

    /**
     * @notice Send an ERC20 to the message sender on the remote domain
=======
        bytes calldata _extraData
    ) public payable {
        _initiateBridgeETH(msg.sender, _to, msg.value, _minGasLimit, _extraData);
    }

    /**
     * @notice Sends ERC20 tokens to the sender's address on the other chain. Note that if the
     *         ERC20 token on the other chain does not recognize the local token as the correct
     *         pair token, the ERC20 bridge will fail and the tokens will be returned to sender on
     *         this chain.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _amount      Amount of local tokens to deposit.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function bridgeERC20(
        address _localToken,
        address _remoteToken,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
=======
        bytes calldata _extraData
>>>>>>> v0.5.23
    ) public virtual onlyEOA {
        _initiateBridgeERC20(
            _localToken,
            _remoteToken,
            msg.sender,
            msg.sender,
            _amount,
            _minGasLimit,
<<<<<<< HEAD
            _data
=======
            _extraData
>>>>>>> v0.5.23
        );
    }

    /**
<<<<<<< HEAD
     * @notice Send an ERC20 to a specified account on the remote domain
=======
     * @notice Sends ERC20 tokens to a receiver's address on the other chain. Note that if the
     *         ERC20 token on the other chain does not recognize the local token as the correct
     *         pair token, the ERC20 bridge will fail and the tokens will be returned to sender on
     *         this chain.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of local tokens to deposit.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function bridgeERC20To(
        address _localToken,
        address _remoteToken,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
=======
        bytes calldata _extraData
>>>>>>> v0.5.23
    ) public virtual {
        _initiateBridgeERC20(
            _localToken,
            _remoteToken,
            msg.sender,
            _to,
            _amount,
            _minGasLimit,
<<<<<<< HEAD
            _data
=======
            _extraData
>>>>>>> v0.5.23
        );
    }

    /**
<<<<<<< HEAD
     * @notice Finalize an ETH sending transaction sent from a remote domain
=======
     * @notice Finalizes an ETH bridge on this chain. Can only be triggered by the other
     *         StandardBridge contract on the remote chain.
     *
     * @param _from      Address of the sender.
     * @param _to        Address of the receiver.
     * @param _amount    Amount of ETH being bridged.
     * @param _extraData Extra data to be sent with the transaction. Note that the recipient will
     *                   not be triggered with this data, but it will be emitted and can be used
     *                   to identify the transaction.
>>>>>>> v0.5.23
     */
    function finalizeBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
<<<<<<< HEAD
        bytes calldata _data
=======
        bytes calldata _extraData
>>>>>>> v0.5.23
    ) public payable onlyOtherBridge {
        require(msg.value == _amount, "StandardBridge: amount sent does not match amount required");
        require(_to != address(this), "StandardBridge: cannot send to self");
        require(_to != address(MESSENGER), "StandardBridge: cannot send to messenger");

<<<<<<< HEAD
        emit ETHBridgeFinalized(_from, _to, _amount, _data);
        (bool success, ) = _to.call{ value: _amount }(new bytes(0));
        require(success, "TransferHelper::safeTransferETH: ETH transfer failed");
    }

    /**
     * @notice Finalize an ERC20 sending transaction sent from a remote domain
=======
        emit ETHBridgeFinalized(_from, _to, _amount, _extraData);

        bool success = SafeCall.call(_to, gasleft(), _amount, hex"");
        require(success, "StandardBridge: ETH transfer failed");
    }

    /**
     * @notice Finalizes an ERC20 bridge on this chain. Can only be triggered by the other
     *         StandardBridge contract on the remote chain.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _from        Address of the sender.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of the ERC20 being bridged.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function finalizeBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
<<<<<<< HEAD
        bytes calldata _data
    ) public onlyOtherBridge {
        try this.completeOutboundTransfer(_localToken, _remoteToken, _to, _amount) {
            emit ERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _data);
        } catch {
            // Something went wrong during the bridging process, return to sender.
            // Can happen if a bridge UI specifies the wrong L2 token.
            _initiateBridgeERC20Unchecked(
                _localToken,
                _remoteToken,
                _from,
                _to,
                _amount,
                0, // _minGasLimit, 0 is fine here
                _data
            );
            emit ERC20BridgeFailed(_localToken, _remoteToken, _from, _to, _amount, _data);
        }
    }

=======
        bytes calldata _extraData
    ) public onlyOtherBridge {
<<<<<<< HEAD
        try this.completeOutboundTransfer(_localToken, _remoteToken, _to, _amount) {
            emit ERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
        } catch {
            // Something went wrong during the bridging process, return to sender.
            // Can happen if a bridge UI specifies the wrong L2 token.
            // We reverse the to and from addresses to make sure the tokens are returned to the
            // sender on the other chain and preserve the accuracy of accounting based on emitted
            // events.
            _initiateBridgeERC20Unchecked(
                _localToken,
                _remoteToken,
                _to,
                _from,
                _amount,
                0, // _minGasLimit, 0 is fine here
                _extraData
            );
            emit ERC20BridgeFailed(_localToken, _remoteToken, _from, _to, _amount, _extraData);
        }
    }

    /**
     * @notice Completes an outbound token transfer. Public function, but can only be called by
     *         this contract. It's security critical that there be absolutely no way for anyone to
     *         trigger this function, except by explicit trigger within this contract. Used as a
     *         simple way to be able to try/catch any type of revert that could occur during an
     *         ERC20 mint/transfer.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of ETH being bridged.
     */
>>>>>>> v0.5.23
    function completeOutboundTransfer(
        address _localToken,
        address _remoteToken,
        address _to,
        uint256 _amount
    ) public onlySelf {
        // Make sure external function calls can't be used to trigger calls to
        // completeOutboundTransfer. We only make external (write) calls to _localToken.
        require(_localToken != address(this), "StandardBridge: local token cannot be self");

=======
>>>>>>> @eth-optimism/l2geth@0.5.27
        if (_isOptimismMintableERC20(_localToken)) {
            require(
                _isCorrectTokenPair(_localToken, _remoteToken),
                "StandardBridge: wrong remote token for Optimism Mintable ERC20 local token"
            );

            OptimismMintableERC20(_localToken).mint(_to, _amount);
        } else {
            deposits[_localToken][_remoteToken] = deposits[_localToken][_remoteToken] - _amount;
            IERC20(_localToken).safeTransfer(_to, _amount);
        }

        emit ERC20BridgeFinalized(_localToken, _remoteToken, _from, _to, _amount, _extraData);
    }

<<<<<<< HEAD
    /**********************
     * Internal Functions *
     **********************/

    /**
     * @notice Initialize the StandardBridge contract with the address of
     * the messenger on the same domain as well as the address of the bridge
     * on the remote domain
     */
    function _initialize(address payable _messenger, address payable _otherBridge) internal {
        require(address(messenger) == address(0), "Contract has already been initialized.");

=======
    /**
<<<<<<< HEAD
     * @notice Initializer.
     *
     * @param _messenger   Address of CrossDomainMessenger on this network.
     * @param _otherBridge Address of the other StandardBridge contract.
     */
    function __StandardBridge_init(address payable _messenger, address payable _otherBridge)
        internal
        onlyInitializing
    {
>>>>>>> v0.5.23
        messenger = CrossDomainMessenger(_messenger);
        otherBridge = StandardBridge(_otherBridge);
    }

    /**
<<<<<<< HEAD
     * @notice Bridge ETH to the remote chain through the messenger
=======
=======
>>>>>>> v0.5.24
     * @notice Initiates a bridge of ETH through the CrossDomainMessenger.
     *
     * @param _from        Address of the sender.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of ETH being bridged.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function _initiateBridgeETH(
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes memory _data
    ) internal {
        emit ETHBridgeInitiated(_from, _to, _amount, _data);

        messenger.sendMessage{ value: _amount }(
            address(otherBridge),
            abi.encodeWithSelector(this.finalizeBridgeETH.selector, _from, _to, _amount, _data),
=======
        bytes memory _extraData
    ) internal {
        require(
            msg.value == _amount,
            "StandardBridge: bridging ETH must include sufficient ETH value"
        );

        emit ETHBridgeInitiated(_from, _to, _amount, _extraData);

        MESSENGER.sendMessage{ value: _amount }(
            address(OTHER_BRIDGE),
            abi.encodeWithSelector(
                this.finalizeBridgeETH.selector,
                _from,
                _to,
                _amount,
                _extraData
            ),
>>>>>>> v0.5.23
            _minGasLimit
        );
    }

    /**
<<<<<<< HEAD
     * @notice Bridge an ERC20 to the remote chain through the messengers
=======
     * @notice Sends ERC20 tokens to a receiver's address on the other chain.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of local tokens to deposit.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function _initiateBridgeERC20(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
=======
        bytes calldata _extraData
>>>>>>> v0.5.23
    ) internal {
        if (_isOptimismMintableERC20(_localToken)) {
            require(
                _isCorrectTokenPair(_localToken, _remoteToken),
                "StandardBridge: wrong remote token for Optimism Mintable ERC20 local token"
            );

<<<<<<< HEAD
            OptimismMintableERC20(_localToken).burn(msg.sender, _amount);
=======
            OptimismMintableERC20(_localToken).burn(_from, _amount);
>>>>>>> v0.5.23
        } else {
            IERC20(_localToken).safeTransferFrom(_from, address(this), _amount);
            deposits[_localToken][_remoteToken] = deposits[_localToken][_remoteToken] + _amount;
        }

<<<<<<< HEAD
        _initiateBridgeERC20Unchecked(
            _localToken,
            _remoteToken,
            _from,
            _to,
            _amount,
            _minGasLimit,
<<<<<<< HEAD
            _data
=======
            _extraData
>>>>>>> v0.5.23
        );
    }

    /**
<<<<<<< HEAD
     * @notice Bridge an ERC20 to the remote chain through the messengers
=======
     * @notice Sends ERC20 tokens to a receiver's address on the other chain WITHOUT doing any
     *         validation. Be EXTREMELY careful when using this function.
     *
     * @param _localToken  Address of the ERC20 on this chain.
     * @param _remoteToken Address of the corresponding token on the remote chain.
     * @param _to          Address of the receiver.
     * @param _amount      Amount of local tokens to deposit.
     * @param _minGasLimit Minimum amount of gas that the bridge can be relayed with.
     * @param _extraData   Extra data to be sent with the transaction. Note that the recipient will
     *                     not be triggered with this data, but it will be emitted and can be used
     *                     to identify the transaction.
>>>>>>> v0.5.23
     */
    function _initiateBridgeERC20Unchecked(
        address _localToken,
        address _remoteToken,
        address _from,
        address _to,
        uint256 _amount,
        uint32 _minGasLimit,
<<<<<<< HEAD
        bytes calldata _data
=======
        bytes calldata _extraData
>>>>>>> v0.5.23
    ) internal {
=======
        emit ERC20BridgeInitiated(_localToken, _remoteToken, _from, _to, _amount, _extraData);

<<<<<<< HEAD
>>>>>>> @eth-optimism/l2geth@0.5.27
        messenger.sendMessage(
            address(otherBridge),
=======
        MESSENGER.sendMessage(
            address(OTHER_BRIDGE),
>>>>>>> @eth-optimism/l2geth@0.5.29
            abi.encodeWithSelector(
                this.finalizeBridgeERC20.selector,
                // Because this call will be executed on the remote chain, we reverse the order of
                // the remote and local token addresses relative to their order in the
                // finalizeBridgeERC20 function.
                _remoteToken,
                _localToken,
                _from,
                _to,
                _amount,
<<<<<<< HEAD
                _data
=======
                _extraData
>>>>>>> v0.5.23
            ),
            _minGasLimit
        );
<<<<<<< HEAD

<<<<<<< HEAD
        emit ERC20BridgeInitiated(_localToken, _remoteToken, _from, _to, _amount, _data);
    }

    /**
     * Checks if a given address is an OptimismMintableERC20. Not perfect, but good enough.
     * Just the way we like it.
     *
     * @param _token Address of the token to check.
     * @return True if the token is an OptimismMintableERC20.
     */
    function _isOptimismMintableERC20(address _token) internal view returns (bool) {
        // 0x1d1d8b63 is mint ^ burn ^ l1Token
        return ERC165Checker.supportsInterface(_token, 0x1d1d8b63);
    }

    /**
     * Checks if the "other token" is the correct pair token for the OptimismMintableERC20.
     *
     * @param _mintableToken OptimismMintableERC20 to check against.
     * @param _otherToken Pair token to check.
=======
        emit ERC20BridgeInitiated(_localToken, _remoteToken, _from, _to, _amount, _extraData);
=======
>>>>>>> @eth-optimism/l2geth@0.5.27
    }

    /**
     * @notice Checks if a given address is an OptimismMintableERC20. Not perfect, but good enough.
     *         Just the way we like it.
     *
     * @param _token Address of the token to check.
     *
     * @return True if the token is an OptimismMintableERC20.
     */
    function _isOptimismMintableERC20(address _token) internal view returns (bool) {
        return
            ERC165Checker.supportsInterface(_token, type(ILegacyMintableERC20).interfaceId) ||
            ERC165Checker.supportsInterface(_token, type(IOptimismMintableERC20).interfaceId);
    }

    /**
     * @notice Checks if the "other token" is the correct pair token for the OptimismMintableERC20.
     *
     * @param _mintableToken OptimismMintableERC20 to check against.
     * @param _otherToken    Pair token to check.
     *
>>>>>>> v0.5.23
     * @return True if the other token is the correct pair token for the OptimismMintableERC20.
     */
    function _isCorrectTokenPair(address _mintableToken, address _otherToken)
        internal
        view
        returns (bool)
    {
        return _otherToken == OptimismMintableERC20(_mintableToken).l1Token();
    }
}
