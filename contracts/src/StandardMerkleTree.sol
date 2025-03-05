// SPDX-License-Identifier: MIT

pragma solidity ^0.8.28;

import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

contract StandardMerkleTree {
    error EmptyProof();
    error InvalidIndex();

    function multiProofVerify(
        bytes32[] calldata proof,
        bool[] calldata proofFlags,
        uint256 offset,
        bytes32 root,
        address[] calldata accounts,
        uint256[] calldata amounts,
        uint256[] calldata indices
    ) external pure returns (bool) {
        bytes32[] memory leaves = new bytes32[](accounts.length);
        for (uint256 i = 0; i < accounts.length; i++) {
            uint256 index = indices[i];
            // Make sure the index is within the acceptable range
            // Not sure if we have to check for duplicate indices as well. Might be fine,
            // so long as the tree wasn't crafted maliciously.
            if (index > offset + accounts.length || index < offset) {
                revert InvalidIndex();
            }
            leaves[i] = _leaf(indices[i], accounts[i], amounts[i]);
        }
        if (leaves.length == 0) {
            revert EmptyProof();
        }
        return MerkleProof.multiProofVerifyCalldata(proof, proofFlags, root, leaves);
    }

    function _leaf(uint256 index, address account, uint256 amount) internal pure returns (bytes32) {
        return keccak256(bytes.concat(keccak256(abi.encode(index, account, amount))));
    }
}
