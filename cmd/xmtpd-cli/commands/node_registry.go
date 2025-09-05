package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/blockchain/migrator"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

func nodeRegistryCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "nodes",
		Short: "Manage Node Registry",
	}

	cmd.PersistentFlags().
		Uint32P("node-id", "n", 0, "Node ID to use")

	cmd.AddCommand(
		registerNodeCmd(),
		canonicalNetworkCmd(),
		getNodeCmd(),
		maxCanonicalCmd(),
		setHttpAddressCmd(),
	)

	return &cmd
}

func registerNodeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "register",
		Short: "Register a node",
		Run:   registerNodeHandler,
		Example: `
Usage: xmtpd-cli nodes register --owner-address <address> --signing-key-pub <key> --http-address <address> [--force]

Register a node:
xmtpd-cli nodes register --owner-address <address> --signing-key-pub <key> --http-address <address>
`,
	}

	cmd.PersistentFlags().
		String("owner-address", "", "owner address to use")

	cmd.PersistentFlags().
		String("signing-key-pub", "", "signing key public key to use")

	cmd.PersistentFlags().
		String("http-address", "", "HTTP address to use")

	_ = cmd.MarkFlagRequired("owner-address")
	_ = cmd.MarkFlagRequired("signing-key-pub")
	_ = cmd.MarkFlagRequired("http-address")

	cmd.PersistentFlags().
		Bool("force", false, "force the registration")

	return &cmd
}

func registerNodeHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	ctx := context.Background()

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		logger.Fatal("could not create registry caller", zap.Error(err))
	}

	nodes, err := migrator.ReadFromRegistry(caller)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	owner, _ := cmd.Flags().GetString("owner-address")
	signingKeyPub, _ := cmd.Flags().GetString("signing-key-pub")
	httpAddress, _ := cmd.Flags().GetString("http-address")
	force, _ := cmd.Flags().GetBool("force")

	ownerAddress := common.HexToAddress(owner)

	if !force {
		for _, node := range nodes {
			if node.SigningKeyPub == signingKeyPub {
				logger.Fatal(
					"signing key public key already registered",
					zap.String("signing-key-pub", signingKeyPub),
				)
			}
		}
	}

	admin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	parsedSigningKeyPub, err := utils.ParseEcdsaPublicKey(signingKeyPub)
	if err != nil {
		logger.Fatal(
			"could not decompress public key",
			zap.Error(err),
			zap.String("key", signingKeyPub),
		)
	}

	nodeID, err := admin.AddNode(ctx, ownerAddress, parsedSigningKeyPub, httpAddress)
	if err != nil {
		logger.Fatal("could not add node", zap.Error(err))
	}

	logger.Info(
		"Node registered",
		zap.String("owner-address", owner),
		zap.Uint32("node-id", nodeID),
	)
}

func canonicalNetworkCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "canonical-network",
		Short: "Manage the canonical network",
		Run:   canonicalNetworkHandler,
		Example: `
Usage: xmtpd-cli nodes canonical-network {--add | --remove} --node-id <node-id>

Add a node to the canonical network:
xmtpd-cli nodes canonical-network --add --node-id <node-id>

Remove a node from the canonical network:
xmtpd-cli nodes canonical-network --remove --node-id <node-id>
`,
	}

	cmd.PersistentFlags().
		Bool("add", false, "add a node to the canonical network")

	cmd.PersistentFlags().
		Bool("remove", false, "remove a node from the canonical network")

	cmd.MarkFlagsMutuallyExclusive("add", "remove")

	return &cmd
}

func canonicalNetworkHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	nodeID, err := cmd.Flags().GetUint32("node-id")
	if err != nil {
		logger.Fatal("could not get node id", zap.Error(err))
	}

	if nodeID == 0 {
		logger.Fatal("node id is required")
	}

	add, err := cmd.Flags().GetBool("add")
	if err != nil {
		logger.Fatal("could not get add flag", zap.Error(err))
	}

	remove, err := cmd.Flags().GetBool("remove")
	if err != nil {
		logger.Fatal("could not get remove flag", zap.Error(err))
	}

	if !add && !remove {
		logger.Fatal("either --add or --remove must be specified")
	}

	ctx := context.Background()

	admin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("failed to create registry admin", zap.Error(err))
	}

	if add {
		err = admin.AddToNetwork(ctx, nodeID)
		if err != nil {
			logger.Fatal("failed to add node to network", zap.Error(err))
		}
	}

	if remove {
		err = admin.RemoveFromNetwork(ctx, nodeID)
		if err != nil {
			logger.Fatal("failed to remove node from network", zap.Error(err))
		}
	}
}

func getNodeCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get and export nodes",
		Run:   getNodeHandler,
		Example: `
Usage: xmtpd-cli nodes get {--all | --node-id <node-id>} [--export <file>]

Get all nodes:
xmtpd-cli nodes get --all

Get a specific node:
xmtpd-cli nodes get --node-id <node-id>

Export all nodesto file:
xmtpd-cli nodes get --all --export <file>
`,
	}

	cmd.PersistentFlags().
		Bool("all", false, "get all nodes")

	cmd.PersistentFlags().
		String("export", "", "export the result to file")

	return &cmd
}

func getNodeHandler(cmd *cobra.Command, _ []string) {
	ctx := context.Background()

	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		logger.Fatal("could not create registry caller", zap.Error(err))
	}

	nodes, err := migrator.ReadFromRegistry(caller)
	if err != nil {
		logger.Fatal("could not retrieve nodes from registry", zap.Error(err))
	}

	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		logger.Fatal("could not get all flag", zap.Error(err))
	}

	nodeID, err := cmd.Flags().GetUint32("node-id")
	if err != nil {
		logger.Fatal("could not get node id", zap.Error(err))
	}

	if nodeID == 0 && !all {
		logger.Fatal("either --node-id or --all must be specified")
	}

	export, err := cmd.Flags().GetString("export")
	if err != nil {
		logger.Fatal("could not get export flag", zap.Error(err))
	}

	switch {
	case all:
		logger.Info("Getting all nodes", zap.Any("nodes", nodes))

		if export != "" {
			err = migrator.DumpNodesToFile(nodes, export)
			if err != nil {
				logger.Fatal("could not dump nodes", zap.Error(err))
			}
		}

	case nodeID != 0:
		var (
			found      bool
			exportNode migrator.SerializableNode
		)

		for _, node := range nodes {
			if node.NodeID == nodeID {
				logger.Info("Got node", zap.Any("node", node))
				found = true
				exportNode = node
			}
		}

		if !found {
			logger.Fatal("node not found", zap.Uint32("node-id", nodeID))
		}

		if export != "" {
			err = migrator.DumpNodesToFile([]migrator.SerializableNode{exportNode}, export)
			if err != nil {
				logger.Fatal("could not dump nodes", zap.Error(err))
			}
		}
	}
}

func maxCanonicalCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "max-canonical",
		Short: "Manage the maximum canonical size",
		Run:   maxCanonicalHandler,
		Example: `
Usage: xmtpd-cli nodes max-canonical [--set <size>]

Set the maximum canonical size:
xmtpd-cli nodes max-canonical --set <size>

Get the current maximum canonical size:
xmtpd-cli nodes max-canonical
`,
	}

	cmd.PersistentFlags().
		Uint8("set", 0, "set the maximum canonical size")

	return &cmd
}

func maxCanonicalHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	ctx := context.Background()

	setVal, err := cmd.Flags().GetUint8("set")
	if err != nil {
		logger.Fatal("could not parse --set flag", zap.Error(err))
	}

	if setVal > 0 {
		admin, err := setupNodeRegistryAdmin(ctx, logger)
		if err != nil {
			logger.Fatal("failed to create registry admin", zap.Error(err))
		}

		err = admin.SetMaxCanonical(ctx, uint8(setVal))
		if err != nil {
			logger.Fatal("failed to set max canonical size", zap.Error(err))
		}

		logger.Info("Set new max canonical size", zap.Uint8("maxCanonicalNodes", setVal))
	}

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		logger.Fatal("failed to create registry caller", zap.Error(err))
	}

	val, err := caller.GetMaxCanonicalNodes(ctx)
	if err != nil {
		logger.Fatal("failed to get max canonical size", zap.Error(err))
	}

	logger.Info("Current max canonical size", zap.Uint8("maxCanonicalNodes", val))
}

func setHttpAddressCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "set-http-address",
		Short: "Set the HTTP address of a node",
		Run:   setHttpAddressHandler,
		Example: `
Usage: xmtpd-cli nodes set-http-address --node-id <node-id> --http-address <address>

Set the HTTP address of a node:
xmtpd-cli nodes set-http-address --node-id <node-id> --http-address <address>
`,
	}

	cmd.PersistentFlags().
		String("http-address", "", "HTTP address to use")

	return &cmd
}

func setHttpAddressHandler(cmd *cobra.Command, _ []string) {
	logger, err := cliLogger()
	if err != nil {
		log.Fatalf("could not build logger: %s", err)
	}

	nodeID, err := cmd.Flags().GetUint32("node-id")
	if err != nil {
		logger.Fatal("could not get node id", zap.Error(err))
	}

	httpAddress, err := cmd.Flags().GetString("http-address")
	if err != nil {
		logger.Fatal("could not get http address", zap.Error(err))
	}

	if nodeID == 0 || httpAddress == "" {
		logger.Fatal("node id and http address are required")
	}

	ctx := context.Background()

	registryAdmin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		logger.Fatal("could not create registry admin", zap.Error(err))
	}

	err = registryAdmin.SetHttpAddress(
		ctx,
		nodeID,
		httpAddress,
	)
	if err != nil {
		logger.Fatal("could not set http address", zap.Error(err))
	}

	logger.Info("Set the HTTP address of a node",
		zap.Uint32("node-id", nodeID),
		zap.String("http-address", httpAddress),
	)
}
