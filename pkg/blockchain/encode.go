package blockchain

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var uint96Size = 12

// Pack helpers (right-aligned, big-endian in bytes32)

func packUint8(v uint8) [32]byte {
	var out [32]byte
	out[31] = v
	return out
}

func packUint16(v uint16) [32]byte {
	var out [32]byte
	out[30] = byte(v >> 8)
	out[31] = byte(v)
	return out
}

func packUint32(v uint32) [32]byte {
	var out [32]byte
	out[28] = byte(v >> 24)
	out[29] = byte(v >> 16)
	out[30] = byte(v >> 8)
	out[31] = byte(v)
	return out
}

func packUint64(v uint64) [32]byte {
	var out [32]byte
	out[24] = byte(v >> 56)
	out[25] = byte(v >> 48)
	out[26] = byte(v >> 40)
	out[27] = byte(v >> 32)
	out[28] = byte(v >> 24)
	out[29] = byte(v >> 16)
	out[30] = byte(v >> 8)
	out[31] = byte(v)
	return out
}

func packUint96Big(v *big.Int) ([32]byte, error) {
	var out [32]byte
	if v == nil {
		return out, fmt.Errorf("uint96: nil value")
	}
	if v.Sign() < 0 {
		return out, fmt.Errorf("uint96: negative value %s", v.String())
	}
	if v.BitLen() > 96 {
		return out, fmt.Errorf("uint96: overflow (%s > 2^96-1)", v.String())
	}
	b := v.Bytes()
	n := len(b)
	if n == 0 {
		return out, nil
	}
	if n > uint96Size {
		return out, fmt.Errorf("uint96: overflow (%d bytes > 12)", n)
	}
	copy(out[32-uint96Size+(uint96Size-n):], b)
	return out, nil
}

func packAddress(a common.Address) [32]byte {
	var out [32]byte
	copy(out[12:], a.Bytes())
	return out
}

func packBool(b bool) [32]byte {
	var out [32]byte
	if b {
		out[31] = 1
	}
	return out
}

// Decode helpers (enforce canonical zero-prefix)

func decodeUint8(val [32]byte) (uint8, error) {
	for i := 0; i < 31; i++ {
		if val[i] != 0 {
			return 0, fmt.Errorf("non-canonical uint8 encoding in bytes32 (non-zero prefix)")
		}
	}
	return val[31], nil
}

func decodeUint16(val [32]byte) (uint16, error) {
	for i := 0; i < 30; i++ {
		if val[i] != 0 {
			return 0, fmt.Errorf("non-canonical uint16 encoding in bytes32 (non-zero prefix)")
		}
	}
	return (uint16(val[30]) << 8) | uint16(val[31]), nil
}

func decodeUint32(val [32]byte) (uint32, error) {
	for i := 0; i < 28; i++ {
		if val[i] != 0 {
			return 0, fmt.Errorf("non-canonical uint32 encoding in bytes32 (non-zero prefix)")
		}
	}
	return (uint32(val[28]) << 24) |
		(uint32(val[29]) << 16) |
		(uint32(val[30]) << 8) |
		uint32(val[31]), nil
}

func decodeUint64(val [32]byte) (uint64, error) {
	for i := 0; i < 24; i++ {
		if val[i] != 0 {
			return 0, fmt.Errorf("non-canonical uint64 encoding in bytes32 (non-zero prefix)")
		}
	}
	return (uint64(val[24]) << 56) |
		(uint64(val[25]) << 48) |
		(uint64(val[26]) << 40) |
		(uint64(val[27]) << 32) |
		(uint64(val[28]) << 24) |
		(uint64(val[29]) << 16) |
		(uint64(val[30]) << 8) |
		uint64(val[31]), nil
}

func decodeUint96Big(val [32]byte) (*big.Int, error) {
	for i := 0; i < 32-uint96Size; i++ {
		if val[i] != 0 {
			return nil, fmt.Errorf("uint96: non-canonical encoding (non-zero prefix)")
		}
	}
	u := new(big.Int).SetBytes(val[32-uint96Size:])
	if u.BitLen() > 96 {
		return nil, fmt.Errorf("uint96: decoded value exceeds 96 bits")
	}
	return u, nil
}

func decodeBool(val [32]byte) (bool, error) {
	v := val[31]
	for i := 0; i < 31; i++ {
		if val[i] != 0 {
			return false, fmt.Errorf("non-canonical bool encoding in bytes32 (non-zero prefix)")
		}
	}
	switch v {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("invalid bool encoding: last byte = %d (want 0 or 1)", v)
	}
}

// Utility

func IsUint96(v *big.Int) bool {
	return v != nil && v.Sign() >= 0 && v.BitLen() <= 96
}
