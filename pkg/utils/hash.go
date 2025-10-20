package utils

import (
	"encoding/binary"
	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
)

func HashPayerSignatureInput(originatorID uint32, unsignedClientEnvelope []byte) []byte {
	targetBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(targetBytes, originatorID)
	return ethcrypto.Keccak256(
		[]byte(constants.TargetOriginatorDomainSeparationLabel),
		targetBytes,
		[]byte(constants.PayerDomainSeparationLabel),
		unsignedClientEnvelope,
	)
}

func HashJWTSignatureInput(textToSign []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.JWTDomainSeparationLabel),
		textToSign,
	)
}

func HashOriginatorSignatureInput(unsignedOriginatorEnvelope []byte) []byte {
	return ethcrypto.Keccak256(
		[]byte(constants.OriginatorDomainSeparationLabel),
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

func PackSortAndHashNodeIDs(nodeIDs []uint32) (common.Hash, error) {
	if !slices.IsSorted(nodeIDs) {
		sortedNodeIDs := slices.Clone(nodeIDs)
		slices.Sort(sortedNodeIDs)
		nodeIDs = sortedNodeIDs
	}

	t, err := abi.NewType("uint32[]", "", nil)
	if err != nil {
		return common.Hash{}, err
	}

	args := abi.Arguments{
		{Type: t},
	}

	encoded, err := args.Pack(nodeIDs)
	if err != nil {
		return common.Hash{}, err
	}

	return ethcrypto.Keccak256Hash(encoded), nil
}
