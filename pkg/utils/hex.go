package utils

import "encoding/hex"

func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}
