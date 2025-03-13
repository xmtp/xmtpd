// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { Test } from "forge-std/src/Test.sol";

import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

import { RatesManager } from "../src/RatesManager.sol";

import { RatesManagerHarness } from "./utils/Harnesses.sol";
import { Utils } from "./utils/Utils.sol";

contract RatesTest is Test, Utils {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 constant RATES_MANAGER_ROLE = keccak256("RATES_MANAGER_ROLE");

    uint256 constant PAGE_SIZE = 50;

    uint64 constant MESSAGE_FEE = 100;
    uint64 constant STORAGE_FEE = 200;
    uint64 constant CONGESTION_FEE = 300;

    address ratesManagerImpImplementation;

    RatesManagerHarness ratesManager;

    address admin = makeAddr("admin");
    address unauthorized = makeAddr("unauthorized");

    function setUp() public {
        ratesManagerImpImplementation = address(new RatesManagerHarness());

        ratesManager = RatesManagerHarness(
            address(
                new ERC1967Proxy(
                    ratesManagerImpImplementation, abi.encodeWithSelector(RatesManager.initialize.selector, admin)
                )
            )
        );
    }

    /* ============ initializer ============ */

    function test_initializer_zeroAdminAddress() public {
        vm.expectRevert(RatesManager.ZeroAdminAddress.selector);

        new ERC1967Proxy(
            ratesManagerImpImplementation, abi.encodeWithSelector(RatesManager.initialize.selector, address(0))
        );
    }

    /* ============ initial state ============ */

    function test_initialState() public view {
        assertEq(_getImplementationFromSlot(address(ratesManager)), ratesManagerImpImplementation);
        assertEq(ratesManager.__getAllRates().length, 0);
    }

    /* ============ addRates ============ */

    function test_addRates_notManager() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, RATES_MANAGER_ROLE
            )
        );

        vm.prank(unauthorized);
        ratesManager.addRates(0, 0, 0, 0);

        // TODO: Test where admin is not the manager.
    }

    function test_addRates_first() public {
        vm.expectEmit(address(ratesManager));
        emit RatesManager.RatesAdded(100, 200, 300, 400);

        vm.prank(admin);
        ratesManager.addRates(100, 200, 300, 400);

        RatesManager.Rates[] memory rates = ratesManager.__getAllRates();

        assertEq(rates.length, 1);

        assertEq(rates[0].messageFee, 100);
        assertEq(rates[0].storageFee, 200);
        assertEq(rates[0].congestionFee, 300);
        assertEq(rates[0].startTime, 400);
    }

    function test_addRates_nth() public {
        ratesManager.__pushRates(0, 0, 0, 0);
        ratesManager.__pushRates(0, 0, 0, 0);
        ratesManager.__pushRates(0, 0, 0, 0);
        ratesManager.__pushRates(0, 0, 0, 0);

        vm.expectEmit(address(ratesManager));
        emit RatesManager.RatesAdded(100, 200, 300, 400);

        vm.prank(admin);
        ratesManager.addRates(100, 200, 300, 400);

        RatesManager.Rates[] memory rates = ratesManager.__getAllRates();

        assertEq(rates.length, 5);

        assertEq(rates[4].messageFee, 100);
        assertEq(rates[4].storageFee, 200);
        assertEq(rates[4].congestionFee, 300);
        assertEq(rates[4].startTime, 400);
    }

    function test_addRates_invalidStartTime() public {
        ratesManager.__pushRates(0, 0, 0, 100);

        vm.expectRevert(RatesManager.InvalidStartTime.selector);

        vm.prank(admin);
        ratesManager.addRates(0, 0, 0, 100);
    }

    /* ============ getRates ============ */

    function test_getRates_emptyArray() public view {
        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);

        assertEq(rates.length, 0);
        assertFalse(hasMore);
    }

    function test_getRates_withinPageSize() public {
        for (uint256 i; i < 3 * PAGE_SIZE; ++i) {
            ratesManager.__pushRates(i, i, i, i);
        }

        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates((3 * PAGE_SIZE) - 10);

        assertEq(rates.length, 10);
        assertFalse(hasMore);

        for (uint256 i; i < rates.length; ++i) {
            assertEq(rates[i].messageFee, i + (3 * PAGE_SIZE) - 10);
            assertEq(rates[i].storageFee, i + (3 * PAGE_SIZE) - 10);
            assertEq(rates[i].congestionFee, i + (3 * PAGE_SIZE) - 10);
            assertEq(rates[i].startTime, i + (3 * PAGE_SIZE) - 10);
        }
    }

    function test_getRates_pagination() public {
        for (uint256 i; i < 3 * PAGE_SIZE; ++i) {
            ratesManager.__pushRates(i, i, i, i);
        }

        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);

        assertEq(rates.length, PAGE_SIZE);
        assertTrue(hasMore);

        for (uint256 i; i < rates.length; ++i) {
            assertEq(rates[i].messageFee, i);
            assertEq(rates[i].storageFee, i);
            assertEq(rates[i].congestionFee, i);
            assertEq(rates[i].startTime, i);
        }
    }

    /* ============ getRatesCount ============ */

    function test_getRatesCount() public {
        assertEq(ratesManager.getRatesCount(), 0);

        for (uint256 i = 1; i <= 1000; ++i) {
            ratesManager.__pushRates(0, 0, 0, 0);
            assertEq(ratesManager.getRatesCount(), i);
        }
    }

    /* ============ pause ============ */

    function test_pause() public {
        vm.expectEmit(address(ratesManager));
        emit PausableUpgradeable.Paused(admin);

        vm.prank(admin);
        ratesManager.pause();

        assertTrue(ratesManager.paused());
    }

    function test_pause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        ratesManager.pause();
    }

    function test_pause_whenPaused() public {
        ratesManager.__pause();

        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);

        vm.prank(admin);
        ratesManager.pause();
    }

    /* ============ unpause ============ */

    function test_unpause() public {
        ratesManager.__pause();

        vm.expectEmit(address(ratesManager));
        emit PausableUpgradeable.Unpaused(admin);

        vm.prank(admin);
        ratesManager.unpause();

        assertFalse(ratesManager.paused());
    }

    function test_unpause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        ratesManager.unpause();
    }

    function test_unpause_whenNotPaused() public {
        vm.expectRevert(PausableUpgradeable.ExpectedPause.selector);

        vm.prank(admin);
        ratesManager.unpause();
    }

    /* ============ upgradeToAndCall ============ */

    function test_upgradeToAndCall_notAdmin() public {
        // Unauthorized upgrade attempts should revert.
        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        ratesManager.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall_zeroImplementationAddress() public {
        vm.expectRevert(RatesManager.ZeroImplementationAddress.selector);

        vm.prank(admin);
        ratesManager.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall() public {
        ratesManager.__pushRates(0, 0, 0, 0);
        ratesManager.__pushRates(1, 1, 1, 1);
        ratesManager.__pushRates(2, 2, 2, 2);

        address newImplementation = address(new RatesManagerHarness());

        // Authorized upgrade should succeed and emit UpgradeAuthorized event.
        vm.expectEmit(address(ratesManager));
        emit RatesManager.UpgradeAuthorized(admin, newImplementation);

        vm.prank(admin);
        ratesManager.upgradeToAndCall(newImplementation, "");

        assertEq(_getImplementationFromSlot(address(ratesManager)), newImplementation);

        RatesManager.Rates[] memory rates = ratesManager.__getAllRates();

        assertEq(rates.length, 3);

        assertEq(rates[0].messageFee, 0);
        assertEq(rates[0].storageFee, 0);
        assertEq(rates[0].congestionFee, 0);
        assertEq(rates[0].startTime, 0);

        assertEq(rates[1].messageFee, 1);
        assertEq(rates[1].storageFee, 1);
        assertEq(rates[1].congestionFee, 1);
        assertEq(rates[1].startTime, 1);

        assertEq(rates[2].messageFee, 2);
        assertEq(rates[2].storageFee, 2);
        assertEq(rates[2].congestionFee, 2);
        assertEq(rates[2].startTime, 2);
    }

    /* ============ helper functions ============ */

    function _getImplementationFromSlot(address proxy) internal view returns (address) {
        // Retrieve the implementation address directly from the proxy storage.
        return address(uint160(uint256(vm.load(proxy, EIP1967_IMPLEMENTATION_SLOT))));
    }
}
