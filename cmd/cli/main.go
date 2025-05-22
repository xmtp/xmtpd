package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xmtp/xmtpd/pkg/currency"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xmtp/xmtpd/pkg/blockchain/migrator"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/stress"

	"github.com/jessevdk/go-flags"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var Version string = "unknown"

type CLI struct {
	config.GlobalOptions
	Command                string
	GetPubKey              config.GetPubKeyOptions
	GenerateKey            config.GenerateKeyOptions
	RegisterNode           config.RegisterNodeOptions
	NetworkAdminOptions    config.NetworkAdminOptions
	GetAllNodes            config.GetAllNodesOptions
	GetNode                config.GetNodeOptions
	SetHttpAddress         config.SetHttpAddressOptions
	MigrateNodes           config.MigrateNodesOptions
	AddRates               config.AddRatesOptions
	GetRates               config.GetRatesOptions
	IdentityUpdatesStress  config.IdentityUpdatesStressOptions
	Watcher                config.WatcherOptions
	GetMaxCanonicalOptions config.GetMaxCanonicalOptions
	SetMaxCanonicalOptions config.SetMaxCanonicalOptions
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
	var networkAdminOptions config.NetworkAdminOptions
	var getPubKeyOptions config.GetPubKeyOptions
	var getAllNodesOptions config.GetAllNodesOptions
	var setHttpAddressOptions config.SetHttpAddressOptions
	var migrateNodesOptions config.MigrateNodesOptions
	var addRatesOptions config.AddRatesOptions
	var getRatesOptions config.GetRatesOptions
	var getNodeOptions config.GetNodeOptions
	var identityUpdatesStressOptions config.IdentityUpdatesStressOptions
	var watcherOptions config.WatcherOptions
	var getMaxCanonicalOptions config.GetMaxCanonicalOptions
	var setMaxCanonicalOptions config.SetMaxCanonicalOptions
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
	if _, err := parser.AddCommand("add-node-to-network", "Add a node to the network", "", &networkAdminOptions); err != nil {
		return nil, fmt.Errorf("could not add add-node-to-network command: %s", err)
	}
	if _, err := parser.AddCommand("remove-node-from-network", "Remove a node from the network", "", &networkAdminOptions); err != nil {
		return nil, fmt.Errorf("could not add remove-node-from-network command: %s", err)
	}
	if _, err := parser.AddCommand("set-max-canonical", "Set the maximum canonical size", "", &setMaxCanonicalOptions); err != nil {
		return nil, fmt.Errorf("could not add set-max-canonical command: %s", err)
	}
	if _, err := parser.AddCommand("migrate-nodes", "Migrate nodes from a file", "", &migrateNodesOptions); err != nil {
		return nil, fmt.Errorf("could not add migrate-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("set-http-address", "Set the HTTP address of a node", "", &setHttpAddressOptions); err != nil {
		return nil, fmt.Errorf("could not add set-http-address command: %s", err)
	}
	if _, err := parser.AddCommand("add-rates", "Add rates of the rates manager", "", &addRatesOptions); err != nil {
		return nil, fmt.Errorf("could not add add-rates command: %s", err)
	}

	// Getter commands
	if _, err := parser.AddCommand("get-all-nodes", "Get all nodes from the registry", "", &getAllNodesOptions); err != nil {
		return nil, fmt.Errorf("could not add get-all-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("get-node", "Get a node from the registry", "", &getNodeOptions); err != nil {
		return nil, fmt.Errorf("could not add get-node command: %s", err)
	}
	if _, err := parser.AddCommand("get-rates", "Get rates of the rates manager", "", &getRatesOptions); err != nil {
		return nil, fmt.Errorf("could not add get-rates command: %s", err)
	}
	if _, err := parser.AddCommand("get-max-canonical-nodes", "Get max canonical nodes for network", "", &getMaxCanonicalOptions); err != nil {
		return nil, fmt.Errorf("could not add get-rates command: %s", err)
	}

	// Dev commands
	if _, err := parser.AddCommand("identity-updates-stress", "Stress the identity updates contract", "", &identityUpdatesStressOptions); err != nil {
		return nil, fmt.Errorf("could not add identity-updates-stress command: %s", err)
	}
	if _, err := parser.AddCommand("start-watcher", "Start the blockchain watcher", "", &watcherOptions); err != nil {
		return nil, fmt.Errorf("could not add start-watcher command: %s", err)
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
		getPubKeyOptions,
		generateKeyOptions,
		registerNodeOptions,
		networkAdminOptions,
		getAllNodesOptions,
		getNodeOptions,
		setHttpAddressOptions,
		migrateNodesOptions,
		addRatesOptions,
		getRatesOptions,
		identityUpdatesStressOptions,
		watcherOptions,
		getMaxCanonicalOptions,
		setMaxCanonicalOptions,
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
	if !options.RegisterNode.Force &&
		isPubKeyAlreadyRegistered(logger, options, options.RegisterNode.SigningKeyPub) {
		logger.Info(
			"provided public key is already registered",
			zap.String("pub-key", options.RegisterNode.SigningKeyPub),
		)
		return
	}

	ctx := context.Background()
	registryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.RegisterNode.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
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
}

func isPubKeyAlreadyRegistered(logger *zap.Logger, options *CLI, pubKey string) bool {
	chainClient, err := blockchain.NewClient(
		context.Background(),
		options.Contracts.SettlementChain.RpcURL,
	)
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
		logger.Fatal(
			"could not retrieve migrated nodes from registry",
			zap.Error(err),
			zap.Any("contracts", options.Contracts),
		)
	}

	for _, node := range nodes {
		if node.SigningKeyPub == pubKey {
			return true
		}
	}

	return false
}

func addNodeToNetwork(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.NetworkAdminOptions.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.AddToNetwork(
		ctx,
		options.NetworkAdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not add node to network", zap.Error(err))
	}
}

func setMaxCanonical(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.SetMaxCanonicalOptions.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetMaxCanonical(
		ctx,
		options.SetMaxCanonicalOptions.Limit,
	)
	if err != nil {
		logger.Fatal("could not set max canonical", zap.Error(err))
	}
}

func removeNodeFromNetwork(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.NetworkAdminOptions.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.RemoveFromNetwork(
		ctx,
		options.NetworkAdminOptions.NodeId,
	)
	if err != nil {
		logger.Fatal("could not remove node from network", zap.Error(err))
	}
}

func migrateNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	nodes, err := migrator.ImportNodesFromFile(options.MigrateNodes.InFile)
	if err != nil {
		logger.Fatal("could not import nodes from file", zap.Error(err))
	}

	newRegistryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.MigrateNodes.AdminOptions.AdminPrivateKey,
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

func getRates(logger *zap.Logger, options *CLI) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.SettlementChain.RpcURL)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	fetcher, err := fees.NewContractRatesFetcher(ctx, chainClient, logger, options.Contracts)
	if err != nil {
		logger.Fatal("could not create rates fetcher", zap.Error(err))
	}

	err = fetcher.Start()
	if err != nil {
		if strings.Contains(err.Error(), "no rates found") {
			logger.Info("no rates found")
			return
		}
		logger.Fatal("could not start rates fetcher", zap.Error(err))
	}

	rates, err := fetcher.GetRates(time.Now())
	if err != nil {
		logger.Fatal("could not get rates", zap.Error(err))
	}

	logger.Info("rates fetched successfully", zap.Any("rates", rates))
}

func addRates(logger *zap.Logger, options *CLI) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()

	registryAdmin, err := setupRateRegistryAdmin(
		ctx,
		logger,
		options.AddRates.AdminOptions.AdminPrivateKey,
		options.Contracts.SettlementChain.ChainID,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	if options.AddRates.MessageFee > math.MaxInt64 {
		logger.Fatal("adding message fee too big", zap.Uint64("fee", options.AddRates.MessageFee))
	}
	if options.AddRates.StorageFee > math.MaxInt64 {
		logger.Fatal("adding storage fee too big", zap.Uint64("fee", options.AddRates.StorageFee))
	}
	if options.AddRates.CongestionFee > math.MaxInt64 {
		logger.Fatal(
			"adding congestion fee too big",
			zap.Uint64("fee", options.AddRates.CongestionFee),
		)
	}

	rates := fees.Rates{
		MessageFee:          currency.PicoDollar(options.AddRates.MessageFee),
		StorageFee:          currency.PicoDollar(options.AddRates.MessageFee),
		CongestionFee:       currency.PicoDollar(options.AddRates.MessageFee),
		TargetRatePerMinute: options.AddRates.TargetRate,
	}

	err = registryAdmin.AddRates(ctx, rates)
	if err != nil {
		logger.Fatal("could not add rates to registry", zap.Error(err))
	}
}

/*
*
Node manager commands
*
*/

func setHttpAddress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	registryAdmin, err := setupNodeRegistryAdmin(
		ctx,
		logger,
		options.SetHttpAddress.NodeManagerOptions.NodePrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup registry admin", zap.Error(err))
	}

	err = registryAdmin.SetHttpAddress(
		ctx,
		options.SetHttpAddress.NodeId,
		options.SetHttpAddress.Address,
	)
	if err != nil {
		logger.Fatal("could not set http address", zap.Error(err))
	}
}

/*
*
Getter commands
*
*/

func getMaxCanonicalNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.SettlementChain.RpcURL)
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

	limit, err := caller.GetMaxCanonicalNodes(ctx)
	if err != nil {
		logger.Fatal("could not get max canonical nodes", zap.Error(err))
	}
	logger.Info("max canonical nodes retrieved successfully", zap.Any("limit", limit))
}

func getAllNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.SettlementChain.RpcURL)
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
	chainClient, err := blockchain.NewClient(ctx, options.Contracts.SettlementChain.RpcURL)
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

	node, err := caller.GetNode(ctx, options.GetNode.NodeId)
	if err != nil {
		logger.Fatal("could not retrieve node from registry", zap.Error(err))
	}

	logger.Info(
		"got nodes",
		zap.Any("node", node),
	)
}

func identityUpdatesStress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()

	logger.Info(
		"creating identity updates",
		zap.Int("count", options.IdentityUpdatesStress.Count),
		zap.String("contract", options.IdentityUpdatesStress.Contract),
	)

	err := stress.StressIdentityUpdates(
		ctx,
		logger,
		options.IdentityUpdatesStress.Count,
		options.IdentityUpdatesStress.Contract,
		options.IdentityUpdatesStress.Rpc,
		options.IdentityUpdatesStress.PrivateKey,
		options.IdentityUpdatesStress.Async,
	)
	if err != nil {
		logger.Fatal("could not create identity updates", zap.Error(err))
	}
}

func startChainWatcher(logger *zap.Logger, options *CLI) {
	ctxwc, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	watcher, err := stress.NewWatcher(
		ctxwc,
		logger,
		options.Watcher.Wss,
		common.HexToAddress(options.Watcher.Contract),
	)
	if err != nil {
		logger.Fatal("could not create watcher", zap.Error(err))
	}

	err = watcher.Listen(ctxwc)
	if err != nil {
		logger.Fatal("could not listen", zap.Error(err))
	}
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
	case "add-node-to-network":
		addNodeToNetwork(logger, options)
		return
	case "remove-node-from-network":
		removeNodeFromNetwork(logger, options)
		return
	case "set-max-canonical":
		setMaxCanonical(logger, options)
		return
	case "migrate-nodes":
		migrateNodes(logger, options)
		return
	case "set-http-address":
		setHttpAddress(logger, options)
		return
	case "get-all-nodes":
		getAllNodes(logger, options)
		return
	case "get-node":
		getNode(logger, options)
		return
	case "get-rates":
		getRates(logger, options)
		return
	case "add-rates":
		addRates(logger, options)
		return
	case "get-max-canonical-nodes":
		getMaxCanonicalNodes(logger, options)
		return
	case "identity-updates-stress":
		identityUpdatesStress(logger, options)
		return
	case "start-watcher":
		startChainWatcher(logger, options)
		return
	}
}

// setupNodeRegistryAdmin creates and returns a node registry admin
func setupNodeRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	options *CLI,
) (blockchain.INodeRegistryAdmin, error) {
	chainClient, err := blockchain.NewClient(
		ctx,
		options.Contracts.SettlementChain.RpcURL,
	)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		options.Contracts.SettlementChain.ChainID,
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

// setupRateRegistryAdmin creates and returns a rate registry admin
func setupRateRegistryAdmin(
	ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	chainID int,
	options *CLI,
) (*blockchain.RatesAdmin, error) {
	chainClient, err := blockchain.NewClient(
		ctx,
		options.Contracts.SettlementChain.RpcURL,
	)
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

	registryAdmin, err := blockchain.NewRatesAdmin(
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
