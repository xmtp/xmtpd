package blockchain

import (
	"context"
	"database/sql"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
	"math/big"
	"sync"
)

type NonceContext struct {
	Nonce   big.Int
	Cancel  func()
	Consume func() error
}
type NonceManager interface {
	GetNonce(ctx context.Context) (*NonceContext, error)
	FastForwardNonce(ctx context.Context, nonce big.Int) error
	Replenish(ctx context.Context, nonce big.Int) error
}

// OpenConnectionsLimiter controls the number of concurrent requests
type OpenConnectionsLimiter struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
}

// MaxConcurrentRequests the blockchain mempool can usually only hold 64 transactions from the same fromAddress
const MaxConcurrentRequests = 64
const BestGuessConcurrency = 32

// NewOpenConnectionsLimiter initializes a OpenConnectionsLimiter with a limit
func NewOpenConnectionsLimiter(maxConcurrent int) *OpenConnectionsLimiter {
	if maxConcurrent > MaxConcurrentRequests {
		maxConcurrent = MaxConcurrentRequests
	}
	return &OpenConnectionsLimiter{
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

type SQLBackedNonceManager struct {
	db      *sql.DB
	logger  *zap.Logger
	limiter *OpenConnectionsLimiter
}

func NewSQLBackedNonceManager(db *sql.DB, logger *zap.Logger) *SQLBackedNonceManager {
	return &SQLBackedNonceManager{
		db:      db,
		logger:  logger.Named("SQLBackedNonceManager"),
		limiter: NewOpenConnectionsLimiter(BestGuessConcurrency),
	}
}

func (s *SQLBackedNonceManager) GetNonce(ctx context.Context) (*NonceContext, error) {
	s.limiter.wg.Add(1)

	// block until there is an available slot in the blockchain rate limiter
	select {
	case s.limiter.semaphore <- struct{}{}:
	case <-ctx.Done():
		s.limiter.wg.Done()
		return nil, ctx.Err()
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
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
		return nil, err
	}

	s.logger.Debug("Generated Nonce", zap.Int64("nonce", nonce))

	ret := &NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			<-s.limiter.semaphore
			s.limiter.wg.Done()
			_ = tx.Rollback()
		},
		Consume: func() error {
			<-s.limiter.semaphore
			s.limiter.wg.Done()
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

func (s *SQLBackedNonceManager) fillNonces(ctx context.Context, startNonce big.Int) error {
	querier := queries.New(s.db)
	return querier.FillNonceSequence(ctx, queries.FillNonceSequenceParams{
		PendingNonce: startNonce.Int64(),
		NumElements:  10000,
	})
}

func (s *SQLBackedNonceManager) abandonNonces(ctx context.Context, endNonce big.Int) error {
	querier := queries.New(s.db)
	_, err := querier.DeleteObsoleteNonces(ctx, endNonce.Int64())
	return err
}

func (s *SQLBackedNonceManager) Replenish(ctx context.Context, nonce big.Int) error {
	s.logger.Debug("Replenishing nonces...", zap.Uint64("nonce", nonce.Uint64()))
	return s.fillNonces(ctx, nonce)
}

func (s *SQLBackedNonceManager) FastForwardNonce(ctx context.Context, nonce big.Int) error {
	err := s.fillNonces(ctx, nonce)
	if err != nil {
		return err
	}
	return s.abandonNonces(ctx, nonce)
}
