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
	messageContents "github.com/xmtp/xmtpd/pkg/proto/mls/message_contents"
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
	signer       *ecdsa.PrivateKey
	originatorID uint32
}

func NewEnvelopesGenerator(
	nodeHTTPAddress string,
	privateKey string,
	originatorID uint32,
) (*EnvelopesGenerator, error) {
	conn, err := buildGRPCConnection(nodeHTTPAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to build gRPC client: %v", err)
	}

	client := message_api.NewReplicationApiClient(conn)

	signer, err := utils.ParseEcdsaPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse payer private key: %v", err)
	}

	return &EnvelopesGenerator{
		client:       client,
		signer:       signer,
		originatorID: originatorID,
	}, nil
}

func (e *EnvelopesGenerator) PublishWelcomeMessageEnvelopes(
	ctx context.Context,
	numEnvelopes int,
	dataSize int,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	clientEnvelopes := make([]*envelopesProto.ClientEnvelope, numEnvelopes)
	for i := 0; i < numEnvelopes; i++ {
		clientEnvelopes[i] = getWelcomeMessageClientEnvelope(dataSize)
	}

	payerEnvelopes := make([]*envelopesProto.PayerEnvelope, numEnvelopes)
	for i := 0; i < numEnvelopes; i++ {
		payerEnvelopes[i], _ = e.buildAndSignPayerEnvelope(
			clientEnvelopes[i],
		)
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

func (e *EnvelopesGenerator) buildAndSignPayerEnvelope(
	protoClientEnvelope *envelopesProto.ClientEnvelope,
) (*envelopesProto.PayerEnvelope, error) {
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
		e.signer,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign client envelope: %w", err)
	}

	return &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     e.originatorID,
		MessageRetentionDays: constants.DEFAULT_STORAGE_DURATION_DAYS,
	}, nil
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

func getWelcomeMessageClientEnvelope(dataSize int) *envelopesProto.ClientEnvelope {
	installationKey := testutils.RandomBytes(32)
	hpk := testutils.RandomBytes(32)
	data := testutils.RandomBytes(dataSize)
	metadata := testutils.RandomBytes(8)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
			WelcomeMessage: &mlsv1.WelcomeMessageInput{
				Version: &mlsv1.WelcomeMessageInput_V1_{
					V1: &mlsv1.WelcomeMessageInput_V1{
						InstallationKey: installationKey,
						Data:            data,
						HpkePublicKey:   hpk,
						WrapperAlgorithm: messageContents.WelcomeWrapperAlgorithm(
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
