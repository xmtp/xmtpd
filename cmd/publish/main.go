package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/pingcap/log"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var originator uint32

func main() {
	addr := os.Args[1]
	originatorID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to parse originator ID: %v", err))
	}
	originator = uint32(originatorID)
	log.Info(fmt.Sprintf("attempting to connect to %s", addr))
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions())
	if err != nil {
		log.Info("Failed to connect to server. Retrying...")
	}
	client := message_api.NewReplicationApiClient(conn)

	_, err = client.PublishEnvelope(context.Background(), &message_api.PublishEnvelopeRequest{
		PayerEnvelope: CreatePayerEnvelope(),
	})
	if err != nil {
		log.Info(fmt.Sprintf("Failed to publish message: %v", err))
	}
}

func Marshal(t *testing.T, msg proto.Message) []byte {
	bytes, err := proto.Marshal(msg)
	require.NoError(t, err)
	return bytes
}

func UnmarshalUnsignedOriginatorEnvelope(
	t *testing.T,
	bytes []byte,
) *message_api.UnsignedOriginatorEnvelope {
	unsignedOriginatorEnvelope := &message_api.UnsignedOriginatorEnvelope{}
	err := proto.Unmarshal(
		bytes,
		unsignedOriginatorEnvelope,
	)
	require.NoError(t, err)
	return unsignedOriginatorEnvelope
}

func CreateClientEnvelope() *message_api.ClientEnvelope {
	return &message_api.ClientEnvelope{
		Payload: nil,
		Aad: &message_api.AuthenticatedData{
			TargetOriginator: originator,
			TargetTopic:      []byte{0x5},
			LastSeen:         &message_api.VectorClock{},
		},
	}
}

func CreatePayerEnvelope(
	clientEnv ...*message_api.ClientEnvelope,
) *message_api.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, CreateClientEnvelope())
	}
	clientEnvBytes, err := proto.Marshal(clientEnv[0])
	if err != nil {
		panic(err)
	}

	return &message_api.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature:         &associations.RecoverableEcdsaSignature{},
	}
}
