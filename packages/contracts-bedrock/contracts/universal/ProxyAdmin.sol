// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Owned } from "@rari-capital/solmate/src/auth/Owned.sol";
import { Proxy } from "./Proxy.sol";
import { AddressManager } from "../legacy/AddressManager.sol";
import { L1ChugSplashProxy } from "../legacy/L1ChugSplashProxy.sol";

/**
 * @title IStaticERC1967Proxy
 * @notice IStaticERC1967Proxy is a static version of the ERC1967 proxy interface.
 */
interface IStaticERC1967Proxy {
    function implementation() external view returns (address);

    function admin() external view returns (address);
}

/**
 * @title IStaticL1ChugSplashProxy
 * @notice IStaticL1ChugSplashProxy is a static version of the ChugSplash proxy interface.
 */
interface IStaticL1ChugSplashProxy {
    function getImplementation() external view returns (address);

    function getOwner() external view returns (address);
}

/**
 * @title ProxyAdmin
 * @notice This is an auxiliary contract meant to be assigned as the admin of an ERC1967 Proxy,
 *         based on the OpenZeppelin implementation. It has backwards compatibility logic to work
 *         with the various types of proxies that have been deployed by Optimism in the past.
 */
contract ProxyAdmin is Owned {
    /**
     * @notice The proxy types that the ProxyAdmin can manage.
     *
     * @custom:value ERC1967    Represents an ERC1967 compliant transparent proxy interface.
     * @custom:value CHUGSPLASH Represents the Chugsplash proxy interface (legacy).
     * @custom:value RESOLVED   Represents the ResolvedDelegate proxy (legacy).
     */
    enum ProxyType {
        ERC1967,
        CHUGSPLASH,
        RESOLVED
    }

    /**
     * @custom:legacy
     * @notice A mapping of proxy types, used for backwards compatibility.
     */
    mapping(address => ProxyType) public proxyType;

    /**
     * @custom:legacy
     * @notice A reverse mapping of addresses to names held in the AddressManager. This must be
     *         manually kept up to date with changes in the AddressManager for this contract
     *         to be able to work as an admin for the ResolvedDelegateProxy type.
     */
    mapping(address => string) public implementationName;

    /**
     * @custom:legacy
     * @notice The address of the address manager, this is required to manage the
     *         ResolvedDelegateProxy type.
     */
    AddressManager public addressManager;

    /**
     * @custom:legacy
     * @notice A legacy upgrading indicator used by the old Chugsplash Proxy.
     */
    bool internal upgrading = false;

    /**
     * @param _owner Address of the initial owner of this contract.
     */
    constructor(address _owner) Owned(_owner) {}

    /**
     * @notice Sets the proxy type for a given address. Only required for non-standard (legacy)
     *         proxy types.
     *
     * @param _address Address of the proxy.
     * @param _type    Type of the proxy.
     */
    function setProxyType(address _address, ProxyType _type) external onlyOwner {
        proxyType[_address] = _type;
    }

    /**
     * @notice Sets the implementation name for a given address. Only required for
     *         ResolvedDelegateProxy type proxies that have an implementation name.
     *
     * @param _address Address of the ResolvedDelegateProxy.
     * @param _name    Name of the implementation for the proxy.
     */
    function setImplementationName(address _address, string memory _name) external onlyOwner {
        implementationName[_address] = _name;
    }

    /**
     * @notice Set the address of the AddressManager. This is required to manage legacy
     *         ResolvedDelegateProxy type proxy contracts.
     *
     * @param _address Address of the AddressManager.
     */
    function setAddressManager(AddressManager _address) external onlyOwner {
        addressManager = _address;
    }

    /**
     * @custom:legacy
     * @notice Set an address in the address manager. Since only the owner of the AddressManager
     *         can directly modify addresses and the ProxyAdmin will own the AddressManager, this
     *         gives the owner of the ProxyAdmin the ability to modify addresses directly.
     *
     * @param _name    Name to set within the AddressManager.
     * @param _address Address to attach to the given name.
     */
    function setAddress(string memory _name, address _address) external onlyOwner {
        addressManager.setAddress(_name, _address);
    }

    /**
     * @custom:legacy
<<<<<<< HEAD
<<<<<<< HEAD
     * @notice Legacy function used by the old Chugsplash proxy to determine if an upgrade is
     *         happening.
     *
     * @return Whether or not there is an upgrade going on
=======
     * @notice Legacy function used to tell ChugSplashProxy contracts if an upgrade is happening.
     *
     * @return Whether or not there is an upgrade going on. May not actually tell you whether an
     *         upgrade is going on, since we don't currently plan to use this variable for anything
     *         other than a legacy indicator to fix a UX bug in the ChugSplash proxy.
>>>>>>> v0.5.23
=======
     * @notice Set the upgrading status for the Chugsplash proxy type.
     *
     * @param _upgrading Whether or not the system is upgrading.
>>>>>>> v0.5.24
     */
    function setUpgrading(bool _upgrading) external onlyOwner {
        upgrading = _upgrading;
    }

    /**
     * @notice Updates the admin of the given proxy address.
     *
     * @param _proxy    Address of the proxy to update.
     * @param _newAdmin Address of the new proxy admin.
     */
    function changeProxyAdmin(address payable _proxy, address _newAdmin) external onlyOwner {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            Proxy(_proxy).changeAdmin(_newAdmin);
        } else if (ptype == ProxyType.CHUGSPLASH) {
            L1ChugSplashProxy(_proxy).setOwner(_newAdmin);
        } else if (ptype == ProxyType.RESOLVED) {
            addressManager.transferOwnership(_newAdmin);
        } else {
            revert("ProxyAdmin: unknown proxy type");
        }
    }

    /**
     * @notice Changes a proxy's implementation contract and delegatecalls the new implementation
     *         with some given data. Useful for atomic upgrade-and-initialize calls.
     *
     * @param _proxy          Address of the proxy to upgrade.
     * @param _implementation Address of the new implementation address.
     * @param _data           Data to trigger the new implementation with.
     */
    function upgradeAndCall(
        address payable _proxy,
        address _implementation,
        bytes memory _data
    ) external payable onlyOwner {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            Proxy(_proxy).upgradeToAndCall{ value: msg.value }(_implementation, _data);
        } else {
            // reverts if proxy type is unknown
            upgrade(_proxy, _implementation);
            (bool success, ) = _proxy.call{ value: msg.value }(_data);
            require(success, "ProxyAdmin: call to proxy after upgrade failed");
        }
    }

    /**
     * @custom:legacy
     * @notice Legacy function used to tell ChugSplashProxy contracts if an upgrade is happening.
     *
     * @return Whether or not there is an upgrade going on. May not actually tell you whether an
     *         upgrade is going on, since we don't currently plan to use this variable for anything
     *         other than a legacy indicator to fix a UX bug in the ChugSplash proxy.
     */
    function isUpgrading() external view returns (bool) {
        return upgrading;
    }

    /**
<<<<<<< HEAD
     * @dev Returns the current implementation of `proxy`.
     *      This contract must be the admin of `proxy`.
     *
     * @param proxy The Proxy to return the implementation of.
     * @return The address of the implementation.
     */
    function getProxyImplementation(Proxy proxy) external view returns (address) {
        ProxyType proxyType = proxyType[address(proxy)];

        // We need to manually run the static call since the getter cannot be flagged as view
        address target;
        bytes memory data;
        if (proxyType == ProxyType.OpenZeppelin) {
            target = address(proxy);
            data = abi.encodeWithSelector(Proxy.implementation.selector);
        } else if (proxyType == ProxyType.Chugsplash) {
            target = address(proxy);
            data = abi.encodeWithSelector(L1ChugSplashProxy.getImplementation.selector);
        } else if (proxyType == ProxyType.ResolvedDelegate) {
            target = address(addressManager);
            data = abi.encodeWithSelector(
                Lib_AddressManager.getAddress.selector,
                implementationName[address(proxy)]
            );
        } else {
            revert("ProxyAdmin: unknown proxy type");
        }

        (bool success, bytes memory returndata) = target.staticcall(data);
        require(success);
        return abi.decode(returndata, (address));
    }

    /**
     * @dev Returns the current admin of `proxy`.
     *      This contract must be the admin of `proxy`.
     *
     * @param proxy The Proxy to return the admin of.
     * @return The address of the admin.
     */
    function getProxyAdmin(Proxy proxy) external view returns (address) {
        ProxyType proxyType = proxyType[address(proxy)];

        // We need to manually run the static call since the getter cannot be flagged as view
        address target;
        bytes memory data;
        if (proxyType == ProxyType.OpenZeppelin) {
            target = address(proxy);
            data = abi.encodeWithSelector(Proxy.admin.selector);
        } else if (proxyType == ProxyType.Chugsplash) {
            target = address(proxy);
            data = abi.encodeWithSelector(L1ChugSplashProxy.getOwner.selector);
        } else if (proxyType == ProxyType.ResolvedDelegate) {
            target = address(addressManager);
            data = abi.encodeWithSignature("owner()");
        } else {
            revert("ProxyAdmin: unknown proxy type");
        }

        (bool success, bytes memory returndata) = target.staticcall(data);
        require(success);
        return abi.decode(returndata, (address));
    }

    /**
     * @dev Changes the admin of `proxy` to `newAdmin`. This contract must be the current admin
     *      of `proxy`.
     *
     * @param proxy    The proxy that will have its admin updated.
     * @param newAdmin The address of the admin to update to.
     */
    function changeProxyAdmin(Proxy proxy, address newAdmin) external onlyOwner {
        ProxyType proxyType = proxyType[address(proxy)];

        if (proxyType == ProxyType.OpenZeppelin) {
            proxy.changeAdmin(newAdmin);
        } else if (proxyType == ProxyType.Chugsplash) {
            L1ChugSplashProxy(payable(proxy)).setOwner(newAdmin);
        } else if (proxyType == ProxyType.ResolvedDelegate) {
            Lib_AddressManager(addressManager).transferOwnership(newAdmin);
=======
     * @notice Returns the implementation of the given proxy address.
     *
     * @param _proxy Address of the proxy to get the implementation of.
     *
     * @return Address of the implementation of the proxy.
     */
    function getProxyImplementation(address _proxy) external view returns (address) {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            return IStaticERC1967Proxy(_proxy).implementation();
        } else if (ptype == ProxyType.CHUGSPLASH) {
            return IStaticL1ChugSplashProxy(_proxy).getImplementation();
        } else if (ptype == ProxyType.RESOLVED) {
            return addressManager.getAddress(implementationName[_proxy]);
        } else {
            revert("ProxyAdmin: unknown proxy type");
        }
    }

    /**
     * @notice Returns the admin of the given proxy address.
     *
     * @param _proxy Address of the proxy to get the admin of.
     *
     * @return Address of the admin of the proxy.
     */
    function getProxyAdmin(address payable _proxy) external view returns (address) {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            return IStaticERC1967Proxy(_proxy).admin();
        } else if (ptype == ProxyType.CHUGSPLASH) {
            return IStaticL1ChugSplashProxy(_proxy).getOwner();
        } else if (ptype == ProxyType.RESOLVED) {
            return addressManager.owner();
        } else {
            revert("ProxyAdmin: unknown proxy type");
        }
    }

    /**
<<<<<<< HEAD
     * @notice Updates the admin of the given proxy address.
     *
     * @param _proxy    Address of the proxy to update.
     * @param _newAdmin Address of the new proxy admin.
     */
    function changeProxyAdmin(address payable _proxy, address _newAdmin) external onlyOwner {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            Proxy(_proxy).changeAdmin(_newAdmin);
        } else if (ptype == ProxyType.CHUGSPLASH) {
            L1ChugSplashProxy(_proxy).setOwner(_newAdmin);
        } else if (ptype == ProxyType.RESOLVED) {
            addressManager.transferOwnership(_newAdmin);
        } else {
            revert("ProxyAdmin: unknown proxy type");
>>>>>>> v0.5.23
        }
    }

    /**
<<<<<<< HEAD
     * @dev Upgrades `proxy` to `implementation`. This contract must be the admin of `proxy`.
     *
     * @param proxy          The address of the proxy.
     * @param implementation The address of the implementation.
     */
    function upgrade(Proxy proxy, address implementation) public onlyOwner {
        ProxyType proxyType = proxyType[address(proxy)];

        if (proxyType == ProxyType.OpenZeppelin) {
            proxy.upgradeTo(implementation);
        } else if (proxyType == ProxyType.Chugsplash) {
            L1ChugSplashProxy(payable(proxy)).setStorage(
                0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc,
                bytes32(uint256(uint160(implementation)))
            );
        } else if (proxyType == ProxyType.ResolvedDelegate) {
            string memory name = implementationName[address(proxy)];
            Lib_AddressManager(addressManager).setAddress(name, implementation);
=======
=======
>>>>>>> v0.5.24
     * @notice Changes a proxy's implementation contract.
     *
     * @param _proxy          Address of the proxy to upgrade.
     * @param _implementation Address of the new implementation address.
     */
    function upgrade(address payable _proxy, address _implementation) public onlyOwner {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            Proxy(_proxy).upgradeTo(_implementation);
        } else if (ptype == ProxyType.CHUGSPLASH) {
            L1ChugSplashProxy(_proxy).setStorage(
                // bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)
                0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc,
                bytes32(uint256(uint160(_implementation)))
            );
        } else if (ptype == ProxyType.RESOLVED) {
            string memory name = implementationName[_proxy];
            addressManager.setAddress(name, _implementation);
        } else {
<<<<<<< HEAD
            revert("ProxyAdmin: unknown proxy type");
>>>>>>> v0.5.23
=======
            // It should not be possible to retrieve a ProxyType value which is not matched by
            // one of the previous conditions.
            assert(false);
>>>>>>> @eth-optimism/l2geth@0.5.27
        }
    }
<<<<<<< HEAD

    /**
<<<<<<< HEAD
     * @dev Upgrades `proxy` to `implementation` and calls a function on the new implementation.
     *      This contract must be the admin of `proxy`.
     *
     * @param proxy           The proxy to call.
     * @param implementation  The implementation to upgrade the proxy to.
     * @param data            The calldata to pass to the implementation.
     */
    function upgradeAndCall(
        Proxy proxy,
        address implementation,
        bytes memory data
    ) external payable onlyOwner {
        ProxyType proxyType = proxyType[address(proxy)];

        if (proxyType == ProxyType.OpenZeppelin) {
            proxy.upgradeToAndCall{ value: msg.value }(implementation, data);
        } else {
            upgrade(proxy, implementation);
            (bool success, ) = address(proxy).call{ value: msg.value }(data);
            require(success);
=======
     * @notice Changes a proxy's implementation contract and delegatecalls the new implementation
     *         with some given data. Useful for atomic upgrade-and-initialize calls.
     *
     * @param _proxy          Address of the proxy to upgrade.
     * @param _implementation Address of the new implementation address.
     * @param _data           Data to trigger the new implementation with.
     */
    function upgradeAndCall(
        address payable _proxy,
        address _implementation,
        bytes memory _data
    ) external payable onlyOwner {
        ProxyType ptype = proxyType[_proxy];
        if (ptype == ProxyType.ERC1967) {
            Proxy(_proxy).upgradeToAndCall{ value: msg.value }(_implementation, _data);
        } else {
            // reverts if proxy type is unknown
            upgrade(_proxy, _implementation);
            (bool success, ) = _proxy.call{ value: msg.value }(_data);
            require(success, "ProxyAdmin: call to proxy after upgrade failed");
>>>>>>> v0.5.23
        }
    }
=======
>>>>>>> v0.5.24
}
