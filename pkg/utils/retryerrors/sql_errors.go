package retryerrors

import (
	"errors"
	"strings"

	"connectrpc.com/connect"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsRetryableSQLError(err error) bool {
	if err == nil {
		return false
	}

	var cerr *connect.Error
	if errors.As(err, &cerr) && cerr.Unwrap() != nil {
		err = cerr.Unwrap()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40P01": // deadlock detected
			return true
		case "40001": // serialization failure
			return true
		case "08006", // connection failure
			"08003", // connection does not exist
			"08000": // connection exception
			return true
		}
	}

	msg := strings.ToLower(err.Error())

	// Fallback matching
	return strings.Contains(msg, "deadlock") ||
		strings.Contains(msg, "serialization failure") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "connection closed") ||
		strings.Contains(msg, "server closed the connection") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "eof") ||
		strings.Contains(msg, "i/o timeout")
}
