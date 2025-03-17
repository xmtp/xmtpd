// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { EnumerableSet } from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

import { GroupMessages } from "../../src/GroupMessages.sol";
import { IdentityUpdates } from "../../src/IdentityUpdates.sol";
import { Nodes } from "../../src/Nodes.sol";
import { RatesManager } from "../../src/RatesManager.sol";

contract GroupMessagesHarness is GroupMessages {
    function __pause() external {
        _pause();
    }

    function __unpause() external {
        _unpause();
    }

    function __setSequenceId(uint64 sequenceId) external {
        _getGroupMessagesStorage().sequenceId = sequenceId;
    }

    function __setMinPayloadSize(uint256 minPayloadSize) external {
        _getGroupMessagesStorage().minPayloadSize = minPayloadSize;
    }

    function __setMaxPayloadSize(uint256 maxPayloadSize) external {
        _getGroupMessagesStorage().maxPayloadSize = maxPayloadSize;
    }

    function __getSequenceId() external view returns (uint64) {
        return _getGroupMessagesStorage().sequenceId;
    }
}

contract IdentityUpdatesHarness is IdentityUpdates {
    function __pause() external {
        _pause();
    }

    function __unpause() external {
        _unpause();
    }

    function __setSequenceId(uint64 sequenceId) external {
        _getIdentityUpdatesStorage().sequenceId = sequenceId;
    }

    function __setMinPayloadSize(uint256 minPayloadSize) external {
        _getIdentityUpdatesStorage().minPayloadSize = minPayloadSize;
    }

    function __setMaxPayloadSize(uint256 maxPayloadSize) external {
        _getIdentityUpdatesStorage().maxPayloadSize = maxPayloadSize;
    }

    function __getSequenceId() external view returns (uint64) {
        return _getIdentityUpdatesStorage().sequenceId;
    }
}

contract NodesHarness is Nodes {
    using EnumerableSet for EnumerableSet.UintSet;

    constructor(address initialAdmin) Nodes(initialAdmin) { }

    function __setNodeCounter(uint256 nodeCounter) external {
        _nodeCounter = uint32(nodeCounter);
    }

    function __setNodeEnabled(uint256 nodeId) external {
        _nodes[nodeId].isDisabled = false;
    }

    function __setNodeDisabled(uint256 nodeId) external {
        _nodes[nodeId].isDisabled = true;
    }

    function __setNode(
        uint256 nodeId,
        bytes calldata signingKeyPub,
        string calldata httpAddress,
        bool isReplicationEnabled,
        bool isApiEnabled,
        bool isDisabled,
        uint256 minMonthlyFeeMicroDollars
    ) external {
        _nodes[nodeId] =
            Node(signingKeyPub, httpAddress, isReplicationEnabled, isApiEnabled, isDisabled, minMonthlyFeeMicroDollars);
    }

    function __setApproval(address to, uint256 tokenId, address authorizer) external {
        _approve(to, tokenId, authorizer);
    }

    function __mint(address to, uint256 nodeId) external {
        _mint(to, nodeId);
    }

    function __addToActiveApiNodesSet(uint256 nodeId) external {
        _activeApiNodes.add(nodeId);
    }

    function __addToActiveReplicationNodesSet(uint256 nodeId) external {
        _activeReplicationNodes.add(nodeId);
    }

    function __activeApiNodesSetContains(uint256 nodeId) external view returns (bool contains) {
        return _activeApiNodes.contains(nodeId);
    }

    function __activeReplicationNodesSetContains(uint256 nodeId) external view returns (bool contains) {
        return _activeReplicationNodes.contains(nodeId);
    }

    function __getNode(uint256 nodeId) external view returns (Node memory node) {
        return _nodes[nodeId];
    }

    function __getOwner(uint256 nodeId) external view returns (address owner) {
        return _ownerOf(nodeId);
    }

    function __getNodeCounter() external view returns (uint32 nodeCounter) {
        return _nodeCounter;
    }

    function __getBaseTokenURI() external view returns (string memory baseTokenURI) {
        return _baseTokenURI;
    }
}

contract RatesManagerHarness is RatesManager {
    function __pause() external {
        _pause();
    }

    function __unpause() external {
        _unpause();
    }

    function __pushRates(uint256 messageFee, uint256 storageFee, uint256 congestionFee, uint256 targetRatePerMinute, uint256 startTime) external {
        _getRatesManagerStorage().allRates.push(
            Rates(uint64(messageFee), uint64(storageFee), uint64(congestionFee), uint64(targetRatePerMinute), uint64(startTime))
        );
    }

    function __getAllRates() external view returns (Rates[] memory) {
        return _getRatesManagerStorage().allRates;
    }
}
