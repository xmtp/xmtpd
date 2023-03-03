package node

import (
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log           zap.Options          `group:"Log options" namespace:"log"`
	P2P           P2POptions           `group:"P2P options" namespace:"p2p"`
	API           apigateway.Options   `group:"API options" namespace:"api"`
	Store         StoreOptions         `group:"Store options" namespace:"store"`
	OpenTelemetry OpenTelemetryOptions `group:"OpenTelemetry options" namespace:"otel"`
}

type P2POptions struct {
	Port        uint   `long:"port" description:"P2P listen port" default:"9000"`
	IdentityKey string `long:"identity-key" description:"Identity private key in hex format" default:""`
}

//nolint:staticcheck // it doesn't like the multiple "choice" struct tags
type StoreOptions struct {
	Type     string                `long:"type" description:"type of storage to use" choice:"mem" choice:"postgres" choice:"bolt" default:"mem"`
	Postgres postgresstore.Options `group:"Store Postgres options" namespace:"postgres"`
	Bolt     bolt.Options          `group:"Store Bolt options" namespace:"bolt"`
}
