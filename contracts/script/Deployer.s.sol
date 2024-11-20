// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import {Script, console} from "forge-std/Script.sol";
import "../src/Nodes.sol";

contract Deployer is Script {
    function setUp() public {}

    function run() public {
        vm.startBroadcast();
        new Nodes();

        vm.broadcast();
    }
}
