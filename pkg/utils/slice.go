package utils

func SliceToArray32(slice []byte) [32]byte {
	var array [32]byte
	copy(array[:], slice)
	return array
}
