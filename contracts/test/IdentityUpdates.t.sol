// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { Test } from "forge-std/src/Test.sol";

import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

import { IdentityUpdates } from "../src/IdentityUpdates.sol";

import { IdentityUpdatesHarness } from "./utils/Harnesses.sol";
import { Utils } from "./utils/Utils.sol";

contract IdentityUpdatesTest is Test, Utils {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;

    uint256 constant ABSOLUTE_MIN_PAYLOAD_SIZE = 78;
    uint256 constant ABSOLUTE_MAX_PAYLOAD_SIZE = 4_194_304;

    address identityUpdatesImplementation;

    IdentityUpdatesHarness identityUpdates;

    address admin = makeAddr("admin");
    address unauthorized = makeAddr("unauthorized");

    function setUp() public {
        identityUpdatesImplementation = address(new IdentityUpdatesHarness());

        identityUpdates = IdentityUpdatesHarness(
            address(
                new ERC1967Proxy(
                    identityUpdatesImplementation, abi.encodeWithSelector(identityUpdates.initialize.selector, admin)
                )
            )
        );
    }

    /* ============ initializer ============ */

    function test_initializer_zeroAdminAddress() public {
        vm.expectRevert(IdentityUpdates.ZeroAdminAddress.selector);

        new ERC1967Proxy(
            identityUpdatesImplementation, abi.encodeWithSelector(IdentityUpdates.initialize.selector, address(0))
        );
    }

    /* ============ initial state ============ */

    function test_initialState() public view {
        assertEq(_getImplementationFromSlot(address(identityUpdates)), identityUpdatesImplementation);
        assertEq(identityUpdates.minPayloadSize(), ABSOLUTE_MIN_PAYLOAD_SIZE);
        assertEq(identityUpdates.maxPayloadSize(), ABSOLUTE_MAX_PAYLOAD_SIZE);
        assertEq(identityUpdates.__getSequenceId(), 0);
    }

    /* ============ addIdentityUpdate ============ */

    function test_addIdentityUpdate_minPayload() public {
        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);

        identityUpdates.addIdentityUpdate(ID, message);

        assertEq(identityUpdates.__getSequenceId(), 1);
    }

    function test_addIdentityUpdate_maxPayload() public {
        bytes memory message = _generatePayload(identityUpdates.maxPayloadSize());

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);

        identityUpdates.addIdentityUpdate(ID, message);

        assertEq(identityUpdates.__getSequenceId(), 1);
    }

    function test_addIdentityUpdate_payloadTooSmall() public {
        bytes memory message = _generatePayload(identityUpdates.minPayloadSize() - 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                IdentityUpdates.InvalidPayloadSize.selector,
                message.length,
                identityUpdates.minPayloadSize(),
                identityUpdates.maxPayloadSize()
            )
        );

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function test_addIdentityUpdate_payloadTooLarge() public {
        bytes memory message = _generatePayload(identityUpdates.maxPayloadSize() + 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                IdentityUpdates.InvalidPayloadSize.selector,
                message.length,
                identityUpdates.minPayloadSize(),
                identityUpdates.maxPayloadSize()
            )
        );

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function test_addIdentityUpdate_whenPaused() public {
        identityUpdates.__pause();

        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function testFuzz_addIdentityUpdate(
        uint256 minPayloadSize,
        uint256 maxPayloadSize,
        uint256 payloadSize,
        uint64 sequenceId,
        bool paused
    ) public {
        minPayloadSize = bound(minPayloadSize, ABSOLUTE_MIN_PAYLOAD_SIZE, ABSOLUTE_MAX_PAYLOAD_SIZE);
        maxPayloadSize = bound(maxPayloadSize, minPayloadSize, ABSOLUTE_MAX_PAYLOAD_SIZE);
        payloadSize = bound(payloadSize, ABSOLUTE_MIN_PAYLOAD_SIZE, ABSOLUTE_MAX_PAYLOAD_SIZE);
        sequenceId = uint64(bound(sequenceId, 0, type(uint64).max - 1));

        identityUpdates.__setSequenceId(sequenceId);
        identityUpdates.__setMinPayloadSize(minPayloadSize);
        identityUpdates.__setMaxPayloadSize(maxPayloadSize);

        if (paused) {
            identityUpdates.__pause();
        }

        bytes memory message = _generatePayload(payloadSize);

        bool shouldFail = (payloadSize < minPayloadSize) || (payloadSize > maxPayloadSize) || paused;

        if (shouldFail) {
            vm.expectRevert();
        } else {
            vm.expectEmit(address(identityUpdates));
            emit IdentityUpdates.IdentityUpdateCreated(ID, message, sequenceId + 1);
        }

        identityUpdates.addIdentityUpdate(ID, message);

        if (shouldFail) return;

        assertEq(identityUpdates.__getSequenceId(), sequenceId + 1);
    }

    /* ============ setMinPayloadSize ============ */

    function test_setMinPayloadSize_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        identityUpdates.setMinPayloadSize(0);
    }

    function test_setMinPayloadSize_requestGreaterThanMax() public {
        identityUpdates.__setMaxPayloadSize(100);

        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMinPayloadSize.selector));

        vm.prank(admin);
        identityUpdates.setMinPayloadSize(101);
    }

    function test_setMinPayloadSize_requestLessThanOrEqualToAbsoluteMin() public {
        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMinPayloadSize.selector));

        vm.prank(admin);
        identityUpdates.setMinPayloadSize(ABSOLUTE_MIN_PAYLOAD_SIZE - 1);
    }

    function test_setMinPayloadSize() public {
        uint256 initialMinSize = identityUpdates.minPayloadSize();
        uint256 newMinSize = initialMinSize + 1;

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.MinPayloadSizeUpdated(initialMinSize, newMinSize);

        vm.prank(admin);
        identityUpdates.setMinPayloadSize(newMinSize);

        assertEq(identityUpdates.minPayloadSize(), newMinSize);
    }

    /* ============ setMaxPayloadSize ============ */

    function test_setMaxPayloadSize_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        identityUpdates.setMaxPayloadSize(0);
    }

    function test_setMaxPayloadSize_requestLessThanMin() public {
        identityUpdates.__setMinPayloadSize(100);

        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMaxPayloadSize.selector));

        vm.prank(admin);
        identityUpdates.setMaxPayloadSize(99);
    }

    function test_setMaxPayloadSize_requestGreaterThanOrEqualToAbsoluteMax() public {
        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMaxPayloadSize.selector));

        vm.prank(admin);
        identityUpdates.setMaxPayloadSize(ABSOLUTE_MAX_PAYLOAD_SIZE + 1);
    }

    function test_setMaxPayloadSize() public {
        uint256 initialMaxSize = identityUpdates.maxPayloadSize();
        uint256 newMaxSize = initialMaxSize - 1;

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.MaxPayloadSizeUpdated(initialMaxSize, newMaxSize);

        vm.prank(admin);
        identityUpdates.setMaxPayloadSize(newMaxSize);

        assertEq(identityUpdates.maxPayloadSize(), newMaxSize);
    }

    /* ============ initialize ============ */

    function test_invalid_reinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        identityUpdates.initialize(admin);
    }

    /* ============ pause ============ */

    function test_pause() public {
        vm.expectEmit(address(identityUpdates));
        emit PausableUpgradeable.Paused(admin);

        vm.prank(admin);
        identityUpdates.pause();

        assertTrue(identityUpdates.paused());
    }

    function test_pause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        identityUpdates.pause();
    }

    function test_pause_whenPaused() public {
        identityUpdates.__pause();

        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);

        vm.prank(admin);
        identityUpdates.pause();
    }

    /* ============ unpause ============ */

    function test_unpause() public {
        identityUpdates.__pause();

        vm.expectEmit(address(identityUpdates));
        emit PausableUpgradeable.Unpaused(admin);

        vm.prank(admin);
        identityUpdates.unpause();

        assertFalse(identityUpdates.paused());
    }

    function test_unpause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        identityUpdates.unpause();
    }

    function test_unpause_whenNotPaused() public {
        vm.expectRevert(PausableUpgradeable.ExpectedPause.selector);

        vm.prank(admin);
        identityUpdates.unpause();
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

        identityUpdates.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall_zeroImplementationAddress() public {
        vm.expectRevert(IdentityUpdates.ZeroImplementationAddress.selector);

        vm.prank(admin);
        identityUpdates.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall() public {
        identityUpdates.__setMaxPayloadSize(100);
        identityUpdates.__setMinPayloadSize(50);
        identityUpdates.__setSequenceId(10);

        address newImplementation = address(new IdentityUpdatesHarness());

        // Authorized upgrade should succeed and emit UpgradeAuthorized event.
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.UpgradeAuthorized(admin, newImplementation);

        vm.prank(admin);
        identityUpdates.upgradeToAndCall(newImplementation, "");

        assertEq(_getImplementationFromSlot(address(identityUpdates)), newImplementation);
        assertEq(identityUpdates.maxPayloadSize(), 100);
        assertEq(identityUpdates.minPayloadSize(), 50);
        assertEq(identityUpdates.__getSequenceId(), 10);
    }

    /* ============ helper functions ============ */

    function _getImplementationFromSlot(address proxy) internal view returns (address) {
        // Retrieve the implementation address directly from the proxy storage.
        return address(uint160(uint256(vm.load(proxy, EIP1967_IMPLEMENTATION_SLOT))));
    }
}
