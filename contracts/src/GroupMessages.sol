// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract GroupMessages {
    event MessageSent(bytes32 groupId, bytes message, uint64 sequenceId);

    uint64 sequenceId;

    function addMessage(bytes32 groupId, bytes memory message) public {
        sequenceId++;

        emit MessageSent(groupId, message, sequenceId);
    }
}
