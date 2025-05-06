package config

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateServerOptions(options ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	// App Chain validation.
	if options.Contracts.AppChain.RpcURL == "" {
		missingSet["--contracts.app-chain.rpc-url"] = struct{}{}
	}

	if options.Contracts.AppChain.ChainID == 0 {
		customSet["--contracts.app-chain.chain-id must be greater than 0"] = struct{}{}
	}

	if options.Contracts.AppChain.GroupMessageBroadcasterAddress == "" {
		missingSet["--contracts.app-chain.group-message-broadcaster-address"] = struct{}{}
	}

	if options.Contracts.AppChain.IdentityUpdateBroadcasterAddress == "" {
		missingSet["--contracts.app-chain.identity-update-broadcaster-address"] = struct{}{}
	}

	if options.Contracts.AppChain.MaxChainDisconnectTime <= 0 {
		customSet["--contracts.app-chain.max-chain-disconnect-time must be greater than 0"] = struct{}{}
	}

	// Settlement Chain validation.
	if options.Contracts.SettlementChain.RpcURL == "" {
		missingSet["--contracts.settlement-chain.rpc-url"] = struct{}{}
	}

	if options.Contracts.SettlementChain.ChainID == 0 {
		customSet["--contracts.settlement-chain.chain-id must be greater than 0"] = struct{}{}
	}

	if options.Contracts.SettlementChain.NodeRegistryAddress == "" {
		missingSet["--contracts.settlement-chain.node-registry-address"] = struct{}{}
	}

	if options.Contracts.SettlementChain.NodeRegistryRefreshInterval <= 0 {
		customSet["--contracts.settlement-chain.node-registry-refresh-interval must be greater than 0"] = struct{}{}
	}

	if options.Contracts.SettlementChain.RateRegistryAddress == "" {
		missingSet["--contracts.settlement-chain.rate-registry-address"] = struct{}{}
	}

	if options.Contracts.SettlementChain.RateRegistryRefreshInterval <= 0 {
		customSet["--contracts.settlement-chain.rate-registry-refresh-interval must be greater than 0"] = struct{}{}
	}

	if options.Payer.Enable {
		if options.DB.WriterConnectionString == "" {
			missingSet["--DB.WriterConnectionString"] = struct{}{}
		}
		if options.Payer.PrivateKey == "" {
			missingSet["--payer.PrivateKey"] = struct{}{}
		}
	}

	if options.Replication.Enable {
		if options.DB.WriterConnectionString == "" {
			missingSet["--DB.WriterConnectionString"] = struct{}{}
		}
		if options.Signer.PrivateKey == "" {
			missingSet["--Signer.PrivateKey"] = struct{}{}
		}
	}

	if options.Sync.Enable {
		if options.DB.WriterConnectionString == "" {
			missingSet["--DB.WriterConnectionString"] = struct{}{}
		}
	}

	if options.Indexer.Enable {
		if options.DB.WriterConnectionString == "" {
			missingSet["--DB.WriterConnectionString"] = struct{}{}
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
