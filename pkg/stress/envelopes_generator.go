package stress

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/rand"
	"time"

	"connectrpc.com/connect"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	contents "github.com/xmtp/xmtpd/pkg/proto/mls/message_contents"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

type EnvelopesGenerator struct {
	client       message_apiconnect.ReplicationApiClient
	cleanup      func()
	privateKey   *ecdsa.PrivateKey
	originatorID uint32
}

type Protocol int

const (
	ProtocolConnect Protocol = iota
	ProtocolConnectGRPC
	ProtocolConnectGRPCWeb
	ProtocolNativeGRPC
)

func NewEnvelopesGenerator(
	nodeHTTPAddress string,
	privateKeyString string,
	originatorID uint32,
	protocol Protocol,
) (*EnvelopesGenerator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	privateKey, err := utils.ParseEcdsaPrivateKey(privateKeyString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse payer private key: %v", err)
	}

	var client message_apiconnect.ReplicationApiClient

	switch protocol {
	case ProtocolConnect:
		client, err = utils.NewConnectReplicationAPIClient(
			ctx,
			nodeHTTPAddress,
			utils.BuildConnectProtocolDialOptions()...,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build replication API client: %w", err)
		}

	case ProtocolConnectGRPC:
		client, err = utils.NewConnectGRPCReplicationAPIClient(
			ctx,
			nodeHTTPAddress,
			utils.BuildGRPCDialOptions()...)
		if err != nil {
			return nil, fmt.Errorf("failed to build replication API client: %w", err)
		}

	case ProtocolConnectGRPCWeb:
		client, err = utils.NewConnectGRPCWebReplicationAPIClient(
			ctx,
			nodeHTTPAddress,
			utils.BuildGRPCWebDialOptions()...)
		if err != nil {
			return nil, fmt.Errorf("failed to build replication API client: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid protocol: %d", protocol)
	}

	return &EnvelopesGenerator{
		client:       client,
		cleanup:      func() { cancel() },
		privateKey:   privateKey,
		originatorID: originatorID,
	}, nil
}

func (e *EnvelopesGenerator) Close() error {
	e.cleanup()
	return nil
}

// PublishWelcomeMessageEnvelopes publishes welcome message envelopes to the XMTPD node.
// The data size can be specified.
func (e *EnvelopesGenerator) PublishWelcomeMessageEnvelopes(
	ctx context.Context,
	numEnvelopes uint,
	dataSize uint,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	if dataSize > math.MaxUint {
		dataSize = uint(math.MaxUint)
	}

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
	size string,
) ([]*envelopesProto.OriginatorEnvelope, error) {
	clientEnvelopes := make([]*envelopesProto.ClientEnvelope, numEnvelopes)
	for i := range clientEnvelopes {
		clientEnvelopes[i] = makeGroupMessageEnvelope(size)
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
		connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: payerEnvelopes,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish payer envelopes: %v", err)
	}

	return r.Msg.OriginatorEnvelopes, nil
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
			MessageRetentionDays: constants.DefaultStorageDurationDays,
		}
	}

	return payerEnvelopes, nil
}

func makeWelcomeMessageClientEnvelope(dataSize uint) *envelopesProto.ClientEnvelope {
	var (
		installationKey  = testutils.RandomBytes(32)
		hpk              = testutils.RandomBytes(32)
		metadata         = testutils.RandomBytes(8)
		wrapperAlgorithm = int32(rand.Intn(3))
		payload          = testutils.RandomBytes(int(dataSize))
	)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
			WelcomeMessage: &mlsv1.WelcomeMessageInput{
				Version: &mlsv1.WelcomeMessageInput_V1_{
					V1: &mlsv1.WelcomeMessageInput_V1{
						InstallationKey: installationKey,
						Data:            payload,
						HpkePublicKey:   hpk,
						WrapperAlgorithm: contents.WelcomeWrapperAlgorithm(
							wrapperAlgorithm,
						),
						WelcomeMetadata: metadata,
					},
				},
			},
		},
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindWelcomeMessagesV1, installationKey[:]).
				Bytes(),
		},
	}
}

func makeKeyPackageClientEnvelope() *envelopesProto.ClientEnvelope {
	installationID := testutils.RandomBytes(32)

	return &envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_UploadKeyPackage{
			UploadKeyPackage: &mlsv1.UploadKeyPackageRequest{
				KeyPackage: &mlsv1.KeyPackageUpload{
					KeyPackageTlsSerialized: keyPackage,
				},
			},
		},
		Aad: &envelopesProto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindKeyPackagesV1, installationID[:]).
				Bytes(),
		},
	}
}

func makeGroupMessageEnvelope(size string) *envelopesProto.ClientEnvelope {
	var (
		groupID    = testutils.RandomGroupID()
		senderHmac = testutils.RandomBytes(32)
		data       []byte
	)

	switch size {
	case "256B":
		data = groupMessage256B
	case "512B":
		data = groupMessage512B
	case "1KB":
		data = groupMessage1KB
	case "5KB":
		data = groupMessage5KB
	default:
		data = groupMessage256B
	}

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
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, groupID[:]).
				Bytes(),
		},
	}
}
