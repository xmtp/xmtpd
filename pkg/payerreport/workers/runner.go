package workers

import (
	"context"
	"fmt"
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

// Build creates a WorkerConfig instance after validating that all required fields are set.
// Returns an error if any required field is nil or invalid.
func (b *WorkerConfigBuilder) Build() (*workerConfig, error) {
	if b.ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if b.logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	if b.registrant == nil {
		return nil, fmt.Errorf("registrant cannot be nil")
	}
	if b.registry == nil {
		return nil, fmt.Errorf("registry cannot be nil")
	}
	if b.reportsManager == nil {
		return nil, fmt.Errorf("reports manager cannot be nil")
	}
	if b.store == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}
	if b.attestationPollInterval <= 0 {
		return nil, fmt.Errorf("attestation poll interval must be greater than 0")
	}
	if b.generateSelfPeriod <= 0 {
		return nil, fmt.Errorf("generate self period must be greater than 0")
	}
	if b.generateOthersPeriod <= 0 {
		return nil, fmt.Errorf("generate others period must be greater than 0")
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
	submissionNotifyCh := make(chan struct{}, 1)

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
	)

	submitterWorker := NewSubmitterWorker(
		cfg.ctx,
		cfg.logger,
		cfg.store,
		cfg.registry,
		cfg.reportsManager,
		cfg.registrant.NodeID(),
		submissionNotifyCh,
	)

	settlementWorker := NewSettlementWorker(
		cfg.ctx,
		cfg.logger,
		cfg.store,
		payerreport.NewPayerReportVerifier(cfg.logger, cfg.store),
		cfg.reportsManager,
		cfg.registrant.NodeID(),
		submissionNotifyCh,
	)

	attestationWorker.Start()
	generatorWorker.Start()
	submitterWorker.Start()
	settlementWorker.Start()

	return &WorkerWrapper{
		workers: []stoppable{attestationWorker, generatorWorker, submitterWorker, settlementWorker},
	}
}
