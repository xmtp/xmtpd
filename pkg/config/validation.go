package config

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateServerOptions(options ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	if options.AppChain.RpcUrl == "" {
		missingSet["--app-chain.rpc-url"] = struct{}{}
	}

	if options.AppChain.MessagesContractAddress == "" {
		missingSet["--app-chain.messages-address"] = struct{}{}
	}

	if options.AppChain.IdentityUpdatesContractAddress == "" {
		missingSet["--app-chain.identity-updates-address"] = struct{}{}
	}

	if options.AppChain.ChainID == 0 {
		customSet["--app-chain.chain-id must be greater than 0"] = struct{}{}
	}

	if options.AppChain.MaxDisconnectTime <= 0 {
		customSet["--app-chain.max-chain-disconnect-time must be greater than 0"] = struct{}{}
	}

	if options.BaseChain.RpcUrl == "" {
		missingSet["--base-chain.rpc-url"] = struct{}{}
	}

	if options.BaseChain.NodesContractAddress == "" {
		missingSet["--base-chain.nodes-address"] = struct{}{}
	}

	if options.BaseChain.ChainID == 0 {
		customSet["--base-chain.chain-id must be greater than 0"] = struct{}{}
	}

	if options.BaseChain.RefreshInterval <= 0 {
		customSet["--base-chain.refresh-interval must be greater than 0"] = struct{}{}
	}

	if options.Payer.Enable {
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
