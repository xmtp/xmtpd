package node

import (
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log           zap.Options          `group:"Log options" namespace:"log"`
	API           apigateway.Options   `group:"API options" namespace:"api"`
	Store         StoreOptions         `group:"Store options" namespace:"store"`
	OpenTelemetry OpenTelemetryOptions `group:"OpenTelemetry options" namespace:"otel"`
}

type StoreOptions struct {
	Type     string                `long:"type" description:"type of storage to use" choice:"mem" choice:"postgres" choice:"bolt" default:"mem"`
	Postgres postgresstore.Options `group:"Store Postgres options" namespace:"postgres"`
	Bolt     bolt.Options          `group:"Store Bolt options" namespace:"bolt"`
}
