// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/extensions/AccessControlDefaultAdminRules.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./interfaces/INodes.sol";

 /// @title XMTP Nodes Registry.
 /// @notice This contract is responsible for minting NFTs and assigning them to node operators.
 /// Each node is minted as an NFT with a unique ID (starting at 100 and increasing by 100 with each new node).
 /// In addition to the standard ERC721 functionality, the contract supports node-specific features,
 /// including node property updates.
 /// 
 /// @dev All nodes on the network periodically check this contract to determine which nodes they should connect to.
 /// The contract owner is responsible for:
 /// - minting and transferring NFTs to node operators.
 /// - updating the node operator's HTTP address and MTLS certificate.
 /// - updating the node operator's minimum monthly fee.
 /// - updating the node operator's API enabled flag.
contract Nodes is AccessControlDefaultAdminRules, ERC721, INodes {
    using EnumerableSet for EnumerableSet.UintSet;

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant NODE_MANAGER_ROLE = keccak256("NODE_MANAGER_ROLE");

    /// @dev The maximum commission percentage that the node operator can receive.
    uint256 public constant MAX_BPS = 10000;

    /// @dev The increment for node IDs.
    uint32 private constant NODE_INCREMENT = 100;

    /// @dev The base URI for the node NFTs.
    string private _baseTokenURI;

    /// @dev Max number of active nodes.
    // slither-disable-next-line constable-states
    uint8 public maxActiveNodes = 20;

    /// @dev The counter for n max IDs.
    // The ERC721 standard expects the tokenID to be uint256 for standard methods unfortunately.
    // slither-disable-next-line constable-states
    uint32 private _nodeCounter = 0;

    /// @dev Mapping of token ID to Node.
    mapping(uint256 => Node) private _nodes;

    /// @dev Nodes with API enabled.
    EnumerableSet.UintSet private _activeApiNodes;

    /// @dev Nodes with replication enabled.
    EnumerableSet.UintSet private _activeReplicationNodes;

    /// @notice The commission percentage that the node operator receives.
    /// @dev This is stored in basis points (1/100th of a percent).
    /// Example: 1% = 100bps, 10% = 1000bps, 100% = 10000bps.
    /// Comission is calculated as (nodeOperatorCommissionPercent * nodeOperatorFee) / MAX_BPS.
    // slither-disable-next-line constable-states
    uint256 public nodeOperatorCommissionPercent;

    constructor(address _initialAdmin)
        ERC721("XMTP Node Operator", "XMTP")
        AccessControlDefaultAdminRules(2 days, _initialAdmin)
    {
        require(_initialAdmin != address(0), InvalidAddress());

        _setRoleAdmin(ADMIN_ROLE, DEFAULT_ADMIN_ROLE);
        _setRoleAdmin(NODE_MANAGER_ROLE, DEFAULT_ADMIN_ROLE);
        // slither-disable-next-line unused-return
        _grantRole(ADMIN_ROLE, _initialAdmin);
        // slither-disable-next-line unused-return
        _grantRole(NODE_MANAGER_ROLE, _initialAdmin);
    }

    // ***************************************************************
    // *                ADMIN-ONLY FUNCTIONS                         *
    // ***************************************************************

    /// @inheritdoc INodes
    function addNode(address to, bytes calldata signingKeyPub, string calldata httpAddress, uint256 minMonthlyFee)
        external
        onlyRole(ADMIN_ROLE)
        returns (uint256)
    {
        require(to != address(0), InvalidAddress());
        require(signingKeyPub.length > 0, InvalidSigningKey());
        require(bytes(httpAddress).length > 0, InvalidHttpAddress());

        // the first node starts with 100
        _nodeCounter++;
        uint32 nodeId = _nodeCounter * NODE_INCREMENT;
        _mint(to, nodeId);
        _nodes[nodeId] = Node(signingKeyPub, httpAddress, false, false, false, minMonthlyFee);
        emit NodeAdded(nodeId, to, signingKeyPub, httpAddress, minMonthlyFee);
        return nodeId;
    }

    /// @inheritdoc INodes
    function disableNode(uint256 nodeId) public onlyRole(ADMIN_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());

        _nodes[nodeId].isDisabled = true;

        // Always remove from active nodes sets when disabled
        _disableApiNode(nodeId);
        _disableReplicationNode(nodeId);

        emit NodeDisabled(nodeId);
    }

    function removeFromApiNodes(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        _disableApiNode(nodeId);
    }

    function removeFromReplicationNodes(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        _disableReplicationNode(nodeId);
    }

    /// @inheritdoc INodes
    function enableNode(uint256 nodeId) external onlyRole(ADMIN_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());

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
        require(bytes(newBaseURI)[bytes(newBaseURI).length - 1] == 0x2f, InvalidURI());
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }

    // ***************************************************************
    // *                NODE MANAGER FUNCTIONS                       *
    // ***************************************************************

    /// @notice Transfers node ownership from one address to another
    /// @dev Only the contract owner may call this. Automatically deactivates the node
    /// @param from The current owner address
    /// @param to The new owner address
    /// @param nodeId The ID of the node being transferred
    function transferFrom(address from, address to, uint256 nodeId) 
        public 
        override(ERC721, IERC721) 
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
        require(_nodeExists(nodeId), NodeDoesNotExist());
        require(bytes(httpAddress).length > 0, InvalidHttpAddress());
        _nodes[nodeId].httpAddress = httpAddress;
        emit HttpAddressUpdated(nodeId, httpAddress);
    }

    /// @inheritdoc INodes
    function setMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) external onlyRole(NODE_MANAGER_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        _nodes[nodeId].minMonthlyFee = minMonthlyFee;
        emit MinMonthlyFeeUpdated(nodeId, minMonthlyFee);
    }

    // ***************************************************************
    // *                NODE OWNER FUNCTION                          *
    // ***************************************************************

    /// @inheritdoc INodes
    function setIsApiEnabled(uint256 nodeId, bool isApiEnabled) external whenNotDisabled(nodeId) {
        require(_ownerOf(nodeId) == msg.sender, Unauthorized());
        if (isApiEnabled) {
            _activateApiNode(nodeId);
        } else {
            _disableApiNode(nodeId);
        }
    }

    /// @inheritdoc INodes
    function setIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) external whenNotDisabled(nodeId) {
        require(_ownerOf(nodeId) == msg.sender, Unauthorized());
        if (isReplicationEnabled) {
            _activateReplicationNode(nodeId);
        } else {
            _disableReplicationNode(nodeId);
        }
    }

    /// @inheritdoc INodes
    function getAllNodes() public view returns (NodeWithId[] memory allNodes) {
        allNodes = new NodeWithId[](_nodeCounter);
        for (uint32 i = 0; i < _nodeCounter; i++) {
            uint32 nodeId = NODE_INCREMENT * (i + 1);
            if (_nodeExists(nodeId)) {
                allNodes[i] = NodeWithId({nodeId: nodeId, node: _nodes[nodeId]});
            }
        }
        return allNodes;
    }

    /// @inheritdoc INodes
    function getAllNodesCount() public view returns (uint256 nodeCount) {
        return _nodeCounter;
    }

    /// @inheritdoc INodes
    function getNode(uint256 nodeId) public view returns (Node memory node) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        return _nodes[nodeId];
    }

    /// @inheritdoc INodes
    function getActiveApiNodes() external view returns (NodeWithId[] memory activeNodes) {
        activeNodes = new NodeWithId[](_activeApiNodes.length());
        for (uint32 i = 0; i < _activeApiNodes.length(); i++) {
            uint256 nodeId = _activeApiNodes.at(i);
            if (_nodeExists(nodeId)) {
                activeNodes[i] = NodeWithId({nodeId: nodeId, node: _nodes[nodeId]});
            }
        }
        return activeNodes;
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
        for (uint32 i = 0; i < _activeReplicationNodes.length(); i++) {
            uint256 nodeId = _activeReplicationNodes.at(i);
            if (_nodeExists(nodeId)) {
                activeNodes[i] = NodeWithId({nodeId: nodeId, node: _nodes[nodeId]});
            }
        }
        return activeNodes;
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

    /// @dev Modifier to check if a node is not disabled.
    modifier whenNotDisabled(uint256 nodeId) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        require(!_nodes[nodeId].isDisabled, NodeIsDisabled());
        _;
    }

    /// @dev Checks if a node exists.
    /// @param nodeId The ID of the node to check.
    /// @return True if the node exists, false otherwise.
    function _nodeExists(uint256 nodeId) private view returns (bool) {
        return _ownerOf(nodeId) != address(0);
    }

    /// @inheritdoc ERC721
    function _baseURI() internal view virtual override returns (string memory) {
        return _baseTokenURI;
    }

    /// @dev Helper function to add a node to the active API nodes set.
    function _activateApiNode(uint256 nodeId) private {
        require(_activeApiNodes.length() < maxActiveNodes, MaxActiveNodesReached());
        _nodes[nodeId].isApiEnabled = true;
        if (!_activeApiNodes.contains(nodeId)) {
            // slither-disable-next-line unused-return
            _activeApiNodes.add(nodeId);
        }
        emit ApiEnabled(nodeId);
    }

    /// @dev Helper function to remove a node from the active API nodes set.
    function _disableApiNode(uint256 nodeId) private {
        _nodes[nodeId].isApiEnabled = false;
        if (_activeApiNodes.contains(nodeId)) {
            // slither-disable-next-line unused-return
            _activeApiNodes.remove(nodeId);
        }
        emit ApiDisabled(nodeId);
    }

    /// @dev Helper function to add a node to the active replication nodes set.
    function _activateReplicationNode(uint256 nodeId) private {
        require(_activeReplicationNodes.length() < maxActiveNodes, MaxActiveNodesReached());
        _nodes[nodeId].isReplicationEnabled = true;
        if (!_activeReplicationNodes.contains(nodeId)) {
            // slither-disable-next-line unused-return
            _activeReplicationNodes.add(nodeId);
        }
        emit ReplicationEnabled(nodeId);
    }

    /// @dev Helper function to remove a node from the active replication nodes set.
    function _disableReplicationNode(uint256 nodeId) private {
        _nodes[nodeId].isReplicationEnabled = false;
        if (_activeReplicationNodes.contains(nodeId)) {
            // slither-disable-next-line unused-return
            _activeReplicationNodes.remove(nodeId);
        }
        emit ReplicationDisabled(nodeId);
    }
    

    /// @dev Override to support ERC721, IERC165, and AccessControlEnumerable.
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, IERC165, AccessControlDefaultAdminRules)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
