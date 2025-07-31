package utils

import "encoding/binary"

func EncodeUint64ToBytes32(v uint64) [32]byte {
	var b [32]byte
	binary.BigEndian.PutUint64(b[24:], v)
	return b
}

func DecodeBytes32ToUint64(b [32]byte) uint64 {
	return binary.BigEndian.Uint64(b[24:]) // last 8 bytes
}
