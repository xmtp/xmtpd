package commands

import (
	"crypto/ecdsa"
	"log"

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
		Run:   generateKeyHandler,
		Example: `
Usage: xmtpd-cli keys generate

Example:
xmtpd-cli keys generate
`,
	}

	return &cmd
}

func generateKeyHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	privKey, err := utils.GenerateEcdsaPrivateKey()
	if err != nil {
		logger.Fatal("could not generate private key", zap.Error(err))
	}

	logger.Info(
		"generated private key",
		zap.String("private-key", utils.EcdsaPrivateKeyToString(privKey)),
		zap.String("public-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)
}

func getPubKeyCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get-public-key",
		Short: "Get the public key for a private key",
		Run:   getPubKeyHandler,
		Example: `
Usage: xmtpd-cli keys get-public-key --private-key <private-key>

Example:
xmtpd-cli keys get-public-key --private-key <private-key>
`,
	}

	return &cmd
}

func getPubKeyHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	privateKey := viper.GetString("private-key")
	if privateKey == "" {
		logger.Fatal("private key is not set")
	}

	privKey, err := utils.ParseEcdsaPrivateKey(privateKey)
	if err != nil {
		logger.Fatal("could not parse private key", zap.Error(err))
	}

	logger.Info(
		"parsed private key",
		zap.String("pub-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)

	privKey.Public()
}
