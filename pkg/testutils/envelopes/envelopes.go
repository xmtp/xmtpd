package testutils

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/constants"

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
	if len(aad) == 0 {
		aad = append(aad, &envelopes.AuthenticatedData{
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
	groupID [16]byte,
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

func CreatePayerReportClientEnvelope(
	originatorID uint32,
) *envelopes.ClientEnvelope {
	return &envelopes.ClientEnvelope{
		Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_PAYER_REPORTS_V1, utils.Uint32ToBytes(originatorID)).
				Bytes(),
		},
		Payload: &envelopes.ClientEnvelope_PayerReport{
			PayerReport: &envelopes.PayerReport{},
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
		},
		Payload: &envelopes.ClientEnvelope_IdentityUpdate{
			IdentityUpdate: update,
		},
	}
}

func CreatePayerEnvelopeWithSigner(
	t *testing.T,
	nodeID uint32,
	signer *ecdsa.PrivateKey,
	expirationDays uint32,
	clientEnv *envelopes.ClientEnvelope,
) *envelopes.PayerEnvelope {
	clientEnvBytes, err := proto.Marshal(clientEnv)
	require.NoError(t, err)

	payerSignature, err := utils.SignClientEnvelope(nodeID, clientEnvBytes, signer)
	require.NoError(t, err)

	return &envelopes.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     nodeID,
		MessageRetentionDays: expirationDays,
	}
}

func CreatePayerEnvelopeWithExpiration(
	t *testing.T,
	nodeID uint32,
	expirationDays uint32,
	clientEnv ...*envelopes.ClientEnvelope,
) *envelopes.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, CreateClientEnvelope())
	}

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	return CreatePayerEnvelopeWithSigner(t, nodeID, key, expirationDays, clientEnv[0])
}

func CreatePayerEnvelope(
	t *testing.T,
	nodeID uint32,
	clientEnv ...*envelopes.ClientEnvelope,
) *envelopes.PayerEnvelope {
	return CreatePayerEnvelopeWithExpiration(
		t,
		nodeID,
		constants.DEFAULT_STORAGE_DURATION_DAYS,
		clientEnv...)
}

func CreateOriginatorEnvelopeWithTimestamp(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	timestamp time.Time,
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
		OriginatorNs:         timestamp.UnixNano(),
		PayerEnvelopeBytes:   marshaledPayerEnv,
	}

	unsignedBytes, err := proto.Marshal(unsignedEnv)
	require.NoError(t, err)

	return &envelopes.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof:                      nil,
	}
}

func CreateOriginatorEnvelope(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	payerEnv ...*envelopes.PayerEnvelope,
) *envelopes.OriginatorEnvelope {
	return CreateOriginatorEnvelopeWithTimestamp(
		t,
		originatorNodeID,
		originatorSequenceID,
		time.Unix(0, 0),
		payerEnv...)
}

func CreateOriginatorEnvelopeWithTopic(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	topic []byte,
) *envelopes.OriginatorEnvelope {
	payerEnv := CreatePayerEnvelope(t, originatorNodeID, CreateClientEnvelope(
		&envelopes.AuthenticatedData{
			TargetTopic: topic,
			DependsOn:   nil,
		},
	))

	return CreateOriginatorEnvelope(t, originatorNodeID, originatorSequenceID, payerEnv)
}
