package testutils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	envelopes "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

const DefaultClientEnvelopeNodeId = uint32(100)

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
	nodeId := DefaultClientEnvelopeNodeId
	if len(aad) == 0 {
		aad = append(aad, &envelopes.AuthenticatedData{
			TargetOriginator: &nodeId,
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte{1, 2, 3}).
				Bytes(),
			DependsOn: &envelopes.Cursor{},
		})
	}
	return &envelopes.ClientEnvelope{
		Payload: &envelopes.ClientEnvelope_GroupMessage{},
		Aad:     aad[0],
	}
}

func CreateGroupMessageClientEnvelope(
	groupID [32]byte,
	message []byte,
) *envelopes.ClientEnvelope {
	return &envelopes.ClientEnvelope{
		Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupID[:]).
				Bytes(),
		},
		Payload: &envelopes.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data: message,
					},
				},
			},
		},
	}
}

func CreateIdentityUpdateClientEnvelope(
	inboxID [32]byte,
	update *associations.IdentityUpdate,
) *envelopes.ClientEnvelope {
	return &envelopes.ClientEnvelope{
		Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, inboxID[:]).
				Bytes(),
			TargetOriginator: new(uint32),
		},
		Payload: &envelopes.ClientEnvelope_IdentityUpdate{
			IdentityUpdate: update,
		},
	}
}

func CreatePayerEnvelope(
	t *testing.T,
	nodeID uint32,
	clientEnv ...*envelopes.ClientEnvelope,
) *envelopes.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, CreateClientEnvelope())
	}

	clientEnvBytes, err := proto.Marshal(clientEnv[0])
	require.NoError(t, err)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	payerSignature, err := utils.SignClientEnvelope(nodeID, clientEnvBytes, key)
	require.NoError(t, err)

	return &envelopes.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator: nodeID,
	}
}

func CreateOriginatorEnvelope(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	payerEnv ...*envelopes.PayerEnvelope,
) *envelopes.OriginatorEnvelope {
	if len(payerEnv) == 0 {
		payerEnv = append(payerEnv, CreatePayerEnvelope(t, originatorNodeID))
	}

	marshaledPayerEnv, err := proto.Marshal(payerEnv[0])
	require.NoError(t, err)

	unsignedEnv := &envelopes.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     originatorNodeID,
		OriginatorSequenceId: originatorSequenceID,
		OriginatorNs:         0,
		PayerEnvelopeBytes:   marshaledPayerEnv,
	}

	unsignedBytes, err := proto.Marshal(unsignedEnv)
	require.NoError(t, err)

	return &envelopes.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof:                      nil,
	}
}

func CreateOriginatorEnvelopeWithTopic(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	topic []byte,
) *envelopes.OriginatorEnvelope {
	payerEnv := CreatePayerEnvelope(t, originatorNodeID, CreateClientEnvelope(
		&envelopes.AuthenticatedData{
			TargetTopic:      topic,
			TargetOriginator: &originatorNodeID,
			DependsOn:        nil,
		},
	))

	return CreateOriginatorEnvelope(t, originatorNodeID, originatorSequenceID, payerEnv)
}
