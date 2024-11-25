// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

contract IdentityUpdates {
    event IdentityUpdateCreated(bytes32 inboxId, bytes update, uint64 sequenceId);

    uint64 sequenceId;

    function addIdentityUpdate(bytes32 inboxId, bytes memory update) public {
        sequenceId++;

        emit IdentityUpdateCreated(inboxId, update, sequenceId);
    }
}
