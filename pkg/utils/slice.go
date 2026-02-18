package utils

import (
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/common"
)

var ErrSliceNot32Bytes = errors.New("slice must be 32 bytes long")

func SliceToArray32(slice []byte) ([32]byte, error) {
	if len(slice) != 32 {
		return [32]byte{}, ErrSliceNot32Bytes
	}

	var array [32]byte
	copy(array[:], slice)
	return array, nil
}

func Uint32SliceToInt32Slice(slice []uint32) ([]int32, error) {
	intSlice := make([]int32, len(slice))
	for i, v := range slice {
		if v > math.MaxInt32 {
			return nil, ErrIntOverflow
		}
		intSlice[i] = int32(v)
	}

	return intSlice, nil
}

func Int32SliceToUint32Slice(slice []int32) []uint32 {
	uintSlice := make([]uint32, len(slice))
	for i, v := range slice {
		uintSlice[i] = uint32(v)
	}

	return uintSlice
}

// ChunkSlice splits a slice into chunks of at most size elements.
func ChunkSlice[T any](s []T, size int) [][]T {
	if size <= 0 {
		return [][]T{s}
	}
	var chunks [][]T
	for i := 0; i < len(s); i += size {
		end := min(i+size, len(s))
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func AddressTo32Slice(addr common.Address) [32]byte {
	var result [32]byte
	copy(result[32-len(addr.Bytes()):], addr.Bytes())
	return result
}

func SliceToSet[T comparable](vals []T) map[T]struct{} {
	set := make(map[T]struct{}, len(vals))
	for _, v := range vals {
		set[v] = struct{}{}
	}
	return set
}
