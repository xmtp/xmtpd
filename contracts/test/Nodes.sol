// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Test} from "forge-std/src/Test.sol";
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

    /// @dev Use _addNodeSet to populate addresses and IDs.
    address node1Operator;
    address node2Operator;
    address node3Operator;
    uint256 node1Id;
    uint256 node2Id;
    uint256 node3Id;

    function setUp() public {
        nodes = new Nodes(admin);
        nodes.grantRole(nodes.NODE_MANAGER_ROLE(), manager);
    }

    function test_addNode() public {
        INodes.Node memory node = _randomNode();

        address operatorAddress = vm.randomAddress();

        uint256 nodeId = nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress, node.minMonthlyFee);

        vm.assertEq(nodes.ownerOf(nodeId), operatorAddress);
        vm.assertEq(nodes.getNode(nodeId).signingKeyPub, node.signingKeyPub);
        vm.assertEq(nodes.getNode(nodeId).httpAddress, node.httpAddress);
        vm.assertEq(nodes.getNode(nodeId).isActive, false);
        vm.assertEq(nodes.getNode(nodeId).isApiEnabled, false);
        vm.assertEq(nodes.getNode(nodeId).isReplicationEnabled, false);
        vm.assertEq(nodes.getNode(nodeId).minMonthlyFee, node.minMonthlyFee);
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

    function test_transferFrom() public {
        _addNodeSet();

        address newOwner = vm.randomAddress();

        // Note: The NFT holder must approve manager to transfer the NFT.
        // A transfer cannot happen without approval.
        vm.prank(node1Operator);
        ERC721(address(nodes)).approve(manager, node1Id);

        vm.prank(manager);
        nodes.transferFrom(node1Operator, newOwner, node1Id);
        vm.assertEq(nodes.ownerOf(node1Id), newOwner);
    }

    function test_RevertWhen_transferFromUnauthorized() public {
        _addNodeSet();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                unauthorized,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(unauthorized);
        nodes.transferFrom(node1Operator, newOwner, node1Id);
    }

    function test_RevertWhen_transferFromNotApproved() public {
        _addNodeSet();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IERC721Errors.ERC721InsufficientApproval.selector,
                manager,
                node1Id
            )
        );
        vm.prank(manager);
        nodes.transferFrom(node1Operator, newOwner, node1Id);
    }

    function test_RevertWhen_ownerCannotTransfer() public {
        _addNodeSet();

        address newOwner = vm.randomAddress();

        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector,
                node1Operator,
                nodes.NODE_MANAGER_ROLE()
            )
        );
        vm.prank(node1Operator);
        nodes.transferFrom(node1Operator, newOwner, node1Id);
    }

    // TODO: test safeTransfer and safeTransferFrom

    function test_allNodes() public {
        _addNodeSet();

        INodes.NodeWithId[] memory allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 3);
    }

    function test_activeNodes() public {
        _addNodeSet();

        nodes.updateActive(node1Id, true);
        nodes.updateActive(node2Id, true);
        nodes.updateActive(node3Id, true);

        uint256[] memory activeNodesIDs = nodes.getActiveNodesIDs();
        vm.assertTrue(activeNodesIDs.length == 3);
    }

    function _addNodeSet() internal {
        INodes.Node memory node1 = _randomNode();
        INodes.Node memory node2 = _randomNode();
        INodes.Node memory node3 = _randomNode();

        node1Operator = vm.randomAddress();
        node2Operator = vm.randomAddress();
        node3Operator = vm.randomAddress();

        node1Id = nodes.addNode(node1Operator, node1.signingKeyPub, node1.httpAddress, node1.minMonthlyFee);
        node2Id = nodes.addNode(node2Operator, node2.signingKeyPub, node2.httpAddress, node2.minMonthlyFee);
        node3Id = nodes.addNode(node3Operator, node3.signingKeyPub, node3.httpAddress, node3.minMonthlyFee);
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
