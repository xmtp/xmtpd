package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ValidateServerOptions(options *ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	err := ParseJSONConfig(&options.Contracts)
	if err != nil {
		return err
	}

	validateBlockchainConfig(options, missingSet, customSet)

	validateMigrationOptions(options, missingSet, customSet)

	validateField(
		options.DB.WriterConnectionString,
		"db.writer-connection-string",
		missingSet,
	)

	if options.Replication.Enable {
		validateField(options.Signer.PrivateKey, "signer.private-key", missingSet)
		validateField(options.MlsValidation.GrpcAddress, "mls-validation.grpc-address", missingSet)
	}

	if options.Indexer.Enable {
		validateField(options.MlsValidation.GrpcAddress, "mls-validation.grpc-address", missingSet)
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

func ValidatePruneOptions(options PruneOptions) error {
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

func validateMigrationOptions(
	opts *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	if opts.MigrationServer.Enable && opts.MigrationClient.Enable {
		missingSet["--migration-server.enable and --migration-client.enable cannot be used together"] = struct{}{}
	}

	if opts.MigrationServer.Enable {
		validateField(
			opts.MigrationServer.PayerPrivateKey,
			"migration-server.payer-private-key",
			missingSet,
		)
		validateField(
			opts.MigrationServer.NodeSigningKey,
			"migration-server.node-signing-key",
			missingSet,
		)
		validateField(
			opts.Signer.PrivateKey,
			"signer.private-key",
			missingSet,
		)
		validateField(
			opts.MigrationServer.ReaderConnectionString,
			"migration-server.reader-connection-string",
			missingSet,
		)
		validateField(
			opts.MigrationServer.ReaderTimeout,
			"migration-server.reader-timeout",
			missingSet,
		)
		validateAppChainConfig(opts, missingSet, customSet)
	}

	if opts.MigrationClient.Enable {
		validateField(
			opts.MigrationClient.FromNodeID,
			"migration-client.from-node-id",
			missingSet,
		)
	}
}

func ContractOptionsFromEnv(filePath string) (ContractsOptions, error) {
	if filePath == "" {
		return ContractsOptions{}, errors.New("config file path is not set")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ContractsOptions{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return ContractsOptions{}, err
	}

	var config ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ContractsOptions{}, err
	}

	// Set default values for missing options in environment json.
	return ContractsOptions{
		SettlementChain: SettlementChainOptions{
			NodeRegistryAddress:         config.NodeRegistry,
			RateRegistryAddress:         config.RateRegistry,
			ParameterRegistryAddress:    config.SettlementChainParameterRegistry,
			PayerRegistryAddress:        config.PayerRegistry,
			PayerReportManagerAddress:   config.PayerReportManager,
			ChainID:                     config.SettlementChainID,
			NodeRegistryRefreshInterval: 60 * time.Second,
			RateRegistryRefreshInterval: 300 * time.Second,
			MaxChainDisconnectTime:      300 * time.Second,
			BackfillBlockPageSize:       500,
		},
		AppChain: AppChainOptions{
			GroupMessageBroadcasterAddress:   config.GroupMessageBroadcaster,
			IdentityUpdateBroadcasterAddress: config.IdentityUpdateBroadcaster,
			ChainID:                          config.AppChainID,
			MaxChainDisconnectTime:           300 * time.Second,
			BackfillBlockPageSize:            500,
			ParameterRegistryAddress:         config.AppChainParameterRegistry,
		},
	}, nil
}

func ParseJSONConfig(options *ContractsOptions) error {
	if options.ConfigFilePath != "" && options.ConfigJSON != "" {
		return errors.New("--config-file and --config-json cannot be used together")
	}

	if options.ConfigFilePath != "" {
		file, err := os.Open(options.ConfigFilePath)
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		// Unmarshal JSON into the Config struct
		var config ChainConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		fillConfigFromJSON(options, &config)
	}

	if options.ConfigJSON != "" {
		// Unmarshal JSON into the Config struct
		var config ChainConfig
		if err := json.Unmarshal([]byte(options.ConfigJSON), &config); err != nil {
			return err
		}

		fillConfigFromJSON(options, &config)
	}

	return nil
}

func fillConfigFromJSON(options *ContractsOptions, config *ChainConfig) {
	// Explicitly specified ENV variables in options take precedence!
	// Only fill in values from the JSON if the relevant fields in options are empty or zero.

	// AppChainOptions
	if options.AppChain.GroupMessageBroadcasterAddress == "" {
		options.AppChain.GroupMessageBroadcasterAddress = config.GroupMessageBroadcaster
	}
	if options.AppChain.IdentityUpdateBroadcasterAddress == "" {
		options.AppChain.IdentityUpdateBroadcasterAddress = config.IdentityUpdateBroadcaster
	}
	if options.AppChain.ChainID == 0 || options.AppChain.ChainID == 31337 {
		options.AppChain.ChainID = config.AppChainID
	}
	if options.AppChain.GatewayAddress == "" {
		options.AppChain.GatewayAddress = config.AppChainGateway
	}
	if options.AppChain.DeploymentBlock == 0 {
		options.AppChain.DeploymentBlock = uint64(config.AppChainDeploymentBlock)
	}
	if options.AppChain.ParameterRegistryAddress == "" {
		options.AppChain.ParameterRegistryAddress = config.AppChainParameterRegistry
	}

	// SettlementChainOptions
	if options.SettlementChain.NodeRegistryAddress == "" {
		options.SettlementChain.NodeRegistryAddress = config.NodeRegistry
	}
	if options.SettlementChain.RateRegistryAddress == "" {
		options.SettlementChain.RateRegistryAddress = config.RateRegistry
	}
	if options.SettlementChain.ParameterRegistryAddress == "" {
		options.SettlementChain.ParameterRegistryAddress = config.SettlementChainParameterRegistry
	}
	if options.SettlementChain.PayerRegistryAddress == "" {
		options.SettlementChain.PayerRegistryAddress = config.PayerRegistry
	}
	if options.SettlementChain.PayerReportManagerAddress == "" {
		options.SettlementChain.PayerReportManagerAddress = config.PayerReportManager
	}
	if options.SettlementChain.ChainID == 0 || options.SettlementChain.ChainID == 31337 {
		options.SettlementChain.ChainID = config.SettlementChainID
	}
	if options.SettlementChain.DeploymentBlock == 0 {
		options.SettlementChain.DeploymentBlock = uint64(config.SettlementChainDeploymentBlock)
	}
	if options.SettlementChain.GatewayAddress == "" {
		options.SettlementChain.GatewayAddress = config.SettlementChainGateway
	}
	if options.SettlementChain.DistributionManagerAddress == "" {
		options.SettlementChain.DistributionManagerAddress = config.DistributionManager
	}
}

func validateBlockchainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	validateAppChainConfig(options, missingSet, customSet)
	validateSettlementChainConfig(options, missingSet, customSet)

	if options.Replication.Enable || options.Sync.Enable {
		validateHexAddress(
			options.Contracts.SettlementChain.RateRegistryAddress,
			"contracts.settlement-chain.rate-registry-address",
			missingSet,
		)
		validateField(
			options.Contracts.SettlementChain.RateRegistryRefreshInterval,
			"contracts.settlement-chain.rate-registry-refresh-interval",
			customSet,
		)
	}
}

func validateAppChainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	validateField(
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.chain-id",
		customSet,
	)
	validateRPCURL(
		options.Contracts.AppChain.RPCURL,
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.rpc-url",
		missingSet,
	)
	validateWebsocketURL(
		options.Contracts.AppChain.WssURL,
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.wss-url",
		missingSet,
	)
	validateHexAddress(
		options.Contracts.AppChain.GroupMessageBroadcasterAddress,
		"contracts.app-chain.group-message-broadcaster-address",
		missingSet,
	)
	validateHexAddress(
		options.Contracts.AppChain.IdentityUpdateBroadcasterAddress,
		"contracts.app-chain.identity-update-broadcaster-address",
		missingSet,
	)
	validateField(
		options.Contracts.AppChain.MaxChainDisconnectTime,
		"contracts.app-chain.max-chain-disconnect-time",
		customSet,
	)
}

func validateSettlementChainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	validateField(
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.chain-id",
		customSet,
	)
	validateRPCURL(
		options.Contracts.SettlementChain.RPCURL,
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.rpc-url",
		missingSet,
	)
	validateWebsocketURL(
		options.Contracts.SettlementChain.WssURL,
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.wss-url",
		missingSet,
	)
	validateHexAddress(
		options.Contracts.SettlementChain.NodeRegistryAddress,
		"contracts.settlement-chain.node-registry-address",
		missingSet,
	)
	validateField(
		options.Contracts.SettlementChain.NodeRegistryRefreshInterval,
		"contracts.settlement-chain.node-registry-refresh-interval",
		customSet,
	)
	validateHexAddress(
		options.Contracts.SettlementChain.PayerRegistryAddress,
		"contracts.settlement-chain.payer-registry-address",
		missingSet,
	)
	validateHexAddress(
		options.Contracts.SettlementChain.PayerReportManagerAddress,
		"contracts.settlement-chain.payer-report-manager-address",
		missingSet,
	)
	validateField(
		options.Contracts.SettlementChain.MaxChainDisconnectTime,
		"contracts.settlement-chain.max-chain-disconnect-time",
		customSet,
	)
}

// validateField checks if a field meets the validation requirements and adds appropriate errors.
func validateField(value interface{}, fieldName string, set map[string]struct{}) {
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

func validateHexAddress(address string, fieldName string, set map[string]struct{}) {
	if address == "" {
		set[fmt.Sprintf("--%s is required", fieldName)] = struct{}{}
	}
	if !common.IsHexAddress(address) || common.HexToAddress(address) == (common.Address{}) {
		set[fmt.Sprintf("--%s is invalid", fieldName)] = struct{}{}
	}
}

func validateRPCURL(rpcURL string, chainID int, fieldName string, set map[string]struct{}) {
	u, err := url.Parse(rpcURL)
	if err != nil {
		set[fmt.Sprintf("--%s is an invalid URL, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		set[fmt.Sprintf("--%s is invalid, expected http or https, got %s", fieldName, u.Scheme)] = struct{}{}
		return
	}

	validateChainID(rpcURL, chainID, fieldName, set)
}

func validateWebsocketURL(wsURL string, chainID int, fieldName string, set map[string]struct{}) {
	u, err := url.Parse(wsURL)
	if err != nil {
		set[fmt.Sprintf("--%s is an invalid URL, %s", fieldName, err.Error())] = struct{}{}
		return
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		set[fmt.Sprintf("--%s is invalid, expected ws or wss, got %s", fieldName, u.Scheme)] = struct{}{}
		return
	}

	validateChainID(wsURL, chainID, fieldName, set)
}

func validateChainID(url string, expectedChainID int, fieldName string, set map[string]struct{}) {
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
