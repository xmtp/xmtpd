// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { Test } from "forge-std/src/Test.sol";

import { IAccessControl } from "@openzeppelin/contracts/access/IAccessControl.sol";
import { IAccessControlDefaultAdminRules } from
    "@openzeppelin/contracts/access/extensions/IAccessControlDefaultAdminRules.sol";
import { IERC721 } from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import { IERC721Errors } from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";
import { IERC165 } from "@openzeppelin/contracts/interfaces/IERC165.sol";

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

import { INodes, INodesEvents, INodesErrors } from "../src/interfaces/INodes.sol";

import { NodesHarness } from "./utils/Harnesses.sol";
import { Utils } from "./utils/Utils.sol";

contract NodesTest is Test, Utils {
    bytes32 constant DEFAULT_ADMIN_ROLE = 0x00;
    bytes32 constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 constant NODE_MANAGER_ROLE = keccak256("NODE_MANAGER_ROLE");

    uint32 constant NODE_INCREMENT = 100;

    uint256 public constant MAX_BPS = 10_000;

    NodesHarness nodes;

    address admin = makeAddr("admin");
    address manager = makeAddr("manager");
    address unauthorized = makeAddr("unauthorized");

    address alice = makeAddr("alice");
    address bob = makeAddr("bob");

    function setUp() public {
        nodes = new NodesHarness(admin);

        vm.prank(admin);
        nodes.grantRole(NODE_MANAGER_ROLE, manager);
    }

    /* ============ initial state ============ */

    function test_initialState() public view {
        assertEq(nodes.maxActiveNodes(), 20);
    }

    /* ============ addNode ============ */

    function test_addNode_first() public {
        INodes.Node memory node = _getRandomNode();

        address operatorAddress = vm.randomAddress();

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeAdded(
            NODE_INCREMENT, operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars
        );

        vm.prank(admin);
        uint256 nodeId =
            nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);

        assertEq(nodeId, NODE_INCREMENT);

        assertEq(nodes.__getOwner(nodeId), operatorAddress);

        assertEq(nodes.__getNode(nodeId).signingKeyPub, node.signingKeyPub);
        assertEq(nodes.__getNode(nodeId).httpAddress, node.httpAddress);
        assertEq(nodes.__getNode(nodeId).isDisabled, false);
        assertEq(nodes.__getNode(nodeId).isApiEnabled, false);
        assertEq(nodes.__getNode(nodeId).isReplicationEnabled, false);
        assertEq(nodes.__getNode(nodeId).minMonthlyFeeMicroDollars, node.minMonthlyFeeMicroDollars);

        assertEq(nodes.__getNodeCounter(), 1);
    }

    function test_addNode_nth() public {
        INodes.Node memory node = _getRandomNode();

        address operatorAddress = vm.randomAddress();

        nodes.__setNodeCounter(11);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeAdded(
            12 * NODE_INCREMENT, operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars
        );

        vm.prank(admin);
        uint256 nodeId =
            nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);

        assertEq(nodeId, 12 * NODE_INCREMENT);

        assertEq(nodes.__getOwner(nodeId), operatorAddress);

        assertEq(nodes.__getNode(nodeId).signingKeyPub, node.signingKeyPub);
        assertEq(nodes.__getNode(nodeId).httpAddress, node.httpAddress);
        assertEq(nodes.__getNode(nodeId).isDisabled, false);
        assertEq(nodes.__getNode(nodeId).isApiEnabled, false);
        assertEq(nodes.__getNode(nodeId).isReplicationEnabled, false);
        assertEq(nodes.__getNode(nodeId).minMonthlyFeeMicroDollars, node.minMonthlyFeeMicroDollars);

        assertEq(nodes.__getNodeCounter(), 12);
    }

    function test_addNode_invalidAddress() public {
        INodes.Node memory node = _getRandomNode();

        vm.expectRevert(INodesErrors.InvalidAddress.selector);

        vm.prank(admin);
        nodes.addNode(address(0), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
    }

    function test_addNode_invalidSigningKey() public {
        INodes.Node memory node = _getRandomNode();

        vm.expectRevert(INodesErrors.InvalidSigningKey.selector);

        vm.prank(admin);
        nodes.addNode(vm.randomAddress(), bytes(""), node.httpAddress, node.minMonthlyFeeMicroDollars);
    }

    function test_addNode_invalidHttpAddress() public {
        INodes.Node memory node = _getRandomNode();

        vm.expectRevert(INodesErrors.InvalidHttpAddress.selector);

        vm.prank(admin);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, "", node.minMonthlyFeeMicroDollars);
    }

    function test_addNode_notAdmin() public {
        INodes.Node memory node = _getRandomNode();

        // Addresses without DEFAULT_ADMIN_ROLE cannot add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);

        // NODE_MANAGER_ROLE is not authorized to add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, manager, ADMIN_ROLE)
        );

        vm.prank(manager);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
    }

    /* ============ enableNode ============ */

    function test_enableNode() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeEnabled(1);

        vm.prank(admin);
        nodes.enableNode(1);
    }

    function test_enableNode_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(admin);
        nodes.enableNode(1);
    }

    function test_enableNode_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.enableNode(0);
    }

    /* ============ disableNode ============ */

    function test_disableNode() public {
        _addNode(1, alice, "", "", true, true, false, 0);
        nodes.__addToActiveApiNodesSet(1);
        nodes.__addToActiveReplicationNodesSet(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiDisabled(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationDisabled(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeDisabled(1);

        vm.prank(admin);
        nodes.disableNode(1);

        assertFalse(nodes.__getNode(1).isReplicationEnabled);
        assertFalse(nodes.__getNode(1).isApiEnabled);
        assertTrue(nodes.__getNode(1).isDisabled);

        assertFalse(nodes.__activeApiNodesSetContains(1));
        assertFalse(nodes.__activeReplicationNodesSetContains(1));
    }

    function test_disableNode_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(admin);
        nodes.disableNode(1);
    }

    function test_disableNode_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.disableNode(0);
    }

    /* ============ removeFromApiNodes ============ */

    function test_removeFromApiNodes() public {
        _addNode(1, alice, "", "", false, true, false, 0);
        nodes.__addToActiveApiNodesSet(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiDisabled(1);

        vm.prank(admin);
        nodes.removeFromApiNodes(1);

        assertFalse(nodes.__getNode(1).isApiEnabled);
        assertFalse(nodes.__activeApiNodesSetContains(1));
    }

    function test_removeFromApiNodes_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(admin);
        nodes.removeFromApiNodes(1);
    }

    function test_removeFromApiNodes_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.removeFromApiNodes(0);
    }

    /* ============ removeFromReplicationNodes ============ */

    function test_removeFromReplicationNodes() public {
        _addNode(1, alice, "", "", true, false, false, 0);
        nodes.__addToActiveReplicationNodesSet(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationDisabled(1);

        vm.prank(admin);
        nodes.removeFromReplicationNodes(1);

        assertFalse(nodes.__getNode(1).isApiEnabled);
        assertFalse(nodes.__activeReplicationNodesSetContains(1));
    }

    function test_removeFromReplicationNodes_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(admin);
        nodes.removeFromReplicationNodes(1);
    }

    function test_removeFromReplicationNodes_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.removeFromReplicationNodes(0);
    }

    /* ============ transferFrom ============ */

    function test_transferFrom() public {
        _addNode(1, alice, "", "", false, false, false, 0);
        nodes.__addToActiveApiNodesSet(1);
        nodes.__addToActiveReplicationNodesSet(1);

        nodes.__setApproval(manager, 1, alice);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiDisabled(1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationDisabled(1);

        vm.expectEmit(address(nodes));
        emit IERC721.Transfer(alice, bob, 1);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeTransferred(1, alice, bob);

        vm.prank(manager);
        nodes.transferFrom(alice, bob, 1);

        assertFalse(nodes.__getNode(1).isApiEnabled);
        assertFalse(nodes.__getNode(1).isReplicationEnabled);

        assertFalse(nodes.__activeApiNodesSetContains(1));
        assertFalse(nodes.__activeReplicationNodesSetContains(1));

        assertEq(nodes.ownerOf(1), bob);
    }

    function test_transferFrom_unauthorized() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, NODE_MANAGER_ROLE
            )
        );

        vm.prank(unauthorized);
        nodes.transferFrom(alice, bob, 1);

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, alice, NODE_MANAGER_ROLE)
        );

        vm.prank(alice);
        nodes.transferFrom(alice, bob, 1);
    }

    function test_transferFrom_insufficientApproval() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(abi.encodeWithSelector(IERC721Errors.ERC721InsufficientApproval.selector, manager, 1));

        vm.prank(manager);
        nodes.transferFrom(alice, bob, 1);
    }

    /* ============ setHttpAddress ============ */

    function test_setHttpAddress() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectEmit(address(nodes));

        emit INodesEvents.HttpAddressUpdated(1, "http://example.com");

        vm.prank(manager);
        nodes.setHttpAddress(1, "http://example.com");

        assertEq(nodes.__getNode(1).httpAddress, "http://example.com");
    }

    function test_setHttpAddress_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(manager);
        nodes.setHttpAddress(1, "");
    }

    function test_setHttpAddress_invalidHttpAddress() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(INodesErrors.InvalidHttpAddress.selector);

        vm.prank(manager);
        nodes.setHttpAddress(1, "");
    }

    function test_setHttpAddress_notManager() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, NODE_MANAGER_ROLE
            )
        );

        vm.prank(unauthorized);
        nodes.setHttpAddress(1, "");

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, alice, NODE_MANAGER_ROLE)
        );

        vm.prank(alice);
        nodes.setHttpAddress(1, "");
    }

    /* ============ setIsApiEnabled ============ */

    function test_setIsApiEnabled() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiEnabled(1);

        vm.prank(alice);
        nodes.setIsApiEnabled(1, true);

        assertTrue(nodes.__getNode(1).isApiEnabled);
        assertTrue(nodes.__activeApiNodesSetContains(1));

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiDisabled(1);

        vm.prank(alice);
        nodes.setIsApiEnabled(1, false);

        assertFalse(nodes.__getNode(1).isApiEnabled);
        assertFalse(nodes.__activeApiNodesSetContains(1));
    }

    function test_setIsApiEnabled_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setIsApiEnabled(1, true);
    }

    function test_setIsApiEnabled_nodeIsDisabled() public {
        _addNode(1, alice, "", "", false, false, true, 0);

        vm.expectRevert(INodesErrors.NodeIsDisabled.selector);

        nodes.setIsApiEnabled(1, true);
    }

    function test_setIsApiEnabled_notOwner() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(unauthorized);
        nodes.setIsApiEnabled(1, true);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(admin);
        nodes.setIsApiEnabled(1, true);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(manager);
        nodes.setIsApiEnabled(1, true);
    }

    /* ============ setIsReplicationEnabled ============ */

    function test_setIsReplicationEnabled() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationEnabled(1);

        vm.prank(alice);
        nodes.setIsReplicationEnabled(1, true);

        assertTrue(nodes.__getNode(1).isReplicationEnabled);
        assertTrue(nodes.__activeReplicationNodesSetContains(1));

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationDisabled(1);

        vm.prank(alice);
        nodes.setIsReplicationEnabled(1, false);

        assertFalse(nodes.__getNode(1).isReplicationEnabled);
        assertFalse(nodes.__activeReplicationNodesSetContains(1));
    }

    function test_setIsReplicationEnabled_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setIsReplicationEnabled(1, true);
    }

    function test_setIsReplicationEnabled_nodeIsDisabled() public {
        _addNode(1, alice, "", "", false, false, true, 0);

        vm.expectRevert(INodesErrors.NodeIsDisabled.selector);

        nodes.setIsReplicationEnabled(1, true);
    }

    function test_setIsReplicationEnabled_notOwner() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(unauthorized);
        nodes.setIsReplicationEnabled(1, true);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(admin);
        nodes.setIsReplicationEnabled(1, true);

        vm.expectRevert(INodesErrors.Unauthorized.selector);

        vm.prank(manager);
        nodes.setIsReplicationEnabled(1, true);
    }

    /* ============ setMinMonthlyFee ============ */

    function test_setMinMonthlyFee() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectEmit(address(nodes));
        emit INodesEvents.MinMonthlyFeeUpdated(1, 1000);

        vm.prank(manager);
        nodes.setMinMonthlyFee(1, 1000);

        assertEq(nodes.__getNode(1).minMonthlyFeeMicroDollars, 1000);
    }

    function test_setMinMonthlyFee_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);

        vm.prank(manager);
        nodes.setMinMonthlyFee(1, 0);
    }

    function test_setMinMonthlyFee_notManager() public {
        _addNode(1, alice, "", "", false, false, false, 0);

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, NODE_MANAGER_ROLE
            )
        );

        vm.prank(unauthorized);
        nodes.setMinMonthlyFee(0, 0);

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, alice, NODE_MANAGER_ROLE)
        );

        vm.prank(alice);
        nodes.setMinMonthlyFee(0, 0);
    }

    /* ============ setMaxActiveNodes ============ */

    function test_setMaxActiveNodes() public {
        vm.expectEmit(address(nodes));
        emit INodesEvents.MaxActiveNodesUpdated(10);

        vm.prank(admin);
        nodes.setMaxActiveNodes(10);

        assertEq(nodes.maxActiveNodes(), 10);
    }

    function test_setMaxActiveNodes_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.setMaxActiveNodes(0);
    }

    function test_setMaxActiveNodes_lessThanActiveApiNodesLength() public {
        nodes.__addToActiveApiNodesSet(1);

        vm.expectRevert(INodesErrors.MaxActiveNodesBelowCurrentCount.selector);

        vm.prank(admin);
        nodes.setMaxActiveNodes(0);
    }

    function test_setMaxActiveNodes_lessThanReplicationApiNodesLength() public {
        nodes.__addToActiveReplicationNodesSet(1);

        vm.expectRevert(INodesErrors.MaxActiveNodesBelowCurrentCount.selector);

        vm.prank(admin);
        nodes.setMaxActiveNodes(0);
    }

    /* ============ setNodeOperatorCommissionPercent ============ */

    function test_setNodeOperatorCommissionPercent() public {
        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeOperatorCommissionPercentUpdated(1000);

        vm.prank(admin);
        nodes.setNodeOperatorCommissionPercent(1000);

        assertEq(nodes.nodeOperatorCommissionPercent(), 1000);
    }

    function test_setNodeOperatorCommissionPercent_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.setNodeOperatorCommissionPercent(0);
    }

    function test_setNodeOperatorCommissionPercent_invalidCommissionPercent() public {
        vm.expectRevert(INodesErrors.InvalidCommissionPercent.selector);

        vm.prank(admin);
        nodes.setNodeOperatorCommissionPercent(MAX_BPS + 1);
    }

    /* ============ setBaseURI ============ */

    function test_setBaseURI() public {
        vm.expectEmit(address(nodes));
        emit INodesEvents.BaseURIUpdated("http://example.com/");

        vm.prank(admin);
        nodes.setBaseURI("http://example.com/");

        assertEq(nodes.__getBaseTokenURI(), "http://example.com/");
    }

    function test_setBaseURI_notAdmin() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, ADMIN_ROLE)
        );

        vm.prank(unauthorized);
        nodes.setBaseURI("");

        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, manager, ADMIN_ROLE)
        );

        vm.prank(manager);
        nodes.setBaseURI("");
    }

    function test_setBaseURI_emptyURI() public {
        vm.expectRevert(INodesErrors.InvalidURI.selector);

        vm.prank(admin);
        nodes.setBaseURI("");
    }

    function test_setBaseURI_noTrailingSlash() public {
        vm.expectRevert(INodesErrors.InvalidURI.selector);

        vm.prank(admin);
        nodes.setBaseURI("http://example.com");
    }

    /* ============ getAllNodes ============ */

    function test_getAllNodes() public {
        INodes.NodeWithId[] memory allNodes;

        _addNode(NODE_INCREMENT, alice, "", "", false, false, false, 0);
        nodes.__setNodeCounter(1);

        allNodes = nodes.getAllNodes();

        assertEq(allNodes.length, 1);
        assertEq(allNodes[0].nodeId, NODE_INCREMENT);

        _addNode(NODE_INCREMENT * 2, alice, "", "", false, false, false, 0);
        nodes.__setNodeCounter(2);

        allNodes = nodes.getAllNodes();

        assertEq(allNodes.length, 2);
        assertEq(allNodes[0].nodeId, NODE_INCREMENT);
        assertEq(allNodes[1].nodeId, NODE_INCREMENT * 2);

        // NOTE: NodeIds not divisible by `NODE_INCREMENT` are not included, but affect length.
        _addNode(NODE_INCREMENT - 1, alice, "", "", false, false, false, 0);
        nodes.__setNodeCounter(3);

        allNodes = nodes.getAllNodes();

        assertEq(allNodes.length, 3);
        assertEq(allNodes[0].nodeId, NODE_INCREMENT);
        assertEq(allNodes[1].nodeId, NODE_INCREMENT * 2);
        assertEq(allNodes[2].nodeId, 0);

        // NOTE: Nodes that do not exist are not included, but affect length.
        _addNode(NODE_INCREMENT * 3, alice, "", "", false, false, false, 0);
        nodes.__setNodeCounter(4);
        nodes.__burn(NODE_INCREMENT * 3);

        allNodes = nodes.getAllNodes();

        assertEq(allNodes.length, 4);
        assertEq(allNodes[0].nodeId, NODE_INCREMENT);
        assertEq(allNodes[1].nodeId, NODE_INCREMENT * 2);
        assertEq(allNodes[2].nodeId, 0);
        assertEq(allNodes[3].nodeId, 0);
    }

    /* ============ getAllNodesCount ============ */

    function test_getAllNodesCount() public {
        nodes.__setNodeCounter(1);

        assertEq(nodes.getAllNodesCount(), 1);

        nodes.__setNodeCounter(2);

        assertEq(nodes.getAllNodesCount(), 2);
    }

    /* ============ getNode ============ */

    function test_getNode() public {
        _addNode(1, alice, hex"1F1F1F", "httpAddress", true, true, true, 1000);

        INodes.Node memory node = nodes.__getNode(1);

        assertEq(node.signingKeyPub, hex"1F1F1F");
        assertEq(node.httpAddress, "httpAddress");
        assertTrue(node.isReplicationEnabled);
        assertTrue(node.isApiEnabled);
        assertTrue(node.isDisabled);
        assertEq(node.minMonthlyFeeMicroDollars, 1000);
    }

    function test_getNode_nodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.getNode(1);
    }

    /* ============ getActiveApiNodes ============ */

    function test_getActiveApiNodes() public {
        INodes.NodeWithId[] memory activeNodes;

        _addNode(1, alice, "", "", false, false, false, 0);
        nodes.__addToActiveApiNodesSet(1);

        activeNodes = nodes.getActiveApiNodes();

        assertEq(activeNodes.length, 1);
        assertEq(activeNodes[0].nodeId, 1);

        _addNode(2, alice, "", "", false, false, false, 0);
        nodes.__addToActiveApiNodesSet(2);

        activeNodes = nodes.getActiveApiNodes();

        assertEq(activeNodes.length, 2);
        assertEq(activeNodes[0].nodeId, 1);
        assertEq(activeNodes[1].nodeId, 2);

        // NOTE: Nodes that do not exist are not included, but affect length.
        _addNode(3, alice, "", "", false, false, false, 0);
        nodes.__addToActiveApiNodesSet(3);
        nodes.__burn(3);

        activeNodes = nodes.getActiveApiNodes();

        assertEq(activeNodes.length, 3);
        assertEq(activeNodes[0].nodeId, 1);
        assertEq(activeNodes[1].nodeId, 2);
        assertEq(activeNodes[2].nodeId, 0);
    }

    /* ============ getActiveReplicationNodes ============ */

    function test_getActiveReplicationNodes() public {
        INodes.NodeWithId[] memory activeNodes;

        _addNode(1, alice, "", "", false, false, false, 0);
        nodes.__addToActiveReplicationNodesSet(1);

        activeNodes = nodes.getActiveReplicationNodes();

        assertEq(activeNodes.length, 1);
        assertEq(activeNodes[0].nodeId, 1);

        _addNode(2, alice, "", "", false, false, false, 0);
        nodes.__addToActiveReplicationNodesSet(2);

        activeNodes = nodes.getActiveReplicationNodes();

        assertEq(activeNodes.length, 2);
        assertEq(activeNodes[0].nodeId, 1);
        assertEq(activeNodes[1].nodeId, 2);

        // NOTE: Nodes that do not exist are not included, but affect length.
        _addNode(3, alice, "", "", false, false, false, 0);
        nodes.__addToActiveReplicationNodesSet(3);
        nodes.__burn(3);

        activeNodes = nodes.getActiveReplicationNodes();

        assertEq(activeNodes.length, 3);
        assertEq(activeNodes[0].nodeId, 1);
        assertEq(activeNodes[1].nodeId, 2);
        assertEq(activeNodes[2].nodeId, 0);
    }

    /* ============ getActiveApiNodesIDs ============ */

    function test_getActiveApiNodesIDs() public {
        nodes.__addToActiveApiNodesSet(1);
        nodes.__addToActiveApiNodesSet(2);
        nodes.__addToActiveApiNodesSet(3);

        uint256[] memory nodeIds = nodes.getActiveApiNodesIDs();

        assertEq(nodeIds.length, 3);
        assertEq(nodeIds[0], 1);
        assertEq(nodeIds[1], 2);
        assertEq(nodeIds[2], 3);
    }

    /* ============ getActiveReplicationNodesIDs ============ */

    function test_getActiveReplicationNodesIDs() public {
        nodes.__addToActiveReplicationNodesSet(1);
        nodes.__addToActiveReplicationNodesSet(2);
        nodes.__addToActiveReplicationNodesSet(3);

        uint256[] memory nodeIds = nodes.getActiveReplicationNodesIDs();

        assertEq(nodeIds.length, 3);
        assertEq(nodeIds[0], 1);
        assertEq(nodeIds[1], 2);
        assertEq(nodeIds[2], 3);
    }

    /* ============ getActiveApiNodesCount ============ */

    function test_getActiveApiNodesCount() public {
        nodes.__addToActiveApiNodesSet(1);
        nodes.__addToActiveApiNodesSet(2);
        nodes.__addToActiveApiNodesSet(3);

        assertEq(nodes.getActiveApiNodesCount(), 3);
    }

    /* ============ getActiveReplicationNodesCount ============ */

    function test_getActiveReplicationNodesCount() public {
        nodes.__addToActiveReplicationNodesSet(1);
        nodes.__addToActiveReplicationNodesSet(2);
        nodes.__addToActiveReplicationNodesSet(3);

        assertEq(nodes.getActiveReplicationNodesCount(), 3);
    }

    /* ============ getApiNodeIsActive ============ */

    function test_getApiNodeIsActive() public {
        nodes.__addToActiveApiNodesSet(1);
        nodes.__addToActiveApiNodesSet(2);
        nodes.__addToActiveApiNodesSet(3);

        assertTrue(nodes.getApiNodeIsActive(1));
        assertTrue(nodes.getApiNodeIsActive(2));
        assertTrue(nodes.getApiNodeIsActive(3));
        assertFalse(nodes.getApiNodeIsActive(4));
    }

    /* ============ getReplicationNodeIsActive ============ */

    function test_getReplicationNodeIsActive() public {
        nodes.__addToActiveReplicationNodesSet(1);
        nodes.__addToActiveReplicationNodesSet(2);
        nodes.__addToActiveReplicationNodesSet(3);

        assertTrue(nodes.getReplicationNodeIsActive(1));
        assertTrue(nodes.getReplicationNodeIsActive(2));
        assertTrue(nodes.getReplicationNodeIsActive(3));
        assertFalse(nodes.getReplicationNodeIsActive(4));
    }

    /* ============ supportsInterface ============ */

    function test_supportsInterface() public view {
        assertTrue(nodes.supportsInterface(type(IERC721).interfaceId));
        assertTrue(nodes.supportsInterface(type(IERC165).interfaceId));
        assertTrue(nodes.supportsInterface(type(IAccessControl).interfaceId));
        assertTrue(nodes.supportsInterface(type(IAccessControlDefaultAdminRules).interfaceId));
    }

    /* ============ revokeRole ============ */

    function test_revokeRole_revokeDefaultAdminRole() public {
        vm.expectRevert(IAccessControlDefaultAdminRules.AccessControlEnforcedDefaultAdminRules.selector);
        nodes.revokeRole(DEFAULT_ADMIN_ROLE, admin);
    }

    /* ============ renounceRole ============ */

    function test_renounceRole_withinDelay() public {
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControlDefaultAdminRules.AccessControlEnforcedDefaultAdminDelay.selector, 0)
        );

        nodes.renounceRole(DEFAULT_ADMIN_ROLE, admin);
    }

    /* ============ helper functions ============ */

    function _addNode(
        uint256 nodeId,
        address nodeOperator,
        bytes memory signingKeyPub,
        string memory httpAddress,
        bool isReplicationEnabled,
        bool isApiEnabled,
        bool isDisabled,
        uint256 minMonthlyFeeMicroDollars
    ) internal {
        nodes.__setNode(
            nodeId,
            signingKeyPub,
            httpAddress,
            isReplicationEnabled,
            isApiEnabled,
            isDisabled,
            minMonthlyFeeMicroDollars
        );
        nodes.__mint(nodeOperator, nodeId);
    }

    function _getRandomNode() internal view returns (INodes.Node memory) {
        return INodes.Node({
            signingKeyPub: _genBytes(32),
            httpAddress: _genString(32),
            isReplicationEnabled: false,
            isApiEnabled: false,
            isDisabled: false,
            minMonthlyFeeMicroDollars: _genRandomInt(100, 10_000)
        });
    }
}
