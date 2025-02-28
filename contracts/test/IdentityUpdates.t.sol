// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Vm.sol";
import {Test} from "forge-std/src/Test.sol";
import {Utils} from "test/utils/Utils.sol";
import {IdentityUpdates} from "src/IdentityUpdates.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {Initializable} from "@openzeppelin-contracts-upgradeable/proxy/utils/Initializable.sol";
import {PausableUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

contract IdentityUpdatesTest is Test, IdentityUpdates, Utils {
    IdentityUpdates identityUpdatesImpl;
    ERC1967Proxy proxy;
    IdentityUpdates identityUpdates;

    address admin = address(this);
    address unauthorized = address(0x1);

    function setUp() public {
        identityUpdatesImpl = new IdentityUpdates();

        proxy = new ERC1967Proxy(
            address(identityUpdatesImpl), abi.encodeWithSelector(identityUpdates.initialize.selector, admin)
        );

        identityUpdates = IdentityUpdates(address(proxy));
    }

    function testAddIdentityUpdateValid() public {
        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function testAddIdentityUpdateWithMaxPayload() public {
        bytes memory message = _generatePayload(identityUpdates.maxPayloadSize());

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function testAddIdentityUpdateTooSmall() public {
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

    function testAddIdentityUpdateTooBig() public {
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

    function testAddIdentityUpdateWhenPaused() public {
        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        identityUpdates.pause();
        assertTrue(identityUpdates.paused());

        vm.expectRevert(abi.encodeWithSelector(PausableUpgradeable.EnforcedPause.selector));

        identityUpdates.addIdentityUpdate(ID, message);
    }

    function testSetMinPayloadSize() public {
        // Store initial min payload size
        uint256 initialMinSize = identityUpdates.minPayloadSize();
        uint256 newMinSize = initialMinSize + 1;

        // Test unauthorized access
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        identityUpdates.setMinPayloadSize(newMinSize);

        // Test authorized access
        identityUpdates.setMinPayloadSize(newMinSize);
        assertEq(identityUpdates.minPayloadSize(), newMinSize);

        // Verify that messages with old minPayloadSize now fail
        bytes memory message = _generatePayload(initialMinSize);
        vm.expectRevert(
            abi.encodeWithSelector(
                IdentityUpdates.InvalidPayloadSize.selector,
                message.length,
                newMinSize,
                identityUpdates.maxPayloadSize()
            )
        );
        identityUpdates.addIdentityUpdate(ID, message);

        // Verify that messages with new minPayloadSize succeed
        bytes memory validMessage = _generatePayload(newMinSize);
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, validMessage, 1);
        identityUpdates.addIdentityUpdate(ID, validMessage);

        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMinPayloadSize.selector));
        identityUpdates.setMinPayloadSize(0);

        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMinPayloadSize.selector));
        identityUpdates.setMinPayloadSize(4194304);
    }

    function testSetMaxPayloadSize() public {
        // Store initial max payload size
        uint256 initialMaxSize = identityUpdates.maxPayloadSize();
        uint256 newMaxSize = 1000;

        // Test unauthorized access
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        identityUpdates.setMaxPayloadSize(newMaxSize);

        // Test authorized access
        identityUpdates.setMaxPayloadSize(newMaxSize);
        assertEq(identityUpdates.maxPayloadSize(), newMaxSize);

        // Verify that messages with old maxPayloadSize now fail
        bytes memory message = _generatePayload(initialMaxSize);
        vm.expectRevert(
            abi.encodeWithSelector(
                IdentityUpdates.InvalidPayloadSize.selector,
                message.length,
                identityUpdates.minPayloadSize(),
                newMaxSize
            )
        );
        identityUpdates.addIdentityUpdate(ID, message);

        // Verify that messages with new maxPayloadSize succeed
        bytes memory validMessage = _generatePayload(newMaxSize);
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, validMessage, 1);
        identityUpdates.addIdentityUpdate(ID, validMessage);

        // Max size should always be greater than min size
        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMaxPayloadSize.selector));
        identityUpdates.setMaxPayloadSize(78);

        // Test setting max size above maxPayloadSize (should fail)
        vm.expectRevert(abi.encodeWithSelector(IdentityUpdates.InvalidMaxPayloadSize.selector));
        identityUpdates.setMaxPayloadSize(4_194_305);
    }

    function testSequenceIdIncrement() public {
        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);
        identityUpdates.addIdentityUpdate(ID, message);

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 2);
        identityUpdates.addIdentityUpdate(ID, message);

        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 3);
        identityUpdates.addIdentityUpdate(ID, message);
    }

    function testInvalidReinitialization() public {
        vm.expectRevert(Initializable.InvalidInitialization.selector);
        identityUpdates.initialize(admin);
    }

    function testPauseUnpause() public {
        identityUpdates.pause();
        assertTrue(identityUpdates.paused());

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        identityUpdates.unpause();

        identityUpdates.unpause();
        assertFalse(identityUpdates.paused());

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        identityUpdates.pause();
    }

    function testRoles() public {
        identityUpdates.grantRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.startPrank(unauthorized);
        identityUpdates.pause();
        identityUpdates.unpause();
        vm.stopPrank();

        identityUpdates.revokeRole(DEFAULT_ADMIN_ROLE, unauthorized);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        identityUpdates.pause();

        identityUpdates.renounceRole(DEFAULT_ADMIN_ROLE, admin);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, admin, identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        identityUpdates.pause();
    }

    function testUpgradeImplementation() public {
        IdentityUpdates newIdentityUpdatesImpl = new IdentityUpdates();
        address newImplAddress = address(newIdentityUpdatesImpl);
        address oldImplAddress = address(identityUpdatesImpl);

        bytes memory message = _generatePayload(identityUpdates.minPayloadSize());

        // Retrieve the implementation address directly from the proxy storage.
        bytes32 rawImplAddress = vm.load(address(identityUpdates), EIP1967_IMPL_SLOT);
        address implementationAddress = address(uint160(uint256(rawImplAddress)));
        assertEq(implementationAddress, oldImplAddress);

        // Initialize sequenceId to 1. The state should be preserved between upgrades.
        vm.expectEmit(address(identityUpdates));
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 1);
        identityUpdates.addIdentityUpdate(ID, message);

        // Unauthorized upgrade attempts should revert.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                identityUpdates.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
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
        emit IdentityUpdates.IdentityUpdateCreated(ID, message, 2);
        identityUpdates.addIdentityUpdate(ID, message);
    }
}
