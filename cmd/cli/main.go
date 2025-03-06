package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xmtp/xmtpd/contracts/pkg/ratesmanager"
	"github.com/xmtp/xmtpd/pkg/blockchain/migrator"
	"github.com/xmtp/xmtpd/pkg/config"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var Version string = "unknown"

type CLI struct {
	config.GlobalOptions
	Command       string
	GetPubKey     config.GetPubKeyOptions
	GenerateKey   config.GenerateKeyOptions
	RegisterNode  config.RegisterNodeOptions
	GetAllNodes   config.GetAllNodesOptions
	UpdateHealth  config.UpdateHealthOptions
	UpdateAddress config.UpdateAddressOptions
	MigrateNodes  config.MigrateNodesOptions
	AddRates      config.AddRatesOptions
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
	var options config.GlobalOptions
	var generateKeyOptions config.GenerateKeyOptions
	var registerNodeOptions config.RegisterNodeOptions
	var getPubKeyOptions config.GetPubKeyOptions
	var getAllNodesOptions config.GetAllNodesOptions
	var updateHealthOptions config.UpdateHealthOptions
	var updateAddressOptions config.UpdateAddressOptions
	var migrateNodesOptions config.MigrateNodesOptions
	var addRatesOptions config.AddRatesOptions

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
	if _, err := parser.AddCommand("mark-healthy", "Mark a node as healthy in the registry", "", &updateHealthOptions); err != nil {
		return nil, fmt.Errorf("Could not add mark-healthy command: %s", err)
	}
	if _, err := parser.AddCommand("mark-unhealthy", "Mark a node as unhealthy in the registry", "", &updateHealthOptions); err != nil {
		return nil, fmt.Errorf("Could not add mark-unhealthy command: %s", err)
	}
	if _, err := parser.AddCommand("update-address", "Update HTTP address of a node", "", &updateAddressOptions); err != nil {
		return nil, fmt.Errorf("Could not add update-address command: %s", err)
	}
	if _, err := parser.AddCommand("migrate-nodes", "Migrate nodes from a file", "", &migrateNodesOptions); err != nil {
		return nil, fmt.Errorf("Could not add dump nodes command: %s", err)
	}
	if _, err := parser.AddCommand("add-rates", "Add rates to the rates manager", "", &addRatesOptions); err != nil {
		return nil, fmt.Errorf("Could not add add-rates command: %s", err)
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
		updateHealthOptions,
		updateAddressOptions,
		migrateNodesOptions,
		addRatesOptions,
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
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)
	privKey.Public()
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
		blockchain.RegistryAdminV1,
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
		zap.String("node-owner-address", options.RegisterNode.OwnerAddress),
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
		zap.String("address", utils.EcdsaPublicKeyToAddress(privKey.Public().(*ecdsa.PublicKey))),
	)
}

func getAllNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		options.Contracts,
		blockchain.RegistryCallerV1,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	nodes, err := migrator.ReadFromRegistry[migrator.SerializableNodeV1](caller)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	logger.Info(
		"got nodes",
		zap.Int("size", len(nodes)),
		zap.Any("nodes", nodes),
	)

	if options.GetAllNodes.OutFile != "" {
		err = migrator.DumpNodesToFile(nodes, options.GetAllNodes.OutFile)
		if err != nil {
			logger.Fatal("could not dump nodes", zap.Error(err))
		}
	}
}

func updateHealth(logger *zap.Logger, options *CLI, health bool) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.UpdateHealth.AdminPrivateKey,
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
		blockchain.RegistryAdminV1,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	err = registryAdmin.UpdateHealth(ctx, options.UpdateHealth.NodeId, health)
	if err != nil {
		logger.Fatal("could not update node health in registry", zap.Error(err))
	}
}

func updateAddress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.UpdateAddress.PrivateKey,
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
		blockchain.RegistryAdminV1,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	err = registryAdmin.UpdateHttpAddress(
		ctx,
		options.UpdateAddress.NodeId,
		options.UpdateAddress.Address,
	)
	if err != nil {
		logger.Fatal("could not update node address in registry", zap.Error(err))
	}
}

func migrateNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	nodes, err := migrator.ImportNodesFromFile(options.MigrateNodes.InFile)
	if err != nil {
		logger.Fatal("could not import nodes from file", zap.Error(err))
	}

	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.MigrateNodes.AdminPrivateKey,
		options.Contracts.ChainID,
	)
	if err != nil {
		logger.Fatal("could not create signer", zap.Error(err))
	}

	registryAdminV2, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
		blockchain.RegistryAdminV2,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	err = migrator.WriteToRegistryV2(logger, nodes, registryAdminV2)
	if err != nil {
		logger.Fatal("could not write nodes to registry", zap.Error(err))
	}
}

func addRates(logger *zap.Logger, options *CLI) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		options.AddRates.AdminPrivateKey,
		options.Contracts.ChainID,
	)
	if err != nil {
		logger.Fatal("could not create signer", zap.Error(err))
	}

	ratesManager, err := blockchain.NewRatesAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create rates admin", zap.Error(err))
	}

	startTime := time.Now().Add(time.Duration(options.AddRates.DelayDays) * 24 * time.Hour)

	rates := ratesmanager.RatesManagerRates{
		MessageFee:    options.AddRates.MessageFee,
		StorageFee:    options.AddRates.StorageFee,
		CongestionFee: options.AddRates.CongestionFee,
		StartTime:     uint64(startTime.Unix()),
	}

	if err = ratesManager.AddRates(ctx, rates); err != nil {
		logger.Fatal("could not add rates", zap.Error(err))
	}

	logger.Info("added rates", zap.Any("rates", rates))
}
func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Version: %s\n", Version)
			return
		}
	}

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
	case "mark-healthy":
		updateHealth(logger, options, true)
		return
	case "mark-unhealthy":
		updateHealth(logger, options, false)
		return
	case "update-address":
		updateAddress(logger, options)
		return
	case "migrate-nodes":
		migrateNodes(logger, options)
		return
	case "add-rates":
		addRates(logger, options)
		return
	}
}
