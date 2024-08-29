package config

import (
	"time"
)

type ApiOptions struct {
	Port int `short:"p" long:"port" description:"Port to listen on" default:"5050"`
}

type ContractsOptions struct {
	RpcUrl                  string        `long:"rpc-url"          description:"Blockchain RPC URL"`
	NodesContractAddress    string        `long:"nodes-address"    description:"Node contract address"`
	MessagesContractAddress string        `long:"messages-address" description:"Message contract address"`
	ChainID                 int           `long:"chain-id"         description:"Chain ID for the appchain"               default:"31337"`
	RefreshInterval         time.Duration `long:"refresh-interval" description:"Refresh interval for the nodes registry" default:"60s"`
}

type DbOptions struct {
	ReaderConnectionString string        `long:"reader-connection-string" description:"Reader connection string"`
	WriterConnectionString string        `long:"writer-connection-string" description:"Writer connection string"                       required:"true"`
	ReadTimeout            time.Duration `long:"read-timeout"             description:"Timeout for reading from the database"                          default:"10s"`
	WriteTimeout           time.Duration `long:"write-timeout"            description:"Timeout for writing to the database"                            default:"10s"`
	MaxOpenConns           int           `long:"max-open-conns"           description:"Maximum number of open connections"                             default:"80"`
	WaitForDB              time.Duration `long:"wait-for"                 description:"wait for DB on start, up to specified duration"`
}

type PayerOptions struct {
	PrivateKey string `long:"private-key" description:"Private key used to sign blockchain transactions"`
}

type ServerOptions struct {
	LogLevel string `short:"l" long:"log-level"    description:"Define the logging level, supported strings are: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL, and their lower-case forms." default:"INFO"`
	//nolint:staticcheck
	LogEncoding string `          long:"log-encoding" description:"Log encoding format. Either console or json"                                                                                  default:"console" choice:"console"`

	SignerPrivateKey string `long:"signer-private-key" description:"Private key used to sign messages"`

	API       ApiOptions       `group:"API Options"       namespace:"api"`
	DB        DbOptions        `group:"Database Options"  namespace:"db"`
	Contracts ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Payer     PayerOptions     `group:"Payer Options"     namespace:"payer"`
}
