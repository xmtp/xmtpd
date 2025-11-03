// Package config implements the configuration options for the xmtpd server.
package config

import (
	"time"

	"github.com/xmtp/xmtpd/pkg/config/environments"
)

type APIOptions struct {
	Enable                bool          `long:"enable"                   env:"XMTPD_API_ENABLE"                   description:"Enable the client API"`
	SendKeepAliveInterval time.Duration `long:"send-keep-alive-interval" env:"XMTPD_API_SEND_KEEP_ALIVE_INTERVAL" description:"Send empty application level keepalive package interval" default:"30s"`
	Port                  int           `long:"port"                     env:"XMTPD_API_PORT"                     description:"Port to listen on"                                       default:"5050" short:"p"`
}

type ContractsOptions struct {
	AppChain        AppChainOptions        `group:"Application Chain Options" namespace:"app-chain"`
	SettlementChain SettlementChainOptions `group:"Settlement Chain Options"  namespace:"settlement-chain"`

	ConfigFilePath string                                `long:"config-file-path" env:"XMTPD_CONTRACTS_CONFIG_FILE_PATH" description:"Path to the JSON contracts config file"`
	ConfigJSON     string                                `long:"config-json"      env:"XMTPD_CONTRACTS_CONFIG_JSON"      description:"JSON contracts config"`
	Environment    environments.SmartContractEnvironment `long:"environment"      env:"XMTPD_CONTRACTS_ENVIRONMENT"      description:"Deployed environment to load contracts config for"`
}

type AppChainOptions struct {
	RPCURL                           string        `long:"rpc-url"                             env:"XMTPD_APP_CHAIN_RPC_URL"                           description:"Blockchain RPC URL"`
	WssURL                           string        `long:"wss-url"                             env:"XMTPD_APP_CHAIN_WSS_URL"                           description:"Blockchain WSS URL"`
	ChainID                          int64         `long:"chain-id"                            env:"XMTPD_APP_CHAIN_CHAIN_ID"                          description:"Chain ID for the application chain"                           default:"31337"`
	MaxChainDisconnectTime           time.Duration `long:"max-chain-disconnect-time"           env:"XMTPD_APP_CHAIN_MAX_CHAIN_DISCONNECT_TIME"         description:"Maximum time to allow the node to operate while disconnected" default:"60s"`
	BackfillBlockPageSize            uint64        `long:"backfill-block-page-size"            env:"XMTPD_APP_CHAIN_BACKFILL_BLOCK_PAGE_SIZE"          description:"Maximal size of a backfill block page"                        default:"500"`
	GroupMessageBroadcasterAddress   string        `long:"group-message-broadcaster-address"   env:"XMTPD_APP_CHAIN_GROUP_MESSAGE_BROADCAST_ADDRESS"   description:"Group message broadcaster contract address"`
	IdentityUpdateBroadcasterAddress string        `long:"identity-update-broadcaster-address" env:"XMTPD_APP_CHAIN_IDENTITY_UPDATE_BROADCAST_ADDRESS" description:"Identity update broadcaster contract address"`
	GatewayAddress                   string        `long:"gateway-address"                     env:"XMTPD_APP_CHAIN_GATEWAY_ADDRESS"                   description:"App Chain Gateway contract address"`
	ParameterRegistryAddress         string        `long:"parameter-registry-address"          env:"XMTPD_APP_CHAIN_PARAMETER_REGISTRY_ADDRESS"        description:"Parameter Registry contract address"`
	DeploymentBlock                  uint64        `long:"deployment-block"                    env:"XMTPD_APP_CHAIN_DEPLOYMENT_BLOCK"                  description:"Deployment block for the application chain"                   default:"0"`
	MaxBlockchainPayloadSize         uint64        `long:"max-blockchain-payload-size"         env:"XMTPD_APP_CHAIN_MAX_BLOCKCHAIN_PAYLOAD_SIZE"       description:"Max payload size for the App Chain in bytes"                  default:"200000"`
}

type SettlementChainOptions struct {
	RPCURL                      string        `long:"rpc-url"                        env:"XMTPD_SETTLEMENT_CHAIN_RPC_URL"                        description:"Blockchain RPC URL"`
	WssURL                      string        `long:"wss-url"                        env:"XMTPD_SETTLEMENT_CHAIN_WSS_URL"                        description:"Blockchain WSS URL"`
	ChainID                     int64         `long:"chain-id"                       env:"XMTPD_SETTLEMENT_CHAIN_CHAIN_ID"                       description:"Chain ID for the settlement chain"                            default:"31337"`
	MaxChainDisconnectTime      time.Duration `long:"max-chain-disconnect-time"      env:"XMTPD_SETTLEMENT_CHAIN_MAX_CHAIN_DISCONNECT_TIME"      description:"Maximum time to allow the node to operate while disconnected" default:"300s"`
	BackfillBlockPageSize       uint64        `long:"backfill-block-page-size"       env:"XMTPD_SETTLEMENT_CHAIN_BACKFILL_BLOCK_PAGE_SIZE"       description:"Maximal size of a backfill block page"                        default:"500"`
	NodeRegistryAddress         string        `long:"node-registry-address"          env:"XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_ADDRESS"          description:"Node Registry contract address"`
	NodeRegistryRefreshInterval time.Duration `long:"node-registry-refresh-interval" env:"XMTPD_SETTLEMENT_CHAIN_NODE_REGISTRY_REFRESH_INTERVAL" description:"Refresh interval for the nodes registry"                      default:"60s"`
	RateRegistryAddress         string        `long:"rate-registry-address"          env:"XMTPD_SETTLEMENT_CHAIN_RATE_REGISTRY_ADDRESS"          description:"Rate registry contract address"`
	RateRegistryRefreshInterval time.Duration `long:"rate-registry-refresh-interval" env:"XMTPD_SETTLEMENT_CHAIN_RATE_REGISTRY_REFRESH_INTERVAL" description:"Refresh interval for the rate registry"                       default:"300s"`
	ParameterRegistryAddress    string        `long:"parameter-registry-address"     env:"XMTPD_SETTLEMENT_CHAIN_PARAMETER_REGISTRY_ADDRESS"     description:"Parameter Registry contract address"`
	PayerRegistryAddress        string        `long:"payer-registry-address"         env:"XMTPD_SETTLEMENT_CHAIN_PAYER_REGISTRY_ADDRESS"         description:"Payer Registry contract address"`
	PayerReportManagerAddress   string        `long:"payer-report-manager-address"   env:"XMTPD_SETTLEMENT_CHAIN_PAYER_REPORT_MANAGER_ADDRESS"   description:"Payer Report Manager contract address"`
	GatewayAddress              string        `long:"gateway-address"                env:"XMTPD_SETTLEMENT_CHAIN_GATEWAY_ADDRESS"                description:"Settlement Chain Gateway contract address"`
	DistributionManagerAddress  string        `long:"distribution-manager-address"   env:"XMTPD_SETTLEMENT_CHAIN_DISTRIBUTION_MANAGER_ADDRESS"   description:"Distribution Manager contract address"`
	DeploymentBlock             uint64        `long:"deployment-block"               env:"XMTPD_SETTLEMENT_CHAIN_DEPLOYMENT_BLOCK"               description:"Deployment block for the settlement chain"                    default:"0"`
	UnderlyingFeeToken          string        `long:"underlying-fee-token-address"   env:"XMTPD_SETTLEMENT_CHAIN_UNDERLYING_FEE_TOKEN_ADDRESS"   description:"Underlying fee token address"`
	FeeToken                    string        `long:"fee-token-address"              env:"XMTPD_SETTLEMENT_CHAIN_FEE_TOKEN_ADDRESS"              description:"Fee token address"`
}

type DBOptions struct {
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

type DebugOptions struct {
	Enable bool `long:"enable" env:"XMTPD_DEBUG_ENABLE" description:"Enable the pprof debug server"`
	Port   int  `long:"port"   env:"XMTPD_DEBUG_PORT"   description:"Port to listen on"             default:"6060"`
}

type PayerOptions struct {
	PrivateKey string `long:"private-key" env:"XMTPD_PAYER_PRIVATE_KEY" description:"Private key used to sign blockchain transactions"`
	Enable     bool   `long:"enable"      env:"XMTPD_PAYER_ENABLE"      description:"Enable the payer API"`
}

type ReplicationOptions struct {
	// Deprecated: use API instead
	Enable bool `long:"enable" env:"XMTPD_REPLICATION_ENABLE" description:"Enable the replication API"`
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

type PayerReportOptions struct {
	Enable                        bool          `long:"enable"                           env:"XMTPD_PAYER_REPORT_RUN_WORKERS"                      description:"Run workers to submit and attest to payer reports"`
	AttestationWorkerPollInterval time.Duration `long:"attestation-worker-poll-interval" env:"XMTPD_PAYER_REPORT_ATTESTATION_WORKER_POLL_INTERVAL" description:"Interval between checks for new reports to attest to"                            default:"1m"`
	GenerateReportSelfPeriod      time.Duration `long:"generate-report-self-period"      env:"XMTPD_PAYER_REPORT_GENERATE_REPORT_SELF_PERIOD"      description:"Period of time we should use to generate payer reports about self"               default:"6h"`
	GenerateReportOthersPeriod    time.Duration `long:"generate-report-others-period"    env:"XMTPD_PAYER_REPORT_GENERATE_REPORT_OTHERS_PERIOD"    description:"Period of time we should use to generate backup payer reports about other nodes" default:"12h"`
}

type ServerOptions struct {
	API             APIOptions             `group:"API Options"              namespace:"api"`
	Contracts       ContractsOptions       `group:"Contracts Options"        namespace:"contracts"`
	DB              DBOptions              `group:"Database Options"         namespace:"db"`
	Log             LogOptions             `group:"Log Options"              namespace:"log"`
	Indexer         IndexerOptions         `group:"Indexer Options"          namespace:"indexer"`
	Metrics         MetricsOptions         `group:"Metrics Options"          namespace:"metrics"`
	MlsValidation   MlsValidationOptions   `group:"MLS Validation Options"   namespace:"mls-validation"`
	PayerReport     PayerReportOptions     `group:"Payer Report Options"     namespace:"payer-report"`
	Reflection      ReflectionOptions      `group:"Reflection Options"       namespace:"reflection"`
	Replication     ReplicationOptions     `group:"API Options"              namespace:"replication"`
	Signer          SignerOptions          `group:"Signer Options"           namespace:"signer"`
	Sync            SyncOptions            `group:"Sync Options"             namespace:"sync"`
	Tracing         TracingOptions         `group:"DD APM Tracing Options"   namespace:"tracing"`
	MigrationServer MigrationServerOptions `group:"Migration Server Options" namespace:"migration-server"`
	MigrationClient MigrationClientOptions `group:"Migration Client Options" namespace:"migration-client"`
	Debug           DebugOptions           `group:"Debug Options"            namespace:"debug"`
	Version         bool                   `                                                              short:"v" long:"version" description:"Output binary version and exit"`
}
