// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Script, console} from "forge-std/src/Script.sol";
import {Environment} from "./utils/Environment.sol";
import {Utils} from "./utils/Utils.sol";
import "src/interfaces/INodes.sol";
import "src/Nodes.sol";

contract DeployXMTPNodeRegistry is Script, Environment, Utils {
    Nodes nodes;

    address admin;
    address deployer;

    function run() public {
        admin = vm.envAddress("XMTP_NODE_REGISTRY_ADMIN_ADDRESS");
        require(admin != address(0), "XMTP_NODE_REGISTRY_ADMIN_ADDRESS not set");

        uint256 privateKey = vm.envUint("PRIVATE_KEY");
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

        string memory addressesOutput;

        addressesOutput = vm.serializeAddress(addresses, "XMTPNodeRegistryDeployer", deployer);
        addressesOutput = vm.serializeAddress(addresses, "XMTPNodeRegistryInitialAdmin", admin);
        addressesOutput = vm.serializeAddress(addresses, "XMTPNodeRegistry", address(nodes));

        string memory finalJson;
        finalJson = vm.serializeString(parent_object, addresses, addressesOutput);
        finalJson = vm.serializeUint(parent_object, "deploymentBlock", block.number);
        finalJson = vm.serializeUint(parent_object, "latestUpgradeBlock", block.number);

        writeOutput(finalJson, XMTP_NODE_REGISTRY_OUTPUT_JSON);
    }
}
