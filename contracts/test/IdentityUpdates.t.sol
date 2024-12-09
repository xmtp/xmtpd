// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import {Test, console} from "forge-std-1.9.4/src/Test.sol";
import {IdentityUpdates} from "../src/IdentityUpdates.sol";

contract IdentityUpdatesTest is Test {
    IdentityUpdates public identityUpdates;

    function setUp() public {
        identityUpdates = new IdentityUpdates();
    }

    function test_AddIdentityUpdate1k() public {
        bytes32 inboxId = bytes32(0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef);
        bytes memory message = new bytes(1024);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256)); // Set each byte to its index modulo 256
        }

        identityUpdates.addIdentityUpdate(inboxId, message);
    }
}
