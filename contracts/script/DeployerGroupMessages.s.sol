// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "src/GroupMessages.sol";
import { ERC1967Proxy } from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployProxiedGroupMessages is Script {
    function run() external {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(privateKey);

        // Step 1: Deploy the implementation contract
        GroupMessages groupMessagesImpl = new GroupMessages();

        // Step 2: Deploy the proxy contract
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(groupMessagesImpl),
            abi.encodeWithSelector(GroupMessages.initialize.selector)
        );

        // Log the deployed contract addresses
        console.log("Implementation Address:", address(groupMessagesImpl));
        console.log("Proxy Address:", address(proxy));

        vm.stopBroadcast();
    }
}
