// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "./interfaces/INodes.sol";


 /// @title XMTP Node Registry.
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
contract Nodes is ERC721, INodes, AccessControl {
    using EnumerableSet for EnumerableSet.UintSet;

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

    /// @dev Active Node Operators IDs set.
    EnumerableSet.UintSet private _activeNodes;

    /// @notice The commission percentage that the node operator receives.
    /// @dev This is stored in basis points (1/100th of a percent).
    /// Example: 1% = 100bps, 10% = 1000bps, 100% = 10000bps.
    /// Comission is calculated as (nodeOperatorCommissionPercent * nodeOperatorFee) / MAX_BPS.
    // slither-disable-next-line constable-states
    uint256 public nodeOperatorCommissionPercent;

    constructor(address _initialAdmin) ERC721("XMTP Node Operator", "XMTP") {
        require(_initialAdmin != address(0), InvalidAddress());

        _grantRole(DEFAULT_ADMIN_ROLE, _initialAdmin);
        _setRoleAdmin(NODE_MANAGER_ROLE, DEFAULT_ADMIN_ROLE);
        _grantRole(NODE_MANAGER_ROLE, _initialAdmin);
    }

    /// @inheritdoc INodes
    function addNode(address to, bytes calldata signingKeyPub, string calldata httpAddress, uint256 minMonthlyFee)
        external
        onlyRole(DEFAULT_ADMIN_ROLE)
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
        _deactivateNode(nodeId);
        super.transferFrom(from, to, nodeId);
        emit NodeTransferred(nodeId, from, to);
    }

    /// @inheritdoc INodes
    function updateHttpAddress(uint256 nodeId, string calldata httpAddress) external onlyRole(NODE_MANAGER_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        require(bytes(httpAddress).length > 0, InvalidHttpAddress());
        _nodes[nodeId].httpAddress = httpAddress;
        emit HttpAddressUpdated(nodeId, httpAddress);
    }

    /// @inheritdoc INodes
    function updateIsReplicationEnabled(uint256 nodeId, bool isReplicationEnabled) external onlyRole(NODE_MANAGER_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        _nodes[nodeId].isReplicationEnabled = isReplicationEnabled;
        emit ReplicationEnabledUpdated(nodeId, isReplicationEnabled);
    }

    /// @inheritdoc INodes
    function updateMinMonthlyFee(uint256 nodeId, uint256 minMonthlyFee) external onlyRole(NODE_MANAGER_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        _nodes[nodeId].minMonthlyFee = minMonthlyFee;
        emit MinMonthlyFeeUpdated(nodeId, minMonthlyFee);
    }

    /// @inheritdoc INodes
    function updateActive(uint256 nodeId, bool isActive) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        if (isActive) {
            require(_activeNodes.length() < maxActiveNodes, MaxActiveNodesReached());
            require(_activeNodes.add(nodeId), NodeAlreadyActive());
        } else {
            require(_activeNodes.remove(nodeId), NodeAlreadyInactive());
        }
        _nodes[nodeId].isActive = isActive;
        emit NodeActivateUpdated(nodeId, isActive);
    }

    /// @inheritdoc INodes
    function batchUpdateActive(uint256[] calldata nodeIds, bool[] calldata isActive) 
        external 
        onlyRole(DEFAULT_ADMIN_ROLE) 
    {
        require(nodeIds.length == isActive.length);
        for (uint256 i = 0; i < nodeIds.length; i++) {
            updateActive(nodeIds[i], isActive[i]);
        }
    }

    /// @inheritdoc INodes
    function updateMaxActiveNodes(uint8 newMaxActiveNodes) external onlyRole(DEFAULT_ADMIN_ROLE) {
        maxActiveNodes = newMaxActiveNodes;
        emit MaxActiveNodesUpdated(newMaxActiveNodes);
    }

    /// @inheritdoc INodes
    function updateNodeOperatorCommissionPercent(uint256 newCommissionPercent) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newCommissionPercent <= MAX_BPS, InvalidCommissionPercent());
        nodeOperatorCommissionPercent = newCommissionPercent;
        emit NodeOperatorCommissionPercentUpdated(newCommissionPercent);
    }   

    /// @inheritdoc INodes
    function setBaseURI(string calldata newBaseURI) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(bytes(newBaseURI).length > 0, "Empty URI not allowed");
        require(bytes(newBaseURI)[bytes(newBaseURI).length - 1] == 0x2f, "URI must end with /");
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
    }

    /// @inheritdoc INodes
    function updateIsApiEnabled(uint256 nodeId) external {
        require(_ownerOf(nodeId) == msg.sender, Unauthorized());
        _nodes[nodeId].isApiEnabled = !_nodes[nodeId].isApiEnabled;
        emit ApiEnabledUpdated(nodeId, _nodes[nodeId].isApiEnabled);
    }

    /// @inheritdoc INodes
    function allNodes() public view returns (NodeWithId[] memory) {
        NodeWithId[] memory allNodesList = new NodeWithId[](_nodeCounter);
        for (uint32 i = 0; i < _nodeCounter; i++) {
            uint32 nodeId = NODE_INCREMENT * (i + 1);
            if (_nodeExists(nodeId)) {
                allNodesList[i] = NodeWithId({nodeId: nodeId, node: _nodes[nodeId]});
            }
        }
        return allNodesList;
    }

    /// @inheritdoc INodes
    function getNode(uint256 nodeId) public view returns (Node memory) {
        require(_nodeExists(nodeId), NodeDoesNotExist());
        return _nodes[nodeId];
    }

    /// @inheritdoc INodes
    function getActiveNodes() external view returns (Node[] memory activeNodes) {
        activeNodes = new Node[](_activeNodes.length());
        for (uint32 i = 0; i < _activeNodes.length(); i++) {
            activeNodes[i] = _nodes[_activeNodes.at(i)];
        }
        return activeNodes;
    }

    /// @inheritdoc INodes
    function getActiveNodesIDs() external view returns (uint256[] memory activeNodesIDs) {
        return _activeNodes.values();
    }

    /// @inheritdoc INodes
    function nodeIsActive(uint256 nodeId) external view returns (bool) {
        return _activeNodes.contains(nodeId);
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

    /// @dev Helper function to deactivate a node
    function _deactivateNode(uint256 nodeId) private {
        if (_activeNodes.contains(nodeId)) {
            _activeNodes.remove(nodeId);
            _nodes[nodeId].isActive = false;
            emit NodeActivateUpdated(nodeId, false);
        }
    }

    /// @dev Required override for AccessControl
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, IERC165, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
