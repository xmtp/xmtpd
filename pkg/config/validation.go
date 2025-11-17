package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config/environments"
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

	if options.API.Enable {
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

	var data []byte
	// Try to parse as URL. If it fails, treat as local path.
	if u, err := url.Parse(filePath); err == nil &&
		(u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "file" || u.Scheme == "config") {
		switch u.Scheme {
		case "config":
			data, err = environments.GetEnvironmentConfig(
				environments.SmartContractEnvironment(u.Host),
			)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf("unknown config environment %s", u.Host)
			}
		case "http", "https":
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Get(filePath)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf("fetching %s: %w", filePath, err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return ContractsOptions{}, fmt.Errorf(
					"fetching %s: http %d",
					filePath,
					resp.StatusCode,
				)
			}
			// Guard against large config files (10KB cap).
			limited := io.LimitReader(resp.Body, 10<<10)
			data, err = io.ReadAll(limited)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf("reading %s: %w", filePath, err)
			}
		case "file":
			// file:// URLs may have URL-encoded paths
			localPath := u.Path
			if u.Host != "" {
				localPath = "//" + u.Host + u.Path
			}
			if localPath == "" {
				// Handle cases like file:///absolute/path
				localPath = strings.TrimPrefix(filePath, "file://")
			}
			localPath, err = url.PathUnescape(localPath)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf(
					"invalid file URL path %q: %w",
					localPath,
					err,
				)
			}

			f, err := os.Open(localPath)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf("open %s: %w", localPath, err)
			}
			defer func() {
				_ = f.Close()
			}()
			limited := io.LimitReader(f, 10<<10)
			data, err = io.ReadAll(limited)
			if err != nil {
				return ContractsOptions{}, fmt.Errorf("read %s: %w", localPath, err)
			}
		default:
			return ContractsOptions{}, fmt.Errorf("unsupported URL scheme %q", u.Scheme)
		}
	} else {
		// Local filesystem path
		f, err := os.Open(filePath)
		if err != nil {
			return ContractsOptions{}, fmt.Errorf("open %s: %w", filePath, err)
		}
		defer func() {
			_ = f.Close()
		}()
		r := io.LimitReader(f, 10<<10)
		data, err = io.ReadAll(r)
		if err != nil {
			return ContractsOptions{}, fmt.Errorf("read %s: %w", filePath, err)
		}
	}

	var config ChainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ContractsOptions{}, fmt.Errorf("unmarshal config: %w", err)
	}

	// Set default values for missing options in environment json.
	return ContractsOptions{
		SettlementChain: SettlementChainOptions{
			NodeRegistryAddress:         config.NodeRegistry,
			RateRegistryAddress:         config.RateRegistry,
			ParameterRegistryAddress:    config.SettlementChainParameterRegistry,
			PayerRegistryAddress:        config.PayerRegistry,
			PayerReportManagerAddress:   config.PayerReportManager,
			ChainID:                     int64(config.SettlementChainID),
			DeploymentBlock:             uint64(config.SettlementChainDeploymentBlock),
			UnderlyingFeeToken:          config.UnderlyingFeeToken,
			FeeToken:                    config.FeeToken,
			NodeRegistryRefreshInterval: 60 * time.Second,
			RateRegistryRefreshInterval: 300 * time.Second,
			MaxChainDisconnectTime:      300 * time.Second,
			BackfillBlockPageSize:       500,
			GatewayAddress:              config.SettlementChainGateway,
			DistributionManagerAddress:  config.DistributionManager,
		},
		AppChain: AppChainOptions{
			GroupMessageBroadcasterAddress:   config.GroupMessageBroadcaster,
			IdentityUpdateBroadcasterAddress: config.IdentityUpdateBroadcaster,
			ChainID:                          int64(config.AppChainID),
			MaxChainDisconnectTime:           300 * time.Second,
			BackfillBlockPageSize:            500,
			GatewayAddress:                   config.AppChainGateway,
			DeploymentBlock:                  uint64(config.AppChainDeploymentBlock),
			ParameterRegistryAddress:         config.AppChainParameterRegistry,
		},
	}, nil
}

func ParseJSONConfig(options *ContractsOptions) error {
	if options.Environment != "" && (options.ConfigFilePath != "" || options.ConfigJSON != "") {
		return errors.New(
			"--contracts.environment cannot be used with --contracts.config-file or --contracts.config-json",
		)
	}

	if options.ConfigFilePath != "" && options.ConfigJSON != "" {
		return errors.New("--config-file and --config-json cannot be used together")
	}

	if options.Environment != "" {
		fmt.Printf("Environment is %s\n", options.Environment)
		data, err := environments.GetEnvironmentConfig(
			options.Environment,
		)
		if err != nil {
			return err
		}

		var config ChainConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}
		fmt.Printf("Chain config %v\n", config)

		fillConfigFromJSON(options, &config)
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
		options.AppChain.ChainID = int64(config.AppChainID)
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
		options.SettlementChain.ChainID = int64(config.SettlementChainID)
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

	if options.SettlementChain.UnderlyingFeeToken == "" {
		options.SettlementChain.UnderlyingFeeToken = config.UnderlyingFeeToken
	}

	if options.SettlementChain.FeeToken == "" {
		options.SettlementChain.FeeToken = config.FeeToken
	}
}

func validateBlockchainConfig(
	options *ServerOptions,
	missingSet map[string]struct{},
	customSet map[string]struct{},
) {
	validateAppChainConfig(options, missingSet, customSet)
	validateSettlementChainConfig(options, missingSet, customSet)

	if options.API.Enable || options.Sync.Enable {
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
func validateField(value any, fieldName string, set map[string]struct{}) {
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

func validateRPCURL(rpcURL string, chainID int64, fieldName string, set map[string]struct{}) {
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

func validateWebsocketURL(wsURL string, chainID int64, fieldName string, set map[string]struct{}) {
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

func validateChainID(url string, expectedChainID int64, fieldName string, set map[string]struct{}) {
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
