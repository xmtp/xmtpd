// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import {Test, console} from "forge-std/src/Test.sol";
import {Utils} from "./utils/Utils.sol";
import {GroupMessages} from "../src/GroupMessages.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {PausableUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

contract GroupMessagesTest is Test, GroupMessages, Utils {
    GroupMessages groupMessagesImpl;
    ERC1967Proxy proxy;
    GroupMessages groupMessages;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        groupMessagesImpl = new GroupMessages();

        proxy = new ERC1967Proxy(
            address(groupMessagesImpl), abi.encodeWithSelector(GroupMessages.initialize.selector, admin)
        );

        groupMessages = GroupMessages(address(proxy));
    }

    function testAddMessageValid() public {
        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);

        groupMessages.addMessage(ID, message);
    }

    function testAddMessageWithMaxPayload() public {
        bytes memory message = _generatePayload(groupMessages.maxPayloadSize());

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);

        groupMessages.addMessage(ID, message);
    }

    function testAddMessageTooSmall() public {
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

    function testAddMessageTooBig() public {
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

    function testAddMessageWhenPaused() public {
        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        groupMessages.pause();
        assertTrue(groupMessages.paused());

        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

        groupMessages.addMessage(ID, message);
    }

    function testSetMinPayloadSize() public {
        // Store initial min payload size
        uint256 initialMinSize = groupMessages.minPayloadSize();
        uint256 newMinSize = initialMinSize + 1;

        // Test unauthorized access
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                groupMessages.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        groupMessages.setMinPayloadSize(newMinSize);

        // Test authorized access
        groupMessages.setMinPayloadSize(newMinSize);
        assertEq(groupMessages.minPayloadSize(), newMinSize);

        // Verify that messages with old minPayloadSize now fail
        bytes memory message = _generatePayload(initialMinSize);
        vm.expectRevert(
            abi.encodeWithSelector(
                GroupMessages.InvalidPayloadSize.selector,
                message.length,
                newMinSize,
                groupMessages.maxPayloadSize()
            )
        );
        groupMessages.addMessage(ID, message);

        // Verify that messages with new minPayloadSize succeed
        bytes memory validMessage = _generatePayload(newMinSize);
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, validMessage, 1);
        groupMessages.addMessage(ID, validMessage);

        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMinPayloadSize.selector));
        groupMessages.setMinPayloadSize(0);

        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMinPayloadSize.selector));
        groupMessages.setMinPayloadSize(4194304);
    }

    function testSetmaxPayloadSize() public {
        // Store initial max payload size
        uint256 initialMaxSize = groupMessages.maxPayloadSize();
        uint256 newMaxSize = 1000;

        // Test unauthorized access
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                groupMessages.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        groupMessages.setMaxPayloadSize(newMaxSize);

        // Test authorized access
        groupMessages.setMaxPayloadSize(newMaxSize);
        assertEq(groupMessages.maxPayloadSize(), newMaxSize);

        // Verify that messages with old maxPayloadSize now fail
        bytes memory message = _generatePayload(initialMaxSize);
        vm.expectRevert(
            abi.encodeWithSelector(
                GroupMessages.InvalidPayloadSize.selector,
                message.length,
                groupMessages.minPayloadSize(),
                newMaxSize
            )
        );
        groupMessages.addMessage(ID, message);

        // Verify that messages with new maxPayloadSize succeed
        bytes memory validMessage = _generatePayload(newMaxSize);
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, validMessage, 1);
        groupMessages.addMessage(ID, validMessage);

        // Max size should always be greater than min size
        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMaxPayloadSize.selector));
        groupMessages.setMaxPayloadSize(78);

        // Test setting max size above maxPayloadSize (should fail)
        vm.expectRevert(abi.encodeWithSelector(GroupMessages.InvalidMaxPayloadSize.selector));
        groupMessages.setMaxPayloadSize(4_194_305);
    }

    function testSequenceIdIncrement() public {
        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);
        groupMessages.addMessage(ID, message);

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 2);
        groupMessages.addMessage(ID, message);

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 3);
        groupMessages.addMessage(ID, message);
    }

    function testInvalidReinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        groupMessages.initialize(admin);
    }

    function testPauseUnpause() public {
        groupMessages.pause();
        assertTrue(groupMessages.paused());

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );
        groupMessages.unpause();

        groupMessages.unpause();
        assertFalse(groupMessages.paused());

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );
        groupMessages.pause();
    }

    function testRoles() public {
        groupMessages.grantRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.startPrank(unauthorized);
        groupMessages.pause();
        groupMessages.unpause();
        vm.stopPrank();

        groupMessages.revokeRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );
        groupMessages.pause();

        groupMessages.renounceRole(DEFAULT_ADMIN_ROLE, admin);
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, admin, DEFAULT_ADMIN_ROLE)
        );
        groupMessages.pause();
    }

    function testUpgradeImplementation() public {
        GroupMessages newGroupMessagesImpl = new GroupMessages();
        address newImplAddress = address(newGroupMessagesImpl);
        address oldImplAddress = address(groupMessagesImpl);

        bytes memory message = _generatePayload(groupMessages.minPayloadSize());

        // Retrieve the implementation address directly from the proxy storage.
        bytes32 rawImplAddress = vm.load(address(groupMessages), EIP1967_IMPL_SLOT);
        address implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, oldImplAddress);

        // Initialize sequenceId to 1. The state should be preserved between upgrades.
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 1);
        groupMessages.addMessage(ID, message);

        // Unauthorized upgrade attempts should revert.
        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, DEFAULT_ADMIN_ROLE
            )
        );
        groupMessages.upgradeToAndCall(address(newGroupMessagesImpl), "");

        // Authorized upgrade should succeed and emit UpgradeAuthorized event.
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.UpgradeAuthorized(address(this), address(newGroupMessagesImpl));
        groupMessages.upgradeToAndCall(address(newGroupMessagesImpl), "");

        // Retrieve the new implementation address directly from the proxy storage.
        rawImplAddress = vm.load(address(groupMessages), EIP1967_IMPL_SLOT);
        implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, newImplAddress);

        // Next sequenceId should be 2.
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(ID, message, 2);
        groupMessages.addMessage(ID, message);
    }
}
