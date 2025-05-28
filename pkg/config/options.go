package config

import (
	"time"
)

type ApiOptions struct {
	Port     int `short:"p" long:"port"      description:"Port to listen on"      env:"XMTPD_API_PORT"      default:"5050"`
	HTTPPort int `          long:"http-port" description:"HTTP Port to listen on" env:"XMTPD_HTTP_API_PORT" default:"5055"`
}

type NodesContractOption struct {
	NodesContractAddress string `long:"nodes-address" env:"XMTPD_CONTRACTS_NODES_ADDRESS" description:"Node contract address"`
}

type ContractsOptions struct {
	// Deprecated: Use AppChainOptions and SettlementChainOptions instead.
	// TODO: Remove.
	RpcUrl                         string        `long:"rpc-url"                   env:"XMTPD_CONTRACTS_RPC_URL"                   description:"Blockchain RPC URL"`
	MessagesContractAddress        string        `long:"messages-address"          env:"XMTPD_CONTRACTS_MESSAGES_ADDRESS"          description:"Message contract address"`
	IdentityUpdatesContractAddress string        `long:"identity-updates-address"  env:"XMTPD_CONTRACTS_IDENTITY_UPDATES_ADDRESS"  description:"Identity updates contract address"`
	RateRegistryContractAddress    string        `long:"rates-registry-address"    env:"XMTPD_CONTRACTS_RATES_MANAGER_ADDRESS"     description:"Rates manager contract address"`
	ChainID                        int           `long:"chain-id"                  env:"XMTPD_CONTRACTS_CHAIN_ID"                  description:"Chain ID for the appchain"                                    default:"31337"`
	RegistryRefreshInterval        time.Duration `long:"registry-refresh-interval" env:"XMTPD_CONTRACTS_REGISTRY_REFRESH_INTERVAL" description:"Refresh interval for the nodes registry"                      default:"60s"`
	RatesRefreshInterval           time.Duration `long:"rates-refresh-interval"    env:"XMTPD_CONTRACTS_RATES_REFRESH_INTERVAL"    description:"Refresh interval for the rates manager"                       default:"300s"`
	MaxChainDisconnectTime         time.Duration `long:"max-chain-disconnect-time" env:"XMTPD_CONTRACTS_MAX_CHAIN_DISCONNECT_TIME" description:"Maximum time to allow the node to operate while disconnected" default:"300s"`
	NodesContract                  NodesContractOption

	BackfillBlockSize uint64 `long:"backfill-block-size" env:"XMTPD_CONTRACTS_BACKFILL_BLOCK_SIZE" description:"Maximal size of a backfill block" default:"500"`

	// New options for app and settlement chains multi-chain deployments.
	AppChain        AppChainOptions        `group:"Application Chain Options" namespace:"app-chain"`
	SettlementChain SettlementChainOptions `group:"Settlement Chain Options"  namespace:"settlement-chain"`
}

type AppChainOptions struct {
	RpcURL                           string        `long:"rpc-url"                             env:"XMTPD_APP_CHAIN_RPC_URL"                           description:"Blockchain RPC URL"`
	ChainID                          int           `long:"chain-id"                            env:"XMTPD_APP_CHAIN_CHAIN_ID"                          description:"Chain ID for the application chain"                           default:"31337"`
	MaxChainDisconnectTime           time.Duration `long:"max-chain-disconnect-time"           env:"XMTPD_APP_CHAIN_MAX_CHAIN_DISCONNECT_TIME"         description:"Maximum time to allow the node to operate while disconnected" default:"300s"`
	GroupMessageBroadcasterAddress   string        `long:"group-message-broadcaster-address"   env:"XMTPD_APP_CHAIN_GROUP_MESSAGE_BROADCAST_ADDRESS"   description:"Group message broadcaster contract address"`
	IdentityUpdateBroadcasterAddress string        `long:"identity-update-broadcaster-address" env:"XMTPD_APP_CHAIN_IDENTITY_UPDATE_BROADCAST_ADDRESS" description:"Identity update broadcaster contract address"`
}

type SettlementChainOptions struct {
	RpcURL                      string        `long:"rpc-url"                        env:"XMTPD_SETTLEMENT_CHAIN_RPC_URL"                        description:"Blockchain RPC URL"`
	ChainID                     int           `long:"chain-id"                       env:"XMTPD_SETTLEMENT_CHAIN_CHAIN_ID"                       description:"Chain ID for the settlement chain"                            default:"31337"`
	NodeRegistryAddress         string        `long:"node-registry-address"          env:"XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS"          description:"Node Registry contract address"`
	NodeRegistryRefreshInterval time.Duration `long:"node-registry-refresh-interval" env:"XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_REFRESH_INTERVAL" description:"Refresh interval for the nodes registry"                      default:"60s"`
	RateRegistryAddress         string        `long:"rate-registry-address"          env:"XMTPD_SETTLEMENT_CHAIN_RATE_REGISTRY_ADDRESS"          description:"Rate registry contract address"`
	RateRegistryRefreshInterval time.Duration `long:"rate-registry-refresh-interval" env:"XMTPD_SETTLEMENT_CHAIN_RATE_REGISTRY_REFRESH_INTERVAL" description:"Refresh interval for the rate registry"                       default:"300s"`
	ParameterRegistryAddress    string        `long:"parameter-registry-address"     env:"XMTPD_SETTLEMENT_CHAIN_PARAMETER_REGISTRY_ADDRESS"     description:"Parameter Registry contract address"`
	PayerRegistryAddress        string        `long:"payer-registry-address"         env:"XMTPD_SETTLEMENT_CHAIN_PAYER_REGISTRY_ADDRESS"         description:"Payer Registry contract address"`
	PayerReportManagerAddress   string        `long:"payer-report-manager-address"   env:"XMTPD_SETTLEMENT_CHAIN_PAYER_REPORT_MANAGER_ADDRESS"   description:"Payer Report Manager contract address"`
	MaxChainDisconnectTime      time.Duration `long:"max-chain-disconnect-time"      env:"XMTPD_SETTLEMENT_CHAIN_MAX_CHAIN_DISCONNECT_TIME"      description:"Maximum time to allow the node to operate while disconnected" default:"300s"`
}

type DbOptions struct {
	ReaderConnectionString string        `long:"reader-connection-string" env:"XMTPD_DB_READER_CONNECTION_STRING" description:"Reader connection string"`
	WriterConnectionString string        `long:"writer-connection-string" env:"XMTPD_DB_WRITER_CONNECTION_STRING" description:"Writer connection string"`
	ReadTimeout            time.Duration `long:"read-timeout"             env:"XMTPD_DB_READ_TIMEOUT"             description:"Timeout for reading from the database"          default:"10s"`
	WriteTimeout           time.Duration `long:"write-timeout"            env:"XMTPD_DB_WRITE_TIMEOUT"            description:"Timeout for writing to the database"            default:"10s"`
	MaxOpenConns           int           `long:"max-open-conns"           env:"XMTPD_DB_MAX_OPEN_CONNS"           description:"Maximum number of open connections"             default:"80"`
	WaitForDB              time.Duration `long:"wait-for"                 env:"XMTPD_DB_WAIT_FOR"                 description:"wait for DB on start, up to specified duration" default:"30s"`
	NameOverride           string        `long:"name-override"            env:"XMTPD_DB_NAME_OVERRIDE"            description:"Override the automatically generated DB name"                 hidden:"true"`
}

type IndexerOptions struct {
	Enable bool `long:"enable" env:"XMTPD_INDEXER_ENABLE" description:"Enable the indexer"`
}

// MetricsOptions are settings used to start a prometheus server
type MetricsOptions struct {
	Enable  bool   `long:"enable"          env:"XMTPD_METRICS_ENABLE"          description:"Enable the metrics server"`
	Address string `long:"metrics-address" env:"XMTPD_METRICS_METRICS_ADDRESS" description:"Listening address of the metrics server"   default:"127.0.0.1"`
	Port    int    `long:"metrics-port"    env:"XMTPD_METRICS_METRICS_PORT"    description:"Listening HTTP port of the metrics server" default:"8008"`
}

type PayerOptions struct {
	PrivateKey string `long:"private-key" env:"XMTPD_PAYER_PRIVATE_KEY" description:"Private key used to sign blockchain transactions"`
	Enable     bool   `long:"enable"      env:"XMTPD_PAYER_ENABLE"      description:"Enable the payer API"`
}

type ReplicationOptions struct {
	Enable                bool          `long:"enable"                   env:"XMTPD_REPLICATION_ENABLE"           description:"Enable the replication API"`
	SendKeepAliveInterval time.Duration `long:"send-keep-alive-interval" env:"XMTPD_API_SEND_KEEP_ALIVE_INTERVAL" description:"Send empty application level keepalive package interval" default:"30s"`
}

type SyncOptions struct {
	Enable bool `long:"enable" env:"XMTPD_SYNC_ENABLE" description:"Enable the sync server"`
}

type MlsValidationOptions struct {
	GrpcAddress string `long:"grpc-address" env:"XMTPD_MLS_VALIDATION_GRPC_ADDRESS" description:"Address of the MLS validation service"`
}

// TracingOptions are settings controlling collection of DD APM traces and error tracking.
type TracingOptions struct {
	Enable bool `long:"enable" env:"XMTPD_TRACING_ENABLE" description:"Enable DD APM trace collection"`
}

// ReflectionOptions are settings controlling collection of GRPC reflection settings.
type ReflectionOptions struct {
	Enable bool `long:"enable" env:"XMTPD_REFLECTION_ENABLE" description:"Enable GRPC reflection"`
}

type LogOptions struct {
	LogLevel    string `short:"l" long:"log-level"    env:"XMTPD_LOG_LEVEL"    description:"Define the logging level, supported strings are: DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL, and their lower-case forms." default:"INFO"`
	LogEncoding string `          long:"log-encoding" env:"XMTPD_LOG_ENCODING" description:"Log encoding format. Either console or json"                                                                                  default:"console"`
}

type SignerOptions struct {
	PrivateKey string `long:"private-key" env:"XMTPD_SIGNER_PRIVATE_KEY" description:"Private key used to sign messages"`
}

type ServerOptions struct {
	API           ApiOptions           `group:"API Options"            namespace:"api"`
	Contracts     ContractsOptions     `group:"Contracts Options"      namespace:"contracts"`
	DB            DbOptions            `group:"Database Options"       namespace:"db"`
	Log           LogOptions           `group:"Log Options"            namespace:"log"`
	Indexer       IndexerOptions       `group:"Indexer Options"        namespace:"indexer"`
	Metrics       MetricsOptions       `group:"Metrics Options"        namespace:"metrics"`
	MlsValidation MlsValidationOptions `group:"MLS Validation Options" namespace:"mls-validation"`
	Payer         PayerOptions         `group:"Payer Options"          namespace:"payer"`
	Reflection    ReflectionOptions    `group:"Reflection Options"     namespace:"reflection"`
	Replication   ReplicationOptions   `group:"Replication Options"    namespace:"replication"`
	Signer        SignerOptions        `group:"Signer Options"         namespace:"signer"`
	Sync          SyncOptions          `group:"Sync Options"           namespace:"sync"`
	Tracing       TracingOptions       `group:"DD APM Tracing Options" namespace:"tracing"`
	Version       bool                 `                                                          short:"v" long:"version" description:"Output binary version and exit"`
}
