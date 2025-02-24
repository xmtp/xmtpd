// SPDX-License-Identifier: MIT
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

    /// @dev This is NOT cryptographically secure. Just good enough for testing.
    function _genRandomInt(uint256 min, uint256 max) internal view returns (uint256) {
        return min + uint256(
            keccak256(
            abi.encodePacked(
                    block.timestamp,
                    block.prevrandao,
                    block.number,
                    msg.sender
                )
            )
        ) % (max - min + 1);
    }

    function _genBytes(uint32 length) internal pure returns (bytes memory) {
        bytes memory message = new bytes(length);
        for (uint256 i = 0; i < length; i++) {
            message[i] = bytes1(uint8(i % 256));
        }

        return message;
    }

    function _genString(uint32 length) internal pure returns (string memory) {
        return string(_genBytes(length));
    }
}
