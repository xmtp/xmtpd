// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "@openzeppelin/contracts/access/Ownable.sol";

contract Nodes is Ownable {
    constructor() Ownable(msg.sender) {}

    struct Node {
        string httpAddress;
        uint256 originatorId;
        bool isHealthy;
        // Maybe we want a TLS cert separate from the public key for MTLS authenticated connections?
    }

    event NodeUpdate(
        bytes publicKey,
        string httpAddress,
        uint256 originatorId,
        bool isHealthy
    );

    // List of public keys
    bytes[] publicKeys;

    // Mapping of publicKey to node
    mapping(bytes => Node) public nodes;

    /**
    Add a node to the network
     */
    function addNode(
        bytes calldata publicKey,
        string calldata httpAddress
    ) public onlyOwner {
        require(
            bytes(nodes[publicKey].httpAddress).length == 0,
            "Node already exists"
        );

        require(bytes(httpAddress).length != 0, "HTTP address is required");

        nodes[publicKey] = Node({
            httpAddress: httpAddress,
            originatorId: publicKeys.length + 1,
            isHealthy: true
        });

        publicKeys.push(publicKey);

        emit NodeUpdate(publicKey, httpAddress, publicKeys.length, true);
    }

    /**
    The contract owner can use this function to mark a node as unhealthy
    triggering all other nodes to stop replicating to/from this node
     */
    function markNodeUnhealthy(bytes calldata publicKey) public onlyOwner {
        require(
            bytes(nodes[publicKey].httpAddress).length != 0,
            "Node does not exist"
        );
        nodes[publicKey].isHealthy = false;

        emit NodeUpdate(
            publicKey,
            nodes[publicKey].httpAddress,
            nodes[publicKey].originatorId,
            false
        );
    }

    /**
    The contract owner can use this function to mark a node as healthy
    triggering all other nodes to 
     */
    function markNodeHealthy(bytes calldata publicKey) public onlyOwner {
        require(
            bytes(nodes[publicKey].httpAddress).length != 0,
            "Node does not exist"
        );
        nodes[publicKey].isHealthy = true;

        emit NodeUpdate(
            publicKey,
            nodes[publicKey].httpAddress,
            nodes[publicKey].originatorId,
            true
        );
    }
}
