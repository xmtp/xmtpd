// Package sync implements the sync server.
package sync

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type MigrationConfig struct {
	Enable     bool
	FromNodeID uint32
}

type SyncServerConfig struct {
	Ctx                        context.Context
	ClientMetrics              *grpcprom.ClientMetrics
	DB                         *sql.DB
	FeeCalculator              fees.IFeeCalculator
	Logger                     *zap.Logger
	Migration                  MigrationConfig
	NodeRegistry               registry.NodeRegistry
	PayerReportDomainSeparator common.Hash
	PayerReportStore           payerreport.IPayerReportStore
	Registrant                 *registrant.Registrant
}

type SyncServerOption func(*SyncServerConfig)

func WithContext(ctx context.Context) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Ctx = ctx }
}

func WithClientMetrics(metrics *grpcprom.ClientMetrics) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.ClientMetrics = metrics }
}

func WithLogger(logger *zap.Logger) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Logger = logger }
}

func WithNodeRegistry(reg registry.NodeRegistry) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.NodeRegistry = reg }
}

func WithRegistrant(r *registrant.Registrant) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Registrant = r }
}

func WithDB(db *sql.DB) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.DB = db }
}

func WithFeeCalculator(calc fees.IFeeCalculator) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.FeeCalculator = calc }
}

func WithPayerReportStore(store payerreport.IPayerReportStore) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.PayerReportStore = store }
}

func WithPayerReportDomainSeparator(domainSeparator common.Hash) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.PayerReportDomainSeparator = domainSeparator }
}

func WithMigration(migration MigrationConfig) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Migration = migration }
}

type SyncServer struct {
	ctx        context.Context
	logger     *zap.Logger
	registrant *registrant.Registrant
	store      *sql.DB
	worker     *syncWorker
}

func NewSyncServer(opts ...SyncServerOption) (*SyncServer, error) {
	cfg := &SyncServerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Ctx == nil {
		return nil, errors.New("syncserver: context is required")
	}

	if cfg.ClientMetrics == nil {
		return nil, errors.New("syncserver: client metrics is required")
	}

	if cfg.Logger == nil {
		return nil, errors.New("syncserver: logger is required")
	}

	if cfg.NodeRegistry == nil {
		return nil, errors.New("syncserver: node registry is required")
	}

	if cfg.Registrant == nil {
		return nil, errors.New("syncserver: registrant is required")
	}

	if cfg.DB == nil {
		return nil, errors.New("syncserver: DB is required")
	}

	if cfg.FeeCalculator == nil {
		return nil, errors.New("syncserver: fee calculator is required")
	}

	worker, err := startSyncWorker(cfg)
	if err != nil {
		return nil, err
	}

	return &SyncServer{
		ctx:        cfg.Ctx,
		logger:     cfg.Logger.Named(utils.SyncLoggerName),
		registrant: cfg.Registrant,
		store:      cfg.DB,
		worker:     worker,
	}, nil
}

func (s *SyncServer) Close() {
	s.logger.Debug("closing")
	s.worker.close()
	s.logger.Debug("closed")
}
