// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { Proxy } from "../universal/Proxy.sol";
import { ProxyAdmin } from "../universal/ProxyAdmin.sol";
import { SimpleStorage } from "./Proxy.t.sol";
import { L1ChugSplashProxy } from "../legacy/L1ChugSplashProxy.sol";
<<<<<<< HEAD
import { Lib_ResolvedDelegateProxy } from "../legacy/Lib_ResolvedDelegateProxy.sol";
import { Lib_AddressManager } from "../legacy/Lib_AddressManager.sol";
=======
import { ResolvedDelegateProxy } from "../legacy/ResolvedDelegateProxy.sol";
import { AddressManager } from "../legacy/AddressManager.sol";
>>>>>>> v0.5.23

contract ProxyAdmin_Test is Test {
    address alice = address(64);

    Proxy proxy;
    L1ChugSplashProxy chugsplash;
<<<<<<< HEAD
    Lib_ResolvedDelegateProxy resolved;

    Lib_AddressManager addressManager;
=======
    ResolvedDelegateProxy resolved;

    AddressManager addressManager;
>>>>>>> v0.5.23

    ProxyAdmin admin;

    SimpleStorage implementation;

    function setUp() external {
        // Deploy the proxy admin
        admin = new ProxyAdmin(alice);
        // Deploy the standard proxy
        proxy = new Proxy(address(admin));

        // Deploy the legacy L1ChugSplashProxy with the admin as the owner
        chugsplash = new L1ChugSplashProxy(address(admin));

<<<<<<< HEAD
        // Deploy the legacy Lib_AddressManager
        addressManager = new Lib_AddressManager();
        // The proxy admin must be the new owner of the address manager
        addressManager.transferOwnership(address(admin));
        // Deploy a legacy Lib_ResolvedDelegateProxy with the name `a`.
        // Whatever `a` is set to in Lib_AddressManager will be the address
        // that is used for the implementation.
        resolved = new Lib_ResolvedDelegateProxy(address(addressManager), "a");

        // Set the address of the address manager in the admin so that it
        // can resolve the implementation address of legacy
        // Lib_ResolvedDelegateProxy based proxies.
        vm.prank(alice);
        admin.setAddressManager(address(addressManager));
        // Set the reverse lookup of the Lib_ResolvedDelegateProxy
        // proxy
        vm.prank(alice);
        admin.setImplementationName(address(resolved), "a");

        // Set the proxy types
        vm.prank(alice);
        admin.setProxyType(address(chugsplash), ProxyAdmin.ProxyType.Chugsplash);
        vm.prank(alice);
        admin.setProxyType(address(resolved), ProxyAdmin.ProxyType.ResolvedDelegate);
=======
        // Deploy the legacy AddressManager
        addressManager = new AddressManager();
        // The proxy admin must be the new owner of the address manager
        addressManager.transferOwnership(address(admin));
        // Deploy a legacy ResolvedDelegateProxy with the name `a`.
        // Whatever `a` is set to in AddressManager will be the address
        // that is used for the implementation.
        resolved = new ResolvedDelegateProxy(addressManager, "a");

        // Impersonate alice for setting up the admin.
        vm.startPrank(alice);
        // Set the address of the address manager in the admin so that it
        // can resolve the implementation address of legacy
        // ResolvedDelegateProxy based proxies.
        admin.setAddressManager(addressManager);
        // Set the reverse lookup of the ResolvedDelegateProxy
        // proxy
        admin.setImplementationName(address(resolved), "a");

        // Set the proxy types
        admin.setProxyType(address(proxy), ProxyAdmin.ProxyType.ERC1967);
        admin.setProxyType(address(chugsplash), ProxyAdmin.ProxyType.CHUGSPLASH);
        admin.setProxyType(address(resolved), ProxyAdmin.ProxyType.RESOLVED);
        vm.stopPrank();
>>>>>>> v0.5.23

        implementation = new SimpleStorage();
    }

    function test_setImplementationName() external {
        vm.prank(alice);
        admin.setImplementationName(address(1), "foo");
        assertEq(
            admin.implementationName(address(1)),
            "foo"
        );
    }

    function test_onlyOwnerSetAddressManager() external {
        vm.expectRevert("UNAUTHORIZED");
<<<<<<< HEAD
        admin.setAddressManager(address(0));
=======
        admin.setAddressManager(AddressManager((address(0))));
>>>>>>> v0.5.23
    }

    function test_onlyOwnerSetImplementationName() external {
        vm.expectRevert("UNAUTHORIZED");
        admin.setImplementationName(address(0), "foo");
    }

    function test_onlyOwnerSetProxyType() external {
        vm.expectRevert("UNAUTHORIZED");
<<<<<<< HEAD
        admin.setProxyType(address(0), ProxyAdmin.ProxyType.Chugsplash);
=======
        admin.setProxyType(address(0), ProxyAdmin.ProxyType.CHUGSPLASH);
>>>>>>> v0.5.23
    }

    function test_owner() external {
        assertEq(admin.owner(), alice);
    }

    function test_proxyType() external {
        assertEq(
            uint256(admin.proxyType(address(proxy))),
<<<<<<< HEAD
            uint256(ProxyAdmin.ProxyType.OpenZeppelin)
        );
        assertEq(
            uint256(admin.proxyType(address(chugsplash))),
            uint256(ProxyAdmin.ProxyType.Chugsplash)
        );
        assertEq(
            uint256(admin.proxyType(address(resolved))),
            uint256(ProxyAdmin.ProxyType.ResolvedDelegate)
        );
    }

    function test_openZeppelinGetProxyImplementation() external {
        getProxyImplementation(proxy);
    }

    function test_chugsplashGetProxyImplementation() external {
        getProxyImplementation(Proxy(payable(chugsplash)));
    }

    function test_delegateResolvedGetProxyImplementation() external {
        getProxyImplementation(Proxy(payable(resolved)));
    }

    function getProxyImplementation(Proxy _proxy) internal {
=======
            uint256(ProxyAdmin.ProxyType.ERC1967)
        );
        assertEq(
            uint256(admin.proxyType(address(chugsplash))),
            uint256(ProxyAdmin.ProxyType.CHUGSPLASH)
        );
        assertEq(
            uint256(admin.proxyType(address(resolved))),
            uint256(ProxyAdmin.ProxyType.RESOLVED)
        );
    }

    function test_erc1967GetProxyImplementation() external {
        getProxyImplementation(payable(proxy));
    }

    function test_chugsplashGetProxyImplementation() external {
        getProxyImplementation(payable(chugsplash));
    }

    function test_delegateResolvedGetProxyImplementation() external {
        getProxyImplementation(payable(resolved));
    }

    function getProxyImplementation(address payable _proxy) internal {
>>>>>>> v0.5.23
        {
            address impl = admin.getProxyImplementation(_proxy);
            assertEq(impl, address(0));
        }

        vm.prank(alice);
        admin.upgrade(_proxy, address(implementation));

        {
            address impl = admin.getProxyImplementation(_proxy);
            assertEq(impl, address(implementation));
        }
    }

<<<<<<< HEAD
    function test_openZeppelinGetProxyAdmin() external {
        getProxyAdmin(proxy);
    }

    function test_chugsplashGetProxyAdmin() external {
        getProxyAdmin(Proxy(payable(chugsplash)));
    }

    function test_delegateResolvedGetProxyAdmin() external {
        getProxyAdmin(Proxy(payable(resolved)));
    }

    function getProxyAdmin(Proxy _proxy) internal {
=======
    function test_erc1967GetProxyAdmin() external {
        getProxyAdmin(payable(proxy));
    }

    function test_chugsplashGetProxyAdmin() external {
        getProxyAdmin(payable(chugsplash));
    }

    function test_delegateResolvedGetProxyAdmin() external {
        getProxyAdmin(payable(resolved));
    }

    function getProxyAdmin(address payable _proxy) internal {
>>>>>>> v0.5.23
        address owner = admin.getProxyAdmin(_proxy);
        assertEq(owner, address(admin));
    }

<<<<<<< HEAD
    function test_openZeppelinChangeProxyAdmin() external {
        changeProxyAdmin(proxy);
    }

    function test_chugsplashChangeProxyAdmin() external {
        changeProxyAdmin(Proxy(payable(chugsplash)));
    }

    function test_delegateResolvedChangeProxyAdmin() external {
        changeProxyAdmin(Proxy(payable(resolved)));
    }

    function changeProxyAdmin(Proxy _proxy) internal {
=======
    function test_erc1967ChangeProxyAdmin() external {
        changeProxyAdmin(payable(proxy));
    }

    function test_chugsplashChangeProxyAdmin() external {
        changeProxyAdmin(payable(chugsplash));
    }

    function test_delegateResolvedChangeProxyAdmin() external {
        changeProxyAdmin(payable(resolved));
    }

    function changeProxyAdmin(address payable _proxy) internal {
>>>>>>> v0.5.23
        ProxyAdmin.ProxyType proxyType = admin.proxyType(address(_proxy));

        vm.prank(alice);
        admin.changeProxyAdmin(_proxy, address(128));

<<<<<<< HEAD
        // The proxy is not longer the admin and can
        // no longer call the proxy interface except for
        // the ResolvedDelegate type which anybody can call
        // the admin interface
        if (proxyType != ProxyAdmin.ProxyType.ResolvedDelegate) {
            vm.expectRevert();
            admin.getProxyAdmin(_proxy);
=======
        // The proxy is no longer the admin and can
        // no longer call the proxy interface except for
        // the ResolvedDelegate type on which anybody can
        // call the admin interface.
        if (proxyType == ProxyAdmin.ProxyType.ERC1967) {
            vm.expectRevert("Proxy: implementation not initialized");
            admin.getProxyAdmin(_proxy);
        } else if (proxyType == ProxyAdmin.ProxyType.CHUGSPLASH) {
            vm.expectRevert("L1ChugSplashProxy: implementation is not set yet");
            admin.getProxyAdmin(_proxy);
        } else if (proxyType == ProxyAdmin.ProxyType.RESOLVED) {
            // Just an empty block to show that all cases are covered
        } else {
            vm.expectRevert("ProxyAdmin: unknown proxy type");
>>>>>>> v0.5.23
        }

        // Call the proxy contract directly to get the admin.
        // Different proxy types have different interfaces.
        vm.prank(address(128));
<<<<<<< HEAD
        if (proxyType == ProxyAdmin.ProxyType.OpenZeppelin) {
            assertEq(_proxy.admin(), address(128));
        } else if (proxyType == ProxyAdmin.ProxyType.Chugsplash) {
=======
        if (proxyType == ProxyAdmin.ProxyType.ERC1967) {
            assertEq(Proxy(payable(_proxy)).admin(), address(128));
        } else if (proxyType == ProxyAdmin.ProxyType.CHUGSPLASH) {
>>>>>>> v0.5.23
            assertEq(
                L1ChugSplashProxy(payable(_proxy)).getOwner(),
                address(128)
            );
<<<<<<< HEAD
        } else if (proxyType == ProxyAdmin.ProxyType.ResolvedDelegate) {
=======
        } else if (proxyType == ProxyAdmin.ProxyType.RESOLVED) {
>>>>>>> v0.5.23
            assertEq(
                addressManager.owner(),
                address(128)
            );
        } else {
            assert(false);
        }
    }

<<<<<<< HEAD
    function test_openZeppelinUpgrade() external {
        upgrade(proxy);
    }

    function test_chugsplashUpgrade() external {
        upgrade(Proxy(payable(chugsplash)));
    }

    function test_delegateResolvedUpgrade() external {
        upgrade(Proxy(payable(resolved)));
    }

    function upgrade(Proxy _proxy) internal {
=======
    function test_erc1967Upgrade() external {
        upgrade(payable(proxy));
    }

    function test_chugsplashUpgrade() external {
        upgrade(payable(chugsplash));
    }

    function test_delegateResolvedUpgrade() external {
        upgrade(payable(resolved));
    }

    function upgrade(address payable _proxy) internal {
>>>>>>> v0.5.23
        vm.prank(alice);
        admin.upgrade(_proxy, address(implementation));

        address impl = admin.getProxyImplementation(_proxy);
        assertEq(impl, address(implementation));
    }

<<<<<<< HEAD
    function test_openZeppelinUpgradeAndCall() external {
        upgradeAndCall(proxy);
    }

    function test_chugsplashUpgradeAndCall() external {
        upgradeAndCall(Proxy(payable(chugsplash)));
    }

    function test_delegateResolvedUpgradeAndCall() external {
        upgradeAndCall(Proxy(payable(resolved)));
    }

    function upgradeAndCall(Proxy _proxy) internal {
=======
    function test_erc1967UpgradeAndCall() external {
        upgradeAndCall(payable(proxy));
    }

    function test_chugsplashUpgradeAndCall() external {
        upgradeAndCall(payable(chugsplash));
    }

    function test_delegateResolvedUpgradeAndCall() external {
        upgradeAndCall(payable(resolved));
    }

    function upgradeAndCall(address payable _proxy) internal {
>>>>>>> v0.5.23
        vm.prank(alice);
        admin.upgradeAndCall(
            _proxy,
            address(implementation),
            abi.encodeWithSelector(SimpleStorage.set.selector, 1, 1)
        );

        address impl = admin.getProxyImplementation(_proxy);
        assertEq(impl, address(implementation));

        uint256 got = SimpleStorage(address(_proxy)).get(1);
        assertEq(got, 1);
    }

    function test_onlyOwner() external {
        vm.expectRevert("UNAUTHORIZED");
<<<<<<< HEAD
        admin.changeProxyAdmin(proxy, address(0));

        vm.expectRevert("UNAUTHORIZED");
        admin.upgrade(proxy, address(implementation));

        vm.expectRevert("UNAUTHORIZED");
        admin.upgradeAndCall(proxy, address(implementation), hex"");
=======
        admin.changeProxyAdmin(payable(proxy), address(0));

        vm.expectRevert("UNAUTHORIZED");
        admin.upgrade(payable(proxy), address(implementation));

        vm.expectRevert("UNAUTHORIZED");
        admin.upgradeAndCall(payable(proxy), address(implementation), hex"");
>>>>>>> v0.5.23
    }

    function test_isUpgrading() external {
        assertEq(false, admin.isUpgrading());

        vm.prank(alice);
        admin.setUpgrading(true);
        assertEq(true, admin.isUpgrading());
    }
}
