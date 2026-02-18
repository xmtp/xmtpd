package bench

import (
	"context"
	"database/sql"
	"log"
)

func seedLedger(_ context.Context, _ *sql.DB) {
	log.Printf("seeded ledger: 0 rows (stub)")
}
