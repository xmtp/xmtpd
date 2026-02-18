package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedUsage(_ context.Context, _ *sql.DB) {
	log.Printf("seeded usage: 0 rows (stub)")
}
