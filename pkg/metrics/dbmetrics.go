package metrics

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/prometheus/client_golang/prometheus"
)

var QueryDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "db",
		Name:      "query_duration_seconds",
		Help:      "Duration of SQL queries by named statement.",
		Buckets:   []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	},
	[]string{"query", "op"},
)

var QueryErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "db",
		Name:      "query_errors_total",
		Help:      "Total SQL query errors by named statement.",
	},
	[]string{"query", "op"},
)

// Extract /* name:FooBar */ from SQL text
var qnameRE = regexp.MustCompile(`--\s*name:\s*([A-Za-z0-9_]+)`)

func queryName(sql string) (string) {
	if m := qnameRE.FindStringSubmatch(sql); m != nil {
		return strings.TrimSpace(m[1])
	}

	return ""
}

type PromLogger struct{}

func (PromLogger) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	switch msg {
	case "Query", "Exec", "Batch":
		// allow these
	default:
		return
	}

	// pgx tracelog includes "sql" and "time"
	rawSQL, _ := data["sql"].(string)
	dur, _ := data["time"].(time.Duration)
	errVal := data["err"]

	name := queryName(rawSQL)

	if name == "" {
		// unknown statement
		return
	}

	// Skip transaction control statements
	if isTxnStatement(name) {
		return
	}

	QueryDuration.WithLabelValues(name, msg).Observe(dur.Seconds())

	if errVal != nil {
		// errVal might be error or string depending on pgx version
		QueryErrors.WithLabelValues(name, msg).Inc()
	}
}

// isTxnStatement returns true if the SQL is a transaction control command.
func isTxnStatement(sql string) bool {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return false
	}
	// Extract the first word
	first := strings.ToLower(strings.Fields(trimmed)[0])
	switch first {
	case "begin", "start", "commit", "rollback", "savepoint", "release", "prepare", "end":
		return true
	default:
		return false
	}
}
