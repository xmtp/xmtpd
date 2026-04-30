package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"connectrpc.com/connect"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	apiv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api/message_apiconnect"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

// newConnectHTTPClient builds an HTTP/2 client, using TLS unless cfg.Insecure.
func newConnectHTTPClient(cfg *config) (client *http.Client, baseURL string) {
	if cfg.Insecure {
		client = &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
					return (&net.Dialer{}).DialContext(ctx, network, addr)
				},
			},
		}
		baseURL = "http://" + cfg.Addr
	} else {
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
				DialTLSContext: func(
					ctx context.Context, network, addr string, tlsCfg *tls.Config,
				) (net.Conn, error) {
					return (&tls.Dialer{Config: tlsCfg}).DialContext(ctx, network, addr)
				},
			},
		}
		baseURL = "https://" + cfg.Addr
	}
	return client, baseURL
}

// runSubscribeEnvelopesDiagnostic tests SubscribeEnvelopes (older RPC) for live delivery.
func runSubscribeEnvelopesDiagnostic(cfg *config) error {
	fmt.Println("\n=== SubscribeEnvelopes Diagnostic ===")

	httpClient, baseURL := newConnectHTTPClient(cfg)
	opts := []connect.ClientOption{connect.WithGRPC()}

	replicationClient := message_apiconnect.NewReplicationApiClient(httpClient, baseURL, opts...)
	publishClient := message_apiconnect.NewPublishApiClient(httpClient, baseURL, opts...)

	topicID := randomBytes(16)
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)
	topicBytes := tp.Bytes()
	fmt.Printf("Topic: %s\n", tp.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Open SubscribeEnvelopes with topic filter (intentionally testing deprecated API)
	stream, err := replicationClient.SubscribeEnvelopes( //nolint:staticcheck
		ctx,
		connect.NewRequest(
			&message_api.SubscribeEnvelopesRequest{
				Query: &message_api.EnvelopesQuery{
					Topics: [][]byte{topicBytes},
				},
			},
		),
	)
	if err != nil {
		return fmt.Errorf("open SubscribeEnvelopes: %w", err)
	}
	fmt.Println("  Stream opened")

	// Consume initial keepalive (empty response)
	if stream.Receive() {
		envs := stream.Msg().GetEnvelopes()
		if len(envs) == 0 {
			fmt.Println("  Got initial keepalive")
		}
	}

	time.Sleep(500 * time.Millisecond)

	// Publish
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("gen key: %w", err)
	}

	published := 0
	for range 5 {
		aad := &envelopesProto.AuthenticatedData{TargetTopic: topicBytes}
		clientEnv := &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_GroupMessage{
				GroupMessage: &apiv1.GroupMessageInput{
					Version: &apiv1.GroupMessageInput_V1_{
						V1: &apiv1.GroupMessageInput_V1{Data: makeGroupMessagePayload(256)},
					},
				},
			},
		}
		clientEnvBytes, err := proto.Marshal(clientEnv)
		if err != nil {
			continue
		}
		sig, err := utils.SignClientEnvelope(cfg.NodeID, clientEnvBytes, key)
		if err != nil {
			continue
		}
		req := connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopesProto.PayerEnvelope{{
				UnsignedClientEnvelope: clientEnvBytes,
				PayerSignature:         &associations.RecoverableEcdsaSignature{Bytes: sig},
				TargetOriginator:       cfg.NodeID,
				MessageRetentionDays:   constants.DefaultStorageDurationDays,
			}},
		})
		_, err = publishClient.PublishPayerEnvelopes(ctx, req)
		if err != nil {
			fmt.Printf("  Publish error: %v\n", err)
			continue
		}
		published++
	}
	fmt.Printf("  Published %d messages\n", published)

	// Wait for live delivery
	fmt.Println("  Waiting 5s for live delivery...")
	var received atomic.Int64
	recvDone := make(chan struct{})
	go func() {
		defer close(recvDone)
		for stream.Receive() {
			envs := stream.Msg().GetEnvelopes()
			if len(envs) > 0 {
				total := received.Add(int64(len(envs)))
				fmt.Printf("  SubscribeEnvelopes received %d (total %d)\n", len(envs), total)
				if int(total) >= published {
					return
				}
			} else {
				fmt.Println("  SubscribeEnvelopes got keepalive")
			}
		}
	}()

	select {
	case <-recvDone:
	case <-time.After(5 * time.Second):
	}

	fmt.Printf("  SubscribeEnvelopes received: %d\n\n", received.Load())
	return nil
}

// runConnectDiagnostic tests live delivery using Connect-RPC client (same as unit tests).
func runConnectDiagnostic(cfg *config) error {
	fmt.Println("=== Connect-RPC Client Diagnostic ===")

	httpClient, baseURL := newConnectHTTPClient(cfg)
	opts := []connect.ClientOption{connect.WithGRPC()}

	queryClient := message_apiconnect.NewQueryApiClient(httpClient, baseURL, opts...)
	publishClient := message_apiconnect.NewPublishApiClient(httpClient, baseURL, opts...)

	// Create a random topic
	topicID := randomBytes(16)
	tp := topic.NewTopic(topic.TopicKindGroupMessagesV1, topicID)
	topicBytes := tp.Bytes()
	fmt.Printf("Topic: %s (%x)\n", tp.String(), topicBytes)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Step 1: Open SubscribeTopics stream via Connect-RPC gRPC client
	stream, err := queryClient.SubscribeTopics(
		ctx,
		connect.NewRequest(&message_api.SubscribeTopicsRequest{
			Filters: []*message_api.SubscribeTopicsRequest_TopicFilter{
				{Topic: topicBytes},
			},
		}),
	)
	if err != nil {
		return fmt.Errorf("open stream: %w", err)
	}
	fmt.Println("  Stream opened")

	// Drain initial status messages
	for stream.Receive() {
		su := stream.Msg().GetStatusUpdate()
		if su == nil {
			break
		}
		fmt.Printf("  Status: %s\n", su.GetStatus())
		if su.GetStatus() == message_api.SubscribeTopicsResponse_SUBSCRIPTION_STATUS_CATCHUP_COMPLETE {
			break
		}
	}
	if stream.Err() != nil {
		return fmt.Errorf("stream error during init: %w", stream.Err())
	}

	time.Sleep(500 * time.Millisecond)

	// Step 2: Publish 5 messages via Connect-RPC
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("gen key: %w", err)
	}

	published := 0
	for range 5 {
		aad := &envelopesProto.AuthenticatedData{TargetTopic: topicBytes}
		clientEnv := &envelopesProto.ClientEnvelope{
			Aad: aad,
			Payload: &envelopesProto.ClientEnvelope_GroupMessage{
				GroupMessage: &apiv1.GroupMessageInput{
					Version: &apiv1.GroupMessageInput_V1_{
						V1: &apiv1.GroupMessageInput_V1{Data: makeGroupMessagePayload(256)},
					},
				},
			},
		}
		clientEnvBytes, err := proto.Marshal(clientEnv)
		if err != nil {
			continue
		}
		sig, err := utils.SignClientEnvelope(cfg.NodeID, clientEnvBytes, key)
		if err != nil {
			continue
		}
		req := connect.NewRequest(&message_api.PublishPayerEnvelopesRequest{
			PayerEnvelopes: []*envelopesProto.PayerEnvelope{{
				UnsignedClientEnvelope: clientEnvBytes,
				PayerSignature:         &associations.RecoverableEcdsaSignature{Bytes: sig},
				TargetOriginator:       cfg.NodeID,
				MessageRetentionDays:   constants.DefaultStorageDurationDays,
			}},
		})
		_, err = publishClient.PublishPayerEnvelopes(ctx, req)
		if err != nil {
			fmt.Printf("  Publish error: %v\n", err)
			continue
		}
		published++
	}
	fmt.Printf("  Published %d messages\n", published)

	// Step 3: Wait for live delivery
	fmt.Println("  Waiting 5s for live delivery...")
	var received2 atomic.Int64
	recvDone := make(chan struct{})
	go func() {
		defer close(recvDone)
		for stream.Receive() {
			if su := stream.Msg().GetStatusUpdate(); su != nil {
				fmt.Printf("  Stream status: %s\n", su.GetStatus())
				continue
			}
			envs := stream.Msg().GetEnvelopes()
			if envs != nil && len(envs.GetEnvelopes()) > 0 {
				n := int64(len(envs.GetEnvelopes()))
				total := received2.Add(n)
				fmt.Printf("  Received %d envelopes (total %d)\n", n, total)
				if int(total) >= published {
					return
				}
			}
		}
	}()

	select {
	case <-recvDone:
	case <-time.After(5 * time.Second):
	}

	recv := received2.Load()
	fmt.Printf("\n=== CONNECT-RPC RESULTS ===\n")
	fmt.Printf("Published:  %d\n", published)
	fmt.Printf("Received:   %d\n", recv)

	if recv == 0 {
		fmt.Println("\nConnect-RPC client ALSO gets 0 messages — server bug (not protocol)")
	} else {
		fmt.Println("\nConnect-RPC works! — native gRPC protocol issue with server")
	}

	return nil
}
