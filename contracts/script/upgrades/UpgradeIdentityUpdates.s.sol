// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/Vm.sol";
import "../utils/Utils.sol";
import "../utils/Environment.sol";
import "src/IdentityUpdates.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract UpgradeIdentityUpdates is Script, Utils, Environment {
    IdentityUpdates newImplementation;
    IdentityUpdates proxy;

    address upgrader;

    function run() external {
        uint256 privateKey = vm.envUint("PRIVATE_KEY");
        upgrader = vm.addr(privateKey);

        vm.startBroadcast(privateKey);

        _initializeProxy();

        // Deploy the new implementation contract.
        newImplementation = new IdentityUpdates();
        require(address(newImplementation) != address(0), "Implementation deployment failed");

        // Upgrade the proxy pointer to the new implementation.
        proxy.upgradeToAndCall(address(newImplementation), "");

        vm.stopBroadcast();

        _serializeUpgradeData();
    }

    function _initializeProxy() internal {
        string memory fileContent = readOutput(XMTP_IDENTITY_UPDATES_OUTPUT_JSON);
        proxy = IdentityUpdates(stdJson.readAddress(fileContent, ".addresses.identityUpdatesProxy"));
        require(address(proxy) != address(0), "proxy address not set");
        require(
            proxy.hasRole(proxy.DEFAULT_ADMIN_ROLE(), upgrader),
            "Upgrader must have admin role"
        );
    }

    function _serializeUpgradeData() internal {
        vm.writeJson(
            vm.toString(address(newImplementation)),
            getOutputPath(XMTP_IDENTITY_UPDATES_OUTPUT_JSON),
            ".addresses.identityUpdatesImpl"
        );
        vm.writeJson(
            vm.toString(block.number),
            getOutputPath(XMTP_IDENTITY_UPDATES_OUTPUT_JSON),
            ".latestUpgradeBlock"
        );
    }
}
