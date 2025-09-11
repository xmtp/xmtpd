package stress

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	contents "github.com/xmtp/xmtpd/pkg/proto/mls/message_contents"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type EnvelopesGenerator struct {
	client       message_api.ReplicationApiClient
	cleanup      func() error
	privateKey   *ecdsa.PrivateKey
	originatorID uint32
}

func NewEnvelopesGenerator(
	nodeHTTPAddress string,
	privateKeyString string,
	originatorID uint32,
) (*EnvelopesGenerator, error) {
	conn, err := buildGRPCConnection(nodeHTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to build gRPC client: %v", err)
	}

	client := message_api.NewReplicationApiClient(conn)

	privateKey, err := utils.ParseEcdsaPrivateKey(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse payer private key: %v", err)
	}

	return &EnvelopesGenerator{
		client:       client,
		cleanup:      func() error { return conn.Close() },
		privateKey:   privateKey,
		originatorID: originatorID,
	}, nil
}

func (e *EnvelopesGenerator) Close() error {
	return e.cleanup()
}

// PublishWelcomeMessageEnvelopes publishes welcome message envelopes to the XMTPD node.
// The data size can be specified.
func (e *EnvelopesGenerator) PublishWelcomeMessageEnvelopes(
	ctx context.Context,
	numEnvelopes uint,
	dataSize uint,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	clientEnvelopes := make([]*envelopesProto.ClientEnvelope, numEnvelopes)
	for i := range clientEnvelopes {
		clientEnvelopes[i] = makeWelcomeMessageClientEnvelope(dataSize)
	}

	payerEnvelopes, err := e.buildAndSignPayerEnvelopes(clientEnvelopes)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign payer envelopes: %v", err)
	}

	return e.publishPayerEnvelopes(ctx, payerEnvelopes)
}

// PublishKeyPackageEnvelopes publishes key package envelopes to the XMTPD node.
// The data size is hardcoded to 1651 bytes, as expected by the protocol.
func (e *EnvelopesGenerator) PublishKeyPackageEnvelopes(
	ctx context.Context,
	numEnvelopes uint,
	dataSize uint,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	clientEnvelopes := make([]*envelopesProto.ClientEnvelope, numEnvelopes)
	for i := range clientEnvelopes {
		clientEnvelopes[i] = makeKeyPackageClientEnvelope()
	}

	payerEnvelopes, err := e.buildAndSignPayerEnvelopes(clientEnvelopes)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign payer envelopes: %v", err)
	}

	return e.publishPayerEnvelopes(ctx, payerEnvelopes)
}

// PublishGroupMessageEnvelopes publishes group message envelopes to the XMTPD node.
// The data size can be specified.
func (e *EnvelopesGenerator) PublishGroupMessageEnvelopes(
	ctx context.Context,
	numEnvelopes uint,
	dataSize uint,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	clientEnvelopes := make([]*envelopesProto.ClientEnvelope, numEnvelopes)
	for i := range clientEnvelopes {
		clientEnvelopes[i] = makeGroupMessageEnvelope(dataSize)
	}

	payerEnvelopes, err := e.buildAndSignPayerEnvelopes(clientEnvelopes)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign payer envelopes: %v", err)
	}

	return e.publishPayerEnvelopes(ctx, payerEnvelopes)
}

func (e *EnvelopesGenerator) publishPayerEnvelopes(
	ctx context.Context,
	payerEnvelopes []*envelopesProto.PayerEnvelope,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	r, err := e.client.PublishPayerEnvelopes(
		ctx,
		&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: payerEnvelopes,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish payer envelopes: %v", err)
	}

	return r.OriginatorEnvelopes, err
}

func (e *EnvelopesGenerator) buildAndSignPayerEnvelopes(
	protoClientEnvelope []*envelopesProto.ClientEnvelope,
) ([]*envelopesProto.PayerEnvelope, error) {
	payerEnvelopes := make([]*envelopesProto.PayerEnvelope, len(protoClientEnvelope))

	for i, protoClientEnvelope := range protoClientEnvelope {
		clientEnvelope, err := envelopes.NewClientEnvelope(protoClientEnvelope)
		if err != nil {
			return nil, fmt.Errorf("failed to build client envelope: %w", err)
		}

		clientEnvelopeBytes, err := clientEnvelope.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to get client envelope bytes: %w", err)
		}

		payerSignature, err := utils.SignClientEnvelope(
			e.originatorID,
			clientEnvelopeBytes,
			e.privateKey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to sign client envelope: %w", err)
		}

		payerEnvelopes[i] = &envelopesProto.PayerEnvelope{
			UnsignedClientEnvelope: clientEnvelopeBytes,
			PayerSignature: &associations.RecoverableEcdsaSignature{
				Bytes: payerSignature,
			},
			TargetOriginator:     e.originatorID,
			MessageRetentionDays: constants.DEFAULT_STORAGE_DURATION_DAYS,
		}
	}

	return payerEnvelopes, nil
}

func makeWelcomeMessageClientEnvelope(dataSize uint) *envelopesProto.ClientEnvelope {
	installationKey := testutils.RandomBytes(32)
	hpk := testutils.RandomBytes(32)
	data := testutils.RandomBytes(int(dataSize))
	metadata := testutils.RandomBytes(8)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
			WelcomeMessage: &mlsv1.WelcomeMessageInput{
				Version: &mlsv1.WelcomeMessageInput_V1_{
					V1: &mlsv1.WelcomeMessageInput_V1{
						InstallationKey: installationKey,
						Data:            data,
						HpkePublicKey:   hpk,
						WrapperAlgorithm: contents.WelcomeWrapperAlgorithm(
							testutils.RandomInt32(),
						),
						WelcomeMetadata: metadata,
					},
				},
			},
		},
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_WELCOME_MESSAGES_V1, installationKey[:]).
				Bytes(),
		},
	}
}

func makeKeyPackageClientEnvelope() *envelopesProto.ClientEnvelope {
	installationID := testutils.RandomBytes(32)
	keyPackage := testutils.RandomBytes(1651)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_UploadKeyPackage{
			UploadKeyPackage: &mlsv1.UploadKeyPackageRequest{
				KeyPackage: &mlsv1.KeyPackageUpload{
					KeyPackageTlsSerialized: keyPackage,
				},
			},
		},
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, installationID[:]).
				Bytes(),
		},
	}
}

func makeGroupMessageEnvelope(dataSize uint) *envelopesProto.ClientEnvelope {
	groupID := testutils.RandomGroupID()
	data := testutils.RandomBytes(int(dataSize))
	senderHmac := testutils.RandomBytes(32)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data:       data,
						SenderHmac: senderHmac,
						ShouldPush: true,
					},
				},
			},
		},
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupID[:]).
				Bytes(),
		},
	}
}

func buildGRPCConnection(
	nodeHTTPAddress string,
	extraDialOpts ...grpc.DialOption,
) (*grpc.ClientConn, error) {
	target, isTLS, err := utils.HttpAddressToGrpcTarget(nodeHTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTTP address to gRPC target: %v", err)
	}

	creds, err := utils.GetCredentialsForAddress(isTLS)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %v", err)
	}

	dialOpts := append([]grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}, extraDialOpts...)

	conn, err := grpc.NewClient(
		target,
		dialOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel at %s: %v", target, err)
	}

	return conn, nil
}
