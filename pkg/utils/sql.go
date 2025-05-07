package utils

import (
	"database/sql"
	"time"
)

func NewNullTime(ts time.Time) sql.NullTime {
	if ts.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: ts, Valid: true}
}

func NewNullInt16[T ~int16](i *T) sql.NullInt16 {
	if i == nil {
		return sql.NullInt16{Valid: false}
	}
	return sql.NullInt16{Int16: int16(*i), Valid: true}
}

// Generic version that accepts both int64 and uint64 types
func NewNullInt64[T ~int64 | ~uint64](i *T) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

func NewNullInt16Slice[T ~int16 | ~uint16](ints []T) []int16 {
	if ints == nil || len(ints) == 0 {
		return nil
	}

	out := make([]int16, len(ints))
	for idx, val := range ints {
		out[idx] = int16(val)
	}
	return out
}
