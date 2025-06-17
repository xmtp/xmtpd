package utils

import (
	"errors"
	"math"
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
