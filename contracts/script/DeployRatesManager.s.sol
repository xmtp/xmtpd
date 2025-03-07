// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/Vm.sol";
import "./utils/Utils.sol";
import "./utils/Environment.sol";
import "../src/RatesManager.sol";

import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract DeployRatesManager is Script, Utils, Environment {
    RatesManager ratesManagerImpl;
    ERC1967Proxy proxy;

    address admin;
    address deployer;

    function run() external {
        admin = vm.envAddress("XMTP_RATES_MANAGER_ADMIN_ADDRESS");
        require(admin != address(0), "XMTP_RATES_MANAGER_ADMIN_ADDRESS not set");

        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        require(privateKey != 0, "PRIVATE_KEY not set");

        deployer = vm.addr(privateKey);
        vm.startBroadcast(privateKey);

        // Deploy the implementation contract.
        ratesManagerImpl = new RatesManager();
        require(address(ratesManagerImpl) != address(0), "Implementation deployment failed");

        // Deploy the proxy contract.
        proxy =
            new ERC1967Proxy(address(ratesManagerImpl), abi.encodeWithSelector(RatesManager.initialize.selector, admin));

        vm.stopBroadcast();

        _serializeDeploymentData();
    }

    function _serializeDeploymentData() internal {
        string memory parent_object = "parent object";
        string memory addresses = "addresses";

        string memory addressesOutput;

        addressesOutput = vm.serializeAddress(addresses, "ratesManagerDeployer", deployer);
        addressesOutput = vm.serializeAddress(addresses, "ratesManagerProxyAdmin", admin);
        addressesOutput = vm.serializeAddress(addresses, "ratesManagerProxy", address(proxy));
        addressesOutput = vm.serializeAddress(addresses, "ratesManagerImpl", address(ratesManagerImpl));

        string memory finalJson;
        finalJson = vm.serializeString(parent_object, addresses, addressesOutput);
        finalJson = vm.serializeUint(parent_object, "deploymentBlock", block.number);
        finalJson = vm.serializeUint(parent_object, "latestUpgradeBlock", block.number);

        writeOutput(finalJson, XMTP_RATES_MANAGER_OUTPUT_JSON);
    }
}
