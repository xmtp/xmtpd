package config

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type OptionsValidator struct {
	logger *zap.Logger
}

func NewOptionsValidator(logger *zap.Logger) *OptionsValidator {
	return &OptionsValidator{logger: logger}
}

func (v *OptionsValidator) ValidateServerOptions(options *ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	err := v.ParseJSONConfig(&options.Contracts)
	if err != nil {
		return err
	}

	v.validateBlockchainConfig(options, missingSet, customSet)

	v.validateMigrationOptions(options, missingSet, customSet)

	v.validateField(
		options.DB.WriterConnectionString,
		"db.writer-connection-string",
		missingSet,
	)

	if options.API.Enable {
		v.validateField(options.Signer.PrivateKey, "signer.private-key", missingSet)
		v.validateField(
			options.MlsValidation.GrpcAddress,
			"mls-validation.grpc-address",
			missingSet,
		)
	}

	if options.Indexer.Enable {
		v.validateField(
			options.MlsValidation.GrpcAddress,
			"mls-validation.grpc-address",
			missingSet,
		)
	}

	if options.Payer.Enable {
		if err := v.validatePayerOptions(&options.Payer, customSet); err != nil {
			return err
		}
	}

	if len(missingSet) > 0 || len(customSet) > 0 {
		var errs []string
		if len(missingSet) > 0 {

			var errorMessages []string
			for err := range missingSet {
				errorMessages = append(errorMessages, err)
			}
			errs = append(
				errs,
				fmt.Sprintf("Missing required arguments: %s", strings.Join(errorMessages, ", ")),
			)
		}
		if len(customSet) > 0 {
			for err := range customSet {
				errs = append(errs, err)
			}
		}
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}

func (v *OptionsValidator) ValidatePruneOptions(options PruneOptions) error {
	missingSet := make(map[string]struct{})

	if options.DB.WriterConnectionString == "" {
		missingSet["--DB.WriterConnectionString"] = struct{}{}
	}

	if options.Contracts.SettlementChain.NodeRegistryAddress == "" {
		missingSet["--contracts.settlement-chain.node-registry-address"] = struct{}{}
	}

	if options.Signer.PrivateKey == "" {
		missingSet["--signer.private-key"] = struct{}{}
	}

	if len(missingSet) > 0 {
		var errorMessages []string
		for err := range missingSet {
			errorMessages = append(errorMessages, err)
		}

		return fmt.Errorf("missing required arguments: %s", strings.Join(errorMessages, ", "))
	}

	if options.PruneConfig.MaxCycles < 1 {
		return fmt.Errorf("max-cycles must be greater than 0")
	}

	return nil
}

func (v *OptionsValidator) validateMigrationOptions(
	opts *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	if opts.MigrationServer.Enable && opts.MigrationClient.Enable {
		missingSet["--migration-server.enable and --migration-client.enable cannot be used together"] = struct{}{}
	}

	if opts.MigrationServer.Enable {
		v.validateField(
			opts.MigrationServer.PayerPrivateKey,
			"migration-server.payer-private-key",
			missingSet,
		)
		v.validateField(
			opts.MigrationServer.NodeSigningKey,
			"migration-server.node-signing-key",
			missingSet,
		)
		v.validateField(
			opts.Signer.PrivateKey,
			"signer.private-key",
			missingSet,
		)
		v.validateField(
			opts.MigrationServer.ReaderConnectionString,
			"migration-server.reader-connection-string",
			missingSet,
		)
		v.validateField(
			opts.MigrationServer.ReaderTimeout,
			"migration-server.reader-timeout",
			missingSet,
		)
		v.validateAppChainConfig(opts, missingSet, customSet)
	}

	if opts.MigrationClient.Enable {
		v.validateField(
			opts.MigrationClient.FromNodeID,
			"migration-client.from-node-id",
			missingSet,
		)
	}
}

// ContractOptionsFromEnv loads contract options from a file path or URL.
// Deprecated: Use LoadContractsConfig with ContractsSource instead.
func ContractOptionsFromEnv(filePath string) (*ContractsOptions, error) {
	return LoadContractsConfig(ContractsSource{
		FilePath: filePath,
	})
}

func (v *OptionsValidator) ParseJSONConfig(options *ContractsOptions) error {
	// Determine which source to use
	source := ContractsSource{
		Environment: string(options.Environment),
		FilePath:    options.ConfigFilePath,
		JSONData:    options.ConfigJSON,
	}

	// Only load if at least one source is specified
	if source.Environment == "" && source.FilePath == "" && source.JSONData == "" {
		return nil // No config source specified, skip loading
	}

	v.logger.Info("Loading contract configuration",
		zap.String("environment", source.Environment),
		zap.String("filePath", source.FilePath),
		zap.Bool("hasJSON", source.JSONData != ""))

	// Load the configuration using the unified loader
	loadedConfig, err := LoadContractsConfig(source)
	if err != nil {
		return fmt.Errorf("load contract config: %w", err)
	}

	// Merge loaded config with existing options (explicitly set values take precedence)
	v.mergeContractsOptions(options, loadedConfig)

	return nil
}

// mergeContractsOptions merges loaded config with existing options.
// Explicitly specified values in options take precedence over loaded values.
func (v *OptionsValidator) mergeContractsOptions(
	options *ContractsOptions,
	loaded *ContractsOptions,
) {
	// AppChainOptions - only fill in values if not explicitly set
	if options.AppChain.GroupMessageBroadcasterAddress == "" {
		options.AppChain.GroupMessageBroadcasterAddress = loaded.AppChain.GroupMessageBroadcasterAddress
	}
	if options.AppChain.IdentityUpdateBroadcasterAddress == "" {
		options.AppChain.IdentityUpdateBroadcasterAddress = loaded.AppChain.IdentityUpdateBroadcasterAddress
	}
	if options.AppChain.ChainID == 0 || options.AppChain.ChainID == 31337 {
		options.AppChain.ChainID = loaded.AppChain.ChainID
	}
	if options.AppChain.GatewayAddress == "" {
		options.AppChain.GatewayAddress = loaded.AppChain.GatewayAddress
	}
	if options.AppChain.DeploymentBlock == 0 {
		options.AppChain.DeploymentBlock = loaded.AppChain.DeploymentBlock
	}
	if options.AppChain.ParameterRegistryAddress == "" {
		options.AppChain.ParameterRegistryAddress = loaded.AppChain.ParameterRegistryAddress
	}

	// SettlementChainOptions - only fill in values if not explicitly set
	if options.SettlementChain.NodeRegistryAddress == "" {
		options.SettlementChain.NodeRegistryAddress = loaded.SettlementChain.NodeRegistryAddress
	}
	if options.SettlementChain.RateRegistryAddress == "" {
		options.SettlementChain.RateRegistryAddress = loaded.SettlementChain.RateRegistryAddress
	}
	if options.SettlementChain.ParameterRegistryAddress == "" {
		options.SettlementChain.ParameterRegistryAddress = loaded.SettlementChain.ParameterRegistryAddress
	}
	if options.SettlementChain.PayerRegistryAddress == "" {
		options.SettlementChain.PayerRegistryAddress = loaded.SettlementChain.PayerRegistryAddress
	}
	if options.SettlementChain.PayerReportManagerAddress == "" {
		options.SettlementChain.PayerReportManagerAddress = loaded.SettlementChain.PayerReportManagerAddress
	}
	if options.SettlementChain.ChainID == 0 || options.SettlementChain.ChainID == 31337 {
		options.SettlementChain.ChainID = loaded.SettlementChain.ChainID
	}
	if options.SettlementChain.DeploymentBlock == 0 {
		options.SettlementChain.DeploymentBlock = loaded.SettlementChain.DeploymentBlock
	}
	if options.SettlementChain.GatewayAddress == "" {
		options.SettlementChain.GatewayAddress = loaded.SettlementChain.GatewayAddress
	}
	if options.SettlementChain.DistributionManagerAddress == "" {
		options.SettlementChain.DistributionManagerAddress = loaded.SettlementChain.DistributionManagerAddress
	}
	if options.SettlementChain.UnderlyingFeeToken == "" {
		options.SettlementChain.UnderlyingFeeToken = loaded.SettlementChain.UnderlyingFeeToken
	}
	if options.SettlementChain.FeeToken == "" {
		options.SettlementChain.FeeToken = loaded.SettlementChain.FeeToken
	}
}

func (v *OptionsValidator) validateBlockchainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	v.validateAppChainConfig(options, missingSet, customSet)
	v.validateSettlementChainConfig(options, missingSet, customSet)

	if options.API.Enable || options.Sync.Enable {
		v.validateHexAddress(
			options.Contracts.SettlementChain.RateRegistryAddress,
			"contracts.settlement-chain.rate-registry-address",
			missingSet,
		)
		v.validateField(
			options.Contracts.SettlementChain.RateRegistryRefreshInterval,
			"contracts.settlement-chain.rate-registry-refresh-interval",
			customSet,
		)
	}
}

func (v *OptionsValidator) validateAppChainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	v.validateField(
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.chain-id",
		customSet,
	)
	v.validateRPCURL(
		options.Contracts.AppChain.RPCURL,
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.rpc-url",
		missingSet,
	)
	v.validateWebsocketURL(
		options.Contracts.AppChain.WssURL,
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.wss-url",
		missingSet,
	)
	v.validateHexAddress(
		options.Contracts.AppChain.GroupMessageBroadcasterAddress,
		"contracts.app-chain.group-message-broadcaster-address",
		missingSet,
	)
	v.validateHexAddress(
		options.Contracts.AppChain.IdentityUpdateBroadcasterAddress,
		"contracts.app-chain.identity-update-broadcaster-address",
		missingSet,
	)
	v.validateField(
		options.Contracts.AppChain.MaxChainDisconnectTime,
		"contracts.app-chain.max-chain-disconnect-time",
		customSet,
	)
}

func (v *OptionsValidator) validateSettlementChainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	v.validateField(
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.chain-id",
		customSet,
	)
	v.validateRPCURL(
		options.Contracts.SettlementChain.RPCURL,
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.rpc-url",
		missingSet,
	)
	v.validateWebsocketURL(
		options.Contracts.SettlementChain.WssURL,
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.wss-url",
		missingSet,
	)
	v.validateHexAddress(
		options.Contracts.SettlementChain.NodeRegistryAddress,
		"contracts.settlement-chain.node-registry-address",
		missingSet,
	)
	v.validateField(
		options.Contracts.SettlementChain.NodeRegistryRefreshInterval,
		"contracts.settlement-chain.node-registry-refresh-interval",
		customSet,
	)
	v.validateHexAddress(
		options.Contracts.SettlementChain.PayerRegistryAddress,
		"contracts.settlement-chain.payer-registry-address",
		missingSet,
	)
	v.validateHexAddress(
		options.Contracts.SettlementChain.PayerReportManagerAddress,
		"contracts.settlement-chain.payer-report-manager-address",
		missingSet,
	)
	v.validateField(
		options.Contracts.SettlementChain.MaxChainDisconnectTime,
		"contracts.settlement-chain.max-chain-disconnect-time",
		customSet,
	)
}

// validateField checks if a field meets the validation requirements and adds appropriate errors.
func (v *OptionsValidator) validateField(value any, fieldName string, set map[string]struct{}) {
	switch v := value.(type) {
	case string:
		if v == "" {
			set[fmt.Sprintf("--%s", fieldName)] = struct{}{}
		}
	case int:
		if v <= 0 {
			set[fmt.Sprintf("--%s must be greater than 0", fieldName)] = struct{}{}
		}
	case time.Duration:
		if v <= 0 {
			set[fmt.Sprintf("--%s must be greater than 0", fieldName)] = struct{}{}
		}
	}
}

func (v *OptionsValidator) validateHexAddress(
	address string,
	fieldName string,
	set map[string]struct{},
) {
	if address == "" {
		set[fmt.Sprintf("--%s is required", fieldName)] = struct{}{}
	}
	if !common.IsHexAddress(address) || common.HexToAddress(address) == (common.Address{}) {
		set[fmt.Sprintf("--%s is invalid", fieldName)] = struct{}{}
	}
}

func (v *OptionsValidator) validateRPCURL(
	rpcURL string,
	chainID int64,
	fieldName string,
	set map[string]struct{},
) {
	u, err := url.Parse(rpcURL)
	if err != nil {
		set[fmt.Sprintf("--%s is an invalid URL, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		set[fmt.Sprintf("--%s is invalid, expected http or https, got %s", fieldName, u.Scheme)] = struct{}{}
		return
	}

	v.validateChainID(rpcURL, chainID, fieldName, set)
}

func (v *OptionsValidator) validateWebsocketURL(
	wsURL string,
	chainID int64,
	fieldName string,
	set map[string]struct{},
) {
	u, err := url.Parse(wsURL)
	if err != nil {
		set[fmt.Sprintf("--%s is an invalid URL, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		set[fmt.Sprintf("--%s is invalid, expected ws or wss, got %s", fieldName, u.Scheme)] = struct{}{}
		return
	}

	v.validateChainID(wsURL, chainID, fieldName, set)
}

func (v *OptionsValidator) validateChainID(
	url string,
	expectedChainID int64,
	fieldName string,
	set map[string]struct{},
) {
	ctx := context.Background()

	client, err := ethclient.DialContext(ctx, url)
	if err != nil {
		set[fmt.Sprintf("--%s error dialing, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	defer func() {
		if client != nil {
			client.Close()
		}
	}()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		set[fmt.Sprintf("--%s error getting chain ID, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	if chainID.Int64() != int64(expectedChainID) {
		set[fmt.Sprintf("--%s is invalid, expected chain ID %d, got %d", fieldName, expectedChainID, chainID.Int64())] = struct{}{}
	}
}

func (v *OptionsValidator) validatePayerOptions(
	options *PayerOptions,
	customSet map[string]struct{},
) error {
	validStrategies := map[string]bool{
		"stable":  true,
		"manual":  true,
		"ordered": true,
		"random":  true,
		"closest": true,
	}

	if options.NodeSelectorStrategy != "" && !validStrategies[options.NodeSelectorStrategy] {
		return fmt.Errorf(
			"invalid node-selector-strategy: %s (must be one of: stable, manual, ordered, random, closest)",
			options.NodeSelectorStrategy,
		)
	}

	if (options.NodeSelectorStrategy == "manual" || options.NodeSelectorStrategy == "ordered") &&
		len(options.NodeSelectorPreferredNodes) == 0 {
		return fmt.Errorf(
			"strategy %s requires at least one node in node-selector-preferred-nodes",
			options.NodeSelectorStrategy,
		)
	}

	if options.NodeSelectorCacheExpiry <= 0 {
		customSet["--payer.node-selector-cache-expiry must be greater than 0"] = struct{}{}
	}

	if options.NodeSelectorTimeout <= 0 {
		customSet["--payer.node-selector-connect-timeout must be greater than 0"] = struct{}{}
	}

	return nil
}
