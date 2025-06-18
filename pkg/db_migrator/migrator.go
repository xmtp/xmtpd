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
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

// TODO: Better ordering of tables. Custom type.
const (
	addressLogTableName      = "address_log"
	groupMessagesTableName   = "group_messages"
	inboxLogTableName        = "inbox_log"
	installationsTableName   = "installations"
	welcomeMessagesTableName = "welcome_messages"

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
	ctx         context.Context
	log         *zap.Logger
	db          *sql.DB
	transformer DataTransformer
	options     *config.MigratorOptions
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

func WithTransformer(transformer DataTransformer) DBMigratorOption {
	return func(cfg *DBMigratorConfig) {
		cfg.transformer = transformer
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

	src *sql.DB
	dst *sql.DB

	transformer DataTransformer

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

	if cfg.transformer == nil {
		return nil, errors.New("transformer is required")
	}

	if cfg.db == nil {
		return nil, errors.New("destination database is required")
	}

	if cfg.options == nil {
		return nil, errors.New("reader connection string is required")
	}

	ctx, cancel := context.WithCancel(cfg.ctx)

	logger := cfg.log.Named("migration-service")

	srcDB, err := db.ConnectToDB(
		ctx,
		logger,
		cfg.options.ReaderConnectionString,
		cfg.options.Namespace,
		cfg.options.WaitForDB,
		cfg.options.ReadTimeout,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	return &dbMigrator{
		ctx:         ctx,
		cancel:      cancel,
		log:         logger,
		dst:         cfg.db,
		src:         srcDB,
		transformer: cfg.transformer,
	}, nil
}

func (s *dbMigrator) Start(ctx context.Context) error {
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

	if err := s.src.Close(); err != nil {
		s.log.Error("Failed to close source database connection", zap.Error(err))
		return err
	}

	s.log.Info("Migration service stopped")

	return nil
}

// migrationWorker continuously processes migration for a specific table.
func (s *dbMigrator) migrationWorker(tableName string) {
	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		fmt.Sprintf("migration-worker-%s", tableName),
		func(ctx context.Context) {
			logger := s.log.Named("worker").With(zap.String("table", tableName))
			logger.Info("migration worker started")

			ticker := time.NewTicker(s.pollInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					logger.Info("migration worker stopping")
					return

				case <-ticker.C:
					err := s.migrateTableBatch(ctx, logger, tableName)
					if err != nil {
						switch err {
						case ErrNoRows:
							logger.Info("no logs to migrate")
							time.Sleep(sleepTimeOnNoRows)

						default:
							logger.Error("migration batch failed", zap.Error(err))
							time.Sleep(sleepTimeOnError)

							// TODO: Handle error and possible connection retry.
						}
					}
				}
			}
		})
}

// migrateTableBatch processes a batch of records for a specific table.
func (s *dbMigrator) migrateTableBatch(
	ctx context.Context,
	logger *zap.Logger,
	tableName string,
) error {
	dstQueries := queries.New(s.dst)

	// Get migration progress for current table.
	lastMigratedID, err := dstQueries.GetMigrationProgress(ctx, tableName)
	if err != nil {
		return fmt.Errorf("failed to get migration progress: %w", err)
	}

	// Get next batch of records from source database.
	records, newLastID, err := s.getNextBatch(ctx, logger, tableName, lastMigratedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRows
		}

		return fmt.Errorf("failed to fetch batch from source database: %w", err)
	}

	if len(records) == 0 {
		return ErrNoRows
	}

	// Transform records to the xmtpd originator envelope format,
	// and inserts them into the destination database.
	processed := 0
	for _, record := range records {
		if err := s.processRecord(ctx, tableName, record); err != nil {
			logger.Error("Failed to process record", zap.Error(err), zap.Any("record", record))
			continue
		}
		processed++
	}

	// Update migration progress.
	if processed > 0 {
		if err := dstQueries.UpdateMigrationProgress(ctx, queries.UpdateMigrationProgressParams{
			LastMigratedID: newLastID,
			SourceTable:    tableName,
		}); err != nil {
			return fmt.Errorf("failed to update migration progress: %w", err)
		}

		logger.Info("Migration batch completed",
			zap.Int("processed", processed),
			zap.Int("total_fetched", len(records)),
			zap.Int64("new_last_id", newLastID))
	}

	return nil
}

// processRecord transforms and inserts a single record.
// TODO: Divide processRecord into producer (transform) and consumer (insert) functions?
func (s *dbMigrator) processRecord(
	ctx context.Context,
	tableName string,
	record interface{},
) error {
	var (
		envelope *envelopes.OriginatorEnvelope
		err      error
	)

	// Transform the record based on its type.
	switch tableName {
	case "address_log":
		addressLog, ok := record.(AddressLog)
		if !ok {
			return fmt.Errorf("invalid record type for address_log")
		}
		envelope, err = s.transformer.TransformAddressLog(&addressLog)

	case "group_messages":
		groupMessage, ok := record.(GroupMessage)
		if !ok {
			return fmt.Errorf("invalid record type for group_messages")
		}
		envelope, err = s.transformer.TransformGroupMessage(&groupMessage)

	case "inbox_log":
		inboxLog, ok := record.(InboxLog)
		if !ok {
			return fmt.Errorf("invalid record type for inbox_log")
		}
		envelope, err = s.transformer.TransformInboxLog(&inboxLog)

	case "installations":
		installation, ok := record.(Installation)
		if !ok {
			return fmt.Errorf("invalid record type for installations")
		}
		envelope, err = s.transformer.TransformInstallation(&installation)

	case "welcome_messages":
		welcomeMessage, ok := record.(WelcomeMessage)
		if !ok {
			return fmt.Errorf("invalid record type for welcome_messages")
		}
		envelope, err = s.transformer.TransformWelcomeMessage(&welcomeMessage)

	default:
		return fmt.Errorf("unknown table: %s", tableName)
	}

	if err != nil {
		return fmt.Errorf("failed to transform record: %w", err)
	}

	if envelope == nil {
		return nil
	}

	// Insert the envelope into the destination database.
	return s.insertOriginatorEnvelope(ctx, envelope)
}
