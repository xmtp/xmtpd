package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ws "github.com/gorilla/websocket"
)

func ValidateServerOptions(options *ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	err := ParseJSONConfig(&options.Contracts)
	if err != nil {
		return err
	}

	validateBlockchainConfig(options, missingSet, customSet)

	validateField(
		options.DB.WriterConnectionString,
		"db.writer-connection-string",
		missingSet,
	)

	if options.Payer.Enable {
		validateField(options.Payer.PrivateKey, "payer.private-key", missingSet)
	}

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

func ParseJSONConfig(options *ContractsOptions) error {
	if options.ConfigFilePath != "" && options.ConfigJson != "" {
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

		fmt.Println("data", string(data))

		// Unmarshal JSON into the Config struct
		var config ChainConfig
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}

		fillConfigFromJson(options, &config)
	}

	if options.ConfigJson != "" {
		// Unmarshal JSON into the Config struct
		var config ChainConfig
		if err := json.Unmarshal([]byte(options.ConfigJson), &config); err != nil {
			return err
		}

		fillConfigFromJson(options, &config)
	}

	return nil
}

func fillConfigFromJson(options *ContractsOptions, config *ChainConfig) {
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
	// TODO: For now, we only validate RpcURL, until deployments are migrated to WssURL.
	validateWebsocketURL(
		options.Contracts.AppChain.RpcURL,
		"contracts.app-chain.rpc-url",
		missingSet,
	)
	validateField(
		options.Contracts.AppChain.ChainID,
		"contracts.app-chain.chain-id",
		customSet,
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
	// TODO: For now, we only validate RpcURL, until deployments are migrated to WssURL.
	validateWebsocketURL(
		options.Contracts.SettlementChain.RpcURL,
		"contracts.settlement-chain.rpc-url",
		missingSet,
	)
	validateField(
		options.Contracts.SettlementChain.ChainID,
		"contracts.settlement-chain.chain-id",
		customSet,
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

func validateWebsocketURL(url string, fieldName string, set map[string]struct{}) {
	dialer := &ws.Dialer{
		HandshakeTimeout: 15 * time.Second,
	}

	// Dial returns an error if the URL is invalid, or if the connection fails.
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		set[fmt.Sprintf("--%s is invalid", fieldName)] = struct{}{}
	}

	_ = conn.Close()
}
