// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console} from "forge-std/Test.sol";
import {GroupMessages} from "../src/GroupMessages.sol";

contract GroupMessagesTest is Test {
    GroupMessages public groupMessages;

    function setUp() public {
        groupMessages = new GroupMessages();
    }

    function test_AddMessage2kb() public {
        bytes32 groupId = bytes32(
            0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
        );
        bytes memory message = new bytes(1024);
        for (uint256 i = 0; i < message.length; i++) {
            message[i] = bytes1(uint8(i % 256)); // Set each byte to its index modulo 256
        }

        groupMessages.addMessage(groupId, message);
    }
}
