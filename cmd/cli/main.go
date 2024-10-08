package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type globalOptions struct {
	Contracts config.ContractsOptions `group:"Contracts Options" namespace:"contracts"`
	Log       config.LogOptions       `group:"Log Options"       namespace:"log"`
}

type CLI struct {
	globalOptions
	Command      string
	GetPubKey    config.GetPubKeyOptions
	GenerateKey  config.GenerateKeyOptions
	RegisterNode config.RegisterNodeOptions
	GetAllNodes  config.GetAllNodesOptions
}

/*
*
Parse the command line options and return the CLI struct.

Some special care has to be made here to ensure that the required options are only evaluated for the correct command.
We use a wrapper type to scope the parser to only the universal options, allowing us to have required fields on
the options for each subcommand.
*
*/
func parseOptions(args []string) (*CLI, error) {
	var options globalOptions
	var generateKeyOptions config.GenerateKeyOptions
	var registerNodeOptions config.RegisterNodeOptions
	var getPubKeyOptions config.GetPubKeyOptions
	var getAllNodesOptions config.GetAllNodesOptions

	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.AddCommand("generate-key", "Generate a public/private keypair", "", &generateKeyOptions); err != nil {
		return nil, fmt.Errorf("Could not add generate-key command: %s", err)
	}
	if _, err := parser.AddCommand("register-node", "Register a node", "", &registerNodeOptions); err != nil {
		return nil, fmt.Errorf("Could not add register-node command: %s", err)
	}
	if _, err := parser.AddCommand("get-pub-key", "Get the public key for a private key", "", &getPubKeyOptions); err != nil {
		return nil, fmt.Errorf("Could not add get-pub-key command: %s", err)
	}
	if _, err := parser.AddCommand("get-all-nodes", "Get all nodes from the registry", "", &getAllNodesOptions); err != nil {
		return nil, fmt.Errorf("Could not add get-all-nodes command: %s", err)
	}
	if _, err := parser.ParseArgs(args); err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			return nil, fmt.Errorf("Could not parse options: %s", err)
		}
		return nil, nil
	}

	if parser.Active == nil {
		return nil, errors.New("No command provided")
	}

	return &CLI{
		options,
		parser.Active.Name,
		getPubKeyOptions,
		generateKeyOptions,
		registerNodeOptions,
		getAllNodesOptions,
	}, nil
}

func getPubKey(logger *zap.Logger, options *CLI) {
	privKey, err := utils.ParseEcdsaPrivateKey(options.GetPubKey.PrivateKey)
	if err != nil {
		logger.Fatal("could not parse private key", zap.Error(err))
	}
	logger.Info(
		"parsed private key",
		zap.String("pub-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
	)
}
func registerNode(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.RegisterNode.AdminPrivateKey,
		options.Contracts.ChainID,
	)

	if err != nil {
		logger.Fatal("could not create signer", zap.Error(err))
	}

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	signingKeyPub, err := utils.ParseEcdsaPublicKey(options.RegisterNode.SigningKeyPub)
	if err != nil {
		logger.Fatal("could not decompress public key", zap.Error(err))
	}

	err = registryAdmin.AddNode(
		ctx,
		options.RegisterNode.OwnerAddress,
		signingKeyPub,
		options.RegisterNode.HttpAddress,
	)
	if err != nil {
		logger.Fatal("could not add node", zap.Error(err))
	}
	logger.Info(
		"successfully added node",
		zap.String("owner-address", options.RegisterNode.OwnerAddress),
		zap.String("node-http-address", options.RegisterNode.HttpAddress),
		zap.String("node-signing-key-pub", utils.EcdsaPublicKeyToString(signingKeyPub)),
	)
}

func generateKey(logger *zap.Logger) {
	privKey, err := utils.GenerateEcdsaPrivateKey()
	if err != nil {
		logger.Fatal("could not generate private key", zap.Error(err))
	}
	logger.Info(
		"generated private key",
		zap.String("private-key", utils.EcdsaPrivateKeyToString(privKey)),
		zap.String("public-key", utils.EcdsaPublicKeyToString(privKey.Public().(*ecdsa.PublicKey))),
	)
}

func getAllNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.GetAllNodes.AdminPrivateKey,
		options.Contracts.ChainID,
	)

	if err != nil {
		logger.Fatal("could not create signer", zap.Error(err))
	}

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	nodes, err := registryAdmin.GetAllNodes(ctx)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}
	logger.Info(
		"got nodes",
		zap.Int("size", len(nodes)),
		zap.Any("nodes", nodes),
	)
}

func main() {
	options, err := parseOptions(os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse options: %s", err)
	}
	if options == nil {
		return
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		log.Fatalf("Could not build logger: %s", err)
	}
	switch options.Command {
	case "generate-key":
		generateKey(logger)
		return
	case "get-pub-key":
		getPubKey(logger, options)
		return
	case "register-node":
		registerNode(logger, options)
		return
	case "get-all-nodes":
		getAllNodes(logger, options)
		return
	}

}
