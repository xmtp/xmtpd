package utils

import "encoding/hex"

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func AssertHexDecode(s string) []byte {
	data, err := HexDecode(s)
	if err != nil {
		panic(err)
	}
	return data
}
