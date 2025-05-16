package sync

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type SyncServerConfig struct {
	Ctx                        context.Context
	Log                        *zap.Logger
	NodeRegistry               registry.NodeRegistry
	Registrant                 *registrant.Registrant
	DB                         *sql.DB
	FeeCalculator              fees.IFeeCalculator
	PayerReportStore           payerreport.IPayerReportStore
	PayerReportDomainSeparator common.Hash
}

type SyncServerOption func(*SyncServerConfig)

func WithContext(ctx context.Context) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Ctx = ctx }
}

func WithLogger(log *zap.Logger) SyncServerOption {
	return func(cfg *SyncServerConfig) { cfg.Log = log }
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

type SyncServer struct {
	ctx        context.Context
	log        *zap.Logger
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

	if cfg.Log == nil {
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
		log:        cfg.Log.Named("sync"),
		registrant: cfg.Registrant,
		store:      cfg.DB,
		worker:     worker,
	}, nil
}

func (s *SyncServer) Close() {
	s.log.Debug("Closing")
	s.worker.close()
	s.log.Debug("Closed")
}
