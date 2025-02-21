// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/src/Test.sol";
import {Vm} from "forge-std/src/Vm.sol";
import {console2} from "forge-std/src/console2.sol";
import {Utils} from "test/utils/Utils.sol";
import {Nodes} from "src/Nodes.sol";
import {INodes} from "src/interfaces/INodes.sol";
import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {AccessControl} from "@openzeppelin/contracts/access/AccessControl.sol";
import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {IERC721Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";

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
        console2.log("admin address", admin);
        console2.log("manager address", manager);
        console2.log("unauthorized address", unauthorized);
    }

    // ***************************************************************
    // *                        addNodes                             *
    // ***************************************************************

    function test_addNode() public {
        INodes.Node memory node = _randomNode();

        address operatorAddress = vm.randomAddress();

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
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidAddress.selector);
        nodes.addNode(address(0), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
    }

    function test_RevertWhen_AddNodeWithInvalidSigningKey() public {
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidSigningKey.selector);
        nodes.addNode(vm.randomAddress(), bytes(""), node.httpAddress, node.minMonthlyFee);
    }

    function test_RevertWhen_AddNodeWithInvalidHttpAddress() public {
        INodes.Node memory node = _randomNode();
        vm.expectRevert(INodes.InvalidHttpAddress.selector);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, "", node.minMonthlyFee);
    }

    function test_RevertWhen_AddNodeUnauthorized() public {
        INodes.Node memory node = _randomNode();

        // Addresses without DEFAULT_ADMIN_ROLE cannot add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);

        // NODE_MANAGER_ROLE is not authorized to add nodes.
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                manager,
                nodes.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(manager);
        nodes.addNode(vm.randomAddress(), node.signingKeyPub, node.httpAddress, node.minMonthlyFee);
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
        nodes.updateHttpAddress(nodeId, "http://example.com");
        vm.assertEq(nodes.getNode(nodeId).httpAddress, "http://example.com");
    }

    function test_RevertWhen_updateHttpAddressNodeDoesNotExist() public {
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateHttpAddress(1337, "http://example.com");
    }

    function test_RevertWhen_updateHttpAddressInvalidHttpAddress() public {
        _addNode();
        vm.expectRevert(INodes.InvalidHttpAddress.selector);
        nodes.updateHttpAddress(nodeId, "");
    }

    function test_RevertWhen_updateHttpAddressUnauthorized() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateHttpAddress(nodeId, "");
    }

    function test_RevertWhen_updateHttpAddressOwnerCannotUpdate() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateHttpAddress(nodeId, "");
    }

    // ***************************************************************
    // *                  updateIsReplicationEnabled                 *
    // ***************************************************************

    function test_updateIsReplicationEnabled() public {
        _addNode();
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
        nodes.updateIsReplicationEnabled(nodeId, true);
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, true);
    }

    function test_RevertWhen_updateIsReplicationEnabledNodeDoesNotExist() public {
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateIsReplicationEnabled(1337, true);
    }

    function test_RevertWhen_updateIsReplicationEnabledUnauthorized() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateIsReplicationEnabled(nodeId, true);
    }

    function test_RevertWhen_updateIsReplicationEnabledOwnerCannotUpdate() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateIsReplicationEnabled(nodeId, true);
    }

    // ***************************************************************
    // *                        updateMinMonthlyFee                  *
    // ***************************************************************

    function test_updateMinMonthlyFee() public {
        _addNode();
        uint256 initialMonthlyFee = nodes.getNode(nodeId).minMonthlyFee;
        nodes.updateMinMonthlyFee(nodeId, 1000);
        vm.assertEq(nodes.getNode(nodeId).minMonthlyFee, 1000);
        vm.assertNotEq(nodes.getNode(nodeId).minMonthlyFee, initialMonthlyFee);
    }

    function test_RevertWhen_updateMinMonthlyFeeNodeDoesNotExist() public {
        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateMinMonthlyFee(1337, 1000);
    }

    function test_RevertWhen_updateMinMonthlyFeeUnauthorized() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.updateMinMonthlyFee(nodeId, 1000);
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
    // *                        allNodes                             *
    // ***************************************************************

    function test_allNodes() public {
        _addNode();

        INodes.NodeWithId[] memory allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 1);
    }

    // ***************************************************************
    // *            updateActive, batchUpdateActive                  *
    // ***************************************************************

    function test_updateActive() public {
        _addNode();
        vm.expectEmit(address(nodes));
        emit INodes.NodeActivateUpdated(nodeId, true);
        nodes.updateActive(nodeId, true);
        vm.assertEq(nodes.getNode(nodeId).isActive, true);
    }

    function test_RevertWhen_updateActiveNodeDoesNotExist() public {
        vm.recordLogs();

        vm.expectRevert(INodes.NodeDoesNotExist.selector);
        nodes.updateActive(1337, true);

        Vm.Log[] memory logs = vm.getRecordedLogs();
        vm.assertTrue(logs.length == 0, "No logs should be emitted");
    }

    function test_RevertWhen_updateActiveUnauthorized() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                manager,
                nodes.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(manager);
        nodes.updateActive(nodeId, true);
    }

    function test_RevertWhen_updateActiveOwnerCannotUpdate() public {
        _addNode();
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                nodeOperator,
                nodes.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(nodeOperator);
        nodes.updateActive(nodeId, true);
    }

    function test_RevertWhen_updateActiveNodeAlreadyActive() public {
        _addNode();
        nodes.updateActive(nodeId, true);
        vm.expectRevert(INodes.NodeAlreadyActive.selector);
        nodes.updateActive(nodeId, true);
    }

    function test_RevertWhen_updateActiveNodeAlreadyInactive() public {
        _addNode();
        vm.expectRevert(INodes.NodeAlreadyInactive.selector);
        nodes.updateActive(nodeId, false);
    }

    function test_batchUpdateActive() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);

        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        vm.expectEmit(address(nodes));
        emit INodes.NodeActivateUpdated(nodeIds[0], true);
        emit INodes.NodeActivateUpdated(nodeIds[1], true);
        emit INodes.NodeActivateUpdated(nodeIds[2], true);
        nodes.batchUpdateActive(nodeIds, isActive);

        uint256[] memory activeNodesIDs = nodes.getActiveNodesIDs();
        vm.assertTrue(activeNodesIDs.length == 3);
    }

    function test_RevertWhen_batchUpdateActiveUnauthorized() public {
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.DEFAULT_ADMIN_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.batchUpdateActive(nodeIds, isActive);
    }

    function test_RevertWhen_batchUpdateActiveInvalidInputLength() public {
        nodes.updateMaxActiveNodes(2);
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](2);
        for (uint256 i = 0; i < 2; i++) {
            isActive[i] = true;
        }

        vm.expectRevert(INodes.InvalidInputLength.selector);
        nodes.batchUpdateActive(nodeIds, isActive);
    }

    function test_RevertWhen_batchUpdateActiveMaxActiveNodesReached() public {
        nodes.updateMaxActiveNodes(2);
        (address[] memory operators, uint256[] memory nodeIds) = _addMultipleNodes(3);
        bool[] memory isActive = new bool[](3);
        for (uint256 i = 0; i < 3; i++) {
            isActive[i] = true;
        }

        vm.expectRevert(INodes.MaxActiveNodesReached.selector);
        nodes.batchUpdateActive(nodeIds, isActive);
    }

    // ***************************************************************
    // *                        Helper functions                     *
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
}
