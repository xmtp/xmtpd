package utils

import "errors"

var (
	ErrSliceNot32Bytes = errors.New("slice must be 32 bytes long")
)

func SliceToArray32(slice []byte) ([32]byte, error) {
	if len(slice) != 32 {
		return [32]byte{}, ErrSliceNot32Bytes
	}

	var array [32]byte
	copy(array[:], slice)
	return array, nil
}
