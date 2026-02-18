package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedCongestion(_ context.Context, _ *sql.DB) {
	log.Printf("seeded congestion: 0 rows (stub)")
}
