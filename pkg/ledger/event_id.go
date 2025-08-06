package ledger

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/core/types"
)

type EventID [32]byte

func (e EventID) String() string {
	return hex.EncodeToString(e[:])
}

// EventIDs are unique for an event in the context of a specific block
// In the case of reorgs, we want to reverse the original event and then
// if the transaction is included a new block, add that event to the ledger
// with a new event ID.
func BuildEventID(log types.Log) EventID {
	inputs := make([]byte, 68)
	copy(inputs[:32], log.TxHash[:])
	copy(inputs[32:64], log.BlockHash[:])
	copy(inputs[64:], convertIndexToBytes(log.Index))

	return sha256.Sum256(inputs)
}

func convertIndexToBytes(index uint) []byte {
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, uint32(index))
	return indexBytes
}
