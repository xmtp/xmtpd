// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import {Test, console} from "forge-std/src/Test.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Nodes} from "../src/Nodes.sol";

contract NodesTest is Test {
    Nodes public nodes;

    function setUp() public {
        nodes = new Nodes();
    }

    function _genBytes(uint32 length) internal pure returns (bytes memory) {
        bytes memory message = new bytes(length);
        for (uint256 i = 0; i < length; i++) {
            message[i] = bytes1(uint8(i % 256));
        }

        return message;
    }

    function _genString(uint32 length) internal pure returns (string memory) {
        return string(_genBytes(length));
    }

    function _randomNode(bool isHealthy) internal pure returns (Nodes.Node memory) {
        return Nodes.Node({signingKeyPub: _genBytes(32), httpAddress: _genString(32), isHealthy: isHealthy});
    }

    function test_canAddNode() public {
        Nodes.Node memory node = _randomNode(true);

        address operatorAddress = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operatorAddress, node.signingKeyPub, node.httpAddress);

        vm.assertEq(nodes.ownerOf(nodeId), operatorAddress);
        vm.assertEq(nodes.getNode(nodeId).signingKeyPub, node.signingKeyPub);
        vm.assertEq(nodes.getNode(nodeId).httpAddress, node.httpAddress);
        vm.assertEq(nodes.getNode(nodeId).isHealthy, true);
    }

    function test_increments100() public {
        Nodes.Node memory node1 = _randomNode(true);
        Nodes.Node memory node2 = _randomNode(true);
        Nodes.Node memory node3 = _randomNode(true);

        address operator1 = vm.randomAddress();
        address operator2 = vm.randomAddress();
        address operator3 = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator1, node1.signingKeyPub, node1.httpAddress);
        vm.assertTrue(nodeId == 100);

        nodeId = nodes.addNode(operator2, node2.signingKeyPub, node2.httpAddress);
        vm.assertTrue(nodeId == 200);

        nodeId = nodes.addNode(operator3, node3.signingKeyPub, node3.httpAddress);
        vm.assertTrue(nodeId == 300);
    }

    function test_canAddMultiple() public {
        Nodes.Node memory node1 = _randomNode(true);
        Nodes.Node memory node2 = _randomNode(true);
        Nodes.Node memory node3 = _randomNode(true);

        address operator1 = vm.randomAddress();
        address operator2 = vm.randomAddress();
        address operator3 = vm.randomAddress();

        uint32 node1Id = nodes.addNode(operator1, node1.signingKeyPub, node1.httpAddress);
        nodes.addNode(operator2, node2.signingKeyPub, node2.httpAddress);
        nodes.addNode(operator3, node3.signingKeyPub, node3.httpAddress);

        Nodes.NodeWithId[] memory allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 3);

        Nodes.NodeWithId[] memory healthyNodes = nodes.healthyNodes();
        vm.assertTrue(healthyNodes.length == 3);

        nodes.updateHealth(node1Id, false);
        allNodes = nodes.allNodes();
        vm.assertTrue(allNodes.length == 3);
        healthyNodes = nodes.healthyNodes();
        vm.assertTrue(healthyNodes.length == 2);
    }

    function test_canMarkUnhealthy() public {
        Nodes.Node memory node = _randomNode(true);
        address operator = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator, node.signingKeyPub, node.httpAddress);

        nodes.updateHealth(nodeId, false);

        vm.assertEq(nodes.getNode(nodeId).isHealthy, false);
        vm.assertEq(nodes.healthyNodes().length, 0);
    }

    function testFail_ownerCannotUpdateHealth() public {
        vm.expectRevert(Ownable.OwnableUnauthorizedAccount.selector);
        Nodes.Node memory node = _randomNode(true);
        address operator = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator, node.signingKeyPub, node.httpAddress);

        vm.prank(operator);
        nodes.updateHealth(nodeId, false);
    }

    function testFail_ownerCannotTransfer() public {
        Nodes.Node memory node = _randomNode(true);
        address operator = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator, node.signingKeyPub, node.httpAddress);

        vm.prank(operator);
        nodes.safeTransferFrom(operator, vm.randomAddress(), uint256(nodeId));
    }

    function test_canChangeHttpAddress() public {
        Nodes.Node memory node = _randomNode(true);
        address operator = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator, node.signingKeyPub, node.httpAddress);

        vm.prank(operator);
        nodes.updateHttpAddress(nodeId, "new-http-address");

        vm.assertEq(nodes.getNode(nodeId).httpAddress, "new-http-address");
    }

    function testFail_cannotChangeOtherHttpAddress() public {
        vm.expectRevert("Only the owner of the Node NFT can update its http address");

        Nodes.Node memory node = _randomNode(true);
        address operator = vm.randomAddress();

        uint32 nodeId = nodes.addNode(operator, node.signingKeyPub, node.httpAddress);

        vm.prank(vm.randomAddress());
        nodes.updateHttpAddress(nodeId, "new-http-address");
    }
}
