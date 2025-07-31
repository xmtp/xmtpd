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
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/deserializer"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	sleepTimeOnNoRows = 10 * time.Second
	sleepTimeOnError  = 1 * time.Second
)

type DBMigratorConfig struct {
	ctx     context.Context
	log     *zap.Logger
	db      *sql.DB
	options *config.MigrationServerOptions
}

type DBMigratorOption func(*DBMigratorConfig)

func WithContext(ctx context.Context) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.ctx = ctx
	}
}

func WithLogger(log *zap.Logger) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.log = log
	}
}

func WithDestinationDB(db *sql.DB) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.db = db
	}
}

func WithMigratorConfig(options *config.MigrationServerOptions) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.options = options
	}
}

type Migrator struct {
	// Internals.
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
	log    *zap.Logger

	// Data management.
	writer              *sql.DB
	reader              *sql.DB
	readers             map[string]ISourceReader
	transformer         IDataTransformer
	blockchainPublisher blockchain.IBlockchainPublisher

	// Configuration.
	pollInterval time.Duration
	batchSize    int32
	running      bool
}

func NewMigrationService(opts ...DBMigratorOption) (*Migrator, error) {
	cfg := &DBMigratorConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ctx == nil {
		return nil, errors.New("context is required")
	}

	if cfg.log == nil {
		return nil, errors.New("logger is required")
	}

	if cfg.db == nil {
		return nil, errors.New("destination database is required")
	}

	if cfg.options == nil {
		return nil, errors.New("migrator options are required")
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

	payerPrivateKey, err := utils.ParseEcdsaPrivateKey(cfg.options.PayerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse payer private key: %v", err)
	}

	nodeSigningKey, err := utils.ParseEcdsaPrivateKey(cfg.options.NodeSigningKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse node signing key: %v", err)
	}

	logger := cfg.log.Named("migrator")

	reader, err := db.ConnectToDB(
		cfg.ctx,
		logger,
		cfg.options.ReaderConnectionString,
		cfg.options.Namespace,
		cfg.options.WaitForDB,
		cfg.options.ReaderTimeout,
	)
	if err != nil {
		return nil, err
	}

	readers := map[string]ISourceReader{
		groupMessagesTableName:   NewGroupMessageReader(reader),
		inboxLogTableName:        NewInboxLogReader(reader),
		keyPackagesTableName:     NewKeyPackageReader(reader),
		welcomeMessagesTableName: NewWelcomeMessageReader(reader),
	}

	transformer := NewTransformer(payerPrivateKey, nodeSigningKey)

	nonceManager, err := setupNonceManager(cfg.ctx, logger, cfg.options)
	if err != nil {
		return nil, err
	}

	blockchainPublisher, err := setupBlockchainPublisher(cfg.ctx, logger, cfg.options, nonceManager)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(cfg.ctx)

	return &Migrator{
		ctx:                 ctx,
		cancel:              cancel,
		wg:                  sync.WaitGroup{},
		mu:                  sync.RWMutex{},
		log:                 logger,
		writer:              cfg.db,
		reader:              reader,
		readers:             readers,
		transformer:         transformer,
		blockchainPublisher: blockchainPublisher,
		pollInterval:        cfg.options.PollInterval,
		batchSize:           cfg.options.BatchSize,
		running:             false,
	}, nil
}

func (m *Migrator) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("migration service is already running")
	}

	m.running = true

	for tableName := range m.readers {
		m.migrationWorker(tableName)
	}

	m.log.Info("Migration service started")

	return nil
}

func (m *Migrator) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	m.cancel()
	m.wg.Wait()
	m.running = false

	if err := m.reader.Close(); err != nil {
		m.log.Error("failed to close connection to source database", zap.Error(err))
	}

	if err := m.writer.Close(); err != nil {
		m.log.Error("failed to close connection to destination database", zap.Error(err))
	}

	m.log.Info("Migration service stopped")

	return nil
}

// migrationWorker continuously processes migration for a specific table.
func (m *Migrator) migrationWorker(tableName string) {
	recvChan := make(chan ISourceRecord, m.batchSize*2)
	wrtrChan := make(chan *envelopes.OriginatorEnvelope, m.batchSize*2)
	wrtrQueries := queries.New(m.writer)

	tracing.GoPanicWrap(
		m.ctx,
		&m.wg,
		fmt.Sprintf("reader-%s", tableName),
		func(ctx context.Context) {
			defer close(recvChan)

			logger := m.log.Named("reader").With(zap.String("table", tableName))
			logger.Info("started")

			ticker := time.NewTicker(m.pollInterval)
			defer ticker.Stop()

			reader, ok := m.readers[tableName]
			if !ok {
				m.log.Error("unknown table", zap.String("table", tableName))
				return
			}

			lastMigratedID, err := wrtrQueries.GetMigrationProgress(ctx, tableName)
			if err != nil {
				logger.Fatal("failed to get migration progress", zap.Error(err))
			}

			for {
				select {
				case <-ctx.Done():
					logger.Info("context cancelled, stopping")
					return

				case <-ticker.C:
					logger.Debug(
						"getting next batch of records",
						zap.Int64("lastMigratedID", lastMigratedID),
					)

					records, newLastID, err := reader.Fetch(ctx, lastMigratedID, m.batchSize)
					if err != nil {
						switch err {
						case sql.ErrNoRows:
							logger.Info("no more records to migrate for now")
							time.Sleep(sleepTimeOnNoRows)

						default:
							logger.Error(
								"getting next batch of records failed, retrying",
								zap.Error(err),
							)

							select {
							case <-ctx.Done():
								return
							case <-time.After(sleepTimeOnError):
							}
						}

						continue
					}

					if len(records) == 0 {
						logger.Info("no more records to migrate for now")

						select {
						case <-ctx.Done():
							return
						case <-time.After(sleepTimeOnNoRows):
						}

						continue
					}

					// Update migration progress only when we have a batch of records.
					// newLastID would be 0 if there are currently no more records to migrate.
					lastMigratedID = newLastID

					logger.Debug(
						"fetched batch of records",
						zap.Int("count", len(records)),
						zap.Int64("lastID", newLastID),
					)

					for _, record := range records {
						select {
						case <-ctx.Done():
							logger.Info("context cancelled, stopping")
							return

						case recvChan <- record:
						}
					}

					logger.Debug(
						"sent batch to transformer",
						zap.Int("count", len(records)),
						zap.Int64("lastID", newLastID),
					)
				}
			}
		})

	tracing.GoPanicWrap(
		m.ctx,
		&m.wg,
		fmt.Sprintf("transformer-%s", tableName),
		func(ctx context.Context) {
			logger := m.log.Named("transformer").With(zap.String("table", tableName))
			logger.Info("started")

			defer close(wrtrChan)

			for {
				select {
				case <-ctx.Done():
					logger.Info("context cancelled, stopping")
					return

				case record, open := <-recvChan:
					if !open {
						logger.Info("channel closed, stopping")
						return
					}

					envelope, err := m.transformer.Transform(record)
					if err != nil {
						logger.Error(
							"failed to transform record",
							zap.Error(err),
							zap.Any("record", record),
						)

						// TODO: Continue, break, alert, metrics?
						continue
					}

					select {
					case <-ctx.Done():
						logger.Info("context cancelled, stopping")
						return

					case wrtrChan <- envelope:
						logger.Debug("envelope sent to writer", zap.Any("envelope", envelope))
					}
				}
			}
		})

	tracing.GoPanicWrap(
		m.ctx,
		&m.wg,
		fmt.Sprintf("writer-%s", tableName),
		func(ctx context.Context) {
			logger := m.log.Named("writer").With(zap.String("table", tableName))
			logger.Info("started")

			for {
				select {
				case <-ctx.Done():
					logger.Info("context cancelled, stopping")
					return

				case envelope, open := <-wrtrChan:
					if !open {
						logger.Info("channel closed, stopping")
						return
					}

					switch envelope.OriginatorNodeID() {
					case WelcomeMessageOriginatorID, KeyPackagesOriginatorID:
						err := retry(
							ctx,
							logger,
							50*time.Millisecond,
							func() re.RetryableError {
								return m.insertOriginatorEnvelopeDatabase(envelope)
							},
						)
						if err != nil {
							// TODO: Send to failed table? Alerts, metrics?
							logger.Error("failed to insert envelope", zap.Error(err))
						}

					case GroupMessageOriginatorID:
						payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload().(*proto.ClientEnvelope_GroupMessage)

						isCommit, err := deserializer.IsGroupMessageCommit(payload)
						if err != nil {
							logger.Error(
								"failed to check if group message is commit",
								zap.Error(err),
							)

							// TODO: Send to failed table? Alerts, metrics?
							continue
						}

						switch isCommit {
						case true:
							err := m.insertOriginatorEnvelopeBlockchain(envelope)
							if err != nil {
								logger.Error(
									"error publishing group message",
									zap.Error(err),
								)
							}

						case false:
							err := retry(
								ctx,
								logger,
								50*time.Millisecond,
								func() re.RetryableError {
									return m.insertOriginatorEnvelopeDatabase(envelope)
								},
							)
							if err != nil {
								// TODO: Send to failed table? Alerts, metrics?
								logger.Error(
									"error publishing group message",
									zap.Error(err),
								)
							}
						}

					case InboxLogOriginatorID:
						err := m.insertOriginatorEnvelopeBlockchain(envelope)
						if err != nil {
							logger.Error(
								"error publishing identity update",
								zap.Error(err),
							)
						}
					}
				}
			}
		})
}
