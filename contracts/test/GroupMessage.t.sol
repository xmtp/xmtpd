// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import {Test, console} from "forge-std/src/Test.sol";
import {GroupMessages} from "src/GroupMessages.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract GroupMessagesTest is Test {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 constant EIP1967_IMPL_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 constant GROUP_ID = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;

    GroupMessages groupMessagesImpl;
    ERC1967Proxy proxy;
    GroupMessages groupMessages;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        groupMessagesImpl = new GroupMessages();

        proxy = new ERC1967Proxy(
            address(groupMessagesImpl), 
            abi.encodeWithSelector(GroupMessages.initialize.selector, admin)
        );

        groupMessages = GroupMessages(address(proxy));
    }

    function testAddMessageValid() public {
        bytes memory message = new bytes(1024);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256)); // Set each byte to its index modulo 256
        }

        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(GROUP_ID, message, 1);

        groupMessages.addMessage(GROUP_ID, message);
    }

    function testAddMessageInvalid() public {
        bytes memory message = new bytes(77);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256));
    }

        vm.expectRevert(GroupMessages.InvalidMessage.selector);
        groupMessages.addMessage(GROUP_ID, message);
    }

    function testInvalidReinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        groupMessages.initialize(admin);
    }

    function testPauseUnpause() public {
        groupMessages.pause();

        vm.prank(unauthorized);
        vm.expectRevert();
        groupMessages.unpause();

        groupMessages.unpause();

        vm.prank(unauthorized);
        vm.expectRevert();
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
        vm.expectRevert(revertRoleData(unauthorized));
        groupMessages.pause();

        groupMessages.renounceRole(DEFAULT_ADMIN_ROLE, admin);
        vm.expectRevert(revertRoleData(admin));
        groupMessages.pause();
    }

    function testUpgradeImplementation() public {
        GroupMessages newGroupMessagesImpl = new GroupMessages();
        address newImplAddress = address(newGroupMessagesImpl);
        address oldImplAddress = address(groupMessagesImpl);
        
        bytes memory message = new bytes(78);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256));
        }

        // Retrieve the implementation address directly from the proxy storage.
        bytes32 rawImplAddress = vm.load(address(groupMessages), EIP1967_IMPL_SLOT);
        address implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, oldImplAddress);

        // Initialize sequenceId to 1. The state should be preserved between upgrades.
        vm.expectEmit(address(groupMessages));
        emit GroupMessages.MessageSent(GROUP_ID, message, 1);
        groupMessages.addMessage(GROUP_ID, message);

        // Unauthorized upgrade attempts should revert.
        vm.prank(unauthorized);
        vm.expectRevert(revertRoleData(unauthorized));
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
        emit GroupMessages.MessageSent(GROUP_ID, message, 2);
        groupMessages.addMessage(GROUP_ID, message);
    }

    function revertRoleData(address _user) public pure returns (bytes memory) {
        return abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, _user, DEFAULT_ADMIN_ROLE);
    }
}