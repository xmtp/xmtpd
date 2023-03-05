package orbitdbnode

import (
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/otel"
	"github.com/xmtp/xmtpd/pkg/store/bolt"
	postgresstore "github.com/xmtp/xmtpd/pkg/store/postgres"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log           zap.Options        `group:"Log options" namespace:"log"`
	P2P           P2POptions         `group:"P2P options" namespace:"p2p"`
	API           apigateway.Options `group:"API options" namespace:"api"`
	Store         StoreOptions       `group:"Store options" namespace:"store"`
	OpenTelemetry otel.Options       `group:"OpenTelemetry options" namespace:"otel"`
}

type P2POptions struct {
	Port            uint     `long:"port" description:"P2P listen port" default:"0"`
	NodeKey         string   `long:"node-key" env:"XMTP_NODE_KEY" description:"P2P node identity private key in hex format" default:""`
	PersistentPeers []string `long:"persistent-peer" description:"P2P persistent peers"`
}

//nolint:staticcheck // it doesn't like the multiple "choice" struct tags
type StoreOptions struct {
	Type     string                `long:"type" description:"type of storage to use" choice:"mem" choice:"postgres" choice:"bolt" default:"mem"`
	Postgres postgresstore.Options `group:"Store Postgres options" namespace:"postgres"`
	Bolt     bolt.Options          `group:"Store Bolt options" namespace:"bolt"`
}
