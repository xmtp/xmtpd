package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ws "github.com/gorilla/websocket"
)

func ValidateServerOptions(options *ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	if !isMultiChainDeployment(*options) {
		normalizeSingleChainConfig(options)
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

	if len(missingSet) > 0 {
		var errorMessages []string
		for err := range missingSet {
			errorMessages = append(errorMessages, err)
		}

		return fmt.Errorf("missing required arguments: %s", strings.Join(errorMessages, ", "))
	}

	return nil
}

func isMultiChainDeployment(options ServerOptions) bool {
	return options.Contracts.AppChain.RpcURL != "" ||
		options.Contracts.SettlementChain.RpcURL != ""
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
}

// normalizeSingleChainConfig copies values from deprecated fields to new fields for single-chain deployments.
func normalizeSingleChainConfig(options *ServerOptions) {
	if options.Contracts.RpcUrl != "" {
		options.Contracts.AppChain.RpcURL = options.Contracts.RpcUrl
		options.Contracts.SettlementChain.RpcURL = options.Contracts.RpcUrl
	}
	if options.Contracts.NodesContract.NodesContractAddress != "" {
		options.Contracts.SettlementChain.NodeRegistryAddress = options.Contracts.NodesContract.NodesContractAddress
	}
	if options.Contracts.MessagesContractAddress != "" {
		options.Contracts.AppChain.GroupMessageBroadcasterAddress = options.Contracts.MessagesContractAddress
	}
	if options.Contracts.IdentityUpdatesContractAddress != "" {
		options.Contracts.AppChain.IdentityUpdateBroadcasterAddress = options.Contracts.IdentityUpdatesContractAddress
	}
	if options.Contracts.RateRegistryContractAddress != "" {
		options.Contracts.SettlementChain.RateRegistryAddress = options.Contracts.RateRegistryContractAddress
	}
	if options.Contracts.RatesRefreshInterval > 0 {
		options.Contracts.SettlementChain.RateRegistryRefreshInterval = options.Contracts.RatesRefreshInterval
	}
	if options.Contracts.ChainID > 0 {
		options.Contracts.AppChain.ChainID = options.Contracts.ChainID
		options.Contracts.SettlementChain.ChainID = options.Contracts.ChainID
	}
	if options.Contracts.RegistryRefreshInterval > 0 {
		options.Contracts.SettlementChain.NodeRegistryRefreshInterval = options.Contracts.RegistryRefreshInterval
	}
	if options.Contracts.MaxChainDisconnectTime > 0 {
		options.Contracts.AppChain.MaxChainDisconnectTime = options.Contracts.MaxChainDisconnectTime
	}
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

	conn.Close()
}
