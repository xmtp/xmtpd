// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { IERC165 } from "@openzeppelin/contracts/utils/introspection/IERC165.sol";

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import { AccessControlDefaultAdminRules } from
    "@openzeppelin/contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import { INodes } from "./interfaces/INodes.sol";

// TODO: `NodeTransferred` event is redundant to `IERC721.Transfer` event.
// TODO: `MaxActiveNodesBelowCurrentCount` should be split across sets for better error handling.

/**
 * @title XMTP Nodes Registry
 *
 * @notice This contract is responsible for minting NFTs and assigning them to node operators.
 * Each node is minted as an NFT with a unique ID (starting at 100 and increasing by 100 with each new node).
 * In addition to the standard ERC721 functionality, the contract supports node-specific features,
 * including node property updates.
 *
 * @dev All nodes on the network periodically check this contract to determine which nodes they should connect to.
 * The contract owner is responsible for:
 *   - minting and transferring NFTs to node operators.
 *   - updating the node operator's HTTP address and MTLS certificate.
 *   - updating the node operator's minimum monthly fee.
 *   - updating the node operator's API enabled flag.
 */
contract Nodes is INodes, AccessControlDefaultAdminRules, ERC721 {
    using EnumerableSet for EnumerableSet.UintSet;

    /// @inheritdoc INodes
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    /// @inheritdoc INodes
    bytes32 public constant NODE_MANAGER_ROLE = keccak256("NODE_MANAGER_ROLE");

    /// @inheritdoc INodes
    uint256 public constant MAX_BPS = 10_000;

    /// @inheritdoc INodes
    uint32 public constant NODE_INCREMENT = 100;

    uint48 internal constant INITIAL_ACCESS_CONTROL_DELAY = 2 days;

    bytes1 internal constant FORWARD_SLASH = 0x2f;

    /// @dev The base URI for the node NFTs.
    string internal _baseTokenURI;

    /// @inheritdoc INodes
    uint8 public maxActiveNodes = 20;

    /**
     * @dev The counter for n max IDs.
     * The ERC721 standard expects the tokenID to be uint256 for standard methods unfortunately.
     */
    uint32 internal _nodeCounter = 0;

    /// @dev Mapping of token ID to Node.
    mapping(uint256 => Node) internal _nodes;

    /// @dev Nodes with API enabled.
    EnumerableSet.UintSet internal _activeApiNodes;

    /// @dev Nodes with replication enabled.
    EnumerableSet.UintSet internal _activeReplicationNodes;

    /// @inheritdoc INodes
    uint256 public nodeOperatorCommissionPercent;

    constructor(address initialAdmin)
        ERC721("XMTP Node Operator", "XMTP")
        AccessControlDefaultAdminRules(INITIAL_ACCESS_CONTROL_DELAY, initialAdmin)
    {
        require(initialAdmin != address(0), InvalidAddress());

        _setRoleAdmin(ADMIN_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(NODE_MANAGER_ROLE, DEFAULT_ADMIN_ROLE);

        // slither-disable-next-line unused-return
        _grantRole(ADMIN_ROLE, initialAdmin); // Will return false if the role is already granted.

        // slither-disable-next-line unused-return
        _grantRole(NODE_MANAGER_ROLE, initialAdmin); // Will return false if the role is already granted.
    }

    /* ============ Admin-Only Functions ============ */

    /// @inheritdoc INodes
    function addNode(
        address to,
        bytes calldata signingKeyPub,
        string calldata httpAddress,
        uint256 minMonthlyFeeMicroDollars
    ) external onlyRole(ADMIN_ROLE) returns (uint256 nodeId) {
        require(to != address(0), InvalidAddress());
        require(signingKeyPub.length > 0, InvalidSigningKey());
        require(bytes(httpAddress).length > 0, InvalidHttpAddress());

        nodeId = ++_nodeCounter * NODE_INCREMENT; // The first node starts with `nodeId = NODE_INCREMENT`.

        _nodes[nodeId] = Node(signingKeyPub, httpAddress, false, false, false, minMonthlyFeeMicroDollars);

        _mint(to, nodeId);

        emit NodeAdded(nodeId, to, signingKeyPub, httpAddress, minMonthlyFeeMicroDollars);
    }

    /// @inheritdoc INodes
    function disableNode(uint256 nodeId) public onlyRole(ADMIN_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);

        _nodes[nodeId].isDisabled = true;

        // Always remove from active nodes sets when disabled.
        _disableApiNode(nodeId);
        _disableReplicationNode(nodeId);

        emit NodeDisabled(nodeId);
    }

    /// @inheritdoc INodes
    function removeFromApiNodes(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);
        _disableApiNode(nodeId);
    }

    /// @inheritdoc INodes
    function removeFromReplicationNodes(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);
        _disableReplicationNode(nodeId);
    }

    /// @inheritdoc INodes
    function enableNode(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);

        // Re-enabling a node just removes the disabled flag.
        // The rest of the node properties are managed by the node operator.
        _nodes[nodeId].isDisabled = false;

        emit NodeEnabled(nodeId);
    }

    /// @inheritdoc INodes
    function setMaxActiveNodes(uint8 newMaxActiveNodes) external onlyRole(ADMIN_ROLE) {
        if (newMaxActiveNodes < _activeApiNodes.length() || newMaxActiveNodes < _activeReplicationNodes.length()) {
            revert MaxActiveNodesBelowCurrentCount();
        }

        maxActiveNodes = newMaxActiveNodes;
        emit MaxActiveNodesUpdated(newMaxActiveNodes);
    }

    /// @inheritdoc INodes
    function setNodeOperatorCommissionPercent(uint256 newCommissionPercent) external onlyRole(ADMIN_ROLE) {
        require(newCommissionPercent <= MAX_BPS, InvalidCommissionPercent());
        nodeOperatorCommissionPercent = newCommissionPercent;
        emit NodeOperatorCommissionPercentUpdated(newCommissionPercent);
    }

    /// @inheritdoc INodes
    function setBaseURI(string calldata newBaseURI) external onlyRole(ADMIN_ROLE) {
        require(bytes(newBaseURI).length > 0, InvalidURI());
        require(bytes(newBaseURI)[bytes(newBaseURI).length - 1] == FORWARD_SLASH, InvalidURI());
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }

    /* ============ Node Manager Functions ============ */

    /// @inheritdoc INodes
    function transferFrom(address from, address to, uint256 nodeId)
        public
        override(INodes, ERC721)
        onlyRole(NODE_MANAGER_ROLE)
    {
        /// @dev Disable the node before transferring ownership.
        /// It's NOP responsibility to re-enable the node after transfer.
        _disableApiNode(nodeId);
        _disableReplicationNode(nodeId);
        super.transferFrom(from, to, nodeId);
        emit NodeTransferred(nodeId, from, to);
    }

    /// @inheritdoc INodes
    function setHttpAddress(uint256 nodeId, string calldata httpAddress) external onlyRole(NODE_MANAGER_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);
        require(bytes(httpAddress).length > 0, InvalidHttpAddress());
        _nodes[nodeId].httpAddress = httpAddress;
        emit HttpAddressUpdated(nodeId, httpAddress);
    }

    /// @inheritdoc INodes
    function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFeeMicroDollars) external onlyRole(NODE_MANAGER_ROLE) {
        _revertIfNodeDoesNotExist(nodeId);
        _nodes[nodeId].minMonthlyFeeMicroDollars = minMonthlyFeeMicroDollars;
        emit MinMonthlyFeeUpdated(nodeId, minMonthlyFeeMicroDollars);
    }

    /* ============ Node Owner Functions ============ */

    /// @inheritdoc INodes
    function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) external {
        _revertIfNodeDoesNotExist(nodeId);
        _revertIfNodeIsDisabled(nodeId);
        _revertIfCallerIsNotOwner(nodeId);

        if (isApiEnabled) {
            _activateApiNode(nodeId);
        } else {
            _disableApiNode(nodeId);
        }
    }

    /// @inheritdoc INodes
    function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) external {
        _revertIfNodeDoesNotExist(nodeId);
        _revertIfNodeIsDisabled(nodeId);
        _revertIfCallerIsNotOwner(nodeId);

        if (isReplicationEnabled) {
            _activateReplicationNode(nodeId);
        } else {
            _disableReplicationNode(nodeId);
        }
    }

    /// @inheritdoc INodes
    function getAllNodes() public view returns (NodeWithId[] memory allNodes) {
        allNodes = new NodeWithId[](_nodeCounter);

        for (uint32 i; i < _nodeCounter; ++i) {
            uint32 nodeId = NODE_INCREMENT * (i + 1);

            allNodes[i] = NodeWithId({ nodeId: nodeId, node: _nodes[nodeId] });
        }
    }

    /// @inheritdoc INodes
    function getAllNodesCount() public view returns (uint256 nodeCount) {
        return _nodeCounter;
    }

    /// @inheritdoc INodes
    function getNode(uint256 nodeId) public view returns (Node memory node) {
        _revertIfNodeDoesNotExist(nodeId);
        return _nodes[nodeId];
    }

    /// @inheritdoc INodes
    function getActiveApiNodes() external view returns (NodeWithId[] memory activeNodes) {
        activeNodes = new NodeWithId[](_activeApiNodes.length());

        for (uint32 i; i < _activeApiNodes.length(); ++i) {
            uint256 nodeId = _activeApiNodes.at(i);

            activeNodes[i] = NodeWithId({ nodeId: nodeId, node: _nodes[nodeId] });
        }
    }

    /// @inheritdoc INodes
    function getActiveApiNodesIDs() external view returns (uint256[] memory activeNodesIDs) {
        return _activeApiNodes.values();
    }

    /// @inheritdoc INodes
    function getActiveApiNodesCount() external view returns (uint256 activeNodesCount) {
        return _activeApiNodes.length();
    }

    /// @inheritdoc INodes
    function getApiNodeIsActive(uint256 nodeId) external view returns (bool isActive) {
        return _activeApiNodes.contains(nodeId);
    }

    /// @inheritdoc INodes
    function getActiveReplicationNodes() external view returns (NodeWithId[] memory activeNodes) {
        activeNodes = new NodeWithId[](_activeReplicationNodes.length());

        for (uint32 i; i < _activeReplicationNodes.length(); ++i) {
            uint256 nodeId = _activeReplicationNodes.at(i);

            activeNodes[i] = NodeWithId({ nodeId: nodeId, node: _nodes[nodeId] });
        }
    }

    /// @inheritdoc INodes
    function getActiveReplicationNodesIDs() external view returns (uint256[] memory activeNodesIDs) {
        return _activeReplicationNodes.values();
    }

    /// @inheritdoc INodes
    function getActiveReplicationNodesCount() external view returns (uint256 activeNodesCount) {
        return _activeReplicationNodes.length();
    }

    /// @inheritdoc INodes
    function getReplicationNodeIsActive(uint256 nodeId) external view returns (bool isActive) {
        return _activeReplicationNodes.contains(nodeId);
    }

    /// @inheritdoc INodes
    function getNodeOperatorCommissionPercent() external view returns (uint256 commissionPercent) {
        return nodeOperatorCommissionPercent;
    }

    /* ============ Internal Functions ============ */

    /**
     * @dev    Checks if a node exists.
     * @param  nodeId The ID of the node to check.
     * @return exists True if the node exists, false otherwise.
     */
    function _nodeExists(uint256 nodeId) internal view returns (bool exists) {
        return _ownerOf(nodeId) != address(0);
    }

    /// @inheritdoc ERC721
    function _baseURI() internal view virtual override returns (string memory baseURI) {
        return _baseTokenURI;
    }

    /// @dev Helper function to add a node to the active API nodes set.
    function _activateApiNode(uint256 nodeId) internal {
        require(_activeApiNodes.length() < maxActiveNodes, MaxActiveNodesReached());

        if (!_activeApiNodes.add(nodeId)) return;

        _nodes[nodeId].isApiEnabled = true;

        emit ApiEnabled(nodeId);
    }

    /// @dev Helper function to remove a node from the active API nodes set.
    function _disableApiNode(uint256 nodeId) internal {
        if (!_activeApiNodes.remove(nodeId)) return;

        _nodes[nodeId].isApiEnabled = false;

        emit ApiDisabled(nodeId);
    }

    /// @dev Helper function to add a node to the active replication nodes set.
    function _activateReplicationNode(uint256 nodeId) internal {
        require(_activeReplicationNodes.length() < maxActiveNodes, MaxActiveNodesReached());

        if (!_activeReplicationNodes.add(nodeId)) return;

        _nodes[nodeId].isReplicationEnabled = true;

        emit ReplicationEnabled(nodeId);
    }

    /// @dev Helper function to remove a node from the active replication nodes set.
    function _disableReplicationNode(uint256 nodeId) internal {
        if (!_activeReplicationNodes.remove(nodeId)) return;

        _nodes[nodeId].isReplicationEnabled = false;

        emit ReplicationDisabled(nodeId);
    }

    /// @dev Override to support INodes, ERC721, IERC165, and AccessControlEnumerable.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, IERC165, AccessControlDefaultAdminRules)
        returns (bool supported)
    {
        return interfaceId == type(INodes).interfaceId || super.supportsInterface(interfaceId);
    }

    /// @dev Reverts if the node does not exist.
    function _revertIfNodeDoesNotExist(uint256 nodeId) internal view {
        require(_nodeExists(nodeId), NodeDoesNotExist());
    }

    /// @dev Reverts if the node is disabled.
    function _revertIfNodeIsDisabled(uint256 nodeId) internal view {
        require(!_nodes[nodeId].isDisabled, NodeIsDisabled());
    }

    /// @dev Reverts if `msg.sender` is not the owner of the node.
    function _revertIfCallerIsNotOwner(uint256 nodeId) internal view {
        require(_ownerOf(nodeId) == msg.sender, Unauthorized());
    }
}
