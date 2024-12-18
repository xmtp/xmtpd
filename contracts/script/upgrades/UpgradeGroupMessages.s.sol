// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/Vm.sol";
import "src/GroupMessages.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract UpgradeGroupMessages is Script {
    function run() external {
        address proxyAddress = vm.envAddress("XMTP_GROUP_MESSAGES_PROXY_ADDRESS");
        require(proxyAddress != address(0), "XMTP_GROUP_MESSAGES_PROXY_ADDRESS not set");

        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        address upgrader = vm.addr(privateKey);
        vm.startBroadcast(privateKey);

        // Step 1: Deploy the new implementation contract.
        GroupMessages newImplementation = new GroupMessages();

        // Step 2: Initialize the proxy.
        GroupMessages proxy = GroupMessages(proxyAddress);

        // Step 3: Upgrade the proxy pointer to the new implementation.
        proxy.upgradeToAndCall(address(newImplementation), "");

        console.log(
            '{"upgrader":"%s","proxy":"%s","newImplementation":"%s"}', 
            upgrader, 
            address(proxy), 
            address(newImplementation)
        );

        vm.stopBroadcast();
    }
}
