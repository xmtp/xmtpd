package workers

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/payerreport"
	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type stoppable interface {
	Stop()
}

// WorkerConfig contains the configuration for running payer report workers.
// Use NewWorkerConfigBuilder() to create instances of this struct.
type workerConfig struct {
	ctx                     context.Context
	logger                  *zap.Logger
	registrant              registrant.IRegistrant
	registry                registry.NodeRegistry
	reportsManager          blockchain.PayerReportsManager
	store                   payerreport.IPayerReportStore
	domainSeparator         common.Hash
	attestationPollInterval time.Duration
	generateSelfPeriod      time.Duration
	generateOthersPeriod    time.Duration
	expirySelfPeriod        time.Duration
	expiryOthersPeriod      time.Duration
}

// WorkerConfigBuilder provides a builder pattern for creating WorkerConfig instances.
// All fields are required and the Build() method will validate that none are nil.
type WorkerConfigBuilder struct {
	ctx                     context.Context
	logger                  *zap.Logger
	registrant              registrant.IRegistrant
	registry                registry.NodeRegistry
	reportsManager          blockchain.PayerReportsManager
	store                   payerreport.IPayerReportStore
	domainSeparator         common.Hash
	attestationPollInterval time.Duration
	generateSelfPeriod      time.Duration
	generateOthersPeriod    time.Duration
	expirySelfPeriod        time.Duration
	expiryOthersPeriod      time.Duration
}

// NewWorkerConfigBuilder creates a new WorkerConfigBuilder instance.
func NewWorkerConfigBuilder() *WorkerConfigBuilder {
	return &WorkerConfigBuilder{}
}

func (b *WorkerConfigBuilder) WithContext(ctx context.Context) *WorkerConfigBuilder {
	b.ctx = ctx
	return b
}

func (b *WorkerConfigBuilder) WithLogger(logger *zap.Logger) *WorkerConfigBuilder {
	b.logger = logger
	return b
}

func (b *WorkerConfigBuilder) WithRegistrant(
	registrant registrant.IRegistrant,
) *WorkerConfigBuilder {
	b.registrant = registrant
	return b
}

func (b *WorkerConfigBuilder) WithRegistry(registry registry.NodeRegistry) *WorkerConfigBuilder {
	b.registry = registry
	return b
}

func (b *WorkerConfigBuilder) WithReportsManager(
	reportsManager blockchain.PayerReportsManager,
) *WorkerConfigBuilder {
	b.reportsManager = reportsManager
	return b
}

func (b *WorkerConfigBuilder) WithStore(store payerreport.IPayerReportStore) *WorkerConfigBuilder {
	b.store = store
	return b
}

func (b *WorkerConfigBuilder) WithDomainSeparator(
	domainSeparator common.Hash,
) *WorkerConfigBuilder {
	b.domainSeparator = domainSeparator
	return b
}

func (b *WorkerConfigBuilder) WithAttestationPollInterval(
	interval time.Duration,
) *WorkerConfigBuilder {
	b.attestationPollInterval = interval
	return b
}

func (b *WorkerConfigBuilder) WithGenerationSelfPeriod(
	period time.Duration,
) *WorkerConfigBuilder {
	b.generateSelfPeriod = period
	return b
}

func (b *WorkerConfigBuilder) WithGenerationOthersPeriod(
	period time.Duration,
) *WorkerConfigBuilder {
	b.generateOthersPeriod = period
	return b
}

func (b *WorkerConfigBuilder) WithExpirySelfPeriod(
	period time.Duration,
) *WorkerConfigBuilder {
	b.expirySelfPeriod = period
	return b
}

func (b *WorkerConfigBuilder) WithExpiryOthersPeriod(
	period time.Duration,
) *WorkerConfigBuilder {
	b.expiryOthersPeriod = period
	return b
}

// Build creates a WorkerConfig instance after validating that all required fields are set.
// Returns an error if any required field is nil or invalid.
func (b *WorkerConfigBuilder) Build() (*workerConfig, error) {
	if b.ctx == nil {
		return nil, errors.New("context cannot be nil")
	}
	if b.logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if b.registrant == nil {
		return nil, errors.New("registrant cannot be nil")
	}
	if b.registry == nil {
		return nil, errors.New("registry cannot be nil")
	}
	if b.reportsManager == nil {
		return nil, errors.New("reports manager cannot be nil")
	}
	if b.store == nil {
		return nil, errors.New("store cannot be nil")
	}
	if b.attestationPollInterval <= 0 {
		return nil, errors.New("attestation poll interval must be greater than 0")
	}
	if b.generateSelfPeriod <= 0 {
		return nil, errors.New("generate self period must be greater than 0")
	}
	if b.generateOthersPeriod <= 0 {
		return nil, errors.New("generate others period must be greater than 0")
	}

	// Expiration periods must be always longer than generation periods to avoid race conditions.
	if float64(
		b.expirySelfPeriod.Nanoseconds(),
	) < float64(
		b.generateSelfPeriod.Nanoseconds(),
	)*1.5 {
		return nil, errors.New("expiry self period must be at least 1.5x the generate self period")
	}

	if float64(
		b.expiryOthersPeriod.Nanoseconds(),
	) < float64(
		b.generateOthersPeriod.Nanoseconds(),
	)*1.5 {
		return nil, errors.New(
			"expiry others period must be at least 1.5x the generate others period",
		)
	}

	return &workerConfig{
		ctx:                     b.ctx,
		logger:                  b.logger,
		registrant:              b.registrant,
		registry:                b.registry,
		reportsManager:          b.reportsManager,
		store:                   b.store,
		domainSeparator:         b.domainSeparator,
		attestationPollInterval: b.attestationPollInterval,
		generateSelfPeriod:      b.generateSelfPeriod,
		generateOthersPeriod:    b.generateOthersPeriod,
		expirySelfPeriod:        b.expirySelfPeriod,
		expiryOthersPeriod:      b.expiryOthersPeriod,
	}, nil
}

type WorkerWrapper struct {
	workers []stoppable
}

func (w *WorkerWrapper) Stop() {
	for _, worker := range w.workers {
		worker.Stop()
	}
}

// RunWorkers creates and starts all payer report workers with the given configuration.
// The configuration should be created using NewWorkerConfigBuilder().Build().
// Returns a WorkerWrapper that can be used to stop all workers.
func RunWorkers(cfg workerConfig) *WorkerWrapper {
	attestationWorker := NewAttestationWorker(
		cfg.ctx,
		cfg.logger,
		cfg.registrant,
		cfg.store,
		cfg.attestationPollInterval,
		cfg.domainSeparator,
	)

	generatorWorker := NewGeneratorWorker(
		cfg.ctx,
		cfg.logger,
		cfg.store,
		cfg.registry,
		cfg.registrant,
		cfg.domainSeparator,
		cfg.generateSelfPeriod,
		cfg.generateOthersPeriod,
		cfg.expirySelfPeriod,
		cfg.expiryOthersPeriod,
	)

	submitterWorker := NewSubmitterWorker(
		cfg.ctx,
		cfg.logger,
		cfg.store,
		cfg.registry,
		cfg.reportsManager,
		cfg.registrant.NodeID(),
	)

	settlementWorker := NewSettlementWorker(
		cfg.ctx,
		cfg.logger,
		cfg.store,
		payerreport.NewPayerReportVerifier(cfg.logger, cfg.store),
		cfg.reportsManager,
		cfg.registrant.NodeID(),
	)

	attestationWorker.Start()
	generatorWorker.Start()
	submitterWorker.Start()
	settlementWorker.Start()

	return &WorkerWrapper{
		workers: []stoppable{attestationWorker, generatorWorker, submitterWorker, settlementWorker},
	}
}
