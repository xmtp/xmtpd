package config

import (
	"errors"
)

func ValidateServerOptions(options ServerOptions) error {
	if options.Payer.Enable {
		if options.Payer.PrivateKey == "" {
			return errors.New("payer.PrivateKey is required")
		}
	}

	if options.Replication.Enable {
		if options.DB.WriterConnectionString == "" {
			return errors.New("DB.WriterConnectionString is required")
		}
		if options.Signer.PrivateKey == "" {
			return errors.New("Signer.PrivateKey is required")
		}
	}

	if options.Sync.Enable {
		if options.DB.WriterConnectionString == "" {
			return errors.New("DB.WriterConnectionString is required")
		}
	}

	if options.Indexer.Enable {
		if options.DB.WriterConnectionString == "" {
			return errors.New("DB.WriterConnectionString is required")
		}
	}

	return nil
}
