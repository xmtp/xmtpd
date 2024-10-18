package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopes "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"google.golang.org/protobuf/proto"
)

func UnmarshalUnsignedOriginatorEnvelope(
	t *testing.T,
	bytes []byte,
) *envelopes.UnsignedOriginatorEnvelope {
	unsignedOriginatorEnvelope := &envelopes.UnsignedOriginatorEnvelope{}
	err := proto.Unmarshal(
		bytes,
		unsignedOriginatorEnvelope,
	)
	require.NoError(t, err)
	return unsignedOriginatorEnvelope
}

func CreateClientEnvelope(aad ...*envelopes.AuthenticatedData) *envelopes.ClientEnvelope {
	if len(aad) == 0 {
		aad = append(aad, &envelopes.AuthenticatedData{
			TargetOriginator: 1,
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
				Bytes(),
			LastSeen: &envelopes.VectorClock{},
		})
	}
	return &envelopes.ClientEnvelope{
		Payload: &envelopes.ClientEnvelope_GroupMessage{},
		Aad:     aad[0],
	}
}

func CreatePayerEnvelope(
	t *testing.T,
	clientEnv ...*envelopes.ClientEnvelope,
) *envelopes.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, CreateClientEnvelope())
	}
	clientEnvBytes, err := proto.Marshal(clientEnv[0])
	require.NoError(t, err)

	return &envelopes.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature:         &associations.RecoverableEcdsaSignature{},
	}
}

func CreateOriginatorEnvelope(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	payerEnv ...*envelopes.PayerEnvelope,
) *envelopes.OriginatorEnvelope {
	if len(payerEnv) == 0 {
		payerEnv = append(payerEnv, CreatePayerEnvelope(t))
	}

	unsignedEnv := &envelopes.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     originatorNodeID,
		OriginatorSequenceId: originatorSequenceID,
		OriginatorNs:         0,
		PayerEnvelope:        payerEnv[0],
	}

	unsignedBytes, err := proto.Marshal(unsignedEnv)
	require.NoError(t, err)

	return &envelopes.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof:                      nil,
	}
}
