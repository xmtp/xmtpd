package utils

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
)

var nodeIdArrayType = abi.Arguments{
	{
		Name: "activeNodeIDs",
		Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.UintTy, Size: 32}},
	},
}

func HashPayerSignatureInput(originatorID uint32, unsignedClientEnvelope []byte) []byte {
	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, originatorID)
	return ethcrypto.Keccak256(
		[]byte(constants.TARGET_ORIGINATOR_SEPARATION_LABEL),
		targetBytes,
		[]byte(constants.PAYER_DOMAIN_SEPARATION_LABEL),
		unsignedClientEnvelope,
	)
}

func HashJWTSignatureInput(textToSign []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.JWT_DOMAIN_SEPARATION_LABEL),
		textToSign,
	)
}

func HashOriginatorSignatureInput(unsignedOriginatorEnvelope []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.ORIGINATOR_DOMAIN_SEPARATION_LABEL),
		unsignedOriginatorEnvelope,
	)
}

func HashPayerReportInput(packedBytes []byte, domainSeparator common.Hash) common.Hash {
	return common.BytesToHash(ethcrypto.Keccak256(
		[]byte("\x19\x01"),
		domainSeparator[:],
		ethcrypto.Keccak256(packedBytes),
	))
}

// Takes a slice of uint32 and returns them packed as 32 BYTE elements that are left padded
func encodePackedUint32Slice(input []uint32) []byte {
	// Convert each uint32 to a 32-byte element, left-padded with zeros
	result := make([]byte, len(input)*32)

	for i, val := range input {
		// Create a 32-byte element with the uint32 value left-padded
		offset := i * 32
		// The uint32 value will be placed at the end of the 32-byte element
		binary.BigEndian.PutUint32(result[offset+28:offset+32], val)
		// The rest of the bytes are already initialized to zero
	}

	return result
}

func PackAndHashNodeIDs(nodeIDs []uint32) common.Hash {
	packed := encodePackedUint32Slice(nodeIDs)

	return common.BytesToHash(ethcrypto.Keccak256(packed))
}
