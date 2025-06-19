package db_migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const (
	sleepTimeOnNoRows = 10 * time.Second
	sleepTimeOnError  = 1 * time.Second
)

var (
	tables = []string{
		addressLogTableName,
		groupMessagesTableName,
		inboxLogTableName,
		installationsTableName,
		welcomeMessagesTableName,
	}

	ErrNoRows = errors.New("no logs to migrate")
)

type DBMigratorConfig struct {
	ctx     context.Context
	log     *zap.Logger
	db      *sql.DB
	options *config.MigratorOptions
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

func WithMigratorConfig(options *config.MigratorOptions) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.options = options
	}
}

type dbMigrator struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	log    *zap.Logger

	writer  *sql.DB
	reader  *sql.DB
	readers map[string]Reader

	transformer Transformer

	running      bool
	batchSize    int32
	pollInterval time.Duration
	mu           sync.RWMutex
}

func NewMigrationService(opts ...DBMigratorOption) (*dbMigrator, error) {
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

	logger := cfg.log.Named("migrator")

	source, err := db.ConnectToDB(
		cfg.ctx,
		logger,
		cfg.options.ReaderConnectionString,
		cfg.options.Namespace,
		cfg.options.WaitForDB,
		cfg.options.ReadTimeout,
	)
	if err != nil {
		return nil, err
	}

	sources := map[string]Reader{
		addressLogTableName:      NewAddressLogReader(source),
		groupMessagesTableName:   NewGroupMessageReader(source),
		inboxLogTableName:        NewInboxLogReader(source),
		installationsTableName:   NewInstallationReader(source),
		welcomeMessagesTableName: NewWelcomeMessageReader(source),
	}

	ctx, cancel := context.WithCancel(cfg.ctx)

	return &dbMigrator{
		ctx:          ctx,
		cancel:       cancel,
		log:          logger,
		writer:       cfg.db,
		reader:       source,
		readers:      sources,
		transformer:  NewTransformer(),
		batchSize:    cfg.options.BatchSize,
		pollInterval: cfg.options.PollInterval,
	}, nil
}

func (s *dbMigrator) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("migration service is already running")
	}

	s.running = true

	for _, table := range tables {
		s.migrationWorker(table)
	}

	s.log.Info("Migration service started", zap.Strings("tables", tables))

	return nil
}

func (s *dbMigrator) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.cancel()
	s.wg.Wait()
	s.running = false

	if err := s.reader.Close(); err != nil {
		s.log.Error("failed to close connection to source database", zap.Error(err))
	}

	if err := s.writer.Close(); err != nil {
		s.log.Error("failed to close connection to destination database", zap.Error(err))
	}

	s.log.Info("Migration service stopped")

	return nil
}

// migrationWorker continuously processes migration for a specific table.
func (s *dbMigrator) migrationWorker(tableName string) {
	recvChan := make(chan Record, s.batchSize*2)
	wrtrChan := make(chan *envelopes.OriginatorEnvelope, s.batchSize*2)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("reader-%s", tableName),
		func(ctx context.Context) {
			defer close(recvChan)

			logger := s.log.Named("reader").With(zap.String("table", tableName))
			logger.Info("started")

			wrtrQueries := queries.New(s.writer)

			ticker := time.NewTicker(s.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					logger.Info("context cancelled, stopping")
					return

				case <-ticker.C:
					records, newLastID, err := s.nextRecords(ctx, logger, wrtrQueries, tableName)
					if err != nil {
						switch err {
						case ErrNoRows:
							logger.Info(ErrNoRows.Error())
							time.Sleep(sleepTimeOnNoRows)

						default:
							logger.Error("migration batch failed", zap.Error(err))
							time.Sleep(sleepTimeOnError)
						}

						continue
					}

					for _, record := range records {
						select {
						case <-ctx.Done():
							logger.Info("context cancelled, stopping")
							return

						case recvChan <- record:
						}
					}

					logger.Debug("batch sent to transformer",
						zap.Int("total_fetched", len(records)),
						zap.Int64("new_last_id", newLastID))

					if err := wrtrQueries.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
						LastMigratedID: newLastID,
						SourceTable:    tableName,
					}); err != nil {
						logger.Error("failed to update migration progress", zap.Error(err))
					}
				}
			}
		})

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("transformer-%s", tableName),
		func(ctx context.Context) {
			logger := s.log.Named("transformer").With(zap.String("table", tableName))
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

					envelope, err := s.transformer.Transform(record)
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
		s.ctx,
		&s.wg,
		fmt.Sprintf("writer-%s", tableName),
		func(ctx context.Context) {
			logger := s.log.Named("writer").With(zap.String("table", tableName))
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

					err := retry(
						ctx,
						logger,
						50*time.Millisecond,
						func() re.RetryableError {
							return s.insertOriginatorEnvelope(ctx, envelope)
						},
					)
					if err != nil {
						logger.Error("failed to insert envelope", zap.Error(err))
						continue
					}
				}
			}
		})
}
