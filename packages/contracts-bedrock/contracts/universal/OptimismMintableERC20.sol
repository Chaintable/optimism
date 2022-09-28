// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
<<<<<<< HEAD

/**
 * @title OptimismMintableERC20
 * This contract represents the remote representation
 * of an ERC20 token. It is linked to the address of
 * a token in another domain and tokens can be locked
 * in the StandardBridge which will mint tokens in the
 * other domain.
 */
contract OptimismMintableERC20 is ERC20 {
    event Mint(address indexed _account, uint256 _amount);
    event Burn(address indexed _account, uint256 _amount);

    /**
     * @notice The address of the token in the remote domain
=======
import "./SupportedInterfaces.sol";

/**
 * @title OptimismMintableERC20
 * @notice OptimismMintableERC20 is a standard extension of the base ERC20 token contract designed
 *         to allow the StandardBridge contracts to mint and burn tokens. This makes it possible to
 *         use an OptimismMintablERC20 as the L2 representation of an L1 token, or vice-versa.
 *         Designed to be backwards compatible with the older StandardL2ERC20 token which was only
 *         meant for use on L2.
 */
contract OptimismMintableERC20 is ERC20 {
    /**
     * @notice Address of the corresponding version of this token on the remote chain.
     */
    address public remoteToken;

    /**
     * @notice Address of the StandardBridge on this network.
     */
    address public bridge;

    /**
     * @notice Emitted whenever tokens are minted for an account.
     *
     * @param account Address of the account tokens are being minted for.
     * @param amount  Amount of tokens minted.
     */
    event Mint(address indexed account, uint256 amount);

    /**
     * @notice Emitted whenever tokens are burned from an account.
     *
     * @param account Address of the account tokens are being burned from.
     * @param amount  Amount of tokens burned.
     */
    event Burn(address indexed account, uint256 amount);

    /**
<<<<<<< HEAD
     * @notice Address of the corresponding version of this token on the remote chain.
>>>>>>> v0.5.23
     */
    address public remoteToken;

    /**
<<<<<<< HEAD
     * @notice The address of the bridge responsible for
     * minting. It is in the same domain.
=======
     * @notice Address of the StandardBridge on this network.
>>>>>>> v0.5.23
=======
     * @notice A modifier that only allows the bridge to call
>>>>>>> v0.5.24
     */
    modifier onlyBridge() {
        require(msg.sender == bridge, "OptimismMintableERC20: only bridge can mint and burn");
        _;
    }

    /**
<<<<<<< HEAD
     * @param _bridge Address of the L2 standard bridge.
     * @param _remoteToken Address of the corresponding L1 token.
     * @param _name ERC20 name.
     * @param _symbol ERC20 symbol.
=======
     * @param _bridge      Address of the L2 standard bridge.
     * @param _remoteToken Address of the corresponding L1 token.
     * @param _name        ERC20 name.
     * @param _symbol      ERC20 symbol.
>>>>>>> v0.5.23
     */
    constructor(
        address _bridge,
        address _remoteToken,
        string memory _name,
        string memory _symbol
    ) ERC20(_name, _symbol) {
        remoteToken = _remoteToken;
        bridge = _bridge;
    }

    /**
<<<<<<< HEAD
<<<<<<< HEAD
     * @notice Returns the corresponding L1 token address.
     * This is a legacy function and wraps the remoteToken value.
=======
     * @custom:legacy
     * @notice Legacy getter for the remote token. Use remoteToken going forward.
>>>>>>> v0.5.23
     */
    function l1Token() public view returns (address) {
        return remoteToken;
    }

    /**
<<<<<<< HEAD
     * @notice The address of the bridge contract
     * responsible for minting tokens. This is a legacy
     * getter function
=======
     * @custom:legacy
     * @notice Legacy getter for the bridge. Use bridge going forward.
>>>>>>> v0.5.23
=======
     * @notice Allows the StandardBridge on this network to mint tokens.
     *
     * @param _to     Address to mint tokens to.
     * @param _amount Amount of tokens to mint.
>>>>>>> v0.5.24
     */
    function mint(address _to, uint256 _amount) external virtual onlyBridge {
        _mint(_to, _amount);
        emit Mint(_to, _amount);
    }

    /**
     * @notice Allows the StandardBridge on this network to burn tokens.
     *
     * @param _from   Address to burn tokens from.
     * @param _amount Amount of tokens to burn.
     */
<<<<<<< HEAD
    modifier onlyBridge() {
<<<<<<< HEAD
        require(msg.sender == bridge, "Only L2 Bridge can mint and burn");
=======
        require(msg.sender == bridge, "OptimismMintableERC20: only bridge can mint and burn");
>>>>>>> v0.5.23
        _;
=======
    function burn(address _from, uint256 _amount) external virtual onlyBridge {
        _burn(_from, _amount);
        emit Burn(_from, _amount);
>>>>>>> v0.5.24
    }

    /**
<<<<<<< HEAD
     * @notice ERC165
     */
    // slither-disable-next-line external-function
    function supportsInterface(bytes4 _interfaceId) public pure returns (bool) {
        bytes4 iface1 = bytes4(keccak256("supportsInterface(bytes4)")); // ERC165
        bytes4 iface2 = this.l1Token.selector ^ this.mint.selector ^ this.burn.selector;
        bytes4 iface3 = this.remoteToken.selector ^ this.mint.selector ^ this.burn.selector;
        return _interfaceId == iface1 || _interfaceId == iface3 || _interfaceId == iface2;
    }

    /**
     * @notice The bridge can mint tokens
     */
    // slither-disable-next-line external-function
    function mint(address _to, uint256 _amount) public virtual onlyBridge {
        _mint(_to, _amount);

=======
     * @notice ERC165 interface check function.
     *
     * @param _interfaceId Interface ID to check.
     *
     * @return Whether or not the interface is supported by this contract.
     */
    function supportsInterface(bytes4 _interfaceId) external pure returns (bool) {
        bytes4 iface1 = type(IERC165).interfaceId;
        bytes4 iface2 = type(IL1Token).interfaceId;
        bytes4 iface3 = type(IRemoteToken).interfaceId;
        return _interfaceId == iface1 || _interfaceId == iface2 || _interfaceId == iface3;
    }

    /**
     * @custom:legacy
     * @notice Legacy getter for the remote token. Use remoteToken going forward.
     */
<<<<<<< HEAD
    function mint(address _to, uint256 _amount) external virtual onlyBridge {
        _mint(_to, _amount);
>>>>>>> v0.5.23
        emit Mint(_to, _amount);
    }

    /**
<<<<<<< HEAD
     * @notice The bridge can burn tokens
     */
    // slither-disable-next-line external-function
    function burn(address _from, uint256 _amount) public virtual onlyBridge {
        _burn(_from, _amount);

=======
     * @notice Allows the StandardBridge on this network to burn tokens.
     *
     * @param _from   Address to burn tokens from.
     * @param _amount Amount of tokens to burn.
     */
    function burn(address _from, uint256 _amount) external virtual onlyBridge {
        _burn(_from, _amount);
>>>>>>> v0.5.23
        emit Burn(_from, _amount);
=======
    function l1Token() public view returns (address) {
        return remoteToken;
    }

    /**
     * @custom:legacy
     * @notice Legacy getter for the bridge. Use bridge going forward.
     */
    function l2Bridge() public view returns (address) {
        return bridge;
>>>>>>> v0.5.24
    }
}
