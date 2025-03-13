// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { Test } from "forge-std/src/Test.sol";

import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";

import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { Initializable } from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import { PausableUpgradeable } from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

import { GroupMessages } from "../src/GroupMessages.sol";

import { GroupMessagesHarness } from "./utils/Harnesses.sol";
import { Utils } from "./utils/Utils.sol";

contract GroupMessagesTest is Test, Utils {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;

    uint256 constant ABSOLUTE_MIN_PAYLOAD_SIZE = 78;
    uint256 constant ABSOLUTE_MAX_PAYLOAD_SIZE = 4_194_304;

    address groupMessagesImplementation;

    GroupMessagesHarness groupMessages;

    address admin = makeAddr("admin");
    address unauthorized = makeAddr("unauthorized");

    function setUp() public {
        groupMessagesImplementation = address(new GroupMessagesHarness());

        groupMessages = GroupMessagesHarness(
            address(
                new ERC1967Proxy(
                    groupMessagesImplementation, abi.encodeWithSelector(GroupMessages.initialize.selector, admin)
                )
            )
        );
    }

    /* ============ initializer ============ */

    function test_initializer_zeroAdminAddress() public {
        vm.expectRevert(GroupMessages.ZeroAdminAddress.selector);

        new ERC1967Proxy(
            groupMessagesImplementation, abi.encodeWithSelector(GroupMessages.initialize.selector, address(0))
        );
    }

    /* ============ initial state ============ */

    function test_initialState() public view {
        assertEq(_getImplementationFromSlot(address(groupMessages)), groupMessagesImplementation);
        assertEq(groupMessages.minPayloadSize(), ABSOLUTE_MIN_PAYLOAD_SIZE);
        assertEq(groupMessages.maxPayloadSize(), ABSOLUTE_MAX_PAYLOAD_SIZE);
        assertEq(groupMessages.__getSequenceId(), 0);
    }

    /* ============ addMessage ============ */

    function test_addMessage_minPayload() public {
        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);

        groupMessages.addMessage(ID, message);

        assertEq(groupMessages.__getSequenceId(), 1);
    }

    function test_addMessage_maxPayload() public {
        bytes memory message = _generatePayload(groupMessages.maxPayloadSize());

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);

        groupMessages.addMessage(ID, message);

        assertEq(groupMessages.__getSequenceId(), 1);
    }

    function test_addMessage_payloadTooSmall() public {
        bytes memory message = _generatePayload(groupMessages.minPayloadSize() - 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                GroupMessages.InvalidPayloadSize.selector,
                message.length,
                groupMessages.minPayloadSize(),
                groupMessages.maxPayloadSize()
            )
        );

        groupMessages.addMessage(ID, message);
    }

    function test_addMessage_payloadTooLarge() public {
        bytes memory message = _generatePayload(groupMessages.maxPayloadSize() + 1);

        vm.expectRevert(
            abi.encodeWithSelector(
                GroupMessages.InvalidPayloadSize.selector,
                message.length,
                groupMessages.minPayloadSize(),
                groupMessages.maxPayloadSize()
            )
        );

        groupMessages.addMessage(ID, message);
    }

    function test_addMessage_whenPaused() public {
        groupMessages.__pause();

        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

        groupMessages.addMessage(ID, message);
    }

    function testFuzz_addMessage(
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

        groupMessages.__setSequenceId(sequenceId);
        groupMessages.__setMinPayloadSize(minPayloadSize);
        groupMessages.__setMaxPayloadSize(maxPayloadSize);

        if (paused) {
            groupMessages.__pause();
        }

        bytes memory message = _generatePayload(payloadSize);

        bool shouldFail = (payloadSize < minPayloadSize) || (payloadSize > maxPayloadSize) || paused;

        if (shouldFail) {
            vm.expectRevert();
        } else {
            vm.expectEmit(address(groupMessages));
            emit GroupMessages.MessageSent(ID, message, sequenceId + 1);
        }

        groupMessages.addMessage(ID, message);

        if (shouldFail) return;

        assertEq(groupMessages.__getSequenceId(), sequenceId + 1);
    }

    /* ============ setMinPayloadSize ============ */

    function test_setMinPayloadSize_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        groupMessages.setMinPayloadSize(0);
    }

    function test_setMinPayloadSize_requestGreaterThanMax() public {
        groupMessages.__setMaxPayloadSize(100);

        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMinPayloadSize.selector));

        vm.prank(admin);
        groupMessages.setMinPayloadSize(101);
    }

    function test_setMinPayloadSize_requestLessThanOrEqualToAbsoluteMin() public {
        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMinPayloadSize.selector));

        vm.prank(admin);
        groupMessages.setMinPayloadSize(ABSOLUTE_MIN_PAYLOAD_SIZE - 1);
    }

    function test_setMinPayloadSize() public {
        uint256 initialMinSize = groupMessages.minPayloadSize();
        uint256 newMinSize = initialMinSize + 1;

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MinPayloadSizeUpdated(initialMinSize, newMinSize);

        vm.prank(admin);
        groupMessages.setMinPayloadSize(newMinSize);

        assertEq(groupMessages.minPayloadSize(), newMinSize);
    }

    /* ============ setMaxPayloadSize ============ */

    function test_setMaxPayloadSize_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        groupMessages.setMaxPayloadSize(0);
    }

    function test_setMaxPayloadSize_requestLessThanMin() public {
        groupMessages.__setMinPayloadSize(100);

        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMaxPayloadSize.selector));

        vm.prank(admin);
        groupMessages.setMaxPayloadSize(99);
    }

    function test_setMaxPayloadSize_requestGreaterThanOrEqualToAbsoluteMax() public {
        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMaxPayloadSize.selector));

        vm.prank(admin);
        groupMessages.setMaxPayloadSize(ABSOLUTE_MAX_PAYLOAD_SIZE + 1);
    }

    function test_setMaxPayloadSize() public {
        uint256 initialMaxSize = groupMessages.maxPayloadSize();
        uint256 newMaxSize = initialMaxSize - 1;

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MaxPayloadSizeUpdated(initialMaxSize, newMaxSize);

        vm.prank(admin);
        groupMessages.setMaxPayloadSize(newMaxSize);

        assertEq(groupMessages.maxPayloadSize(), newMaxSize);
    }

    /* ============ initialize ============ */

    function test_invalid_reinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        groupMessages.initialize(admin);
    }

    /* ============ pause ============ */

    function test_pause() public {
        vm.expectEmit(address(groupMessages));
        emit PausableUpgradeable.Paused(admin);

        vm.prank(admin);
        groupMessages.pause();

        assertTrue(groupMessages.paused());
    }

    function test_pause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        groupMessages.pause();
    }

    function test_pause_whenPaused() public {
        groupMessages.__pause();

        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);

        vm.prank(admin);
        groupMessages.pause();
    }

    /* ============ unpause ============ */

    function test_unpause() public {
        groupMessages.__pause();

        vm.expectEmit(address(groupMessages));
        emit PausableUpgradeable.Unpaused(admin);

        vm.prank(admin);
        groupMessages.unpause();

        assertFalse(groupMessages.paused());
    }

    function test_unpause_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );

        vm.prank(unauthorized);
        groupMessages.unpause();
    }

    function test_unpause_whenNotPaused() public {
        vm.expectRevert(PausableUpgradeable.ExpectedPause.selector);

        vm.prank(admin);
        groupMessages.unpause();
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

        groupMessages.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall_zeroImplementationAddress() public {
        vm.expectRevert(GroupMessages.ZeroImplementationAddress.selector);

        vm.prank(admin);
        groupMessages.upgradeToAndCall(address(0), "");
    }

    function test_upgradeToAndCall() public {
        groupMessages.__setMaxPayloadSize(100);
        groupMessages.__setMinPayloadSize(50);
        groupMessages.__setSequenceId(10);

        address newImplementation = address(new GroupMessagesHarness());

        // Authorized upgrade should succeed and emit UpgradeAuthorized event.
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.UpgradeAuthorized(admin, newImplementation);

        vm.prank(admin);
        groupMessages.upgradeToAndCall(newImplementation, "");

        assertEq(_getImplementationFromSlot(address(groupMessages)), newImplementation);
        assertEq(groupMessages.maxPayloadSize(), 100);
        assertEq(groupMessages.minPayloadSize(), 50);
        assertEq(groupMessages.__getSequenceId(), 10);
    }

    /* ============ helper functions ============ */

    function _getImplementationFromSlot(address proxy) internal view returns (address) {
        // Retrieve the implementation address directly from the proxy storage.
        return address(uint160(uint256(vm.load(proxy, EIP1967_IMPLEMENTATION_SLOT))));
    }
}
