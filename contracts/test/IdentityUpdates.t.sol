// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import {Test, console} from "forge-std/src/Test.sol";
import {IdentityUpdates} from "src/IdentityUpdates.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract IdentityUpdatesTest is Test {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 constant EIP1967_IMPL_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 constant INBOX_ID = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;

    IdentityUpdates identityUpdatesImpl;
    ERC1967Proxy proxy;
    IdentityUpdates identityUpdates;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        identityUpdatesImpl = new IdentityUpdates();

        proxy = new ERC1967Proxy(
            address(identityUpdatesImpl), 
            abi.encodeWithSelector(identityUpdates.initialize.selector, admin)
        );

        identityUpdates = IdentityUpdates(address(proxy));
    }

    function test_AddIdentityUpdateValid() public {
        bytes memory message = new bytes(1024);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256)); // Set each byte to its index modulo 256
        }

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(INBOX_ID, message, 1);
        identityUpdates.addIdentityUpdate(INBOX_ID, message);
    }

    function testAddMessageInvalid() public {
        bytes memory message = new bytes(103);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256));
        }

        vm.expectRevert(IdentityUpdates.InvalidIdentityUpdate.selector);
        identityUpdates.addIdentityUpdate(INBOX_ID, message);
    }

    function testInvalidReinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        identityUpdates.initialize(admin);
    }

    function testPauseUnpause() public {
        identityUpdates.pause();
        assertTrue(identityUpdates.paused());

        vm.prank(unauthorized);
        vm.expectRevert();
        identityUpdates.unpause();

        identityUpdates.unpause();
        assertFalse(identityUpdates.paused());

        vm.prank(unauthorized);
        vm.expectRevert();
        identityUpdates.pause();
    }

    function testRoles() public {
        identityUpdates.grantRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.startPrank(unauthorized);
        identityUpdates.pause();
        identityUpdates.unpause();
        vm.stopPrank();

        identityUpdates.revokeRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.prank(unauthorized);
        vm.expectRevert(revertRoleData(unauthorized));
        identityUpdates.pause();

        identityUpdates.renounceRole(DEFAULT_ADMIN_ROLE, admin);
        vm.expectRevert(revertRoleData(admin));
        identityUpdates.pause();
    }

    function testUpgradeImplementation() public {
        IdentityUpdates newIdentityUpdatesImpl = new IdentityUpdates();
        address newImplAddress = address(newIdentityUpdatesImpl);
        address oldImplAddress = address(identityUpdatesImpl);
        
        bytes memory message = new bytes(104);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256));
        }

        // Retrieve the implementation address directly from the proxy storage.
        bytes32 rawImplAddress = vm.load(address(identityUpdates), EIP1967_IMPL_SLOT);
        address implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, oldImplAddress);

        // Initialize sequenceId to 1. The state should be preserved between upgrades.
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(INBOX_ID, message, 1);
        identityUpdates.addIdentityUpdate(INBOX_ID, message);

        // Unauthorized upgrade attempts should revert.
        vm.prank(unauthorized);
        vm.expectRevert(revertRoleData(unauthorized));
        identityUpdates.upgradeToAndCall(address(newIdentityUpdatesImpl), "");

        // Authorized upgrade should succeed and emit UpgradeAuthorized event.
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.UpgradeAuthorized(address(this), address(newIdentityUpdatesImpl));
        identityUpdates.upgradeToAndCall(address(newIdentityUpdatesImpl), "");

        // Retrieve the new implementation address directly from the proxy storage.
        rawImplAddress = vm.load(address(identityUpdates), EIP1967_IMPL_SLOT);
        implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, newImplAddress);

        // Next sequenceId should be 2.
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(INBOX_ID, message, 2);
        identityUpdates.addIdentityUpdate(INBOX_ID, message);
    }

    function revertRoleData(address _user) public pure returns (bytes memory) {
        return abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, _user, DEFAULT_ADMIN_ROLE);
    }
}
