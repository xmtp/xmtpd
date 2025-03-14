// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import { stdJson } from "forge-std/src/StdJson.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import { GroupMessages } from "../../src/GroupMessages.sol";

import { Utils } from "../utils/Utils.sol";
import { Environment } from "../utils/Environment.sol";

contract UpgradeGroupMessages is Utils, Environment {
    GroupMessages newImplementation;
    GroupMessages proxy;

    address upgrader;

    function run() external {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        upgrader = vm.addr(privateKey);

        vm.startBroadcast(privateKey);

        _initializeProxy();

        // Deploy the new implementation contract.
        newImplementation = new GroupMessages();
        require(address(newImplementation) != address(0), "Implementation deployment failed");

        // Upgrade the proxy pointer to the new implementation.
        proxy.upgradeToAndCall(address(newImplementation), "");

        vm.stopBroadcast();

        _serializeUpgradeData();
    }

    function _initializeProxy() internal {
        string memory fileContent = readOutput(XMTP_GROUP_MESSAGES_OUTPUT_JSON);
        proxy = GroupMessages(stdJson.readAddress(fileContent, ".addresses.proxy"));
        require(address(proxy) != address(0), "proxy address not set");
        require(proxy.hasRole(proxy.DEFAULT_ADMIN_ROLE(), upgrader), "Upgrader must have admin role");
    }

    function _serializeUpgradeData() internal {
        vm.writeJson(
            vm.toString(address(newImplementation)),
            getOutputPath(XMTP_GROUP_MESSAGES_OUTPUT_JSON),
            ".addresses.implementation"
        );
        vm.writeJson(vm.toString(block.number), getOutputPath(XMTP_GROUP_MESSAGES_OUTPUT_JSON), ".latestUpgradeBlock");
    }
}
