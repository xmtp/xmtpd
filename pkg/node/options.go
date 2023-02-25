package node

import (
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log           zap.Options          `group:"Log options" namespace:"log"`
	API           apigateway.Options   `group:"API options" namespace:"api"`
	Store         StoreOptions         `group:"Store options" namespace:"store"`
	OpenTelemetry OpenTelemetryOptions `group:"OpenTelemetry options" namespace:"otel"`
}

type StoreOptions struct {
	Postgres PostgresOptions `group:"Postgres options" namespace:"postgres"`
}

type PostgresOptions struct {
	DSN string `long:"dsn" env:"POSTGRES_DSN" description:"Postgres connection string" default:""`
}
