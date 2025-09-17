package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/stress"
	"go.uber.org/zap"
)

func generateCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "generate",
		Short: "Generate and publish payloads",
	}
	cmd.AddCommand(
		welcomeMessageCmd(),
		groupMessageCmd(),
		keyPackageCmd(),
	)
	return &cmd
}

func welcomeMessageCmd() *cobra.Command {
	var (
		numEnvelopes     uint
		dataSize         uint
		grpcAddress      string
		privateKeyString string
		originatorID     uint32
	)

	cmd := &cobra.Command{
		Use:   "welcome-message",
		Short: "Publish welcome message envelopes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return welcomeMessageHandler(
				grpcAddress,
				privateKeyString,
				originatorID,
				numEnvelopes,
				dataSize,
			)
		},
		Example: `
Usage: xmtpd-cli generate welcome-message --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num> --data-size <size>

Example:
xmtpd-cli generate welcome-message --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num> --data-size <size>
`,
	}

	cmd.Flags().StringVarP(&grpcAddress, "grpc-address", "a", "", "xmtpd node gRPC address")
	_ = cmd.MarkFlagRequired("grpc-address")

	cmd.Flags().StringVar(&privateKeyString, "private-key", "", "payer private key")
	_ = cmd.MarkFlagRequired("private-key")

	cmd.Flags().Uint32VarP(&originatorID, "originator-id", "o", 100, "xmtpd node originator ID")

	// Note that num-envelopes doesn't currently work, until we implement batch publishing.
	cmd.Flags().UintVar(&numEnvelopes, "num-envelopes", 1, "number of envelopes to generate")

	cmd.Flags().UintVar(&dataSize, "data-size", 256, "data size in bytes")

	return cmd
}

func welcomeMessageHandler(
	nodeHTTPAddress string,
	privateKeyString string,
	originatorID uint32,
	numEnvelopes uint,
	dataSize uint,
) error {
	logger, err := cliLogger()
	if err != nil {
		return err
	}

	ctx := context.Background()

	generator, err := stress.NewEnvelopesGenerator(nodeHTTPAddress, privateKeyString, originatorID)
	if err != nil {
		return err
	}

	defer func() {
		err := generator.Close()
		if err != nil {
			logger.Error("error closing generator", zap.Error(err))
		}
	}()

	var envelopes []*envelopes.OriginatorEnvelope

	envelopes, err = generator.PublishWelcomeMessageEnvelopes(ctx, numEnvelopes, dataSize)
	if err != nil {
		logger.Error("error publishing welcome message", zap.Error(err))
		return err
	}

	for _, envelope := range envelopes {
		logger.Info(
			"welcome message published",
			zap.Int(
				"unsigned_originator_envelope_size",
				len(envelope.UnsignedOriginatorEnvelope),
			),
			zap.Any("proof", envelope.Proof),
		)
	}

	return nil
}

func groupMessageCmd() *cobra.Command {
	var (
		numEnvelopes     uint
		dataSize         string
		grpcAddress      string
		privateKeyString string
		originatorID     uint32
	)

	cmd := &cobra.Command{
		Use:   "group-message",
		Short: "Publish group message envelopes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return groupMessageHandler(
				grpcAddress,
				privateKeyString,
				originatorID,
				numEnvelopes,
				dataSize,
			)
		},
		Example: `
Usage: xmtpd-cli generate group-message --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num> --data-size [256B|512B|1KB|5KB]

Example:
xmtpd-cli generate group-message --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num> --data-size [256B|512B|1KB|5KB]
`,
	}

	cmd.Flags().StringVarP(&grpcAddress, "grpc-address", "a", "", "xmtpd node gRPC address")
	_ = cmd.MarkFlagRequired("grpc-address")

	cmd.Flags().StringVar(&privateKeyString, "private-key", "", "payer private key")
	_ = cmd.MarkFlagRequired("private-key")

	cmd.Flags().Uint32VarP(&originatorID, "originator-id", "o", 100, "xmtpd node originator ID")

	// Note that num-envelopes doesn't currently work, until we implement batch publishing.
	cmd.Flags().UintVar(&numEnvelopes, "num-envelopes", 1, "number of envelopes to generate")

	cmd.Flags().
		StringVar(&dataSize, "data-size", "256B", "data size in bytes, options: 256B, 512B, 1KB, or 5KB")

	return cmd
}

func groupMessageHandler(
	nodeHTTPAddress string,
	privateKeyString string,
	originatorID uint32,
	numEnvelopes uint,
	dataSize string,
) error {
	logger, err := cliLogger()
	if err != nil {
		return err
	}

	ctx := context.Background()

	generator, err := stress.NewEnvelopesGenerator(nodeHTTPAddress, privateKeyString, originatorID)
	if err != nil {
		return err
	}

	defer func() {
		err := generator.Close()
		if err != nil {
			logger.Error("error closing generator", zap.Error(err))
		}
	}()

	var envelopes []*envelopes.OriginatorEnvelope

	envelopes, err = generator.PublishGroupMessageEnvelopes(ctx, numEnvelopes, dataSize)
	if err != nil {
		logger.Error("error publishing group message", zap.Error(err))
		return err
	}

	for _, envelope := range envelopes {
		logger.Info(
			"group message published",
			zap.Int(
				"unsigned_originator_envelope_size",
				len(envelope.UnsignedOriginatorEnvelope),
			),
			zap.Any("proof", envelope.Proof),
		)
	}

	return nil
}

func keyPackageCmd() *cobra.Command {
	var (
		numEnvelopes     uint
		grpcAddress      string
		privateKeyString string
		originatorID     uint32
	)

	cmd := &cobra.Command{
		Use:   "key-package",
		Short: "Publish key package envelopes",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return keyPackageHandler(
				grpcAddress,
				privateKeyString,
				originatorID,
				numEnvelopes,
			)
		},
		Example: `
Usage: xmtpd-cli generate key-package --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num>

Example:
xmtpd-cli generate key-package --grpc-address <address> --private-key <key> --originator-id <id> --num-envelopes <num>
`,
	}

	cmd.Flags().StringVarP(&grpcAddress, "grpc-address", "a", "", "xmtpd node gRPC address")
	_ = cmd.MarkFlagRequired("grpc-address")

	cmd.Flags().StringVar(&privateKeyString, "private-key", "", "payer private key")
	_ = cmd.MarkFlagRequired("private-key")

	cmd.Flags().Uint32VarP(&originatorID, "originator-id", "o", 100, "xmtpd node originator ID")

	// Note that num-envelopes doesn't currently work, until we implement batch publishing.
	cmd.Flags().UintVar(&numEnvelopes, "num-envelopes", 1, "number of envelopes to generate")

	return cmd
}

func keyPackageHandler(
	nodeHTTPAddress string,
	privateKeyString string,
	originatorID uint32,
	numEnvelopes uint,
) error {
	logger, err := cliLogger()
	if err != nil {
		return err
	}

	ctx := context.Background()

	generator, err := stress.NewEnvelopesGenerator(nodeHTTPAddress, privateKeyString, originatorID)
	if err != nil {
		return err
	}

	defer func() {
		err := generator.Close()
		if err != nil {
			logger.Error("error closing generator", zap.Error(err))
		}
	}()

	var envelopes []*envelopes.OriginatorEnvelope

	envelopes, err = generator.PublishKeyPackageEnvelopes(ctx, numEnvelopes)
	if err != nil {
		logger.Error("error publishing key package", zap.Error(err))
		return err
	}

	for _, envelope := range envelopes {
		logger.Info(
			"key package published",
			zap.Int(
				"unsigned_originator_envelope_size",
				len(envelope.UnsignedOriginatorEnvelope),
			),
			zap.Any("proof", envelope.Proof),
		)
	}

	return nil
}
