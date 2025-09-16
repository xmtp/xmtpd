package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xmtp/xmtpd/cmd/xmtpd-cli/options"
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
		messageCmd(),
	)
	return &cmd
}

func messageCmd() *cobra.Command {
	var (
		messageType      options.MessageType
		numEnvelopes     uint
		dataSize         uint
		grpcAddress      string
		privateKeyString string
		originatorID     uint32
	)

	cmd := &cobra.Command{
		Use:   "message",
		Short: "Generate and publish message payloads",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return messageSetHandler(
				messageType,
				grpcAddress,
				privateKeyString,
				originatorID,
				numEnvelopes,
				dataSize,
			)
		},
	}

	cmd.Flags().Var(&messageType, "type", "welcome | key-package | group-message")
	_ = cmd.MarkFlagRequired("type")

	cmd.Flags().StringVarP(&grpcAddress, "grpc-address", "a", "", "xmtpd node gRPC address")
	_ = cmd.MarkFlagRequired("grpc-address")

	cmd.Flags().StringVar(&privateKeyString, "private-key", "", "payer private key")
	_ = cmd.MarkFlagRequired("private-key")

	cmd.Flags().Uint32VarP(&originatorID, "originator-id", "o", 100, "xmtpd node originator ID")

	// Note that num-envelopes doesn't currently work, until we implement batch publishing.
	cmd.Flags().UintVar(&numEnvelopes, "num-envelopes", 1, "number of envelopes to generate")
	cmd.Flags().UintVar(&dataSize, "data-size", 100, "data size in bytes")

	return cmd
}

func messageSetHandler(
	target options.MessageType,
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

	var envelopes []*envelopes.OriginatorEnvelope

	switch target {
	case options.MessageTypeWelcome:
		envelopes, err = generator.PublishWelcomeMessageEnvelopes(ctx, numEnvelopes, dataSize)
		if err != nil {
			logger.Fatal("write", zap.Error(err))
		}

		logger.Info(
			"welcome message published",
			zap.Uint("num_envelopes", numEnvelopes),
			zap.Uint("data_size", dataSize),
		)
	case options.MessageTypeKeyPackage:
		envelopes, err = generator.PublishKeyPackageEnvelopes(ctx, numEnvelopes)
		if err != nil {
			logger.Fatal("write", zap.Error(err))
		}

		logger.Info(
			"key package published",
			zap.Uint("num_envelopes", numEnvelopes),
		)
	case options.MessageTypeGroupMessage:
		envelopes, err = generator.PublishGroupMessageEnvelopes(ctx, numEnvelopes, dataSize)
		if err != nil {
			logger.Fatal("write", zap.Error(err))
		}

		logger.Info(
			"group message published",
			zap.Uint("num_envelopes", numEnvelopes),
			zap.Uint("data_size", dataSize),
		)
	default:
		logger.Fatal("target must be welcome | key-package | group-message")
	}

	for _, envelope := range envelopes {
		logger.Info(
			"envelope",
			zap.Any("envelope", envelope.Proof),
		)
	}

	return nil
}
