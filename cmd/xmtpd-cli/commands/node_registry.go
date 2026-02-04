package commands

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
	"github.com/xmtp/xmtpd/pkg/blockchain/migrator"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

func nodeRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "nodes",
		Short:        "Manage Node Registry",
		SilenceUsage: true,
	}

	// Shared flag (read by some subcommands)
	cmd.PersistentFlags().Uint32P("node-id", "n", 0, "Node ID to use")

	cmd.AddCommand(
		registerNodeCmd(),
		canonicalNetworkCmd(),
		getNodeCmd(),
		maxCanonicalCmd(),
		setHTTPAddressCmd(),
		healthCheckCmd(),
	)

	return cmd
}

//------------------------------------------------------------------------------
// register
//------------------------------------------------------------------------------

func registerNodeCmd() *cobra.Command {
	var owner options.AddressFlag
	var signingKeyPub string
	var httpAddress string
	var force bool

	cmd := &cobra.Command{
		Use:          "register",
		Short:        "Register a node",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes register --owner-address <address> --signing-key-pub <key> --http-address <address> [--force]

Register a node:
xmtpd-cli nodes register --owner-address <address> --signing-key-pub <key> --http-address <address>
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return registerNodeHandler(signingKeyPub, httpAddress, owner.Address, force)
		},
	}

	cmd.Flags().Var(&owner, "owner-address", "Owner address to use")
	_ = cmd.MarkFlagRequired("owner-address")

	cmd.Flags().StringVar(&signingKeyPub, "signing-key-pub", "", "signing key public key to use")
	_ = cmd.MarkFlagRequired("signing-key-pub")

	cmd.Flags().StringVar(&httpAddress, "http-address", "", "HTTP address to use")
	_ = cmd.MarkFlagRequired("http-address")

	cmd.Flags().BoolVar(&force, "force", false, "force the registration")

	return cmd
}

func registerNodeHandler(
	signingKeyPub, httpAddress string,
	owner common.Address,
	force bool,
) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not create registry caller: %w", err)
	}

	nodes, err := migrator.ReadFromRegistry(caller)
	if err != nil {
		return fmt.Errorf("could not retrieve nodes from registry: %w", err)
	}

	if !force {
		for _, node := range nodes {
			if node.SigningKeyPub == signingKeyPub {
				logger.Warn(
					"signing key public key already registered",
					zap.String("signing-key-pub", signingKeyPub),
				)
				return nil
			}
		}
	}

	admin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not create registry admin: %w", err)
	}

	parsedSigningKeyPub, err := utils.ParseEcdsaPublicKey(signingKeyPub)
	if err != nil {
		logger.Error(
			"could not decompress public key",
			zap.Error(err),
			zap.String("key", signingKeyPub),
		)
		return fmt.Errorf("could not decompress public key: %w", err)
	}

	nodeID, err := admin.AddNode(ctx, owner, parsedSigningKeyPub, httpAddress)
	if err != nil {
		return fmt.Errorf("could not add node: %w", err)
	}

	logger.Info(
		"node registered",
		zap.String("owner-address", owner.Hex()),
		zap.Uint32("node-id", nodeID),
	)
	return nil
}

//------------------------------------------------------------------------------
// canonical-network
//------------------------------------------------------------------------------

func canonicalNetworkCmd() *cobra.Command {
	var add bool
	var remove bool

	cmd := &cobra.Command{
		Use:          "canonical-network",
		Short:        "Manage the canonical network",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes canonical-network {--add | --remove} --node-id <node-id>

Add a node to the canonical network:
xmtpd-cli nodes canonical-network --add --node-id <node-id>

Remove a node from the canonical network:
xmtpd-cli nodes canonical-network --remove --node-id <node-id>
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			nodeID, err := cmd.Flags().GetUint32("node-id")
			if err != nil {
				return fmt.Errorf("could not get node id: %w", err)
			}
			return canonicalNetworkHandler(add, remove, nodeID)
		},
	}

	cmd.Flags().BoolVar(&add, "add", false, "add a node to the canonical network")
	cmd.Flags().BoolVar(&remove, "remove", false, "remove a node from the canonical network")
	cmd.MarkFlagsMutuallyExclusive("add", "remove")

	return cmd
}

func canonicalNetworkHandler(add, remove bool, nodeID uint32) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	if nodeID == 0 {
		return fmt.Errorf("node id is required")
	}
	if !add && !remove {
		return fmt.Errorf("either --add or --remove must be specified")
	}

	ctx := context.Background()
	admin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to create registry admin: %w", err)
	}

	if add {
		if err := admin.AddToNetwork(ctx, nodeID); err != nil {
			return fmt.Errorf("failed to add node to network: %w", err)
		}
		logger.Info("added node to canonical network", zap.Uint32("node-id", nodeID))
	}

	if remove {
		if err := admin.RemoveFromNetwork(ctx, nodeID); err != nil {
			return fmt.Errorf("failed to remove node from network: %w", err)
		}
		logger.Info("removed node from canonical network", zap.Uint32("node-id", nodeID))
	}

	return nil
}

//------------------------------------------------------------------------------
// get
//------------------------------------------------------------------------------

func getNodeCmd() *cobra.Command {
	var all bool
	var exportPath string

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get and export nodes",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes get {--all | --node-id <node-id>} [--export <file>]

Get all nodes:
xmtpd-cli nodes get --all

Get a specific node:
xmtpd-cli nodes get --node-id <node-id>

Export all nodes to file:
xmtpd-cli nodes get --all --export <file>
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			nodeID, err := cmd.Flags().GetUint32("node-id")
			if err != nil {
				return fmt.Errorf("could not get node id: %w", err)
			}
			return getNodeHandler(all, nodeID, exportPath)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "get all nodes")
	cmd.Flags().StringVar(&exportPath, "export", "", "export the result to file")

	return cmd
}

func getNodeHandler(all bool, nodeID uint32, exportPath string) error {
	ctx := context.Background()

	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not create registry caller: %w", err)
	}

	nodes, err := migrator.ReadFromRegistry(caller)
	if err != nil {
		return fmt.Errorf("could not retrieve nodes from registry: %w", err)
	}

	if nodeID == 0 && !all {
		return fmt.Errorf("either --node-id or --all must be specified")
	}

	switch {
	case all:
		logger.Info("getting all nodes", zap.Any("nodes", nodes))
		if exportPath != "" {
			if err := migrator.DumpNodesToFile(nodes, exportPath); err != nil {
				return fmt.Errorf("could not dump nodes: %w", err)
			}
		}
	default:
		var (
			found      bool
			exportNode migrator.SerializableNode
		)

		for _, node := range nodes {
			if node.NodeID == nodeID {
				logger.Info("got node", zap.Any("node", node))
				found = true
				exportNode = node
				break
			}
		}
		if !found {
			return fmt.Errorf("node not found: %d", nodeID)
		}
		if exportPath != "" {
			if err := migrator.DumpNodesToFile([]migrator.SerializableNode{exportNode}, exportPath); err != nil {
				return fmt.Errorf("could not dump nodes: %w", err)
			}
		}
	}

	return nil
}

//------------------------------------------------------------------------------
// max-canonical
//------------------------------------------------------------------------------

func maxCanonicalCmd() *cobra.Command {
	var setVal uint8

	cmd := &cobra.Command{
		Use:          "max-canonical",
		Short:        "Manage the maximum canonical size",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes max-canonical [--set <size>]

Set the maximum canonical size:
xmtpd-cli nodes max-canonical --set <size>

Get the current maximum canonical size:
xmtpd-cli nodes max-canonical
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return maxCanonicalHandler(setVal)
		},
	}

	cmd.Flags().Uint8Var(&setVal, "set", 0, "set the maximum canonical size")
	return cmd
}

func maxCanonicalHandler(setVal uint8) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()

	if setVal > 0 {
		admin, err := setupNodeRegistryAdmin(ctx, logger)
		if err != nil {
			return fmt.Errorf("failed to create registry admin: %w", err)
		}
		if err := admin.SetMaxCanonical(ctx, setVal); err != nil {
			return fmt.Errorf("failed to set max canonical size: %w", err)
		}
		logger.Info("set new max canonical size", zap.Uint8("max_canonical_nodes", setVal))
	}

	caller, err := setupNodeRegistryCaller(ctx, logger)
	if err != nil {
		return fmt.Errorf("failed to create registry caller: %w", err)
	}

	val, err := caller.GetMaxCanonicalNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get max canonical size: %w", err)
	}

	logger.Info("current max canonical size", zap.Uint8("max_canonical_nodes", val))
	return nil
}

//------------------------------------------------------------------------------
// set-http-address
//------------------------------------------------------------------------------

func setHTTPAddressCmd() *cobra.Command {
	var httpAddress string

	cmd := &cobra.Command{
		Use:          "set-http-address",
		Short:        "Set the HTTP address of a node",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes set-http-address --node-id <node-id> --http-address <address>

Set the HTTP address of a node:
xmtpd-cli nodes set-http-address --node-id <node-id> --http-address <address>
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			nodeID, err := cmd.Flags().GetUint32("node-id")
			if err != nil {
				return fmt.Errorf("could not get node id: %w", err)
			}
			return setHTTPAddressHandler(nodeID, httpAddress)
		},
	}

	cmd.Flags().StringVar(&httpAddress, "http-address", "", "HTTP address to use")
	_ = cmd.MarkFlagRequired("http-address")

	return cmd
}

func setHTTPAddressHandler(nodeID uint32, httpAddress string) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	if nodeID == 0 || httpAddress == "" {
		return fmt.Errorf("node id and http address are required")
	}

	ctx := context.Background()

	registryAdmin, err := setupNodeRegistryAdmin(ctx, logger)
	if err != nil {
		return fmt.Errorf("could not create registry admin: %w", err)
	}

	if err := registryAdmin.SetHTTPAddress(ctx, nodeID, httpAddress); err != nil {
		return fmt.Errorf("could not set http address: %w", err)
	}

	logger.Info("set the HTTP address of a node",
		zap.Uint32("node-id", nodeID),
		zap.String("http-address", httpAddress),
	)
	return nil
}

func healthCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "health-check",
		Short:        "Check the health of a node",
		SilenceUsage: true,
		Example: `
Usage: xmtpd-cli nodes health-check [node url]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("at least one node url is required")
			}

			return healthCheckHandler(args)
		},
	}

	return cmd
}

func healthCheckHandler(addresses []string) error {
	logger, err := cliLogger()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}

	ctx := context.Background()

	for _, address := range addresses {
		logger.Info("checking health of node", zap.String("address", address))

		healthy := true

		/* Connect-RPC */

		logger.Info("checking Connect-RPC is not exposed", zap.String("address", address))

		connectRPCClient, err := utils.NewConnectMetadataAPIClient(
			ctx,
			address,
		)
		if err != nil {
			logger.Error(
				"could not create metadata api client for address",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		_, err = connectRPCClient.GetVersion(
			ctx,
			&connect.Request[metadata_api.GetVersionRequest]{
				Msg: &metadata_api.GetVersionRequest{},
			},
		)
		if err == nil {
			logger.Error(
				"❌ Connect-RPC is exposed",
				zap.String("address", address),
			)
			healthy = false
		}

		/* gRPC */

		logger.Info("checking gRPC is supported", zap.String("address", address))

		gRPCClient, err := utils.NewConnectMetadataAPIClient(
			ctx,
			address,
			utils.BuildGRPCDialOptions()...,
		)
		if err != nil {
			logger.Error(
				"could not create metadata api client for address",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		_, err = gRPCClient.GetVersion(
			ctx,
			&connect.Request[metadata_api.GetVersionRequest]{
				Msg: &metadata_api.GetVersionRequest{},
			},
		)
		if err != nil {
			logger.Error(
				"❌ should accept gRPC requests",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		/* gRPC-Web */

		logger.Info("checking gRPC-Web is supported", zap.String("address", address))

		gRPCWebClient, err := utils.NewConnectMetadataAPIClient(
			ctx,
			address,
			utils.BuildGRPCWebDialOptions()...,
		)
		if err != nil {
			logger.Error(
				"could not create metadata api client for address",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		_, err = gRPCWebClient.GetVersion(
			ctx,
			&connect.Request[metadata_api.GetVersionRequest]{
				Msg: &metadata_api.GetVersionRequest{},
			},
		)
		if err != nil {
			logger.Error(
				"❌ should accept gRPC-Web requests",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		/* CORS: OPTIONS request should be supported */

		logger.Info("checking CORS: OPTIONS request is supported", zap.String("address", address))

		req, err := http.NewRequest(
			http.MethodOptions,
			address+"/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion",
			bytes.NewReader([]byte{0, 0, 0, 0, 0}),
		)
		if err != nil {
			logger.Error(
				"could not create request",
				zap.String("address", address),
				zap.Error(err),
			)
			continue
		}

		// These headers are what browsers send in preflight requests.
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Error(
				"could not make http request",
				zap.String("address", address),
				zap.Error(err),
			)
			healthy = false
		}

		// OPTIONS request should return a 204 No Content response.
		if resp.StatusCode != http.StatusNoContent {
			logger.Error(
				"❌ OPTIONS requests not supported",
				zap.String("address", address),
				zap.Int("status-code", resp.StatusCode),
			)
			healthy = false
		}

		_ = resp.Body.Close()

		if healthy {
			logger.Info("✅ node is healthy", zap.String("address", address))
		} else {
			logger.Error("❌ node is unhealthy", zap.String("address", address))
		}
	}

	return nil
}
