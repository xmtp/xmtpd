// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

// Import the standard ERC721 interface.
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/// @title INodes
/// @notice This interface defines the ERC721-based registry for “nodes” in the system.
/// Each node is minted as an NFT with a unique ID (starting at 100 and increasing by 100 with each new node).
/// In addition to the standard ERC721 functionality, the contract supports node-specific features,
/// including node property updates.
///
/// Note: The staking contract reference has been removed per updated requirements.
interface INodes is IERC721 {

    /// @notice Error thrown when an unauthorized address attempts to call a function.
    error Unauthorized();

    /// @notice Error thrown when an invalid address is provided.
    error InvalidAddress();

    /// @notice Error thrown when an invalid signing key is provided.
    error InvalidSigningKey();

    /// @notice Error thrown when an invalid HTTP address is provided.
    error InvalidHttpAddress();

    /// @notice Error thrown when a node does not exist.
    error NodeDoesNotExist();

    /// @notice Error thrown when a node is already active.
    error NodeAlreadyActive();

    /// @notice Error thrown when a node is already inactive.
    error NodeAlreadyInactive();

    /// @notice Error thrown when an invalid commission percentage is provided.
    error InvalidCommissionPercent();

    /// @notice Error thrown when the maximum number of active nodes is reached.
    error MaxActiveNodesReached();

    /// @notice Error thrown when an invalid URI is provided.
    error InvalidURI();

    /// @notice Struct representing a node in the registry.
    /// @param signingKeyPub The public key used for node signing/verification.
    /// @param httpAddress The HTTP endpoint address for the node.
    /// @param isReplicationEnabled A flag indicating whether the node supports replication.
    /// @param isApiEnabled A flag indicating whether the node has its API enabled.
    /// @param isActive A flag indicating whether the node is actively participating in the network.
    /// @param minMonthlyFee The minimum monthly fee collected by the node operator.
    struct Node {
        bytes signingKeyPub;
        string httpAddress;
        bool isReplicationEnabled;
        bool isApiEnabled;
        bool isActive;
        uint256 minMonthlyFee;
    }

    /// @notice Struct representing a node with its ID
    /// @param nodeId The unique identifier for the node
    /// @param node The node struct
    struct NodeWithId {
        uint256 nodeId;
        Node node;
    }

    // ***************************************************************
    // *                        EVENTS                             *
    // ***************************************************************

    /// @notice Emitted when a new node is added and its NFT minted.
    /// @param nodeId The unique identifier for the node (starts at 100, increments by 100).
    /// @param owner The address that receives the new node NFT.
    /// @param signingKeyPub The node’s signing key public value.
    /// @param httpAddress The node’s HTTP endpoint.
    /// @param minMonthlyFee The minimum monthly fee for the node.
    event NodeAdded(
        uint256 indexed nodeId,
        address indexed owner,
        bytes signingKeyPub,
        string httpAddress,
        uint256 minMonthlyFee
    );

    /// @notice Emitted when the HTTP address for a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param newHttpAddress The new HTTP address.
    event HttpAddressUpdated(
        uint256 indexed nodeId,
        string newHttpAddress
    );

    /// @notice Emitted when the replication flag for a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param isReplicationEnabled The updated replication flag.
    event ReplicationEnabledUpdated(
        uint256 indexed nodeId,
        bool isReplicationEnabled
    );

    /// @notice Emitted when the minimum monthly fee for a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param minMonthlyFee The updated minimum fee.
    event MinMonthlyFeeUpdated(
        uint256 indexed nodeId,
        uint256 minMonthlyFee
    );

    /// @notice Emitted when the API enabled flag for a node is toggled.
    /// @param nodeId The identifier of the node.
    /// @param isApiEnabled The updated API enabled flag.
    event ApiEnabledUpdated(
        uint256 indexed nodeId,
        bool isApiEnabled
    );

    /// @notice Emitted when the node operator commission percent is updated.
    /// @param newCommissionPercent The new commission percentage.
    event NodeOperatorCommissionPercentUpdated(
        uint256 newCommissionPercent
    );

    /// @notice Emitted when the active status of a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param isActive The updated active status.
    event NodeActivateUpdated(
        uint256 indexed nodeId,
        bool isActive);

    /// @notice Emitted when a node is transferred
    event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to);

    /// @notice Emitted when the maximum number of active nodes is updated.
    /// @param newMaxActiveNodes The new maximum number of active nodes.
    event MaxActiveNodesUpdated(uint8 newMaxActiveNodes);

    /// @notice Emitted when the base URI is updated.
    /// @param newBaseURI The new base URI.
    event BaseURIUpdated(string newBaseURI);

    // ***************************************************************
    // *                ADMIN-ONLY FUNCTIONS                         *
    // ***************************************************************

    /// @notice Adds a new node to the registry and mints its corresponding ERC721 token.
    /// @dev Only the contract owner may call this. Node IDs start at 100 and increase by 100 for each new node.
    /// @param to The address that will own the new node NFT.
    /// @param signingKeyPub The public signing key for the node.
    /// @param httpAddress The node’s HTTP address.
    /// @param minMonthlyFee The minimum monthly fee that the node operator collects.
    /// @return nodeId The unique identifier of the newly added node.
    function addNode(
        address to,
        bytes calldata signingKeyPub,
        string calldata httpAddress,
        uint256 minMonthlyFee
    ) external returns (uint256 nodeId);

    /// @notice Updates the HTTP address of an existing node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param httpAddress The new HTTP address.
    function updateHttpAddress(uint256 nodeId, string calldata httpAddress) external;

    /// @notice Updates whether replication is enabled for a node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param isReplicationEnabled A boolean indicating if replication should be enabled.
    function updateIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) external;

    /// @notice Updates the minimum monthly fee for a node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param minMonthlyFee The new minimum monthly fee.
    function updateMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) external;

    /// @notice Updates the commission percentage that the node operator receives.
    /// @dev Only the contract owner may call this.
    /// @param newCommissionPercent The new commission percentage.
    function updateNodeOperatorCommissionPercent(uint256 newCommissionPercent) external;

    /// @notice Updates the active status of a node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param isActive The new active status.
    function updateActive(uint256 nodeId, bool isActive) external;

    /// @notice Updates multiple nodes active status at once.
    /// @dev Only the contract owner may call this.
    /// @param nodeIds Array of node IDs.
    /// @param isActive Array of active status flags.
    function batchUpdateActive(uint256[] calldata nodeIds, bool[] calldata isActive) external;

    /// @notice Updates the maximum number of active nodes.
    /// @dev Only the contract owner may call this.
    /// @param newMaxActiveNodes The new maximum number of active nodes.
    function updateMaxActiveNodes(uint8 newMaxActiveNodes) external;

    /// @notice Updates the base URI for the node NFTs.
    /// @dev Only the contract owner may call this.
    /// @param newBaseURI The new base URI. Has to end with a trailing slash.
    function setBaseURI(string calldata newBaseURI) external;

    // ***************************************************************
    // *                NODE OWNER FUNCTION                        *
    // ***************************************************************

    /// @notice Toggles the API enabled flag for the node owned by the caller.
    /// @dev Only the owner of the node NFT may call this.
    /// @param nodeId The unique identifier of the node.
    function updateIsApiEnabled(uint256 nodeId) external;

    // ***************************************************************
    // *                     GETTER FUNCTIONS                      *
    // ***************************************************************

    /// @notice Gets all nodes regardless of their health status
    /// @return allNodesList An array of all nodes in the registry
    function allNodes() external view returns (NodeWithId[] memory);

    /// @notice Retrieves the details of a given node.
    /// @param nodeId The unique identifier of the node.
    /// @return The Node struct containing the node's details.
    function getNode(uint256 nodeId) external view returns (Node memory);

    /// @notice Retrieves the current node operator commission percentage.
    /// @return The commission percentage.
    function nodeOperatorCommissionPercent() external view returns (uint256);

    /// @notice Retrieves a list of active nodes.
    /// @dev Active nodes are those with `isActive` set to true.
    /// @return activeNodes An array of Node structs representing active nodes.
    function getActiveNodes() external view returns (Node[] memory activeNodes);

    /// @notice Retrieves a list of active nodes IDs.
    /// @dev Active nodes are those with `isActive` set to true.
    /// @return activeNodesIDs An array of node IDs representing active nodes.
    function getActiveNodesIDs() external view returns (uint256[] memory activeNodesIDs);

    /// @notice Retrieves if a node is active.
    /// @param nodeId The ID of the node NFT.
    /// @return isActive A boolean indicating if the node is active.
    function nodeIsActive(uint256 nodeId) external view returns (bool);
}