package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var migratorE2ELatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_migrator_e2e_latency_seconds",
		Help:    "Time spent migrating a message",
		Buckets: []float64{0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 60.0},
	},
	[]string{"table", "destination"},
)

func EmitMigratorE2ELatency(table, destination string, duration float64) {
	migratorE2ELatency.With(prometheus.Labels{
		"table":       table,
		"destination": destination,
	}).Observe(duration)
}

var migratorDestLastSequenceIDBlockchain = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_migrator_destination_blockchain_last_sequence_id",
		Help: "Last sequence ID published to blockchain",
	},
	[]string{"table"},
)

func EmitMigratorDestLastSequenceIDBlockchain(table string, sequenceID int64) {
	migratorDestLastSequenceIDBlockchain.With(prometheus.Labels{"table": table}).
		Set(float64(sequenceID))
}

var migratorDestLastSequenceIDDatabase = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_migrator_destination_database_last_sequence_id",
		Help: "Last sequence ID persisted in destination database",
	},
	[]string{"table"},
)

func EmitMigratorDestLastSequenceIDDatabase(table string, sequenceID int64) {
	migratorDestLastSequenceIDDatabase.With(prometheus.Labels{"table": table}).
		Set(float64(sequenceID))
}

var migratorReaderErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_migrator_reader_errors_total",
		Help: "Total number of reader errors",
	},
	[]string{"table", "error_type"},
)

func EmitMigratorReaderError(table, errorType string) {
	migratorReaderErrors.With(prometheus.Labels{
		"table":      table,
		"error_type": errorType,
	}).Inc()
}

var migratorReaderFetchDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_migrator_reader_fetch_duration_seconds",
		Help:    "Time spent fetching records from source database",
		Buckets: []float64{0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 60.0},
	},
	[]string{"table"},
)

func EmitMigratorReaderFetchDuration(table string, duration float64) {
	migratorReaderFetchDuration.With(prometheus.Labels{"table": table}).Observe(duration)
}

func MeasureReaderLatency[Return any](table string, fn func() (Return, error)) (Return, error) {
	start := time.Now()
	ret, err := fn()
	if err == nil {
		EmitMigratorReaderFetchDuration(table, time.Since(start).Seconds())
	}
	return ret, err
}

var migratorReaderNumRowsFound = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_migrator_reader_num_rows_found",
		Help: "Number of rows fetched from source database",
	},
	[]string{"table"},
)

func EmitMigratorReaderNumRowsFound(table string, numRows int64) {
	migratorReaderNumRowsFound.With(prometheus.Labels{"table": table}).Add(float64(numRows))
}

var migratorSourceLastSequenceID = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_migrator_source_last_sequence_id",
		Help: "Last sequence ID pulled from source DB",
	},
	[]string{"table"},
)

func EmitMigratorSourceLastSequenceID(table string, sequenceID int64) {
	migratorSourceLastSequenceID.With(prometheus.Labels{"table": table}).Set(float64(sequenceID))
}

var migratorTransformerErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_migrator_transformer_errors_total",
		Help: "Total number of transformation errors",
	},
	[]string{"table"},
)

func EmitMigratorTransformerError(table string) {
	migratorTransformerErrors.With(prometheus.Labels{"table": table}).Inc()
}

var migratorWriterErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_migrator_writer_errors_total",
		Help: "Total number of writer errors by destination and error type",
	},
	[]string{"table", "destination", "error_type"},
)

func EmitMigratorWriterError(table, destination, errorType string) {
	migratorWriterErrors.With(prometheus.Labels{
		"table":       table,
		"destination": destination,
		"error_type":  errorType,
	}).Inc()
}

var migratorWriterLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_migrator_writer_latency_seconds",
		Help:    "Time spent writing to destination",
		Buckets: []float64{0.1, 0.2, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 60.0},
	},
	[]string{"table", "destination"},
)

func EmitMigratorWriterLatency(table, destination string, duration float64) {
	migratorWriterLatency.With(prometheus.Labels{
		"table":       table,
		"destination": destination,
	}).Observe(duration)
}

func MeasureWriterLatency(
	table, destination string,
	fn func() error,
) error {
	start := time.Now()
	err := fn()
	if err == nil {
		EmitMigratorWriterLatency(table, destination, time.Since(start).Seconds())
	}
	return err
}

var migratorWriterRetryAttempts = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_migrator_writer_retry_attempts",
		Help:    "Number of retry attempts before success or failure",
		Buckets: []float64{0, 1, 2, 3, 4, 5, 10, 20},
	},
	[]string{"table", "destination"},
)

var migratorWriterRowsMigrated = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_migrator_writer_rows_migrated",
		Help: "Total number of rows successfully migrated",
	},
	[]string{"table"},
)

func EmitMigratorWriterRowsMigrated(table string, numRows int64) {
	migratorWriterRowsMigrated.With(prometheus.Labels{"table": table}).Add(float64(numRows))
}

func EmitMigratorWriterRetryAttempts(table, destination string, attempts int) {
	migratorWriterRetryAttempts.With(prometheus.Labels{
		"table":       table,
		"destination": destination,
	}).Observe(float64(attempts))
}
