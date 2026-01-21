package config

import "time"

type MigrationServerOptions struct {
	Enable bool `long:"enable" env:"XMTPD_MIGRATION_SERVER_ENABLE" description:"Enable the migration server"`

	PayerPrivateKey string `long:"payer-private-key" env:"XMTPD_MIGRATION_PAYER_PRIVATE_KEY" description:"Private key used to sign payer envelopes"`
	NodeSigningKey  string `long:"node-signing-key"  env:"XMTPD_MIGRATION_NODE_SIGNING_KEY"  description:"Private key used to sign originator envelopes"`

	ReaderConnectionString string        `long:"reader-connection-string" env:"XMTPD_MIGRATION_DB_READER_CONNECTION_STRING" description:"Reader connection string"`
	ReaderTimeout          time.Duration `long:"reader-timeout"           env:"XMTPD_MIGRATION_DB_READER_TIMEOUT"           description:"Timeout for reading from the database"          default:"10s"`
	WaitForDB              time.Duration `long:"wait-for"                 env:"XMTPD_MIGRATION_DB_WAIT_FOR"                 description:"wait for DB on start, up to specified duration" default:"30s"`
	BatchSize              int32         `long:"batch-size"               env:"XMTPD_MIGRATION_DB_BATCH_SIZE"               description:"Batch size for migration"                       default:"1000"`
	PollInterval           time.Duration `long:"process-interval"         env:"XMTPD_MIGRATION_DB_PROCESS_INTERVAL"         description:"Interval for processing migration"              default:"10s"`
	DatabaseWriterWorkers  int           `long:"database-writer-workers"  env:"XMTPD_MIGRATION_DB_WRITER_WORKERS"           description:"Number of database writer workers"              default:"2"`
	Namespace              string        `long:"namespace"                env:"XMTPD_MIGRATION_DB_NAMESPACE"                description:"Namespace for migration"                        default:""`
	StartDate              time.Time     `long:"start-date"               env:"XMTPD_MIGRATION_START_DATE"                  description:"Start date for migration"                       default:"2025-10-01T00:00:00Z"`
}

type MigrationClientOptions struct {
	Enable     bool   `long:"enable"       env:"XMTPD_MIGRATION_CLIENT_ENABLE"       description:"Enable the migration client"`
	FromNodeID uint32 `long:"from-node-id" env:"XMTPD_MIGRATION_CLIENT_FROM_NODE_ID" description:"Node ID to start migration from" default:"100"`
}
