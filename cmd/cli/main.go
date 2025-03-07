package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"os"

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
	Command                          string
	AdminOptions                     config.AdminOptions
	NodeManagerOptions               config.NodeManagerOptions
	GetPubKey                        config.GetPubKeyOptions
	GenerateKey                      config.GenerateKeyOptions
	RegisterNode                     config.RegisterNodeOptions
	GetAllNodes                      config.GetAllNodesOptions
	SetHttpAddress                   config.SetHttpAddressOptions
	MigrateNodes                     config.MigrateNodesOptions
	NodeOperator                     config.NodeOperatorOptions
	SetMinMonthlyFee                 config.SetMinMonthlyFeeOptions
	SetMaxActiveNodes                config.SetMaxActiveNodesOptions
	SetNodeOperatorCommissionPercent config.SetNodeOperatorCommissionPercentOptions
	GetOptions                       config.GetOptions
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
	var adminOptions config.AdminOptions
	var nodeManagerOptions config.NodeManagerOptions
	var generateKeyOptions config.GenerateKeyOptions
	var registerNodeOptions config.RegisterNodeOptions
	var getPubKeyOptions config.GetPubKeyOptions
	var getAllNodesOptions config.GetAllNodesOptions
	var setHttpAddressOptions config.SetHttpAddressOptions
	var migrateNodesOptions config.MigrateNodesOptions
	var nodeOperatorOptions config.NodeOperatorOptions
	var setMinMonthlyFeeOptions config.SetMinMonthlyFeeOptions
	var setMaxActiveNodesOptions config.SetMaxActiveNodesOptions
	var setNodeOperatorCommissionPercentOptions config.SetNodeOperatorCommissionPercentOptions
	var getOptions config.GetOptions
	parser := flags.NewParser(&options, flags.Default)

	// Admin commands
	if _, err := parser.AddCommand("generate-key", "Generate a public/private keypair", "", &generateKeyOptions); err != nil {
		return nil, fmt.Errorf("could not add generate-key command: %s", err)
	}
	if _, err := parser.AddCommand("get-pub-key", "Get the public key for a private key", "", &getPubKeyOptions); err != nil {
		return nil, fmt.Errorf("could not add get-pub-key command: %s", err)
	}
	if _, err := parser.AddCommand("register-node", "Register a node", "", &registerNodeOptions); err != nil {
		return nil, fmt.Errorf("could not add register-node command: %s", err)
	}
	if _, err := parser.AddCommand("migrate-nodes", "Migrate nodes from a file", "", &migrateNodesOptions); err != nil {
		return nil, fmt.Errorf("could not add migrate-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("disable-node", "Disable a node", "", &adminOptions); err != nil {
		return nil, fmt.Errorf("could not add disable-node command: %s", err)
	}
	if _, err := parser.AddCommand("enable-node", "Enable a node", "", &adminOptions); err != nil {
		return nil, fmt.Errorf("could not add enable-node command: %s", err)
	}
	if _, err := parser.AddCommand("remove-from-api-nodes", "Remove a node from the API nodes", "", &adminOptions); err != nil {
		return nil, fmt.Errorf("could not add remove-from-api-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("remove-from-replication-nodes", "Remove a node from the replication nodes", "", &adminOptions); err != nil {
		return nil, fmt.Errorf("could not add remove-from-replication-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("set-http-address", "Set the HTTP address of a node", "", &setHttpAddressOptions); err != nil {
		return nil, fmt.Errorf("could not add set-http-address command: %s", err)
	}
	if _, err := parser.AddCommand("set-min-monthly-fee", "Set the minimum monthly fee of a node", "", &setMinMonthlyFeeOptions); err != nil {
		return nil, fmt.Errorf("could not add set-min-monthly-fee command: %s", err)
	}
	if _, err := parser.AddCommand("set-max-active-nodes", "Set the maximum number of active nodes", "", &setMaxActiveNodesOptions); err != nil {
		return nil, fmt.Errorf("could not add set-max-active-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("set-node-operator-commission-percent", "Set the node operator commission percent", "", &setNodeOperatorCommissionPercentOptions); err != nil {
		return nil, fmt.Errorf(
			"could not add set-node-operator-commission-percent command: %s",
			err,
		)
	}

	// Node operator commands
	if _, err := parser.AddCommand("set-api-enabled", "Set API enabled for a node", "", &nodeOperatorOptions); err != nil {
		return nil, fmt.Errorf("could not add set-api-enabled command: %s", err)
	}
	if _, err := parser.AddCommand("set-replication-enabled", "Set replication enabled for a node", "", &nodeOperatorOptions); err != nil {
		return nil, fmt.Errorf("could not add set-replication-enabled command: %s", err)
	}

	// Getter commands
	if _, err := parser.AddCommand("get-all-nodes", "Get all nodes from the registry", "", &getAllNodesOptions); err != nil {
		return nil, fmt.Errorf("could not add get-all-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("get-active-api-nodes", "Get all active API nodes from the registry", "", &getOptions); err != nil {
		return nil, fmt.Errorf("could not add get-active-api-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("get-active-replication-nodes", "Get all active replication nodes from the registry", "", &getOptions); err != nil {
		return nil, fmt.Errorf("could not add get-active-replication-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("get-node", "Get a node from the registry", "", &getOptions); err != nil {
		return nil, fmt.Errorf("could not add get-node command: %s", err)
	}
	if _, err := parser.ParseArgs(args); err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			return nil, fmt.Errorf("could not parse options: %s", err)
		}
		return nil, nil
	}

	if parser.Active == nil {
		return nil, errors.New("no command provided")
	}

	return &CLI{
		options,
		parser.Active.Name,
		adminOptions,
		nodeManagerOptions,
		getPubKeyOptions,
		generateKeyOptions,
		registerNodeOptions,
		getAllNodesOptions,
		setHttpAddressOptions,
		migrateNodesOptions,
		nodeOperatorOptions,
		setMinMonthlyFeeOptions,
		setMaxActiveNodesOptions,
		setNodeOperatorCommissionPercentOptions,
		getOptions,
	}, nil
}

/*
*
Key commands
*
*/

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

/*
*
Admin commands
*
*/

func registerNode(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.RegisterNode.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	signingKeyPub, err := utils.ParseEcdsaPublicKey(options.RegisterNode.SigningKeyPub)
	if err != nil {
		logger.Fatal("could not decompress public key", zap.Error(err))
	}

	minMonthlyFee := int64(0)
	if options.RegisterNode.MinMonthlyFee < 0 {
		logger.Fatal("provided negative min monthly fee is not allowed")
	}
	minMonthlyFee = options.RegisterNode.MinMonthlyFee

	err = registryAdmin.AddNode(
		ctx,
		options.RegisterNode.OwnerAddress,
		signingKeyPub,
		options.RegisterNode.HttpAddress,
		minMonthlyFee,
	)
	if err != nil {
		logger.Fatal("could not add node", zap.Error(err))
	}
}

func disableNode(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.DisableNode(
		ctx,
		options.AdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not disable node", zap.Error(err))
	}
}

func enableNode(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.EnableNode(
		ctx,
		options.AdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not enable node", zap.Error(err))
	}
}

func removeFromApiNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.RemoveFromApiNodes(
		ctx,
		options.AdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not remove from api nodes", zap.Error(err))
	}
}

func removeFromReplicationNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.RemoveFromReplicationNodes(
		ctx,
		options.AdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not remove from replication nodes", zap.Error(err))
	}
}

func migrateNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	nodes, err := migrator.ImportNodesFromFile(options.MigrateNodes.InFile)
	if err != nil {
		logger.Fatal("could not import nodes from file", zap.Error(err))
	}

	newRegistryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.MigrateNodes.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = migrator.WriteToRegistry(logger, nodes, newRegistryAdmin)
	if err != nil {
		logger.Fatal("could not write nodes to registry", zap.Error(err))
	}
}

func setMaxActiveNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.SetMaxActiveNodes.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetMaxActiveNodes(
		ctx,
		options.SetMaxActiveNodes.MaxActiveNodes,
	)
	if err != nil {
		logger.Fatal("could not set max active nodes", zap.Error(err))
	}
}

func setNodeOperatorCommissionPercent(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.SetNodeOperatorCommissionPercent.AdminOptions.AdminPrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetNodeOperatorCommissionPercent(
		ctx,
		options.SetNodeOperatorCommissionPercent.CommissionPercent,
	)
	if err != nil {
		logger.Fatal("could not set node operator commission percent", zap.Error(err))
	}
}

/*
*
Node manager commands
*
*/

func setHttpAddress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.SetHttpAddress.NodeManagerOptions.NodePrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetHttpAddress(
		ctx,
		options.NodeManagerOptions.NodeId,
		options.SetHttpAddress.Address,
	)
	if err != nil {
		logger.Fatal("could not set http address", zap.Error(err))
	}
}

func setMinMonthlyFee(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.SetMinMonthlyFee.NodeManagerOptions.NodePrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	if options.SetMinMonthlyFee.MinMonthlyFee < 0 {
		logger.Fatal("invalid negative minMonthlyFee provided")
	}

	err = registryAdmin.SetMinMonthlyFee(
		ctx,
		options.NodeManagerOptions.NodeId,
		options.SetMinMonthlyFee.MinMonthlyFee,
	)
	if err != nil {
		logger.Fatal("could not set min monthly fee", zap.Error(err))
	}
}

/*
*
Node operator commands
*
*/

func setApiEnabled(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.NodeOperator.NodePrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetIsApiEnabled(
		ctx,
		options.NodeOperator.NodeId,
		options.NodeOperator.Enable,
	)
	if err != nil {
		logger.Fatal("could not set API enabled", zap.Error(err))
	}
}

func setReplicationEnabled(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupRegistryAdmin(
		ctx,
		logger,
		options.NodeOperator.NodePrivateKey,
		options.Contracts.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetIsReplicationEnabled(
		ctx,
		options.NodeOperator.NodeId,
		options.NodeOperator.Enable,
	)
	if err != nil {
		logger.Fatal("could not set replication enabled", zap.Error(err))
	}
}

/*
*
Getter commands
*
*/

func getActiveApiNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry caller", zap.Error(err))
	}

	nodes, err := caller.GetActiveApiNodes(ctx)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	logger.Info(
		"got nodes",
		zap.Int("size", len(nodes)),
		zap.Any("nodes", nodes),
	)
}

func getActiveReplicationNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	nodes, err := caller.GetActiveReplicationNodes(ctx)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	logger.Info(
		"got nodes",
		zap.Int("size", len(nodes)),
		zap.Any("nodes", nodes),
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
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	nodes, err := migrator.ReadFromRegistry(caller)
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

func getNode(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	node, err := caller.GetNode(ctx, options.NodeOperator.NodeId)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	logger.Info(
		"got nodes",
		zap.Any("node", node),
	)
}

/*
*
Main Command
*
*/

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("Version: %s\n", Version)
			return
		}
	}

	options, err := parseOptions(os.Args[1:])
	if err != nil {
		log.Fatalf("could not parse options: %s", err)
	}
	if options == nil {
		return
	}

	logger, _, err := utils.BuildLogger(options.Log)
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
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
	case "disable-node":
		disableNode(logger, options)
		return
	case "enable-node":
		enableNode(logger, options)
		return
	case "migrate-nodes":
		migrateNodes(logger, options)
		return
	case "remove-from-api-nodes":
		removeFromApiNodes(logger, options)
		return
	case "remove-from-replication-nodes":
		removeFromReplicationNodes(logger, options)
		return
	case "set-http-address":
		setHttpAddress(logger, options)
		return
	case "set-min-monthly-fee":
		setMinMonthlyFee(logger, options)
		return
	case "set-max-active-nodes":
		setMaxActiveNodes(logger, options)
		return
	case "set-node-operator-commission-percent":
		setNodeOperatorCommissionPercent(logger, options)
		return
	case "set-api-enabled":
		setApiEnabled(logger, options)
		return
	case "set-replication-enabled":
		setReplicationEnabled(logger, options)
		return
	case "get-active-api-nodes":
		getActiveApiNodes(logger, options)
		return
	case "get-active-replication-nodes":
		getActiveReplicationNodes(logger, options)
		return
	case "get-all-nodes":
		getAllNodes(logger, options)
		return
	case "get-node":
		getNode(logger, options)
		return
	}
}

// setupRegistryAdmin creates and returns a node registry admin
func setupRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	chainID int,
	options *CLI,
) (blockchain.INodeRegistryAdmin, error) {
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.RpcUrl)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		chainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry admin: %w", err)
	}

	return registryAdmin, nil
}
