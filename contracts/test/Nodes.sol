// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/src/Test.sol";
import {Vm} from "forge-std/src/Vm.sol";
import {Utils} from "test/utils/Utils.sol";
import {Nodes} from "src/Nodes.sol";
import {INodes} from "src/interfaces/INodes.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {IAccessControlDefaultAdminRules} from "@openzeppelin/contracts/access/extensions/IAccessControlDefaultAdminRules.sol";
import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {IERC721Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";
import {IERC165} from "@openzeppelin/contracts/interfaces/IERC165.sol";

contract NodesTest is Test, Utils {
    Nodes public nodes;

    address admin = address(this);
    address manager = vm.randomAddress();
    address unauthorized = address(0x1);

    /// @dev Use _addNode to populate addresses and IDs.
    /// @dev Use _addMultipleNodes to populate multiple nodes.
    address nodeOperator;
    uint256 nodeId;

    function setUp() public {
        nodes = new Nodes(admin);
        nodes.grantRole(nodes.NODE_MANAGER_ROLE(), manager);
    }

    // ***************************************************************
    // *                        addNodes                             *
    // ***************************************************************

    function test_addNode() public {
        INodes.Node memory node = _randomNode();

        address operatorAddress = vm.randomAddress();

        vm.expectEmit(address(nodes));
        emit INodes.NodeAdded(100, operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        uint256 tmpNodeId = nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFee);

        vm.assertEq(nodes.ownerOf(tmpNodeId), operatorAddress);
        vm.assertEq(nodes.getNode(tmpNodeId).signingKeyPub, node.signingKeyPub);
        vm.assertEq(nodes.getNode(tmpNodeId).httpAddress, node.httpAddress);
        vm.assertEq(nodes.getNode(tmpNodeId).isActive, false);
        vm.assertEq(nodes.getNode(tmpNodeId).isApiEnabled, false);
        vm.assertEq(nodes.getNode(tmpNodeId).isReplicationEnabled, false);
        vm.assertEq(nodes.getNode(tmpNodeId).minMonthlyFee, node.minMonthlyFee);
    }

    function test_RevertWhen_AddNodeWithZeroAddress() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidAddress.selector);
        nodes.addNode(address(0), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_AddNodeWithInvalidSigningKey() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidSigningKey.selector);
        nodes.addNode(vm.randomAddress(), bytes(""), node.httpAddress, node.minMonthlyFee);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_AddNodeWithInvalidHttpAddress() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidHttpAddress.selector);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, "", node.minMonthlyFee);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_AddNodeUnauthorized() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();

        // Addresses without DEFAULT_ADMIN_ROLE cannot add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);

        // NODE_MANAGER_ROLE is not authorized to add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                manager,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(manager);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        _checkNoLogsEmitted();
    }

    function test_NodeIdIncrementsByHundred() public {
        INodes.Node memory node = _randomNode();
        uint256 firstId = nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        uint256 secondId = nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        assertEq(secondId - firstId, 100);
    }

    // ***************************************************************
    // *                        transferFrom                         *
    // ***************************************************************

    function test_transferFrom() public {
        _addNode();

        address newOwner = vm.randomAddress();

        // Note: The NFT holder must approve manager to transfer the NFT.
        // A transfer cannot happen without approval.
        vm.prank(nodeOperator);
        ERC721(address(nodes)).approve(manager, nodeId);

        vm.prank(manager);
        nodes.transferFrom(nodeOperator, newOwner, nodeId);
        vm.assertEq(nodes.ownerOf(nodeId), newOwner);
    }

    function test_RevertWhen_transferFromUnauthorized() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.transferFrom(nodeOperator, newOwner, nodeId);
    }

    function test_RevertWhen_transferFromNotApproved() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IERC721Errors.ERC721InsufficientApproval.selector,
                manager,
                nodeId
            )
        );
        vm.prank(manager);
        nodes.transferFrom(nodeOperator, newOwner, nodeId);
    }

    function test_RevertWhen_ownerCannotTransfer() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.transferFrom(nodeOperator, newOwner, nodeId);
    }

    // ***************************************************************
    // *                        safeTransferFrom                     *
    // ***************************************************************

    /// @dev Nodes.sol does not override safeTransferFrom.
    /// While this is not necessary, as safeTransferFrom uses transferFrom,
    /// we should test that it works as expected.
    function test_safeTransferFrom() public {
        _addNode();

        address newOwner = vm.randomAddress();

        // Note: The NFT holder must approve manager to transfer the NFT.
        // A transfer cannot happen without approval.
        vm.prank(nodeOperator);
        ERC721(address(nodes)).approve(manager, nodeId);

        vm.prank(manager);
        nodes.safeTransferFrom(nodeOperator, newOwner, nodeId);
        vm.assertEq(nodes.ownerOf(nodeId), newOwner);
    }

    function test_RevertWhen_safeTransferFromUnauthorized() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.safeTransferFrom(nodeOperator, newOwner, nodeId);
    }

    function test_RevertWhen_safeTransferFromNotApproved() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IERC721Errors.ERC721InsufficientApproval.selector,
                manager,
                nodeId
            )
        );
        vm.prank(manager);
        nodes.safeTransferFrom(nodeOperator, newOwner, nodeId);
    }

    function test_RevertWhen_ownerCannotSafeTransferFrom() public {
        _addNode();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.safeTransferFrom(nodeOperator, newOwner, nodeId);
    }

    // ***************************************************************
    // *                        updateHttpAddress                    *
    // ***************************************************************

    function test_updateHttpAddress() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodes.HttpAddressUpdated(nodeId, "http://example.com");
        nodes.updateHttpAddress(nodeId, "http://example.com");
        vm.assertEq(nodes.getNode(nodeId).httpAddress, "http://example.com");
    }

    function test_RevertWhen_updateHttpAddressNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateHttpAddress(1337, "http://example.com");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateHttpAddressInvalidHttpAddress() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodes.InvalidHttpAddress.selector);
        nodes.updateHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateHttpAddressUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateHttpAddressOwnerCannotUpdate() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                  updateIsReplicationEnabled                 *
    // ***************************************************************

    function test_updateIsReplicationEnabled() public {
        _addNode();
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
        vm.expectEmit(address(nodes));
        emit INodes.ReplicationEnabledUpdated(nodeId, true);
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, true);
    }

    function test_RevertWhen_updateIsReplicationEnabledNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateIsReplicationEnabled(1337, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateIsReplicationEnabledUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateIsReplicationEnabled(nodeId, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateIsReplicationEnabledOwnerCannotUpdate() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateIsReplicationEnabled(nodeId, true);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                        updateMinMonthlyFee                  *
    // ***************************************************************

    function test_updateMinMonthlyFee() public {
        _addNode();
        uint256 initialMonthlyFee = nodes.getNode(nodeId).minMonthlyFee;
        vm.expectEmit(address(nodes));
        emit INodes.MinMonthlyFeeUpdated(nodeId, 1000);
        nodes.updateMinMonthlyFee(nodeId, 1000);
        vm.assertEq(nodes.getNode(nodeId).minMonthlyFee, 1000);
        vm.assertNotEq(nodes.getNode(nodeId).minMonthlyFee, initialMonthlyFee);
    }

    function test_RevertWhen_updateMinMonthlyFeeNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateMinMonthlyFee(1337, 1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateMinMonthlyFeeUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateMinMonthlyFee(nodeId, 1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateMinMonthlyFeeOwnerCannotUpdate() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateMinMonthlyFee(nodeId, 1000);
    }

    // ***************************************************************
    // *            updateActive, batchUpdateActive                  *
    // ***************************************************************

    function test_updateActive() public {
        _addNode();
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.prank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);
        vm.expectEmit(address(nodes));
        emit INodes.NodeActivateUpdated(nodeId, true);
        nodes.updateActive(nodeId, true);
        vm.assertEq(nodes.getNode(nodeId).isActive, true);
    }

    function test_RevertWhen_updateActiveInvalidNodeConfig() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodes.InvalidNodeConfig.selector);
        nodes.updateActive(nodeId, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateActiveNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateActive(1337, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateActiveUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                manager,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(manager);
        nodes.updateActive(nodeId, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateActiveOwnerCannotUpdate() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateActive(nodeId, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateActiveNodeAlreadyActive() public {
        _addNode();
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.prank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);

        nodes.updateActive(nodeId, true);
        vm.recordLogs();
        vm.expectRevert(INodes.NodeAlreadyActive.selector);
        nodes.updateActive(nodeId, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateActiveNodeAlreadyInactive() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodes.NodeAlreadyInactive.selector);
        nodes.updateActive(nodeId, false);
        _checkNoLogsEmitted();
    }

    function test_batchUpdateActive() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);

        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        _enableNodes(operators, nodeIds);

        vm.expectEmit(address(nodes));
        emit INodes.NodeActivateUpdated(nodeIds[0], true);
        emit INodes.NodeActivateUpdated(nodeIds[1], true);
        emit INodes.NodeActivateUpdated(nodeIds[2], true);
        nodes.batchUpdateActive(nodeIds, isActive);

        uint256[] memory activeNodesIDs = nodes.getActiveNodesIDs();
        vm.assertTrue(activeNodesIDs.length == 3);
    }

    function test_RevertWhen_batchUpdateActiveUnauthorized() public {
        (, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.batchUpdateActive(nodeIds, isActive);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_batchUpdateActiveInvalidInputLength() public {
        nodes.updateMaxActiveNodes(2);
        (, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](2);
        for (uint256 i = 0; i < 2; i++) {
            isActive[i] = true;
        }

        vm.recordLogs();
        vm.expectRevert(INodes.InvalidInputLength.selector);
        nodes.batchUpdateActive(nodeIds, isActive);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_batchUpdateActiveMaxActiveNodesReached() public {
        nodes.updateMaxActiveNodes(2);
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        _enableNodes(operators, nodeIds);

        vm.recordLogs();
        vm.expectRevert(INodes.MaxActiveNodesReached.selector);
        nodes.batchUpdateActive(nodeIds, isActive);
        Vm.Log[] memory logs = vm.getRecordedLogs();
        vm.assertTrue(logs.length == 2, "Only 2 nodes should be added");
    }

    // ***************************************************************
    // *                    updateMaxActiveNodes                     *
    // ***************************************************************

    function test_updateMaxActiveNodes() public {
        vm.expectEmit(address(nodes));
        emit INodes.MaxActiveNodesUpdated(10);
        nodes.updateMaxActiveNodes(10);
        vm.assertEq(nodes.maxActiveNodes(), 10);
    }

    function test_RevertWhen_updateMaxActiveNodesUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateMaxActiveNodes(10);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateMaxActiveNodesDeactivateNodes() public {
        _addNode();
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.prank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);

        nodes.updateMaxActiveNodes(1);
        nodes.updateActive(nodeId, true);
        vm.recordLogs();
        vm.expectRevert(INodes.MaxActiveNodesBelowCurrentCount.selector);
        nodes.updateMaxActiveNodes(0);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *             updateNodeOperatorCommissionPercent             *
    // ***************************************************************  

    function test_updateNodeOperatorCommissionPercent() public {
        vm.expectEmit(address(nodes));
        emit INodes.NodeOperatorCommissionPercentUpdated(1000);
        nodes.updateNodeOperatorCommissionPercent(1000);
        vm.assertEq(nodes.nodeOperatorCommissionPercent(), 1000);
    }

    function test_RevertWhen_updateNodeOperatorCommissionPercentUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateNodeOperatorCommissionPercent(1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_updateNodeOperatorCommissionPercentInvalidCommissionPercent() public {
        vm.recordLogs();
        vm.expectRevert(INodes.InvalidCommissionPercent.selector);
        nodes.updateNodeOperatorCommissionPercent(10001);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                       setBaseURI                            *
    // ***************************************************************  

    function test_setBaseURI() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodes.BaseURIUpdated("http://example.com/");
        nodes.setBaseURI("http://example.com/");
        vm.assertEq(nodes.tokenURI(100), "http://example.com/100");
    }

    function test_RevertWhen_setBaseURIUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.setBaseURI("http://example.com/");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setBaseURIEmptyURI() public {
        vm.recordLogs();
        vm.expectRevert(INodes.InvalidURI.selector);
        nodes.setBaseURI("");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setBaseURINoTrailingSlash() public {
        vm.recordLogs();
        vm.expectRevert(INodes.InvalidURI.selector);
        nodes.setBaseURI("http://example.com");
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                  updateIsApiEnabled                         *
    // ***************************************************************  

    function test_updateIsApiEnabled() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodes.ApiEnabledUpdated(nodeId, true);
        vm.startPrank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);
        assertEq(nodes.getNode(nodeId).isApiEnabled, true);

        nodes.updateIsApiEnabled(nodeId, false);
        assertEq(nodes.getNode(nodeId).isApiEnabled, false);
        vm.stopPrank();
    }

    function test_RevertWhen_updateIsApiEnabledOnlyNodeOperatorCanUpdate() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodes.Unauthorized.selector);
        /// @dev Default user is admin
        nodes.updateIsApiEnabled(nodeId, true);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                        Getters                              *
    // ***************************************************************

    function test_getAllNodes() public {
        _addNode();

        INodes.NodeWithId[] memory allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 1);
    }

    function test_getAllNodesMultiple() public {
        _addMultipleNodes(10);

        INodes.NodeWithId[] memory allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 10);
    }

    function test_getAllNodesCount() public {
        _addMultipleNodes(1);
        vm.assertEq(nodes.getAllNodesCount(), 1);
        _addMultipleNodes(2);
        vm.assertEq(nodes.getAllNodesCount(), 3);
        _addMultipleNodes(3);
        vm.assertEq(nodes.getAllNodesCount(), 6);
    }

    function test_getNode() public {
        _addNode();
        INodes.Node memory node = nodes.getNode(nodeId);
        vm.assertEq(node.signingKeyPub, nodes.getNode(nodeId).signingKeyPub);
        vm.assertEq(node.httpAddress, nodes.getNode(nodeId).httpAddress);
        vm.assertEq(node.isReplicationEnabled, nodes.getNode(nodeId).isReplicationEnabled);
        vm.assertEq(node.isApiEnabled, nodes.getNode(nodeId).isApiEnabled);
        vm.assertEq(node.isActive, nodes.getNode(nodeId).isActive);
        vm.assertEq(node.minMonthlyFee, nodes.getNode(nodeId).minMonthlyFee);
    }

    function test_RevertWhen_getNodeNodeDoesNotExist() public {
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.getNode(1337);
    }

    function test_getActiveNodes() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }
        _enableNodes(operators, nodeIds);
        nodes.batchUpdateActive(nodeIds, isActive);

        INodes.NodeWithId[] memory activeNodes = nodes.getActiveNodes();
        vm.assertTrue(activeNodes.length == 3);
    }

    function test_getActiveNodesIDs() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }
        _enableNodes(operators, nodeIds);
        nodes.batchUpdateActive(nodeIds, isActive);

        uint256[] memory activeNodeIds = nodes.getActiveNodesIDs();
        vm.assertTrue(activeNodeIds.length == 3);
        vm.assertEq(activeNodeIds[0], nodeIds[0]);
        vm.assertEq(activeNodeIds[1], nodeIds[1]);
        vm.assertEq(activeNodeIds[2], nodeIds[2]);
    }

    function test_getActiveNodesCount() public {
        _addNode();
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.prank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);
        nodes.updateActive(nodeId, true);
        vm.assertEq(nodes.getActiveNodesCount(), 1);

        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(2);
        bool[] memory isActive = new bool[](2);
        for (uint256 i = 0; i < 2; i++) {
            isActive[i] = true;
        }
        _enableNodes(operators, nodeIds);
        nodes.batchUpdateActive(nodeIds, isActive);
        vm.assertEq(nodes.getActiveNodesCount(), 3);

        (address[] memory operators2, uint256[] memory nodeIds2) = _addMultipleNodes(3);
        bool[] memory isActive2 = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive2[i] = true;
        }
        _enableNodes(operators2, nodeIds2);
        nodes.batchUpdateActive(nodeIds2, isActive2);
        vm.assertEq(nodes.getActiveNodesCount(), 6);

        nodes.updateActive(nodeIds[0], false);
        vm.assertEq(nodes.getActiveNodesCount(), 5);
    }

    function test_getNodeIsActive() public {
        _addNode();
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.prank(nodeOperator);
        nodes.updateIsApiEnabled(nodeId, true);

        nodes.updateActive(nodeId, true);
        vm.assertEq(nodes.getNodeIsActive(nodeId), true);

        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(2);
        _enableNodes(operators, nodeIds);

        nodes.updateActive(nodeIds[0], true);
        vm.assertEq(nodes.getNodeIsActive(nodeIds[0]), true);
        vm.assertEq(nodes.getNodeIsActive(nodeIds[1]), false);
    }

    // ***************************************************************
    // *                       Security                              *
    // ***************************************************************

    function test_supportsInterface() public view {
        vm.assertEq(nodes.supportsInterface(type(IERC721).interfaceId), true);
        vm.assertEq(nodes.supportsInterface(type(IERC165).interfaceId), true);
        vm.assertEq(nodes.supportsInterface(type(IAccessControl).interfaceId), true);
        vm.assertEq(nodes.supportsInterface(type(IAccessControlDefaultAdminRules).interfaceId), true);
    }

    function test_RevertWhen_revokeRoleDefaultAdminRole() public {
        bytes32 role = nodes.DEFAULT_ADMIN_ROLE();
        vm.expectRevert(IAccessControlDefaultAdminRules.AccessControlEnforcedDefaultAdminRules.selector);
        nodes.revokeRole(role, admin);
    }

    function test_RevertWhen_renounceRoleDefaultAdminRole() public {
        bytes32 role = nodes.DEFAULT_ADMIN_ROLE();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControlDefaultAdminRules.AccessControlEnforcedDefaultAdminDelay.selector,
                0
            )
        );
        nodes.renounceRole(role, admin);
    }

    // ***************************************************************
    // *                     Helper functions                        *
    // ***************************************************************

    function _addNode() internal {
        INodes.Node memory node1 = _randomNode();
        nodeOperator = vm.randomAddress();
        nodeId = nodes.addNode(nodeOperator, node1.signingKeyPub, node1.httpAddress, node1.minMonthlyFee);
    }

    function _addMultipleNodes(uint256 numberOfNodes) internal returns (address[] memory operators, uint256[] memory nodeIds) {
        operators = new address[](numberOfNodes);
        nodeIds = new uint256[](numberOfNodes);
        for (uint256 i = 0; i < numberOfNodes; i++) {
            INodes.Node memory node = _randomNode();
            operators[i] = vm.randomAddress();
            nodeIds[i] = nodes.addNode(operators[i], node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
        }
        return (operators, nodeIds);
    }

    function _enableNodes(address[] memory operators, uint256[] memory nodeIds) internal {
        for (uint256 i = 0; i < nodeIds.length; i++) {
            nodes.updateIsReplicationEnabled(nodeIds[i], true);
            vm.prank(operators[i]);
            nodes.updateIsApiEnabled(nodeIds[i], true);
        }
    }

    function _randomNode() internal view returns (INodes.Node memory) {
        return INodes.Node({
            signingKeyPub: _genBytes(32), 
            httpAddress: _genString(32), 
            isReplicationEnabled: false, 
            isApiEnabled: false, 
            isActive: false,
            minMonthlyFee: _genRandomInt(100, 10000)
        });
    }

    function _checkNoLogsEmitted() internal {
        Vm.Log[] memory logs = vm.getRecordedLogs();
        vm.assertTrue(logs.length == 0, "No logs should be emitted");
    }
}
