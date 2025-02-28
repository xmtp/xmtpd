package blockchain

import (
	"context"
	"database/sql"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

type NonceContext struct {
	Nonce   uint64
	Cancel  func()
	Consume func() error
}
type NonceManager interface {
	GetNonce(ctx context.Context) (*NonceContext, error)
	FastForwardNonce(ctx context.Context, nonce uint64) error
}

type SQLBackedNonceManager struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSQLBackedNonceManager(db *sql.DB, logger *zap.Logger) *SQLBackedNonceManager {
	return &SQLBackedNonceManager{
		db:     db,
		logger: logger.Named("SQLBackedNonceManager"),
	}
}

func (s *SQLBackedNonceManager) GetNonce(ctx context.Context) (*NonceContext, error) {
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

	nonce, err := txQuerier.GetNextAvailablePayerSequence(ctx)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Generated Nonce", zap.Int64("nonce", nonce))

	ret := &NonceContext{
		Nonce: uint64(nonce),
		Cancel: func() {
			_ = tx.Rollback()
		},
		Consume: func() error {
			_, err = txQuerier.DeleteAvailablePayerSequence(ctx, nonce)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			return tx.Commit()
		},
	}

	return ret, nil

}

func (s *SQLBackedNonceManager) FillNonces(ctx context.Context, startNonce uint64) (err error) {
	querier := queries.New(s.db)
	return querier.FillPayerSequence(ctx, queries.FillPayerSequenceParams{
		PendingNonce: int64(startNonce),
		NumElements:  100,
	})
}

func (s *SQLBackedNonceManager) FastForwardNonce(ctx context.Context, nonce uint64) error {
	//querier := queries.New(s.db)
	return nil
}
