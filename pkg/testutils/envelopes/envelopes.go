package testutils

import (
	"crypto/ecdsa"
	"encoding/hex"
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

func getRealisticGroupMessagePayload() []byte {
	b, err := hex.DecodeString(
		"0001000210ed8344ef48284cf0ed2f74fc53b1887a000000000000000101001cf0abcef28080027c94c2ade64f6c3e" +
			"d6c7feaf5d75ed5927d1521ab340cb375f3b8a50520540c7c2cb9f6646812646ec8a9b74868f3049ef66d706e6e9b6" +
			"45014571d67b9483b1af909f5008ebfff94d870e74fd0c2791feb3ef08f92cf55b645d7992103fa18012d8f225b13d" +
			"589ce366ad8d041744f4e18e6b63b90c67325c24cb5a7e1d3e5717df5fa402b52e0e418f671053e10236337ac0e408" +
			"0de124f36e59a6a70dbf9f5d62cdfc60004bb16fbc1f89a289bd8edc08b137ffba4dc948f1867b17ea4962a8740082" +
			"7eccf73a4e8cbf965b2ef7070d0a604cb6fe70a6e52a0c1eb52bfe2273c4",
	)
	if err != nil {
		panic("could not generate bytes")
	}
	return b
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
		Payload: &envelopes.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data: getRealisticGroupMessagePayload(),
					},
				},
			},
		},
		Aad: aad[0],
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
