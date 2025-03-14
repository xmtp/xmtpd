// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

import {Script, console} from "forge-std/src/Script.sol";

/*
 * @dev   Compute the storage slot for a given namespace.
 *        forge script script/ComputeStorageSlot.sol --sig "compute(string)" "xmtp.storage.Payer"
 */
contract ComputeStorageSlot is Script {
    function compute(string calldata namespace) external view {
        bytes32 namespaceHash = keccak256(bytes(namespace));
        uint256 namespaceInt = uint256(namespaceHash);
        bytes32 slot = bytes32((namespaceInt - 1) & ~uint256(0xff));

        console.log("Namespace:", namespace);
        console.logBytes32(namespaceHash);
        console.log("Storage Slot:");
        console.logBytes32(slot);
    }

    function run() external view {
        this.compute("");
    }
}
