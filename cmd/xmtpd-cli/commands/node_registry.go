package commands

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"

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
			if err := migrator.DumpNodesToFile(
				[]migrator.SerializableNode{exportNode},
				exportPath,
			); err != nil {
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
		logger.Info("checking node", zap.String("address", address))

		/* Connect-RPC */

		connectRPCHealthy, err := checkConnectRPC(ctx, address)
		if err != nil || !connectRPCHealthy {
			logger.Error("❌ Connect-RPC should not be exposed")
		} else {
			logger.Info("✅ Connect-RPC is not exposed")
		}

		/* gRPC */

		grpcHealthy, err := checkGRPC(ctx, address)
		if err != nil || !grpcHealthy {
			logger.Error("❌ gRPC should be supported")
		} else {
			logger.Info("✅ gRPC is supported")
		}

		/* gRPC-Web */

		grpcWebHealthy, err := checkGRPCWeb(ctx, address)
		if err != nil || !grpcWebHealthy {
			logger.Error("❌ gRPC-Web should be supported")
		} else {
			logger.Info("✅ gRPC-Web is supported")
		}

		/* CORS: OPTIONS request should be supported */

		corsHealthy, err := checkCORS(address)
		if err != nil || !corsHealthy {
			logger.Error("❌ CORS: OPTIONS request should be supported")
		} else {
			logger.Info("✅ CORS: OPTIONS requests are supported")
		}

		if connectRPCHealthy && grpcHealthy && grpcWebHealthy && corsHealthy {
			logger.Info("✅ node is healthy")
		} else {
			logger.Error("❌ node is unhealthy")
		}
	}

	return nil
}

func checkConnectRPC(ctx context.Context, address string) (bool, error) {
	connectRPCClient, err := utils.NewConnectMetadataAPIClient(
		ctx,
		address,
	)
	if err != nil {
		return false, fmt.Errorf("could not create metadata api client for address: %w", err)
	}

	_, err = connectRPCClient.GetVersion(
		ctx,
		&connect.Request[metadata_api.GetVersionRequest]{
			Msg: &metadata_api.GetVersionRequest{},
		},
	)
	if err == nil {
		return false, fmt.Errorf("Connect-RPC is exposed")
	}

	return true, nil
}

func checkGRPC(ctx context.Context, address string) (bool, error) {
	gRPCClient, err := utils.NewConnectMetadataAPIClient(
		ctx,
		address,
		utils.BuildGRPCDialOptions()...,
	)
	if err != nil {
		return false, fmt.Errorf("could not create metadata api client for address: %w", err)
	}

	_, err = gRPCClient.GetVersion(
		ctx,
		&connect.Request[metadata_api.GetVersionRequest]{
			Msg: &metadata_api.GetVersionRequest{},
		},
	)
	if err != nil {
		return false, fmt.Errorf("should accept gRPC requests: %w", err)
	}

	return true, nil
}

func checkGRPCWeb(ctx context.Context, address string) (bool, error) {
	gRPCWebClient, err := utils.NewConnectMetadataAPIClient(
		ctx,
		address,
		utils.BuildGRPCWebDialOptions()...,
	)
	if err != nil {
		return false, fmt.Errorf("could not create metadata api client for address: %w", err)
	}

	_, err = gRPCWebClient.GetVersion(
		ctx,
		&connect.Request[metadata_api.GetVersionRequest]{
			Msg: &metadata_api.GetVersionRequest{},
		},
	)
	if err != nil {
		return false, fmt.Errorf("should accept gRPC-Web requests: %w", err)
	}

	return true, nil
}

func checkCORS(address string) (bool, error) {
	url := address + "/xmtp.xmtpv4.metadata_api.MetadataApi/GetVersion"

	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return false, fmt.Errorf("could not create request: %w", err)
	}

	// Match your successful curl (browser-like preflight)
	origin := "https://network.xmtp.org"
	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "content-type,x-grpc-web,accept,authorization")

	// Force HTTP/2 by using an http2.Transport directly.
	// (This will only work for https:// URLs.)
	client := &http.Client{
		Transport: &http2.Transport{},
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("could not make http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.ProtoMajor != 2 {
		return false, fmt.Errorf("did not use HTTP/2 (got %s)", resp.Proto)
	}

	// Many proxies use 204 for preflight; allow 200 too just in case.
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf(
			"preflight failed: status=%d proto=%s",
			resp.StatusCode,
			resp.Proto,
		)
	}

	// Validate key CORS headers
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOrigin != origin && allowOrigin != "*" {
		return false, fmt.Errorf("unexpected Access-Control-Allow-Origin: %q", allowOrigin)
	}

	allowMethods := resp.Header.Get("Access-Control-Allow-Methods")
	if !strings.Contains(strings.ToUpper(allowMethods), "POST") {
		return false, fmt.Errorf(
			"POST not allowed by Access-Control-Allow-Methods: %q",
			allowMethods,
		)
	}

	allowHeaders := strings.ToLower(resp.Header.Get("Access-Control-Allow-Headers"))
	for _, h := range []string{"content-type", "x-grpc-web", "accept", "authorization"} {
		if !strings.Contains(allowHeaders, h) {
			return false, fmt.Errorf(
				"missing allowed header %q in Access-Control-Allow-Headers: %q",
				h,
				allowHeaders,
			)
		}
	}

	return true, nil
}
