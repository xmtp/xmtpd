// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.28;

import {IAccessControl} from "@openzeppelin/contracts/access/IAccessControl.sol";
import {PausableUpgradeable} from "@openzeppelin-contracts-upgradeable/utils/PausableUpgradeable.sol";

contract Utils {
    bytes32 public constant EIP1967_IMPL_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 public constant ID = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;

    function _generatePayload(uint256 length) public pure returns (bytes memory) {
        bytes memory payload = new bytes(length);
        for (uint256 i = 0; i < payload.length; i++) {
            payload[i] = bytes1(uint8(i % 256));
        }
        return payload;
    }
}
