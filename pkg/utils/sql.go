package utils

import (
	"database/sql"
	"math"
	"time"
)

func NewNullTime(ts time.Time) sql.NullTime {
	if ts.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: ts, Valid: true}
}

func NewNullInt16[T ~int16 | ~uint16](i *T) sql.NullInt16 {
	if i == nil {
		return sql.NullInt16{Valid: false}
	}
	if *i > math.MaxInt16 {
		return sql.NullInt16{Valid: false}
	}
	return sql.NullInt16{Int16: int16(*i), Valid: true}
}

func NewNullInt32[T ~int32 | ~uint32](i *T) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	if *i > math.MaxInt32 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

// Generic version that accepts both int64 and uint64 types
func NewNullInt64[T ~int64 | ~uint64](i *T) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	if *i > math.MaxInt64 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func NewNullInt16Slice[T ~int16 | ~uint16](ints []T) []int16 {
	if len(ints) == 0 {
		return nil
	}

	out := make([]int16, len(ints))
	for idx, val := range ints {
		out[idx] = int16(val)
	}
	return out
}

func NewNullBytes(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	return b
}
