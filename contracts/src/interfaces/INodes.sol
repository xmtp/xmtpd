// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

// Import the standard ERC721 interface.
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/// @title INodesErrors
/// @notice This interface defines the errors emitted by the INodes contract.
interface INodesErrors {
    /// @notice Error thrown when a node is disabled.
    error NodeIsDisabled();

    /// @notice Error thrown when an invalid address is provided.
    error InvalidAddress();

    /// @notice Error thrown when an invalid commission percentage is provided.
    error InvalidCommissionPercent();

    /// @notice Error thrown when an invalid HTTP address is provided.
    error InvalidHttpAddress();

    /// @notice Error thrown when the input length is invalid.
    error InvalidInputLength();

    /// @notice Error thrown when a node config is invalid.
    error InvalidNodeConfig();

    /// @notice Error thrown when an invalid signing key is provided.
    error InvalidSigningKey();

    /// @notice Error thrown when an invalid URI is provided.
    error InvalidURI();

    /// @notice Error thrown when the maximum number of active nodes is reached.
    error MaxActiveNodesReached();

    /// @notice Error when trying to set max active nodes below current active count.
    error MaxActiveNodesBelowCurrentCount();

    /// @notice Error thrown when a node does not exist.
    error NodeDoesNotExist();

    /// @notice Error thrown when an unauthorized address attempts to call a function.
    error Unauthorized();
}

/// @title INodesEvents
/// @notice This interface defines the events emitted by the INodes contract.
interface INodesEvents {
    /// @notice Emitted when a new node is added and its NFT minted.
    /// @param nodeId The unique identifier for the node (starts at 100, increments by 100).
    /// @param owner The address that receives the new node NFT.
    /// @param signingKeyPub The node’s signing key public value.
    /// @param httpAddress The node’s HTTP endpoint.
    /// @param minMonthlyFeeMicroDollars The minimum monthly fee for the node.
    event NodeAdded(
        uint256 indexed nodeId,
        address indexed owner,
        bytes signingKeyPub,
        string httpAddress,
        uint256 minMonthlyFeeMicroDollars
    );

    /// @notice Emitted when a disabled status is removed from a node.
    /// @param nodeId The identifier of the node.
    event NodeEnabled(uint256 indexed nodeId);

    /// @notice Emitted when a node is disabled by an administrator.
    /// @param nodeId The identifier of the node.
    event NodeDisabled(uint256 indexed nodeId);

    /// @notice Emitted when a node is transferred
    event NodeTransferred(uint256 indexed nodeId, address indexed from, address indexed to);

    /// @notice Emitted when the HTTP address for a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param newHttpAddress The new HTTP address.
    event HttpAddressUpdated(
        uint256 indexed nodeId,
        string newHttpAddress
    );

    /// @notice Emitted when the replication flag for a node is enabled.
    /// @param nodeId The identifier of the node.
    event ReplicationEnabled(uint256 indexed nodeId);

    /// @notice Emitted when the replication flag for a node is disabled.
    /// @param nodeId The identifier of the node.
    event ReplicationDisabled(uint256 indexed nodeId);

    /// @notice Emitted when the API enabled flag for a node is enabled.
    /// @param nodeId The identifier of the node.
    event ApiEnabled(uint256 indexed nodeId);

    /// @notice Emitted when the API enabled flag for a node is disabled.
    /// @param nodeId The identifier of the node.
    event ApiDisabled(uint256 indexed nodeId);

    /// @notice Emitted when the minimum monthly fee for a node is updated.
    /// @param nodeId The identifier of the node.
    /// @param minMonthlyFeeMicroDollars The updated minimum fee.
    event MinMonthlyFeeUpdated(
        uint256 indexed nodeId,
        uint256 minMonthlyFeeMicroDollars
    );

    /// @notice Emitted when the node operator commission percent is updated.
    /// @param newCommissionPercent The new commission percentage.
    event NodeOperatorCommissionPercentUpdated(
        uint256 newCommissionPercent
    );

    /// @notice Emitted when the maximum number of active nodes is updated.
    /// @param newMaxActiveNodes The new maximum number of active nodes.
    event MaxActiveNodesUpdated(uint8 newMaxActiveNodes);

    /// @notice Emitted when the base URI is updated.
    /// @param newBaseURI The new base URI.
    event BaseURIUpdated(string newBaseURI);
}

/// @title INodes
/// @notice This interface defines the ERC721-based registry for “nodes” in the system.
/// Each node is minted as an NFT with a unique ID (starting at 100 and increasing by 100 with each new node).
/// In addition to the standard ERC721 functionality, the contract supports node-specific features,
/// including node property updates.
interface INodes is IERC721, INodesErrors, INodesEvents {
    /// @notice Struct representing a node in the registry.
    /// @param signingKeyPub The public key used for node signing/verification.
    /// @param httpAddress The HTTP endpoint address for the node.
    /// @param isReplicationEnabled A flag indicating whether the node supports replication.
    /// @param isApiEnabled A flag indicating whether the node has its API enabled.
    /// @param isDisabled A flag indicating whether the node was disabled by an administrator.
    /// @param minMonthlyFeeMicroDollars The minimum monthly fee collected by the node operator.
    struct Node {
        bytes signingKeyPub;
        string httpAddress;
        bool isReplicationEnabled;
        bool isApiEnabled;
        bool isDisabled;
        uint256 minMonthlyFeeMicroDollars;
    }

    /// @notice Struct representing a node with its ID
    /// @param nodeId The unique identifier for the node
    /// @param node The node struct
    struct NodeWithId {
        uint256 nodeId;
        Node node;
    }

    // ***************************************************************
    // *                ADMIN-ONLY FUNCTIONS                         *
    // ***************************************************************

    /// @notice Adds a new node to the registry and mints its corresponding ERC721 token.
    /// @dev Only the contract owner may call this. Node IDs start at 100 and increase by 100 for each new node.
    /// @param to The address that will own the new node NFT.
    /// @param signingKeyPub The public signing key for the node.
    /// @param httpAddress The node’s HTTP address.
    /// @param minMonthlyFeeMicroDollars The minimum monthly fee that the node operator collects.
    /// @return nodeId The unique identifier of the newly added node.
    function addNode(
        address to,
        bytes calldata signingKeyPub,
        string calldata httpAddress,
        uint256 minMonthlyFeeMicroDollars
    ) external returns (uint256 nodeId);

    /// @notice Disables a node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    function disableNode(uint256 nodeId) external;

    /// @notice Removes a node from the active API nodes set.
    /// @dev Only the contract owner may call this.
    /// enableNode sets isDisabled to false, it does not activate the node.
    /// The node must be activated separately.
    /// @param nodeId The unique identifier of the node.
    function enableNode(uint256 nodeId) external;

    /// @notice Set the HTTP address of an existing node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param httpAddress The new HTTP address.
    function setHttpAddress(uint256 nodeId, string calldata httpAddress) external;

    /// @notice Set the minimum monthly fee for a node.
    /// @dev Only the contract owner may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param minMonthlyFeeMicroDollars The new minimum monthly fee.
    function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) external;

    /// @notice Sets the commission percentage that the node operator receives.
    /// @dev Only the contract owner may call this.
    /// @param newCommissionPercent The new commission percentage.
    function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) external;

    /// @notice Sets the maximum number of active nodes.
    /// @dev Only the contract owner may call this.
    /// @param newMaxActiveNodes The new maximum number of active nodes.
    function setMaxActiveNodes(uint8 newMaxActiveNodes) external;

    /// @notice Set the base URI for the node NFTs.
    /// @dev Only the contract owner may call this.
    /// @param newBaseURI The new base URI. Has to end with a trailing slash.
    function setBaseURI(string calldata newBaseURI) external;

    // ***************************************************************
    // *                NODE OWNER FUNCTION                        *
    // ***************************************************************

    /// @notice Sets the API enabled flag for the node owned by the caller.
    /// @dev Only the owner of the node NFT may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param isApiEnabled The new API enabled flag.
    function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) external;

    /// @notice Sets the replication enabled flag for the node owned by the caller.
    /// @dev Only the owner of the node NFT may call this.
    /// @param nodeId The unique identifier of the node.
    /// @param isReplicationEnabled A boolean indicating if replication should be enabled.
    function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) external;

    // ***************************************************************
    // *                     GETTER FUNCTIONS                      *
    // ***************************************************************

    /// @notice Retrieves the current node operator commission percentage.
    /// @return commissionPercent The commission percentage.
    function getNodeOperatorCommissionPercent() external view returns (uint256 commissionPercent);

    /// @notice Gets all nodes regardless of their health status
    /// @return allNodes An array of all nodes in the registry
    function getAllNodes() external view returns (NodeWithId[] memory allNodes);

    /// @notice Gets the total number of nodes in the registry.
    /// @return nodeCount The total number of nodes.
    function getAllNodesCount() external view returns (uint256 nodeCount);

    /// @notice Retrieves the details of a given node.
    /// @param nodeId The unique identifier of the node.
    /// @return node The Node struct containing the node's details.
    function getNode(uint256 nodeId) external view returns (Node memory node);

    /// @notice Retrieves a list of active API nodes.
    /// @dev Active nodes are those with `isActive` set to true.
    /// @return activeNodes An array of Node structs representing active nodes.
    function getActiveApiNodes() external view returns (NodeWithId[] memory activeNodes);

    /// @notice Retrieves a list of active API nodes IDs.
    /// @return activeNodesIDs An array of node IDs representing active nodes.
    function getActiveApiNodesIDs() external view returns (uint256[] memory activeNodesIDs);

    /// @notice Retrieves the total number of active API nodes.
    /// @return activeNodesCount The total number of active API nodes.
    function getActiveApiNodesCount() external view returns (uint256 activeNodesCount);

    /// @notice Retrieves if a node API is active.
    /// @param nodeId The ID of the node NFT.
    /// @return isActive A boolean indicating if the node is active.
    function getApiNodeIsActive(uint256 nodeId) external view returns (bool isActive);

    /// @notice Retrieves a list of active replication nodes.
    /// @dev Active nodes are those with `isActive` set to true.
    /// @return activeNodes An array of Node structs representing active nodes.
    function getActiveReplicationNodes() external view returns (NodeWithId[] memory activeNodes);

    /// @notice Retrieves a list of active replication nodes IDs.
    /// @return activeNodesIDs An array of node IDs representing active nodes.
    function getActiveReplicationNodesIDs() external view returns (uint256[] memory activeNodesIDs);

    /// @notice Retrieves the total number of active replication nodes.
    /// @return activeNodesCount The total number of active replication nodes.
    function getActiveReplicationNodesCount() external view returns (uint256 activeNodesCount);

    /// @notice Retrieves if a node replication is active.
    /// @param nodeId The ID of the node NFT.
    /// @return isActive A boolean indicating if the node is active.
    function getReplicationNodeIsActive(uint256 nodeId) external view returns (bool isActive);
}
