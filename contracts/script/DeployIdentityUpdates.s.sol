// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/Vm.sol";
import "src/IdentityUpdates.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployIdentityUpdates is Script {
    function run() external {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(privateKey);
        vm.startBroadcast(privateKey);

        // Step 1: Deploy the implementation contract
        IdentityUpdates idUpdatesImpl = new IdentityUpdates();

        // Step 2: Deploy the proxy contract
        ERC1967Proxy proxy =
            new ERC1967Proxy(
                address(idUpdatesImpl), 
                abi.encodeWithSelector(IdentityUpdates.initialize.selector, deployer)
        );

        console.log(
            '{"deployer":"%s","proxy":"%s","implementation":"%s"}', 
            deployer, 
            address(proxy), 
            address(idUpdatesImpl)
        );

        vm.stopBroadcast();
    }
}
