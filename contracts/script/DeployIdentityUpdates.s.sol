// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/Vm.sol";
import "./utils/Utils.sol";
import "./utils/Environment.sol";
import "src/IdentityUpdates.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployIdentityUpdates is Script, Utils, Environment {
    IdentityUpdates idUpdatesImpl;
    ERC1967Proxy proxy;

    address admin;
    address deployer;

    function run() external {
        admin = vm.envAddress("XMTP_IDENTITY_UPDATES_ADMIN_ADDRESS");
        require(admin != address(0), "XMTP_IDENTITY_UPDATES_ADMIN_ADDRESS not set");

        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        require(privateKey != 0, "PRIVATE_KEY not set");

        deployer = vm.addr(privateKey);
        vm.startBroadcast(privateKey);

        // Deploy the implementation contract.
        idUpdatesImpl = new IdentityUpdates();
        require(address(idUpdatesImpl) != address(0), "Implementation deployment failed");

        // Deploy the proxy contract.
        proxy =
            new ERC1967Proxy(address(idUpdatesImpl), abi.encodeWithSelector(IdentityUpdates.initialize.selector, admin));

        vm.stopBroadcast();

        _serializeDeploymentData();
    }

    function _serializeDeploymentData() internal {
        string memory parent_object = "parent object";
        string memory addresses = "addresses";

        string memory addressesOutput;

        addressesOutput = vm.serializeAddress(addresses, "identityUpdatesDeployer", deployer);
        addressesOutput = vm.serializeAddress(addresses, "identityUpdatesProxyAdmin", admin);
        addressesOutput = vm.serializeAddress(addresses, "identityUpdatesProxy", address(proxy));
        addressesOutput = vm.serializeAddress(addresses, "identityUpdatesImpl", address(idUpdatesImpl));

        string memory finalJson;
        finalJson = vm.serializeString(parent_object, addresses, addressesOutput);
        finalJson = vm.serializeUint(parent_object, "deploymentBlock", block.number);
        finalJson = vm.serializeUint(parent_object, "latestUpgradeBlock", block.number);

        writeOutput(finalJson, XMTP_IDENTITY_UPDATES_OUTPUT_JSON);
    }
}
