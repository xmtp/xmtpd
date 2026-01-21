/*
Package migrator implements a service that migrates data from a source database to a destination database.

The data to migrate (source) is expected to be of the types defined in xmtp/proto MLS V1.
https://github.com/xmtp/proto/blob/main/proto/mls/api/v1/mls.proto

Upon migration, the data is transformed and written into the xmtpd database, in the originator_envelopes table.
- The OriginatorEnvelope will have a hardcoded originator ID, based on the type of data. See types.go.
- Original sequence IDs are preserved.
- Expiry (retentionDays) are set based on the type of data.
- Congestion fee and base fee are calculated and set, based on retentionDays.
- Payer and originator envelopes are signed with the payer and node signing keys, respectively.
*/
package migrator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	sleepTimeOnNoRows = 10 * time.Second
	sleepTimeOnError  = 1 * time.Second

	// Migration destinations.
	destinationDatabase   = "database"
	destinationBlockchain = "blockchain"

	// Logging fields.
	idField             = "id"
	tableField          = "table"
	lastMigratedIDField = "last_migrated_id"

	// Logging messages.
	noMoreRecordsToMigrateMessage = "no more records to migrate for now"
	channelClosedMessage          = "channel closed, stopping"
	contextCancelledMessage       = "context cancelled, stopping"
)

type DBMigratorConfig struct {
	ctx           context.Context
	logger        *zap.Logger
	db            *db.Handler
	options       *config.MigrationServerOptions
	contracts     *config.ContractsOptions
	feeCalculator fees.IFeeCalculator
}

type DBMigratorOption func(*DBMigratorConfig)

func WithContext(ctx context.Context) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.ctx = ctx
	}
}

func WithLogger(logger *zap.Logger) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.logger = logger
	}
}

func WithDestinationDB(db *db.Handler) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.db = db
	}
}

func WithMigratorConfig(options *config.MigrationServerOptions) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.options = options
	}
}

func WithContractsOptions(contracts *config.ContractsOptions) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.contracts = contracts
	}
}

func WithFeeCalculator(calc fees.IFeeCalculator) DBMigratorOption {
	return func(cfg *DBMigratorConfig) { cfg.feeCalculator = calc }
}

type Migrator struct {
	// Internals.
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
	logger *zap.Logger

	// Data management.
	target              *db.Handler
	source              *db.Handler
	readers             map[string]ISourceReader
	transformer         IDataTransformer
	blockchainPublisher blockchain.IBlockchainPublisher

	// Configuration.
	pollInterval          time.Duration
	batchSize             int32
	databaseWriterWorkers int
	running               atomic.Bool
}

func NewMigrationService(opts ...DBMigratorOption) (*Migrator, error) {
	cfg := &DBMigratorConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ctx == nil {
		return nil, errors.New("context is required")
	}

	if cfg.logger == nil {
		return nil, errors.New("logger is required")
	}

	if cfg.db == nil {
		return nil, errors.New("destination database is required")
	}

	if cfg.options == nil {
		return nil, errors.New("migrator options are required")
	}

	if cfg.contracts == nil {
		return nil, errors.New("contracts are required")
	}

	if cfg.options.ReaderConnectionString == "" {
		return nil, errors.New("reader connection string is required")
	}

	if cfg.options.PayerPrivateKey == "" {
		return nil, errors.New("payer private key is required")
	}

	if cfg.options.NodeSigningKey == "" {
		return nil, errors.New("node signing key is required")
	}

	if cfg.feeCalculator == nil {
		return nil, errors.New("fee calculator is required")
	}

	payerPrivateKey, err := utils.ParseEcdsaPrivateKey(cfg.options.PayerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse payer private key: %v", err)
	}

	nodeSigningKey, err := utils.ParseEcdsaPrivateKey(cfg.options.NodeSigningKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse node signing key: %v", err)
	}

	logger := cfg.logger.Named(utils.MigratorLoggerName)

	reader, err := db.ConnectToDB(
		cfg.ctx,
		logger,
		cfg.options.ReaderConnectionString,
		cfg.options.Namespace,
		cfg.options.WaitForDB,
		cfg.options.ReaderTimeout,
		nil,
	)
	if err != nil {
		return nil, err
	}

	readDB := db.NewDBHandler(reader, db.WithReadReplica(reader))

	readers := map[string]ISourceReader{
		groupMessagesTableName: NewGroupMessageReader(readDB.DB(), cfg.options.StartDate.Unix()),
		inboxLogTableName:      NewInboxLogReader(readDB.DB(), cfg.options.StartDate.UnixNano()),
		keyPackagesTableName:   NewKeyPackageReader(readDB.DB()),
		welcomeMessagesTableName: NewWelcomeMessageReader(
			readDB.DB(),
			cfg.options.StartDate.Unix(),
		),
		commitMessagesTableName: NewCommitMessageReader(readDB.DB(), cfg.options.StartDate.Unix()),
	}

	transformer := NewTransformer(cfg.feeCalculator, payerPrivateKey, nodeSigningKey)

	blockchainPublisher, err := setupBlockchainPublisher(
		cfg.ctx,
		logger,
		cfg.db,
		cfg.options.PayerPrivateKey,
		cfg.contracts,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(cfg.ctx)

	return &Migrator{
		ctx:                   ctx,
		cancel:                cancel,
		wg:                    sync.WaitGroup{},
		mu:                    sync.RWMutex{},
		logger:                logger,
		target:                cfg.db,
		source:                readDB,
		readers:               readers,
		transformer:           transformer,
		blockchainPublisher:   blockchainPublisher,
		pollInterval:          cfg.options.PollInterval,
		batchSize:             cfg.options.BatchSize,
		databaseWriterWorkers: cfg.options.DatabaseWriterWorkers,
		running:               atomic.Bool{},
	}, nil
}

func (m *Migrator) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running.Swap(true) {
		return fmt.Errorf("migration service is already running")
	}

	if err := m.startKeyPackagesWorker(); err != nil {
		return err
	}

	if err := m.startWelcomeMessagesWorker(); err != nil {
		return err
	}

	if err := m.startGroupMessagesWorker(); err != nil {
		return err
	}

	if err := m.startCommitMessagesWorker(); err != nil {
		return err
	}

	if err := m.startInboxLogWorker(); err != nil {
		return err
	}

	m.logger.Info("migration service started")

	return nil
}

func (m *Migrator) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running.Swap(false) {
		return fmt.Errorf("migration service is not running")
	}

	m.cancel()
	m.wg.Wait()

	if err := m.source.Close(); err != nil {
		m.logger.Error("failed to close connection to source database", zap.Error(err))
	}

	m.logger.Info("migration service stopped")

	return nil
}

func (m *Migrator) startKeyPackagesWorker() error {
	keyPackagesWorker := NewWorker(
		keyPackagesTableName,
		m.batchSize,
		m.target,
		nil,
		m.logger,
		m.pollInterval,
		m.databaseWriterWorkers,
	)

	if err := keyPackagesWorker.StartReader(m.ctx, m.readers[keyPackagesTableName]); err != nil {
		return err
	}

	if err := keyPackagesWorker.StartTransformer(m.ctx, m.transformer); err != nil {
		return err
	}

	if err := keyPackagesWorker.StartDatabaseWriter(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) startWelcomeMessagesWorker() error {
	welcomeMessagesWorker := NewWorker(
		welcomeMessagesTableName,
		m.batchSize,
		m.target,
		nil,
		m.logger,
		m.pollInterval,
		m.databaseWriterWorkers,
	)

	if err := welcomeMessagesWorker.StartReader(m.ctx, m.readers[welcomeMessagesTableName]); err != nil {
		return err
	}

	if err := welcomeMessagesWorker.StartTransformer(m.ctx, m.transformer); err != nil {
		return err
	}

	if err := welcomeMessagesWorker.StartDatabaseWriter(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) startGroupMessagesWorker() error {
	groupMessagesWorker := NewWorker(
		groupMessagesTableName,
		m.batchSize,
		m.target,
		nil,
		m.logger,
		m.pollInterval,
		m.databaseWriterWorkers,
	)

	if err := groupMessagesWorker.StartReader(m.ctx, m.readers[groupMessagesTableName]); err != nil {
		return err
	}

	if err := groupMessagesWorker.StartTransformer(m.ctx, m.transformer); err != nil {
		return err
	}

	if err := groupMessagesWorker.StartDatabaseWriter(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) startCommitMessagesWorker() error {
	commitMessagesWorker := NewWorker(
		commitMessagesTableName,
		m.batchSize,
		m.target,
		m.blockchainPublisher,
		m.logger,
		m.pollInterval,
		m.databaseWriterWorkers,
	)

	if err := commitMessagesWorker.StartReader(m.ctx, m.readers[commitMessagesTableName]); err != nil {
		return err
	}

	if err := commitMessagesWorker.StartTransformer(m.ctx, m.transformer); err != nil {
		return err
	}

	if err := commitMessagesWorker.StartBlockchainWriterBatch(m.ctx); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) startInboxLogWorker() error {
	inboxLogWorker := NewWorker(
		inboxLogTableName,
		m.batchSize,
		m.target,
		m.blockchainPublisher,
		m.logger,
		m.pollInterval,
		m.databaseWriterWorkers,
	)

	if err := inboxLogWorker.StartReader(m.ctx, m.readers[inboxLogTableName]); err != nil {
		return err
	}

	if err := inboxLogWorker.StartTransformer(m.ctx, m.transformer); err != nil {
		return err
	}

	if err := inboxLogWorker.StartBlockchainWriterBatch(m.ctx); err != nil {
		return err
	}

	return nil
}
