package utils

import (
	"errors"
	"math"
)

var ErrIntOverflow = errors.New("overflow in conversion from uint to int")

func Uint64ToInt64(u uint64) (int64, error) {
	if u > math.MaxInt64 {
		return 0, ErrIntOverflow
	}
	return int64(u), nil
}

func Uint32ToInt32(u uint32) (int32, error) {
	if u > math.MaxInt32 {
		return 0, ErrIntOverflow
	}
	return int32(u), nil
}
