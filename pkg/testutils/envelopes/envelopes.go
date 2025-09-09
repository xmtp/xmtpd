// Package envelopes implements the envelopes test utils.
package envelopes

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

const DefaultClientEnvelopeNodeID = uint32(100)

const (
	MinimalCommitPayload      = "0001000210aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000000000000103000000"
	MinimalApplicationPayload = "0001000210aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000000000000101000000"
)

func GetRealisticGroupMessagePayload(makeCommit bool) []byte {
	if makeCommit {
		b, err := hex.DecodeString(MinimalCommitPayload)
		if err != nil {
			panic("could not generate bytes")
		}
		return b
	} else {
		b, err := hex.DecodeString(MinimalApplicationPayload)
		if err != nil {
			panic("could not generate bytes")
		}
		return b
	}
}

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

type ClientEnvelopeOptions struct {
	Aad      *envelopes.AuthenticatedData
	IsCommit bool
}

func CreateClientEnvelope(options ...*ClientEnvelopeOptions) *envelopes.ClientEnvelope {
	var aad *envelopes.AuthenticatedData
	var isCommit bool

	if len(options) == 0 {
		aad = &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3}).
				Bytes(),
			DependsOn: &envelopes.Cursor{},
		}
		isCommit = false
	} else {
		option := options[0]
		if option.IsCommit {
			isCommit = true
		}
		if option.Aad != nil {
			aad = option.Aad
		}
	}
	return &envelopes.ClientEnvelope{
		Payload: &envelopes.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data: GetRealisticGroupMessagePayload(isCommit),
					},
				},
			},
		},
		Aad: aad,
	}
}

func CreateGroupMessageClientEnvelope(
	groupID [16]byte,
	message []byte,
) *envelopes.ClientEnvelope {
	return &envelopes.ClientEnvelope{
		Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, groupID[:]).
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
			TargetTopic: topic.NewTopic(topic.TopicKindPayerReportsV1, utils.Uint32ToBytes(originatorID)).
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
			TargetTopic: topic.NewTopic(topic.TopicKindIdentityUpdatesV1, inboxID[:]).
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
		constants.DefaultStorageDurationDays,
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
		&ClientEnvelopeOptions{Aad: &envelopes.AuthenticatedData{
			TargetTopic: topic,
			DependsOn:   nil,
		}},
	))

	return CreateOriginatorEnvelope(t, originatorNodeID, originatorSequenceID, payerEnv)
}
