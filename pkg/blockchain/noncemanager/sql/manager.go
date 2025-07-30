// Package sql provides a SQL database-backed implementation of the NonceManager interface.
// It uses database transactions and row-level locking to ensure atomic nonce allocation
// even under high concurrency.
package sql

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"go.uber.org/zap"
)

// SQLBackedNonceManager implements NonceManager using a SQL database for persistence.
// It provides thread-safe nonce allocation with configurable concurrency limits.
type SQLBackedNonceManager struct {
	db      *sql.DB
	logger  *zap.Logger
	limiter *noncemanager.OpenConnectionsLimiter
}

// NewSQLBackedNonceManager creates a new SQL-backed nonce manager with default concurrency settings
func NewSQLBackedNonceManager(db *sql.DB, logger *zap.Logger) *SQLBackedNonceManager {
	return &SQLBackedNonceManager{
		db:      db,
		logger:  logger.Named("SQLBackedNonceManager"),
		limiter: noncemanager.NewOpenConnectionsLimiter(noncemanager.BestGuessConcurrency),
	}
}

// GetNonce atomically reserves the next available nonce from the database.
// It uses SELECT FOR UPDATE SKIP LOCKED to ensure no two transactions see the same nonce.
// The transaction stays open until Cancel() (rollback) or Consume() (delete + commit).
func (s *SQLBackedNonceManager) GetNonce(ctx context.Context) (*noncemanager.NonceContext, error) {
	s.limiter.WG.Add(1)

	// Block until there is an available slot in the blockchain rate limiter
	select {
	case s.limiter.Semaphore <- struct{}{}:
	case <-ctx.Done():
		s.limiter.WG.Done()
		return nil, ctx.Err()
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		<-s.limiter.Semaphore
		s.limiter.WG.Done()
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	txQuerier := queries.New(s.db).WithTx(tx)

	nonce, err := txQuerier.GetNextAvailableNonce(ctx)
	if err != nil {
		<-s.limiter.Semaphore
		s.limiter.WG.Done()
		return nil, err
	}

	metrics.EmitPayerCurrentNonce(float64(nonce))

	ret := &noncemanager.NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			<-s.limiter.Semaphore
			s.limiter.WG.Done()
			_ = tx.Rollback()
		},
		Consume: func() error {
			<-s.limiter.Semaphore
			s.limiter.WG.Done()
			_, err = txQuerier.DeleteAvailableNonce(ctx, nonce)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			return tx.Commit()
		},
	}

	return ret, nil
}

// fillNonces generates and stores a batch of sequential nonces starting from the given value.
// This is used internally by both Replenish and FastForwardNonce methods.
func (s *SQLBackedNonceManager) fillNonces(ctx context.Context, startNonce big.Int) (int32, error) {
	querier := queries.New(s.db)
	return querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: startNonce.Int64(),
		NumElements:  10000,
	})
}

// abandonNonces removes all nonces below the given threshold value.
// This is used by FastForwardNonce to clean up obsolete nonces.
func (s *SQLBackedNonceManager) abandonNonces(ctx context.Context, endNonce big.Int) error {
	querier := queries.New(s.db)
	_, err := querier.DeleteObsoleteNonces(ctx, endNonce.Int64())
	return err
}

// Replenish ensures a sufficient number of nonces are available starting from the given nonce.
// It generates up to 10,000 nonces in a single batch operation.
func (s *SQLBackedNonceManager) Replenish(ctx context.Context, nonce big.Int) error {
	cnt, err := s.fillNonces(ctx, nonce)
	if cnt > 0 {
		s.logger.Debug(
			"Replenished nonces...",
			zap.Uint64("starting_nonce", nonce.Uint64()),
			zap.Int32("num_nonces", cnt),
		)
	}
	return err
}

// FastForwardNonce sets the nonce sequence to start from the given value and removes
// all nonces below it. This is typically used when recovering from blockchain state issues.
func (s *SQLBackedNonceManager) FastForwardNonce(ctx context.Context, nonce big.Int) error {
	_, err := s.fillNonces(ctx, nonce)
	if err != nil {
		return err
	}
	return s.abandonNonces(ctx, nonce)
}
