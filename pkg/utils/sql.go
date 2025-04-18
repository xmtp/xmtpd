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
