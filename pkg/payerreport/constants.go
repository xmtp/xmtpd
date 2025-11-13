package payerreport

import (
	"github.com/ethereum/go-ethereum/common"
)

// payerReportDigestTypeHash is the type hash as defined by EIP-712 for structured data hashing and signing.
// Calculated as the following keccak256 hash:
//
//	keccak256("PayerReport(uint32 originatorNodeId,uint64 startSequenceId,uint64 endSequenceId,uint32 endMinuteSinceEpoch,bytes32 payersMerkleRoot,uint32[] nodeIds)")
//
// Reference: https://github.com/xmtp/smart-contracts/blob/main/src/settlement-chain/PayerReportManager.sol#L29
var payerReportDigestTypeHash = common.HexToHash(
	"3737a2cced99bb28fc5aede45aa81d3ce0aa9137c5f417641835d0d71d303346",
)
