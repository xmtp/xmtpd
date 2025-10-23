package utils

import (
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log.go defines the logger names and good practices for logging in the system.
// A good logging system is:
// - Easy to understand and debug using 3rd party tools or debug logging in development.
// - Unifies logging across different components of the system, including fields, errors and subsystems.
//
// Good practices for logging include:
// - Logger instances are named logger, for uniformity across the codebase.
// - Build the logger in a way the name chain has readable and meaningful names, such as "xmtpd.api.publish-worker".
// - Prefer circuit breaking in hot paths, such as "if logger.Core().Enabled(zap.DebugLevel) { ... }".
// - Always circuit break when using zap.Any, as it relies on reflection.
// - As a general thumb rule, prefer Debug on hot paths.
// - For repeating fields across different log messages in the same package, use a const.
// - For repeating fields across different packages, use a field function.
// - For repeating errors, use a constant error message.
// - Fields in zap.Field or errors should be snake_case.
// - Avoid duplicate fields in the log messages.
// - Log and error messages always start with lowercase.

const (
	// Base services.
	BaseLoggerName                 = "xmtpd"
	APILoggerName                  = "api"
	DatabaseSubscriptionLoggerName = "database-subscription"
	NodeRegistryWatchdogLoggerName = "node-registry-watchdog"
	PublishWorkerName              = "publish-worker"
	SubscribeWorkerLoggerName      = "subscribe-worker"
	SyncLoggerName                 = "sync"
	SyncWorkerName                 = "sync-worker"
	ContractRatesFetcherLoggerName = "contract-rates-fetcher"
	MetricsLoggerName              = "metrics"
	MisbehaviorLoggerName          = "misbehavior"

	// Gateway services.
	GatewayLoggerName             = "gateway"
	BlockchainPublisherLoggerName = "blockchain-publisher"

	// Indexer.
	IndexerLoggerName                    = "indexer"
	AppChainIndexerLoggerName            = "app-chain"
	GroupMessageBroadcasterLoggerName    = "group-message-broadcaster"
	IdentityUpdateBroadcasterLoggerName  = "identity-update-broadcaster"
	PayerReportContractLoggerName        = "payer-report"
	PayerReportManagerContractLoggerName = "payer-report-manager"
	RPCLogStreamerLoggerName             = "rpc-log-streamer"
	StorerLoggerName                     = "storer"
	ReorgHandlerLoggerName               = "reorg-handler"
	SettlementChainAdminLoggerName       = "settlement-chain-admin"
	SettlementChainIndexerLoggerName     = "settlement-chain"

	// Nonce managers.
	RedisNonceManagerLoggerName = "redis-nonce-manager"
	SQLNonceManagerLoggerName   = "sql-nonce-manager"

	// Migrator.
	MigratorLoggerName            = "migrator"
	MigratorReaderLoggerName      = "reader"
	MigratorTransformerLoggerName = "transformer"
	MigratorWriterLoggerName      = "writer"

	// On-chain protocol services.
	AppChainAdminLoggerName           = "app-chain-admin"
	FundsAdminLoggerName              = "funds-admin"
	NodeRegistryAdminLoggerName       = "node-registry-admin"
	NodeRegistryCallerLoggerName      = "node-registry-caller"
	ParameterAdminLoggerName          = "parameter-registry-admin"
	RateRegistryAdminLoggerName       = "rate-registry-admin"
	PayerReportManagerAdminLoggerName = "payer-report-manager-admin"

	// Payer report subsystem.
	LedgerLoggerName                       = "ledger"
	PayerReportMainLoggerName              = "payer-report"
	PayerReportStoreLoggerName             = "store"
	PayerReportAttestationWorkerLoggerName = "attestation-worker"
	PayerReportGeneratorWorkerLoggerName   = "generator-worker"
	PayerReportSettlementWorkerLoggerName  = "settlement-worker"
	PayerReportSubmitterWorkerLoggerName   = "submitter-worker"

	// Misc and utilities.
	PrunerLoggerName             = "pruner-worker"
	StressChainWatcherLoggerName = "chain-watcher"
)

/* BuildLogger */

func BuildLogger(options config.LogOptions) (*zap.Logger, *zap.Config, error) {
	atom := zap.NewAtomicLevel()
	level := zapcore.InfoLevel
	err := level.Set(options.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	atom.SetLevel(level)

	cfg := zap.Config{
		Encoding:         options.LogEncoding,
		Level:            atom,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        "time",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			NameKey:        "caller",
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	return logger, &cfg, nil
}

/* Fields */

func AddressField(address string) zap.Field {
	return zap.String("address", address)
}

func AmountField(amount string) zap.Field {
	return zap.String("amount", amount)
}

func AppChainChainIDField(chainID int64) zap.Field {
	return zap.Int64("app_chain_id", chainID)
}

func BalanceField(balance string) zap.Field {
	return zap.String("balance", balance)
}

func BlockNumberField(blockNumber uint64) zap.Field {
	return zap.Uint64("block_number", blockNumber)
}

// BodyField uses reflection to log the body of the request or response.
// Do not use in hot paths unless guarded with "if logger.Core().Enabled(zap.DebugLevel) { ... }"
func BodyField(body interface{}) zap.Field {
	return zap.Any("body", body)
}

func ChainIDField(chainID int64) zap.Field {
	return zap.Int64("chain_id", chainID)
}

func ContractAddressField(contractAddress string) zap.Field {
	return zap.String("contract_address", contractAddress)
}

func CountField(count int64) zap.Field {
	return zap.Int64("count", count)
}

func DurationMsField(waitTime time.Duration) zap.Field {
	return zap.Int64("duration_ms", waitTime.Milliseconds())
}

func EnvelopeIDField(envelopeID int64) zap.Field {
	return zap.Int64("envelope_id", envelopeID)
}

func EventField(event string) zap.Field {
	return zap.String("event", event)
}

func EventIDField(eventID string) zap.Field {
	return zap.String("event_id", eventID)
}

func FromAddressField(fromAddress string) zap.Field {
	return zap.String("from_address", fromAddress)
}

func GroupIDField(groupID string) zap.Field {
	return zap.String("group_id", groupID)
}

func HashField(hash string) zap.Field {
	return zap.String("hash", hash)
}

func InboxIDField(inboxID string) zap.Field {
	return zap.String("inbox_id", inboxID)
}

func LastProcessedField(lastProcessed int64) zap.Field {
	return zap.Int64("last_processed", lastProcessed)
}

func LastSequenceIDField(lastSequenceID int64) zap.Field {
	return zap.Int64("last_sequence_id", lastSequenceID)
}

func LimitField(limit uint8) zap.Field {
	return zap.Uint8("limit", limit)
}

func MethodField(method string) zap.Field {
	return zap.String("method", method)
}

func NodeHTTPAddressField(nodeHTTPAddress string) zap.Field {
	return zap.String("node_http_address", nodeHTTPAddress)
}

func NodeOwnerField(nodeOwner string) zap.Field {
	return zap.String("node_owner", nodeOwner)
}

func NodeSigningPublicKeyField(nodeSigningPublicKey string) zap.Field {
	return zap.String("node_signing_public_key", nodeSigningPublicKey)
}

func NonceField(nonce uint64) zap.Field {
	return zap.Uint64("nonce", nonce)
}

func NumEnvelopesField(numEnvelopes int) zap.Field {
	return zap.Int("num_envelopes", numEnvelopes)
}

func NumNoncesField(numNonces int32) zap.Field {
	return zap.Int32("num_nonces", numNonces)
}

func NumResponsesField(numResponses int) zap.Field {
	return zap.Int("num_responses", numResponses)
}

func NumRowsField(numRows int32) zap.Field {
	return zap.Int32("num_rows", numRows)
}

func NumTopicsField(numTopics int) zap.Field {
	return zap.Int("num_topics", numTopics)
}

func OriginatorIDField(originatorID uint32) zap.Field {
	return zap.Uint32("originator_id", originatorID)
}

func PayerAddressField(payerAddress string) zap.Field {
	return zap.String("payer_address", payerAddress)
}

func PayerIDField(payerID int32) zap.Field {
	return zap.Int32("payer_id", payerID)
}

func PayerInfoField(payerInfo *metadata_api.GetPayerInfoResponse_PayerInfo) zap.Field {
	return zap.Any("payer_info", payerInfo)
}

func PayerReportIDField(reportID string) zap.Field {
	return zap.String("report_id", reportID)
}

func PublicKeyField(publicKey string) zap.Field {
	return zap.String("public_key", publicKey)
}

func RecipientField(recipient string) zap.Field {
	return zap.String("recipient", recipient)
}

func SequenceIDField(sequenceID int64) zap.Field {
	return zap.Int64("sequence_id", sequenceID)
}

func SettlementChainChainIDField(chainID int64) zap.Field {
	return zap.Int64("settlement_chain_id", chainID)
}

func StartingNonceField(startingNonce uint64) zap.Field {
	return zap.Uint64("starting_nonce", startingNonce)
}

func StartSequenceIDField(startSequenceID int64) zap.Field {
	return zap.Int64("start_sequence_id", startSequenceID)
}

func TimeField(time time.Time) zap.Field {
	return zap.Time("time", time)
}

func ToAddressField(toAddress string) zap.Field {
	return zap.String("to", toAddress)
}

func TopicField(topic string) zap.Field {
	return zap.String("topic", topic)
}
