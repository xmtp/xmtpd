// SPDX-License-Identifier: MIT
pragma solidity 0.8.28;

contract Utils {
    bytes32 public constant EIP1967_IMPLEMENTATION_SLOT =
        0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 public constant ID = 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef;

    function _generatePayload(uint256 length) public pure returns (bytes memory payload) {
        payload = new bytes(length);

        for (uint256 i; i < payload.length; ++i) {
            payload[i] = bytes1(uint8(i % 256));
        }
    }

    /// @dev This is NOT cryptographically secure. Just good enough for testing.
    function _genRandomInt(uint256 min, uint256 max) internal view returns (uint256) {
        return min
            + uint256(keccak256(abi.encodePacked(block.timestamp, block.prevrandao, block.number, msg.sender)))
                % (max - min + 1);
    }

    function _genBytes(uint32 length) internal pure returns (bytes memory message) {
        message = new bytes(length);

        for (uint256 i; i < length; ++i) {
            message[i] = bytes1(uint8(i % 256));
        }
    }

    function _genString(uint32 length) internal pure returns (string memory) {
        return string(_genBytes(length));
    }
}
