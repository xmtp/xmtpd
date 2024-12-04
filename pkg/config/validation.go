package config

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateServerOptions(options ServerOptions) error {
	missingSet := make(map[string]struct{})
	customSet := make(map[string]struct{})

	if options.Contracts.RpcUrl == "" {
		missingSet["--contracts.rpc-url"] = struct{}{}
	}

	if options.Contracts.NodesContractAddress == "" {
		missingSet["--contracts.nodes-address"] = struct{}{}
	}

	if options.Contracts.MessagesContractAddress == "" {
		missingSet["--contracts.messages-address"] = struct{}{}
	}

	if options.Contracts.IdentityUpdatesContractAddress == "" {
		missingSet["--contracts.identity-updates-address"] = struct{}{}
	}

	if options.Contracts.ChainID == 0 {
		customSet["--contracts.chain-id must be greater than 0"] = struct{}{}
	}

	if options.Contracts.RefreshInterval <= 0 {
		customSet["--contracts.refresh-interval must be greater than 0"] = struct{}{}
	}

	if options.Contracts.MaxChainDisconnectTime <= 0 {
		customSet["--contracts.max-chain-disconnect-time must be greater than 0"] = struct{}{}
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
