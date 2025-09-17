package commands

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

func keyManagementCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "keys",
		Short: "Manage keys",
	}

	cmd.AddCommand(
		generateKeyCommand(),
		getPubKeyCommand(),
	)

	return &cmd
}

func generateKeyCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate",
		Short: "Generate a new key pair",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return generateKeyHandler()
		},
		Example: `
Usage: xmtpd-cli keys generate

Example:
xmtpd-cli keys generate
`,
	}

	return &cmd
}

func generateKeyHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	privKey, err := utils.GenerateEcdsaPrivateKey()
	if err != nil {
		return fmt.Errorf("could not generate private key: %w", err)
	}

	logger.Info(
		"generated private key",
		zap.String("private-key", utils.EcdsaPrivateKeyToString(privKey)),
		zap.String("public-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)

	return nil
}

func getPubKeyCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get-public-key",
		Short: "Get the public key for a private key",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return getPubKeyHandler()
		},
		Example: `
Usage: xmtpd-cli keys get-public-key --private-key <private-key>

Example:
xmtpd-cli keys get-public-key --private-key <private-key>
`,
	}

	return &cmd
}

func getPubKeyHandler() error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	privateKey := viper.GetString("private-key")
	if privateKey == "" {
		return fmt.Errorf("private key is not set")
	}

	privKey, err := utils.ParseEcdsaPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("could not parse private key: %w", err)
	}

	logger.Info(
		"parsed private key",
		zap.String("pub-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)

	return nil
}
