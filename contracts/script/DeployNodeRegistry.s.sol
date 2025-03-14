// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { Nodes } from "../src/Nodes.sol";

import { Utils } from "./utils/Utils.sol";
import { Environment } from "./utils/Environment.sol";

contract DeployXMTPNodeRegistry is Utils, Environment {
    Nodes nodes;

    address admin;
    address deployer;

    function run() public {
        admin = vm.envAddress("XMTP_NODE_REGISTRY_ADMIN_ADDRESS");
        require(admin != address(0), "XMTP_NODE_REGISTRY_ADMIN_ADDRESS not set");

        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        require(privateKey != 0, "PRIVATE_KEY not set");

        deployer = vm.addr(privateKey);
        vm.startBroadcast(privateKey);

        nodes = new Nodes(admin);
        require(address(nodes) != address(0), "Nodes deployment failed");

        vm.stopBroadcast();

        _serializeDeploymentData();
    }

    function _serializeDeploymentData() internal {
        string memory parent_object = "parent object";
        string memory addresses = "addresses";
        string memory constructorArgs = "constructorArgs";

        string memory addressesOutput;
        addressesOutput = vm.serializeAddress(addresses, "deployer", deployer);
        addressesOutput = vm.serializeAddress(addresses, "implementation", address(nodes));

        string memory constructorArgsOutput = vm.serializeAddress(constructorArgs, "initialAdmin", admin);

        string memory finalJson;
        finalJson = vm.serializeString(parent_object, addresses, addressesOutput);
        finalJson = vm.serializeString(parent_object, constructorArgs, constructorArgsOutput);
        finalJson = vm.serializeUint(parent_object, "deploymentBlock", block.number);

        writeOutput(finalJson, XMTP_NODE_REGISTRY_OUTPUT_JSON);
    }
}
