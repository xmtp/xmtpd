package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api/payer_apiconnect"
	"github.com/xmtp/xmtpd/pkg/topic"
	"golang.org/x/net/http2"
)

// gatewayClient wraps a Connect-RPC PayerApi client for sending unsigned
// ClientEnvelopes through the gateway (which signs them with its payer key).
type gatewayClient struct {
	client payer_apiconnect.PayerApiClient
}

// newGatewayClient creates a Connect-RPC client targeting the gateway's PayerApi.
func newGatewayClient(gatewayAddr string) *gatewayClient {
	httpClient := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
			DialTLSContext: func(ctx context.Context, network, addr string, tlsCfg *tls.Config) (net.Conn, error) {
				return tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, network, addr, tlsCfg)
			},
		},
	}

	baseURL := "https://" + gatewayAddr
	client := payer_apiconnect.NewPayerApiClient(httpClient, baseURL, connect.WithGRPC())

	return &gatewayClient{client: client}
}

// publishClientEnvelope sends an unsigned ClientEnvelope to the gateway.
func (g *gatewayClient) publishClientEnvelope(
	ctx context.Context,
	topicKind topic.TopicKind,
	targetTopic *topic.Topic,
	payloadSize int,
) error {
	tc := testCase{
		TopicKind:   topicKind,
		PayloadSize: payloadSize,
	}
	// Build the unsigned client envelope (same as what SDK sends to gateway)
	clientEnv := buildClientEnvelopeForTopic(targetTopic, tc)

	req := connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
		Envelopes: []*envelopesProto.ClientEnvelope{clientEnv},
	})

	_, err := g.client.PublishClientEnvelopes(ctx, req)
	return err
}

// publishCommitEnvelope sends a commit-type GroupMessage through gateway → blockchain.
func (g *gatewayClient) publishCommitEnvelope(
	ctx context.Context,
	targetTopic *topic.Topic,
	payloadSize int,
) error {
	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: targetTopic.Bytes(),
	}
	clientEnv := &envelopesProto.ClientEnvelope{
		Aad: aad,
		Payload: &envelopesProto.ClientEnvelope_GroupMessage{
			GroupMessage: &apiv1.GroupMessageInput{
				Version: &apiv1.GroupMessageInput_V1_{
					V1: &apiv1.GroupMessageInput_V1{
						Data: makeCommitMessagePayload(payloadSize),
					},
				},
			},
		},
	}

	req := connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
		Envelopes: []*envelopesProto.ClientEnvelope{clientEnv},
	})

	_, err := g.client.PublishClientEnvelopes(ctx, req)
	return err
}

// publishKeyPackageEnvelope sends a key package upload through gateway → node.
func (g *gatewayClient) publishKeyPackageEnvelope(
	ctx context.Context,
) error {
	installationID := randomBytes(32)
	// Minimal valid key package TLS payload (~1651 bytes in real MLS, we use 256B placeholder)
	keyPkgData := randomBytes(256)

	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: topic.NewTopic(topic.TopicKindKeyPackagesV1, installationID).Bytes(),
	}
	clientEnv := &envelopesProto.ClientEnvelope{
		Aad: aad,
		Payload: &envelopesProto.ClientEnvelope_UploadKeyPackage{
			UploadKeyPackage: &apiv1.UploadKeyPackageRequest{
				KeyPackage: &apiv1.KeyPackageUpload{
					KeyPackageTlsSerialized: keyPkgData,
				},
			},
		},
	}

	req := connect.NewRequest(&payer_api.PublishClientEnvelopesRequest{
		Envelopes: []*envelopesProto.ClientEnvelope{clientEnv},
	})

	_, err := g.client.PublishClientEnvelopes(ctx, req)
	return err
}

// buildClientEnvelopeForTopic builds an unsigned ClientEnvelope targeting a specific topic.
func buildClientEnvelopeForTopic(tp *topic.Topic, tc testCase) *envelopesProto.ClientEnvelope {
	aad := &envelopesProto.AuthenticatedData{
		TargetTopic: tp.Bytes(),
	}

	switch tc.TopicKind {
	case topic.TopicKindGroupMessagesV1:
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_GroupMessage{
				GroupMessage: buildGroupMessageInput(tc.PayloadSize),
			},
		}
	case topic.TopicKindWelcomeMessagesV1:
		return &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_WelcomeMessage{
				WelcomeMessage: buildWelcomeMessageInput(tc.PayloadSize),
			},
		}
	default:
		panic("unsupported topic kind")
	}
}

func buildGroupMessageInput(size int) *apiv1.GroupMessageInput {
	return &apiv1.GroupMessageInput{
		Version: &apiv1.GroupMessageInput_V1_{
			V1: &apiv1.GroupMessageInput_V1{
				Data: makeGroupMessagePayload(size),
			},
		},
	}
}

func buildWelcomeMessageInput(size int) *apiv1.WelcomeMessageInput {
	return &apiv1.WelcomeMessageInput{
		Version: &apiv1.WelcomeMessageInput_V1_{
			V1: &apiv1.WelcomeMessageInput_V1{
				Data: randomBytes(size),
			},
		},
	}
}
