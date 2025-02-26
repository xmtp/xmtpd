// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Test.sol";
import {RatesManager} from "../src/RatesManager.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {PausableUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

contract RatesTest is Test {
    RatesManager ratesManagerImpl;
    ERC1967Proxy proxy;
    RatesManager ratesManager;

    address admin = address(this);
    address unauthorized = address(0x1);

    uint64 constant messageFee = 100;
    uint64 constant storageFee = 200;
    uint64 constant congestionFee = 300;

    function setUp() public {
        ratesManagerImpl = new RatesManager();

        proxy =
            new ERC1967Proxy(address(ratesManagerImpl), abi.encodeWithSelector(RatesManager.initialize.selector, admin));

        ratesManager = RatesManager(address(proxy));
    }

    function addRate(uint64 startTime) internal {
        ratesManager.addRates(messageFee, storageFee, congestionFee, startTime);
    }

    function testAddRatesValid() public {
        uint64 startTime = uint64(block.timestamp + 1);

        vm.expectEmit(true, true, true, true);
        emit RatesManager.RatesAdded(messageFee, storageFee, congestionFee, startTime);

        addRate(startTime);

        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);
        assertEq(rates.length, 1);
        assertEq(rates[0].messageFee, messageFee);
        assertEq(rates[0].storageFee, storageFee);
        assertEq(rates[0].congestionFee, congestionFee);
        assertEq(rates[0].startTime, startTime);
        assertFalse(hasMore);
    }

    function testAddRatesUnauthorized() public {
        uint64 startTime = uint64(block.timestamp + 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                ratesManager.RATES_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        ratesManager.addRates(messageFee, storageFee, congestionFee, startTime);
    }

    function testAddRatesChronologicalOrder() public {
        uint64 startTime1 = uint64(block.timestamp + 1);
        uint64 startTime2 = uint64(block.timestamp + 2);

        addRate(startTime1);

        vm.expectRevert(abi.encodeWithSelector(RatesManager.InvalidStartTime.selector));
        addRate(startTime1);
        addRate(startTime2);

        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);
        assertEq(rates.length, 2);
        assertEq(rates[1].startTime, startTime2);
        assertFalse(hasMore);
    }

    function testGetRatesPagination() public {
        for (uint64 i = 0; i < 60; i++) {
            addRate(uint64(block.timestamp + i + 1));
        }

        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);
        assertEq(rates.length, 50);
        assertTrue(hasMore);

        (rates, hasMore) = ratesManager.getRates(50);
        assertEq(rates.length, 10);
        assertFalse(hasMore);
    }

    function testGetRatesEmptyArray() public view {
        // Verify the rates count is zero
        assertEq(ratesManager.getRatesCount(), 0);

        // Query the empty rates list
        (RatesManager.Rates[] memory rates, bool hasMore) = ratesManager.getRates(0);
        assertEq(rates.length, 0);
        assertFalse(hasMore);
    }
}
