// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

/**
 * Interface for the XMTP Node Operators NFT contract.
 */
interface INodes {
    struct Node {
        bytes signingKeyPub;
        string httpAddress;
        bool isHealthy;
        bool isActive;
    }

    struct NodeWithId {
        uint32 nodeId;
        Node node;
    }

    event NodeCreated(uint256 nodeId, Node node);
    event NodeUpdated(uint256 nodeId, Node node);
    event NodeActivated(uint256 nodeId);
    event NodeDeactivated(uint256 nodeId);

    function addNode(address to, bytes calldata signingKeyPub, string calldata httpAddress) external returns (uint32);
    function transferFrom(address from, address to, uint256 tokenId) external;
    function updateHttpAddress(uint256 tokenId, string calldata httpAddress) external;
    function updateHealth(uint256 tokenId, bool isHealthy) external;
    function updateActive(uint256 tokenId, bool isActive) external;
    function healthyNodes() external view returns (NodeWithId[] memory);
    function allNodes() external view returns (NodeWithId[] memory);
    function getNode(uint256 tokenId) external view returns (Node memory);

    // Include required ERC721 interface functions
    function ownerOf(uint256 tokenId) external view returns (address);
} 