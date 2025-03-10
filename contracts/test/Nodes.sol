// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/src/Test.sol";
import {Vm} from "forge-std/src/Vm.sol";
import {Utils} from "test/utils/Utils.sol";
import {Nodes} from "src/Nodes.sol";
import {INodes, INodesEvents, INodesErrors} from "src/interfaces/INodes.sol";
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

    /// @dev Use _addNode to populate nodeOperator and nodeId.
    /// @dev Use _enableNode to enable the node.
    /// @dev Use _addMultipleNodes to populate multiple nodes.
    /// @dev Use _enableNodes to enable multiple nodes.
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
        emit INodesEvents.NodeAdded(100, operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        uint256 tmpNodeId = nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);

        vm.assertEq(nodes.ownerOf(tmpNodeId), operatorAddress);
        vm.assertEq(nodes.getNode(tmpNodeId).signingKeyPub, node.signingKeyPub);
        vm.assertEq(nodes.getNode(tmpNodeId).httpAddress, node.httpAddress);
        vm.assertEq(nodes.getNode(tmpNodeId).isDisabled, false);
        vm.assertEq(nodes.getNode(tmpNodeId).isApiEnabled, false);
        vm.assertEq(nodes.getNode(tmpNodeId).isReplicationEnabled, false);
        vm.assertEq(nodes.getNode(tmpNodeId).minMonthlyFeeMicroDollars, node.minMonthlyFeeMicroDollars);
    }

    function test_RevertWhen_AddNodeWithZeroAddress() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodesErrors.InvalidAddress.selector);
        nodes.addNode(address(0), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_AddNodeWithInvalidSigningKey() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodesErrors.InvalidSigningKey.selector);
        nodes.addNode(vm.randomAddress(), bytes(""), node.httpAddress, node.minMonthlyFeeMicroDollars);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_AddNodeWithInvalidHttpAddress() public {
        vm.recordLogs();
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodesErrors.InvalidHttpAddress.selector);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, "", node.minMonthlyFeeMicroDollars);
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
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);

        // NODE_MANAGER_ROLE is not authorized to add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                manager,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(manager);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        _checkNoLogsEmitted();
    }

    function test_NodeIdIncrementsByHundred() public {
        INodes.Node memory node = _randomNode();
        uint256 firstId = nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        uint256 secondId = nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        assertEq(secondId - firstId, 100);
    }

    // ***************************************************************
    // *                    enable/disableNode                       *
    // ***************************************************************

    function test_disableNode() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);
        assertEq(nodes.getNode(nodeId).isDisabled, false);
        assertEq(nodes.getNode(nodeId).isApiEnabled, true);
        assertEq(nodes.getNode(nodeId).isReplicationEnabled, true);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeDisabled(nodeId);
        nodes.disableNode(nodeId);
        assertEq(nodes.getNode(nodeId).isDisabled, true);
        assertEq(nodes.getNode(nodeId).isApiEnabled, false);
        assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);

        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeEnabled(nodeId);
        nodes.enableNode(nodeId);
        assertEq(nodes.getNode(nodeId).isDisabled, false);
        assertEq(nodes.getNode(nodeId).isApiEnabled, false);
        assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
    }

    function test_RevertWhen_disableNodeNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.disableNode(1337);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_disableNodeUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.disableNode(nodeId);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                    removeFromApiNodes                       *
    // ***************************************************************

    function test_removeFromApiNodes() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);
        nodes.removeFromApiNodes(nodeId);
        assertEq(nodes.getNode(nodeId).isApiEnabled, false);
    }

    function test_RevertWhen_removeFromApiNodesNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.removeFromApiNodes(1337);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_removeFromApiNodesUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.removeFromApiNodes(nodeId);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                  removeFromReplicationNodes                 *
    // ***************************************************************

    function test_removeFromReplicationNodes() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);
        nodes.removeFromReplicationNodes(nodeId);
        assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
    }

    function test_RevertWhen_removeFromReplicationNodesNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.removeFromReplicationNodes(1337);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_removeFromReplicationNodesUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.removeFromReplicationNodes(nodeId);
        _checkNoLogsEmitted();
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
    // *                      setHttpAddress                         *
    // ***************************************************************

    function test_setHttpAddress() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodesEvents.HttpAddressUpdated(nodeId, "http://example.com");
        nodes.setHttpAddress(nodeId, "http://example.com");
        vm.assertEq(nodes.getNode(nodeId).httpAddress, "http://example.com");
    }

    function test_RevertWhen_setHttpAddressNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setHttpAddress(1337, "http://example.com");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setHttpAddressInvalidHttpAddress() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodesErrors.InvalidHttpAddress.selector);
        nodes.setHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setHttpAddressUnauthorized() public {
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
        nodes.setHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setHttpAddressOwnerCannotUpdate() public {
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
        nodes.setHttpAddress(nodeId, "");
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                     setIsApiEnabled                         *
    // ***************************************************************

    function test_setIsApiEnabled() public {
        _addNode();
        vm.assertEq(nodes.getNode(nodeId).isApiEnabled, false);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiEnabled(nodeId);
        vm.startPrank(nodeOperator);
        nodes.setIsApiEnabled(nodeId, true);
        assertEq(nodes.getNode(nodeId).isApiEnabled, true);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ApiDisabled(nodeId);
        nodes.setIsApiEnabled(nodeId, false);
        assertEq(nodes.getNode(nodeId).isApiEnabled, false);
        vm.stopPrank();
    }

    function test_RevertWhen_setIsApiEnabledNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setIsApiEnabled(1337, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setIsApiEnabledNodeIsDisabled() public {
        _addNode();

        nodes.disableNode(nodeId);

        vm.expectRevert(INodesErrors.NodeIsDisabled.selector);
        vm.prank(nodeOperator);
        nodes.setIsApiEnabled(nodeId, true);
    }

    function test_RevertWhen_setIsApiEnabledUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodesErrors.Unauthorized.selector);
        vm.prank(unauthorized);
        nodes.setIsApiEnabled(nodeId, true);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                  setIsReplicationEnabled                    *
    // ***************************************************************

    function test_setIsReplicationEnabled() public {
        _addNode();
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationEnabled(nodeId);
        vm.startPrank(nodeOperator);
        nodes.setIsReplicationEnabled(nodeId, true);
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, true);

        vm.expectEmit(address(nodes));
        emit INodesEvents.ReplicationDisabled(nodeId);
        nodes.setIsReplicationEnabled(nodeId, false);
        assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
        vm.stopPrank();
    }

    function test_RevertWhen_setIsReplicationEnabledNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setIsReplicationEnabled(1337, true);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setIsReplicationEnabledNodeIsDisabled() public {
        _addNode();

        nodes.disableNode(nodeId);

        vm.expectRevert(INodesErrors.NodeIsDisabled.selector);
        vm.prank(nodeOperator);
        nodes.setIsReplicationEnabled(nodeId, true);
    }

    function test_RevertWhen_setIsReplicationEnabledUnauthorized() public {
        _addNode();
        vm.recordLogs();
        vm.expectRevert(INodesErrors.Unauthorized.selector);
        vm.prank(unauthorized);
        nodes.setIsReplicationEnabled(nodeId, true);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                    setMinMonthlyFee                         *
    // ***************************************************************

    function test_setMinMonthlyFee() public {
        _addNode();
        uint256 initialMonthlyFee = nodes.getNode(nodeId).minMonthlyFeeMicroDollars;
        vm.expectEmit(address(nodes));
        emit INodesEvents.MinMonthlyFeeUpdated(nodeId, 1000);
        nodes.setMinMonthlyFee(nodeId, 1000);
        vm.assertEq(nodes.getNode(nodeId).minMonthlyFeeMicroDollars, 1000);
        vm.assertNotEq(nodes.getNode(nodeId).minMonthlyFeeMicroDollars, initialMonthlyFee);
    }

    function test_RevertWhen_setMinMonthlyFeeNodeDoesNotExist() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.setMinMonthlyFee(1337, 1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setMinMonthlyFeeUnauthorized() public {
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
        nodes.setMinMonthlyFee(nodeId, 1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setMinMonthlyFeeOwnerCannotUpdate() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.setMinMonthlyFee(nodeId, 1000);
    }

    // ***************************************************************
    // *                    setMaxActiveNodes                        *
    // ***************************************************************

    function test_setMaxActiveNodes() public {
        vm.expectEmit(address(nodes));
        emit INodesEvents.MaxActiveNodesUpdated(10);
        nodes.setMaxActiveNodes(10);
        vm.assertEq(nodes.maxActiveNodes(), 10);
    }

    function test_RevertWhen_setMaxActiveNodesUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.setMaxActiveNodes(10);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setMaxActiveNodesBelowCurrentCount() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);

        /// @dev It shouldn't fail to set the maxActiveNodes to the current count.
        nodes.setMaxActiveNodes(1);

        vm.recordLogs();
        vm.expectRevert(INodesErrors.MaxActiveNodesBelowCurrentCount.selector);
        nodes.setMaxActiveNodes(0);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *             setNodeOperatorCommissionPercent                *
    // ***************************************************************  

    function test_setNodeOperatorCommissionPercent() public {
        vm.expectEmit(address(nodes));
        emit INodesEvents.NodeOperatorCommissionPercentUpdated(1000);
        nodes.setNodeOperatorCommissionPercent(1000);
        vm.assertEq(nodes.nodeOperatorCommissionPercent(), 1000);
    }

    function test_RevertWhen_setNodeOperatorCommissionPercentUnauthorized() public {
        vm.recordLogs();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.setNodeOperatorCommissionPercent(1000);
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setNodeOperatorCommissionPercentInvalidCommissionPercent() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.InvalidCommissionPercent.selector);
        nodes.setNodeOperatorCommissionPercent(10001);
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                       setBaseURI                            *
    // ***************************************************************  

    function test_setBaseURI() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodesEvents.BaseURIUpdated("http://example.com/");
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
        vm.expectRevert(INodesErrors.InvalidURI.selector);
        nodes.setBaseURI("");
        _checkNoLogsEmitted();
    }

    function test_RevertWhen_setBaseURINoTrailingSlash() public {
        vm.recordLogs();
        vm.expectRevert(INodesErrors.InvalidURI.selector);
        nodes.setBaseURI("http://example.com");
        _checkNoLogsEmitted();
    }

    // ***************************************************************
    // *                  getNode / getAllNodes                      *
    // ***************************************************************

    function test_getAllNodes() public {
        _addNode();

        INodes.NodeWithId[] memory allNodes = nodes.getAllNodes();
        vm.assertTrue(allNodes.length == 1);
    }

    function test_getAllNodesMultiple() public {
        _addMultipleNodes(10);

        INodes.NodeWithId[] memory allNodes = nodes.getAllNodes();
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
        vm.assertEq(node.isDisabled, nodes.getNode(nodeId).isDisabled);
        vm.assertEq(node.minMonthlyFeeMicroDollars, nodes.getNode(nodeId).minMonthlyFeeMicroDollars);
    }

    function test_RevertWhen_getNodeNodeDoesNotExist() public {
        vm.expectRevert(INodesErrors.NodeDoesNotExist.selector);
        nodes.getNode(1337);
    }

    // ***************************************************************
    // *      getActiveReplicationNodes / getActiveApiNodes          *
    // ***************************************************************

    function test_getActiveApiNodes() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        _enableNodes(operators, nodeIds);

        INodes.NodeWithId[] memory activeApiNodes = nodes.getActiveApiNodes();
        vm.assertTrue(activeApiNodes.length == 3);

        INodes.NodeWithId[] memory activeReplicationNodes = nodes.getActiveReplicationNodes();
        vm.assertTrue(activeReplicationNodes.length == 3);
    }

    function test_getActiveApiNodesIDs() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        _enableNodes(operators, nodeIds);

        uint256[] memory activeNodeIds = nodes.getActiveApiNodesIDs();
        vm.assertTrue(activeNodeIds.length == 3);
        vm.assertEq(activeNodeIds[0], nodeIds[0]);
        vm.assertEq(activeNodeIds[1], nodeIds[1]);
        vm.assertEq(activeNodeIds[2], nodeIds[2]);

        INodes.NodeWithId[] memory activeReplicationNodes = nodes.getActiveReplicationNodes();
        vm.assertTrue(activeReplicationNodes.length == 3);
        vm.assertEq(activeReplicationNodes[0].nodeId, nodeIds[0]);
        vm.assertEq(activeReplicationNodes[1].nodeId, nodeIds[1]);
        vm.assertEq(activeReplicationNodes[2].nodeId, nodeIds[2]);
    }

    function test_getActiveApiNodesCount() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);
        vm.assertEq(nodes.getActiveApiNodesCount(), 1);
        vm.assertEq(nodes.getActiveReplicationNodesCount(), 1);

        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(2);
        bool[] memory isActive = new bool[](2);
        for (uint256 i = 0; i < 2; i++) {
            isActive[i] = true;
        }
        _enableNodes(operators, nodeIds);
        vm.assertEq(nodes.getActiveApiNodesCount(), 3);
        vm.assertEq(nodes.getActiveReplicationNodesCount(), 3);

        (address[] memory operators2, uint256[] memory nodeIds2) = _addMultipleNodes(3);
        bool[] memory isActive2 = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive2[i] = true;
        }
        _enableNodes(operators2, nodeIds2);
        vm.assertEq(nodes.getActiveApiNodesCount(), 6);
        vm.assertEq(nodes.getActiveReplicationNodesCount(), 6);
        
        vm.prank(operators2[0]);
        nodes.setIsApiEnabled(nodeIds2[0], false);
        vm.assertEq(nodes.getActiveApiNodesCount(), 5);

        vm.prank(operators2[0]);
        nodes.setIsReplicationEnabled(nodeIds2[0], false);
        vm.assertEq(nodes.getActiveReplicationNodesCount(), 5);
    }

    function test_getNodeIsActive() public {
        _addNode();
        _enableNode(nodeOperator, nodeId);
        vm.assertEq(nodes.getApiNodeIsActive(nodeId), true);
        vm.assertEq(nodes.getReplicationNodeIsActive(nodeId), true);

        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(2);
        _enableNode(operators[0], nodeIds[0]);

        vm.assertEq(nodes.getApiNodeIsActive(nodeIds[0]), true);
        vm.assertEq(nodes.getApiNodeIsActive(nodeIds[1]), false);
        vm.assertEq(nodes.getReplicationNodeIsActive(nodeIds[0]), true);
        vm.assertEq(nodes.getReplicationNodeIsActive(nodeIds[1]), false);
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
        INodes.Node memory node = _randomNode();
        nodeOperator = vm.randomAddress();
        nodeId = nodes.addNode(nodeOperator, node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
    }

    function _addMultipleNodes(uint256 numberOfNodes) internal returns (address[] memory operators, uint256[] memory nodeIds) {
        operators = new address[](numberOfNodes);
        nodeIds = new uint256[](numberOfNodes);
        for (uint256 i = 0; i < numberOfNodes; i++) {
            INodes.Node memory node = _randomNode();
            operators[i] = vm.randomAddress();
            nodeIds[i] = nodes.addNode(operators[i], node.signingKeyPub, node.httpAddress, node.minMonthlyFeeMicroDollars);
        }
        return (operators, nodeIds);
    }

    function _enableNode(address operator, uint256 id) internal {
        vm.startPrank(operator);
        nodes.setIsReplicationEnabled(id, true);
        nodes.setIsApiEnabled(id, true);
        vm.stopPrank();
    }

    function _enableNodes(address[] memory operators, uint256[] memory nodeIds) internal {
        for (uint256 i = 0; i < nodeIds.length; i++) {
            vm.startPrank(operators[i]);
            nodes.setIsReplicationEnabled(nodeIds[i], true);
            nodes.setIsApiEnabled(nodeIds[i], true);
            vm.stopPrank();
        }
    }

    function _randomNode() internal view returns (INodes.Node memory) {
        return INodes.Node({
            signingKeyPub: _genBytes(32), 
            httpAddress: _genString(32), 
            isReplicationEnabled: false, 
            isApiEnabled: false, 
            isDisabled: false,
            minMonthlyFeeMicroDollars: _genRandomInt(100, 10000)
        });
    }

    function _checkNoLogsEmitted() internal {
        Vm.Log[] memory logs = vm.getRecordedLogs();
        vm.assertTrue(logs.length == 0, "No logs should be emitted");
    }
}
