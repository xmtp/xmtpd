// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import "forge-std/src/Script.sol";
import "forge-std/src/StdJson.sol";

contract Utils is Script {
    uint256 constant CHAIN_ID_ANVIL_LOCALNET = 31337;
    uint256 constant CHAIN_ID_XMTP_TESTNET = 241320161;
    uint256 constant CHAIN_ID_BASE_SEPOLIA = 84532;

    string constant OUTPUT_ANVIL_LOCALNET = "anvil_localnet";
    string constant OUTPUT_XMTP_TESTNET = "xmtp_testnet";
    string constant OUTPUT_BASE_SEPOLIA = "base_sepolia";
    string constant OUTPUT_UNKNOWN = "unknown";

    function readInput(string memory inputFileName) internal view returns (string memory) {
        string memory file = getInputPath(inputFileName);
        return vm.readFile(file);
    }

    function getInputPath(string memory inputFileName) internal view returns (string memory) {
        string memory inputDir = string.concat(vm.projectRoot(), "/config/");
        string memory chainDir = string.concat(_resolveChainID(), "/");
        string memory file = string.concat(inputFileName, ".json");
        return string.concat(inputDir, chainDir, file);
    }

    function readOutput(string memory outputFileName) internal view returns (string memory) {
        string memory file = getOutputPath(outputFileName);
        return vm.readFile(file);
    }

    function writeOutput(string memory outputJson, string memory outputFileName) internal {
        string memory outputFilePath = getOutputPath(outputFileName);
        vm.writeJson(outputJson, outputFilePath);
    }

    function getOutputPath(string memory outputFileName) internal view returns (string memory) {
        string memory outputDir = string.concat(vm.projectRoot(), "/config/");
        string memory chainDir = string.concat(_resolveChainID(), "/");
        string memory outputFilePath = string.concat(outputDir, chainDir, outputFileName, ".json");
        return outputFilePath;
    }

    function _resolveChainID() internal view returns (string memory) {
        uint256 chainID = block.chainid;
        if (chainID == CHAIN_ID_ANVIL_LOCALNET) {
            return OUTPUT_ANVIL_LOCALNET;
        } else if (chainID == CHAIN_ID_XMTP_TESTNET) {
            return OUTPUT_XMTP_TESTNET;
        } else if (chainID == CHAIN_ID_BASE_SEPOLIA) {
            return OUTPUT_BASE_SEPOLIA;
        } else {
            return OUTPUT_UNKNOWN;
        }
    }
}
