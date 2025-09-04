package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	GetBootstrapperAddress config.GetBootstrapperAddressOptions
	SetBootstrapperAddress config.SetBootstrapperAddressOptions
	SetPause               config.SetPauseOptions
	GetPause               config.GetPauseOptions
	GetNodeRegistryAdmin   config.GetNodeRegistryAdminOptions
	SetNodeRegistryAdmin   config.SetNodeRegistryAdminOptions
	GetPayloadSize         config.GetPayloadSizeOptions
	SetPayloadSize         config.SetPayloadSizeOptions
	GetDMFeesRecipient     config.GetDistributionManagerProtocolFeesRecipientOptions
	SetDMFeesRecipient     config.SetDistributionManagerProtocolFeesRecipientOptions
	GetPayerMinDeposit     config.GetPayerMinimumDepositOptions
	SetPayerMinDeposit     config.SetPayerMinimumDepositOptions
	GetPayerWithdrawLock   config.GetPayerWithdrawLockPeriodOptions
	SetPayerWithdrawLock   config.SetPayerWithdrawLockPeriodOptions
	GetPRMFeeRate          config.GetPayerReportProtocolFeeRateOptions
	SetPRMFeeRate          config.SetPayerReportProtocolFeeRateOptions
	GetRateMigrator        config.GetRateRegistryMigratorOptions
	SetRateMigrator        config.SetRateRegistryMigratorOptions
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
	var getBootstrapperAddressOptions config.GetBootstrapperAddressOptions
	var setBootstrapperAddressOptions config.SetBootstrapperAddressOptions
	var setPauseOptions config.SetPauseOptions
	var getPauseOptions config.GetPauseOptions

	var getNodeRegistryAdminOptions config.GetNodeRegistryAdminOptions
	var setNodeRegistryAdminOptions config.SetNodeRegistryAdminOptions

	var getPayloadSizeOptions config.GetPayloadSizeOptions
	var setPayloadSizeOptions config.SetPayloadSizeOptions

	var getDMFeesRecipientOptions config.GetDistributionManagerProtocolFeesRecipientOptions
	var setDMFeesRecipientOptions config.SetDistributionManagerProtocolFeesRecipientOptions

	var getPayerMinDepositOptions config.GetPayerMinimumDepositOptions
	var setPayerMinDepositOptions config.SetPayerMinimumDepositOptions

	var getPayerWithdrawLockOptions config.GetPayerWithdrawLockPeriodOptions
	var setPayerWithdrawLockOptions config.SetPayerWithdrawLockPeriodOptions

	var getPRMFeeRateOptions config.GetPayerReportProtocolFeeRateOptions
	var setPRMFeeRateOptions config.SetPayerReportProtocolFeeRateOptions

	var getRateMigratorOptions config.GetRateRegistryMigratorOptions
	var setRateMigratorOptions config.SetRateRegistryMigratorOptions
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
	if _, err := parser.AddCommand("set-bootstrapper-address", "Set bootstrapper address for V3 migration", "", &setBootstrapperAddressOptions); err != nil {
		return nil, fmt.Errorf("could not add set-bootstrapper-address command: %s", err)
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
		return nil, fmt.Errorf("could not add get-max-canonical-nodes command: %s", err)
	}
	if _, err := parser.AddCommand("get-bootstrapper-address", "Get bootstrapper address for V3 migration", "", &getBootstrapperAddressOptions); err != nil {
		return nil, fmt.Errorf("could not add get-bootstrapper-address command: %s", err)
	}

	// Dev commands
	if _, err := parser.AddCommand("identity-updates-stress", "Stress the identity updates contract", "", &identityUpdatesStressOptions); err != nil {
		return nil, fmt.Errorf("could not add identity-updates-stress command: %s", err)
	}
	if _, err := parser.AddCommand("start-watcher", "Start the blockchain watcher", "", &watcherOptions); err != nil {
		return nil, fmt.Errorf("could not add start-watcher command: %s", err)
	}

	// Control pause status (both broadcasters)
	if _, err := parser.AddCommand("set-pause", "Set pause status for a broadcaster (identity/group)", "", &setPauseOptions); err != nil {
		return nil, fmt.Errorf("could not add set-pause command: %s", err)
	}
	if _, err := parser.AddCommand("get-pause", "Get pause status for a broadcaster (identity/group)", "", &getPauseOptions); err != nil {
		return nil, fmt.Errorf("could not add get-pause command: %s", err)
	}

	// --- NodeRegistry admin address
	if _, err := parser.AddCommand("get-node-admin", "Get node registry admin address", "", &getNodeRegistryAdminOptions); err != nil {
		return nil, fmt.Errorf("could not add get-node-admin command: %s", err)
	}
	if _, err := parser.AddCommand("set-node-admin", "Set node registry admin address", "", &setNodeRegistryAdminOptions); err != nil {
		return nil, fmt.Errorf("could not add set-node-admin command: %s", err)
	}

	// --- Payload size (unified for identity|group, min|max)
	if _, err := parser.AddCommand("get-payload-size", "Get payload size bound on a broadcaster (identity/group)", "", &getPayloadSizeOptions); err != nil {
		return nil, fmt.Errorf("could not add get-payload-size command: %s", err)
	}
	if _, err := parser.AddCommand("set-payload-size", "Set payload size bound on a broadcaster (identity/group)", "", &setPayloadSizeOptions); err != nil {
		return nil, fmt.Errorf("could not add set-payload-size command: %s", err)
	}

	// --- DistributionManager: protocolFeesRecipient
	if _, err := parser.AddCommand("get-dm-protocol-fees-recipient", "Get DistributionManager protocol fees recipient", "", &getDMFeesRecipientOptions); err != nil {
		return nil, fmt.Errorf("could not add get-dm-protocol-fees-recipient command: %s", err)
	}
	if _, err := parser.AddCommand("set-dm-protocol-fees-recipient", "Set DistributionManager protocol fees recipient", "", &setDMFeesRecipientOptions); err != nil {
		return nil, fmt.Errorf("could not add set-dm-protocol-fees-recipient command: %s", err)
	}

	// --- PayerRegistry: minimumDeposit (uint96 picodollars)
	if _, err := parser.AddCommand("get-payer-min-deposit", "Get PayerRegistry minimum deposit (uint96 picodollars)", "", &getPayerMinDepositOptions); err != nil {
		return nil, fmt.Errorf("could not add get-payer-min-deposit command: %s", err)
	}
	if _, err := parser.AddCommand("set-payer-min-deposit", "Set PayerRegistry minimum deposit (uint96 picodollars)", "", &setPayerMinDepositOptions); err != nil {
		return nil, fmt.Errorf("could not add set-payer-min-deposit command: %s", err)
	}

	// --- PayerRegistry: withdrawLockPeriod (uint32 seconds)
	if _, err := parser.AddCommand("get-payer-withdraw-lock", "Get PayerRegistry withdraw lock period (seconds)", "", &getPayerWithdrawLockOptions); err != nil {
		return nil, fmt.Errorf("could not add get-payer-withdraw-lock command: %s", err)
	}
	if _, err := parser.AddCommand("set-payer-withdraw-lock", "Set PayerRegistry withdraw lock period (seconds)", "", &setPayerWithdrawLockOptions); err != nil {
		return nil, fmt.Errorf("could not add set-payer-withdraw-lock command: %s", err)
	}

	// --- PayerReportManager: protocolFeeRate (uint16, bps)
	if _, err := parser.AddCommand("get-prm-fee-rate", "Get PayerReportManager protocol fee rate (bps, uint16)", "", &getPRMFeeRateOptions); err != nil {
		return nil, fmt.Errorf("could not add get-prm-fee-rate command: %s", err)
	}
	if _, err := parser.AddCommand("set-prm-fee-rate", "Set PayerReportManager protocol fee rate (bps, uint16)", "", &setPRMFeeRateOptions); err != nil {
		return nil, fmt.Errorf("could not add set-prm-fee-rate command: %s", err)
	}

	// --- RateRegistry: migrator (address)
	if _, err := parser.AddCommand("get-rate-migrator", "Get RateRegistry migrator address", "", &getRateMigratorOptions); err != nil {
		return nil, fmt.Errorf("could not add get-rate-migrator command: %s", err)
	}
	if _, err := parser.AddCommand("set-rate-migrator", "Set RateRegistry migrator address", "", &setRateMigratorOptions); err != nil {
		return nil, fmt.Errorf("could not add set-rate-migrator command: %s", err)
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

	if err := config.ParseJSONConfig(&options.Contracts); err != nil {
		return nil, fmt.Errorf("could not parse contracts JSON config: %s", err)
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
		getBootstrapperAddressOptions,
		setBootstrapperAddressOptions,
		setPauseOptions,
		getPauseOptions,
		getNodeRegistryAdminOptions,
		setNodeRegistryAdminOptions,
		getPayloadSizeOptions,
		setPayloadSizeOptions,
		getDMFeesRecipientOptions,
		setDMFeesRecipientOptions,
		getPayerMinDepositOptions,
		setPayerMinDepositOptions,
		getPayerWithdrawLockOptions,
		setPayerWithdrawLockOptions,
		getPRMFeeRateOptions,
		setPRMFeeRateOptions,
		getRateMigratorOptions,
		setRateMigratorOptions,
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

	_, err = registryAdmin.AddNode(
		ctx,
		options.RegisterNode.OwnerAddress.Address,
		signingKeyPub,
		options.RegisterNode.HttpAddress,
	)
	if err != nil {
		logger.Fatal("could not add node", zap.Error(err))
	}
}

func isPubKeyAlreadyRegistered(logger *zap.Logger, options *CLI, pubKey string) bool {
	chainClient, err := blockchain.NewRPCClient(
		context.Background(),
		options.Contracts.SettlementChain.RPCURL,
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

	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		logger.Fatal("could not create chain client", zap.Error(err))
	}

	newNodes, err := migrator.ImportNodesFromFile(options.MigrateNodes.InFile)
	if err != nil {
		logger.Fatal("could not import nodes from file", zap.Error(err))
	}

	caller, err := blockchain.NewNodeRegistryCaller(
		logger,
		chainClient,
		options.Contracts,
	)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	oldNodes, err := migrator.ReadFromRegistry(caller)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
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

	err = migrator.WriteToRegistry(ctx, newNodes, oldNodes, newRegistryAdmin)
	if err != nil {
		logger.Fatal("could not write nodes to registry", zap.Error(err))
	}
}

func getRates(logger *zap.Logger, options *CLI) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*15))
	defer cancel()
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
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

	rates := fees.Rates{
		MessageFee:          options.AddRates.MessageFee,
		StorageFee:          options.AddRates.StorageFee,
		CongestionFee:       options.AddRates.CongestionFee,
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
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
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

	limit, err := caller.GetMaxCanonicalNodes(ctx)
	if err != nil {
		logger.Fatal("could not get max canonical nodes", zap.Error(err))
	}
	logger.Info("max canonical nodes retrieved successfully", zap.Any("limit", limit))
}

func getBootstrapperAddress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()

	admin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.GetBootstrapperAddress.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup admin", zap.Error(err))
	}

	iuAddr, err := admin.GetIdentityUpdateBootstrapper(ctx)
	if err != nil {
		logger.Fatal("could not get bootstrapper address", zap.Error(err))
	}

	logger.Info(
		"identity update bootstrapper address retrieved successfully",
		zap.String("IU address", iuAddr.String()),
	)

	gmAddr, err := admin.GetGroupMessageBootstrapper(ctx)
	if err != nil {
		logger.Fatal("could not get bootstrapper address", zap.Error(err))
	}

	logger.Info(
		"group message bootstrapper address retrieved successfully",
		zap.String("GM address", gmAddr.String()),
	)

	if iuAddr.String() != gmAddr.String() {
		logger.Warn(
			"identity update bootstrapper address and group message bootstrapper address do not match",
		)
	}
}

func setBootstrapperAddress(logger *zap.Logger, options *CLI) {
	ctx := context.Background()

	admin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.SetBootstrapperAddress.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup admin", zap.Error(err))
	}

	err = admin.SetIdentityUpdateBootstrapper(
		ctx,
		options.SetBootstrapperAddress.Address.Address,
	)
	if err != nil {
		logger.Fatal("could not set identity update bootstrapper address", zap.Error(err))
	}

	err = admin.SetGroupMessageBootstrapper(
		ctx,
		options.SetBootstrapperAddress.Address.Address,
	)
	if err != nil {
		logger.Fatal("could not set group message bootstrapper address", zap.Error(err))
	}
}

func setPause(logger *zap.Logger, options *CLI) {
	ctx := context.Background()

	appChainAdmin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.SetPause.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}

	settlementChainAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetPause.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	switch options.SetPause.Target {
	case config.TargetIdentity:
		if err := appChainAdmin.SetIdentityUpdatePauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set identity pause status", zap.Error(err))
		}
		logger.Info(
			"identity update pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)

	case config.TargetGroup:
		if err := appChainAdmin.SetGroupMessagePauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set group message pause status", zap.Error(err))
		}
		logger.Info(
			"group message pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)
	case config.TargetAppChainGateway:
		if err := appChainAdmin.SetAppChainGatewayPauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set appchain gateway pause status", zap.Error(err))
		}
		logger.Info(
			"appchain gateway pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)
	case config.TargetSettlementChainGateway:
		if err := settlementChainAdmin.SetSettlementChainGatewayPauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set settlement chain gateway pause status", zap.Error(err))
		}
		logger.Info(
			"settlement chain gateway pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)
	case config.TargetPayerRegistry:
		if err := settlementChainAdmin.SetPayerRegistryPauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set payer registry pause status", zap.Error(err))
		}
		logger.Info(
			"payer registry pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)
	case config.TargetDistributionManager:
		if err := settlementChainAdmin.SetDistributionManagerPauseStatus(ctx, options.SetPause.Paused.Bool()); err != nil {
			logger.Fatal("could not set distribution manager pause status", zap.Error(err))
		}
		logger.Info(
			"distribution manager pause status set",
			zap.Bool("paused", options.SetPause.Paused.Bool()),
		)
	}
}

func getPause(logger *zap.Logger, options *CLI) {
	ctx := context.Background()

	appChainAdmin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.GetPause.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}

	settlementChainAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetPause.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	switch options.GetPause.Target {
	case config.TargetIdentity:
		paused, err := appChainAdmin.GetIdentityUpdatePauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get identity pause status", zap.Error(err))
		}
		logger.Info("identity pause status", zap.Bool("paused", paused))

	case config.TargetGroup:
		paused, err := appChainAdmin.GetGroupMessagePauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get group pause status", zap.Error(err))
		}
		logger.Info("group pause status", zap.Bool("paused", paused))
	case config.TargetAppChainGateway:
		paused, err := appChainAdmin.GetAppChainGatewayPauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get appchain gateway pause status", zap.Error(err))
		}
		logger.Info("app-chain gateway pause status", zap.Bool("paused", paused))
	case config.TargetSettlementChainGateway:
		paused, err := settlementChainAdmin.GetSettlementChainGatewayPauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get settlementchain gateway pause status", zap.Error(err))
		}
		logger.Info("settlement-chain gateway pause status", zap.Bool("paused", paused))
	case config.TargetPayerRegistry:
		paused, err := settlementChainAdmin.GetPayerRegistryPauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get payer registry pause status", zap.Error(err))
		}
		logger.Info("payer registry pause status", zap.Bool("paused", paused))
	case config.TargetDistributionManager:
		paused, err := settlementChainAdmin.GetDistributionManagerPauseStatus(ctx)
		if err != nil {
			logger.Fatal("could not get distribution manager pause status", zap.Error(err))
		}
		logger.Info("distribution manager pause status", zap.Bool("paused", paused))
	}
}

func getAllNodes(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
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
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
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
		zap.String("contract", options.IdentityUpdatesStress.Contract.String()),
	)

	err := stress.StressIdentityUpdates(
		ctx,
		logger,
		options.IdentityUpdatesStress.Count,
		options.IdentityUpdatesStress.Contract.Address,
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
		options.Watcher.Contract.Address,
	)
	if err != nil {
		logger.Fatal("could not create watcher", zap.Error(err))
	}

	err = watcher.Listen(ctxwc)
	if err != nil {
		logger.Fatal("could not listen", zap.Error(err))
	}
}

func getNodeRegistryAdmin(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetNodeRegistryAdmin.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	addr, err := scAdmin.GetNodeRegistryAdmin(ctx)
	if err != nil {
		logger.Fatal("could not get node registry admin", zap.Error(err))
	}
	logger.Info("node registry admin", zap.String("address", addr.Hex()))
}

func setNodeRegistryAdmin(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetNodeRegistryAdmin.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	if err := scAdmin.SetNodeRegistryAdmin(ctx, options.SetNodeRegistryAdmin.Address.Address); err != nil {
		logger.Fatal("could not set node registry admin", zap.Error(err))
	}
	logger.Info(
		"node registry admin set",
		zap.String("address", options.SetNodeRegistryAdmin.Address.String()),
	)
}

func getPayloadSize(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	admin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.GetPayloadSize.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}
	var size uint64
	switch options.GetPayloadSize.Target {
	case config.TargetIdentity:
		if options.GetPayloadSize.Bound == config.PayloadMin {
			size, err = admin.GetIdentityUpdateMinPayloadSize(ctx)
		} else {
			size, err = admin.GetIdentityUpdateMaxPayloadSize(ctx)
		}
	case config.TargetGroup:
		if options.GetPayloadSize.Bound == config.PayloadMin {
			size, err = admin.GetGroupMessageMinPayloadSize(ctx)
		} else {
			size, err = admin.GetGroupMessageMaxPayloadSize(ctx)
		}
	default:
		logger.Fatal(
			"payload size only supports target identity|group",
			zap.String("target", string(options.GetPayloadSize.Target)),
		)
	}
	if err != nil {
		logger.Fatal("could not read payload size", zap.Error(err))
	}
	logger.Info(
		"payload size",
		zap.String("target", string(options.GetPayloadSize.Target)),
		zap.String("bound", string(options.GetPayloadSize.Bound)),
		zap.Uint64("bytes", size),
	)
}

func setPayloadSize(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	admin, err := setupAppChainAdmin(
		ctx,
		logger,
		options.SetPayloadSize.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup appchain admin", zap.Error(err))
	}

	sz := options.SetPayloadSize.Size

	switch options.SetPayloadSize.Target {
	case config.TargetIdentity:
		if options.SetPayloadSize.Bound == config.PayloadMin {
			err = admin.SetIdentityUpdateMinPayloadSize(ctx, sz)
		} else {
			err = admin.SetIdentityUpdateMaxPayloadSize(ctx, sz)
		}
	case config.TargetGroup:
		if options.SetPayloadSize.Bound == config.PayloadMin {
			err = admin.SetGroupMessageMinPayloadSize(ctx, sz)
		} else {
			err = admin.SetGroupMessageMaxPayloadSize(ctx, sz)
		}
	default:
		logger.Fatal(
			"payload size only supports target identity|group",
			zap.String("target", string(options.SetPayloadSize.Target)),
		)
	}
	if err != nil {
		logger.Fatal("could not set payload size", zap.Error(err))
	}
	logger.Info("payload size set",
		zap.String("target", string(options.SetPayloadSize.Target)),
		zap.String("bound", string(options.SetPayloadSize.Bound)),
		zap.Uint64("bytes", sz),
	)
}

func getDMProtocolFeesRecipient(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetDMFeesRecipient.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	addr, err := scAdmin.GetDistributionManagerProtocolFeesRecipient(ctx)
	if err != nil {
		logger.Fatal("could not get distribution manager protocol fees recipient", zap.Error(err))
	}
	logger.Info("distribution manager protocol fees recipient", zap.String("address", addr.Hex()))
}

func setDMProtocolFeesRecipient(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetDMFeesRecipient.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}

	err = scAdmin.SetDistributionManagerProtocolFeesRecipient(
		ctx,
		options.SetDMFeesRecipient.Address.Address,
	)
	if err != nil {
		logger.Fatal("could not set distribution manager protocol fees recipient", zap.Error(err))
	}
	logger.Info(
		"distribution manager protocol fees recipient set",
		zap.String("address", options.SetDMFeesRecipient.Address.String()),
	)
}

func getPayerMinDeposit(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetPayerMinDeposit.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	v, perr := scAdmin.GetPayerRegistryMinimumDeposit(ctx)
	if perr != nil {
		logger.Fatal("could not read minimum deposit", zap.Error(perr))
	}
	logger.Info("payer registry minimum deposit (picodollars)", zap.String("value", v.String()))
}

func setPayerMinDeposit(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetPayerMinDeposit.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	// currency.PicoDollar is likely a uint64/newtype; convert to big.Int for uint96
	bi := new(big.Int).SetUint64(uint64(options.SetPayerMinDeposit.Amount))
	if err := scAdmin.SetPayerRegistryMinimumDeposit(ctx, bi); err != nil {
		logger.Fatal("could not set minimum deposit", zap.Error(err))
	}
	logger.Info(
		"payer registry minimum deposit set (picodollars)",
		zap.String("value", bi.String()),
	)
}

func getPayerWithdrawLock(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetPayerWithdrawLock.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	secs, perr := scAdmin.GetPayerRegistryWithdrawLockPeriod(ctx)
	if perr != nil {
		logger.Fatal("could not read withdraw lock period", zap.Error(perr))
	}
	logger.Info("payer registry withdraw lock period", zap.Uint32("seconds", secs))
}

func setPayerWithdrawLock(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetPayerWithdrawLock.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	if err := scAdmin.SetPayerRegistryWithdrawLockPeriod(ctx, options.SetPayerWithdrawLock.Seconds); err != nil {
		logger.Fatal("could not set withdraw lock period", zap.Error(err))
	}
	logger.Info(
		"payer registry withdraw lock period set",
		zap.Uint32("seconds", options.SetPayerWithdrawLock.Seconds),
	)
}

func getPRMFeeRate(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetPRMFeeRate.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	v, perr := scAdmin.GetPayerReportManagerProtocolFeeRate(ctx)
	if perr != nil {
		logger.Fatal("could not read PRM protocol fee rate", zap.Error(perr))
	}
	logger.Info("payer report manager protocol fee rate (bps)", zap.Uint16("bps", v))
}

func setPRMFeeRate(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.SetPRMFeeRate.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	if err := scAdmin.SetPayerReportManagerProtocolFeeRate(ctx, options.SetPRMFeeRate.FeeRateBps); err != nil {
		logger.Fatal("could not set PRM protocol fee rate", zap.Error(err))
	}
	logger.Info(
		"payer report manager protocol fee rate set",
		zap.Uint16("bps", options.SetPRMFeeRate.FeeRateBps),
	)
}

func getRateMigrator(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetRateMigrator.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	addr, perr := scAdmin.GetRateRegistryMigrator(ctx)
	if perr != nil {
		logger.Fatal("could not read rate registry migrator", zap.Error(perr))
	}
	logger.Info("rate registry migrator", zap.String("address", addr.Hex()))
}

func setRateMigrator(logger *zap.Logger, options *CLI) {
	ctx := context.Background()
	scAdmin, err := setupSettlementChainAdmin(
		ctx,
		logger,
		options.GetRateMigrator.AdminOptions.AdminPrivateKey,
		options,
	)
	if err != nil {
		logger.Fatal("could not setup settlement chain admin", zap.Error(err))
	}
	if err := scAdmin.SetRateRegistryMigrator(ctx, options.SetRateMigrator.Address.Address); err != nil {
		logger.Fatal("could not set rate registry migrator", zap.Error(err))
	}
	logger.Info(
		"rate registry migrator set",
		zap.String("address", options.SetRateMigrator.Address.String()),
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
	case "get-bootstrapper-address":
		getBootstrapperAddress(logger, options)
		return
	case "set-bootstrapper-address":
		setBootstrapperAddress(logger, options)
		return
	case "set-pause":
		setPause(logger, options)
		return
	case "get-pause":
		getPause(logger, options)
		return
	case "get-node-admin":
		getNodeRegistryAdmin(logger, options)
		return
	case "set-node-admin":
		setNodeRegistryAdmin(logger, options)
		return
	case "get-payload-size":
		getPayloadSize(logger, options)
		return
	case "set-payload-size":
		setPayloadSize(logger, options)
		return

	case "get-dm-protocol-fees-recipient":
		getDMProtocolFeesRecipient(logger, options)
		return
	case "set-dm-protocol-fees-recipient":
		setDMProtocolFeesRecipient(logger, options)
		return

	case "get-payer-min-deposit":
		getPayerMinDeposit(logger, options)
		return
	case "set-payer-min-deposit":
		setPayerMinDeposit(logger, options)
		return

	case "get-payer-withdraw-lock":
		getPayerWithdrawLock(logger, options)
		return
	case "set-payer-withdraw-lock":
		setPayerWithdrawLock(logger, options)
		return
	case "get-prm-fee-rate":
		getPRMFeeRate(logger, options)
		return
	case "set-prm-fee-rate":
		setPRMFeeRate(logger, options)
		return

	case "get-rate-migrator":
		getRateMigrator(logger, options)
		return
	case "set-rate-migrator":
		setRateMigrator(logger, options)
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
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		options.Contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	parameterAdmin, err := setupParameterAdmin(ctx, logger, privateKey, options)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	registryAdmin, err := blockchain.NewNodeRegistryAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
		parameterAdmin,
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
	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		chainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	paramAdmin, err := blockchain.NewParameterAdmin(logger, chainClient, signer, options.Contracts)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	registryAdmin, err := blockchain.NewRatesAdmin(
		logger,
		paramAdmin,
		chainClient,
		options.Contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create registry admin: %w", err)
	}

	return registryAdmin, nil
}

func setupParameterAdmin(ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	options *CLI,
) (*blockchain.ParameterAdmin, error) {
	if options.Contracts.SettlementChain.RPCURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}
	if options.Contracts.SettlementChain.ChainID == 0 {
		return nil, fmt.Errorf("chain id is required")
	}
	if options.Contracts.SettlementChain.ParameterRegistryAddress == "" {
		return nil, fmt.Errorf("parameter registry address is required")
	}

	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		options.Contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	parameterAdmin, err := blockchain.NewParameterAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	return parameterAdmin, nil
}

// setupAppChainAdmin creates and returns a appchain admin
func setupAppChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	options *CLI,
) (blockchain.IAppChainAdmin, error) {
	if options.Contracts.AppChain.RPCURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}
	if options.Contracts.AppChain.ChainID == 0 {
		return nil, fmt.Errorf("chain id is required")
	}
	if options.Contracts.AppChain.GroupMessageBroadcasterAddress == "" {
		return nil, fmt.Errorf("group message broadcaster address is required")
	}
	if options.Contracts.AppChain.IdentityUpdateBroadcasterAddress == "" {
		return nil, fmt.Errorf("identity update broadcaster address is required")
	}

	// TODO(mkysel) https://github.com/xmtp/smart-contracts/issues/125
	//if options.Contracts.AppChain.GatewayAddress == "" {
	//	return nil, fmt.Errorf("gateway address is required")
	//}

	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.AppChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		options.Contracts.AppChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	parameterAdmin, err := setupParameterAdmin(ctx, logger, privateKey, options)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	appchainAdmin, err := blockchain.NewAppChainAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
		parameterAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create appchain admin: %w", err)
	}

	return appchainAdmin, nil
}

// setupSettlementChainAdmin creates and returns a settlementchain admin
func setupSettlementChainAdmin(
	ctx context.Context,
	logger *zap.Logger,
	privateKey string,
	options *CLI,
) (blockchain.ISettlementChainAdmin, error) {
	if options.Contracts.SettlementChain.RPCURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}
	if options.Contracts.SettlementChain.ChainID == 0 {
		return nil, fmt.Errorf("chain id is required")
	}

	if options.Contracts.SettlementChain.PayerRegistryAddress == "" {
		return nil, fmt.Errorf("payer registry address is required")
	}
	if options.Contracts.SettlementChain.DistributionManagerAddress == "" {
		return nil, fmt.Errorf("distribution manager address is required")
	}

	// TODO(mkysel) https://github.com/xmtp/smart-contracts/issues/125
	//if options.Contracts.SettlementChain.GatewayAddress == "" {
	//	return nil, fmt.Errorf("gateway address is required")
	//}

	chainClient, err := blockchain.NewRPCClient(
		ctx,
		options.Contracts.SettlementChain.RPCURL,
	)
	if err != nil {
		return nil, err
	}

	signer, err := blockchain.NewPrivateKeySigner(
		privateKey,
		options.Contracts.SettlementChain.ChainID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create signer: %w", err)
	}

	parameterAdmin, err := setupParameterAdmin(ctx, logger, privateKey, options)
	if err != nil {
		return nil, fmt.Errorf("could not create parameter admin: %w", err)
	}

	settlementChainAdmin, err := blockchain.NewSettlementChainAdmin(
		logger,
		chainClient,
		signer,
		options.Contracts,
		parameterAdmin,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create settlementchain admin: %w", err)
	}

	return settlementChainAdmin, nil
}
