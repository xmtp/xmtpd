package envelopes

import (
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func TestValidOriginatorEnvelope(t *testing.T) {
	originatorNodeID := uint32(1)
	originatorSequenceID := uint64(1)

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(t, originatorNodeID, clientEnv)
	originatorEnvelope := envelopeTestUtils.CreateOriginatorEnvelope(
		t,
		originatorNodeID,
		originatorSequenceID,
		payerEnvelope,
	)

	processed, err := NewOriginatorEnvelope(originatorEnvelope)
	require.NoError(t, err)
	require.Equal(t, originatorNodeID, processed.UnsignedOriginatorEnvelope.OriginatorNodeID())
	require.Equal(
		t,
		originatorSequenceID,
		processed.UnsignedOriginatorEnvelope.OriginatorSequenceID(),
	)

	serializedClientEnv, err := proto.Marshal(clientEnv)
	require.NoError(t, err)
	serializedClientEnvAfterParse, err := processed.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	require.NoError(t, err)
	require.Equal(t, serializedClientEnv, serializedClientEnvAfterParse)
}

func TestSerialize(t *testing.T) {
	originatorNodeID := uint32(1)
	originatorSequenceID := uint64(1)

	clientEnv := envelopeTestUtils.CreateClientEnvelope()
	payerEnvelope := envelopeTestUtils.CreatePayerEnvelope(t, originatorNodeID, clientEnv)
	originatorEnvelope := envelopeTestUtils.CreateOriginatorEnvelope(
		t,
		originatorNodeID,
		originatorSequenceID,
		payerEnvelope,
	)

	serializedFromProto, err := proto.Marshal(originatorEnvelope)
	require.NoError(t, err)

	originatorStruct, err := NewOriginatorEnvelope(originatorEnvelope)
	require.NoError(t, err)
	serializedFromStruct, err := originatorStruct.Bytes()
	require.NoError(t, err)

	require.Equal(t, serializedFromProto, serializedFromStruct)
}

func TestInvalidOriginatorEnvelope(t *testing.T) {
	_, err := NewOriginatorEnvelope(nil)
	require.Error(t, err)

	empty := &envelopesProto.OriginatorEnvelope{}
	_, err = NewOriginatorEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidUnsignedOriginatorEnvelope(t *testing.T) {
	_, err := NewUnsignedOriginatorEnvelope(nil)
	require.Error(t, err)

	empty := &envelopesProto.UnsignedOriginatorEnvelope{}
	_, err = NewUnsignedOriginatorEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidPayerEnvelope(t *testing.T) {
	_, err := NewPayerEnvelope(nil)
	require.Error(t, err)

	empty := &envelopesProto.PayerEnvelope{}
	_, err = NewPayerEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidClientEnvelope(t *testing.T) {
	_, err := NewClientEnvelope(nil)
	require.Error(t, err)

	empty := &envelopesProto.ClientEnvelope{}
	_, err = NewClientEnvelope(empty)
	require.Error(t, err)
}

func buildAad(topic *topic.Topic) *envelopesProto.AuthenticatedData {
	return &envelopesProto.AuthenticatedData{
		TargetTopic: topic.Bytes(),
		DependsOn:   &envelopesProto.Cursor{},
	}
}

func TestPayloadType(t *testing.T) {
	// Group Message envelope with matching topic
	clientEnvelope, err := NewClientEnvelope(&envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_GroupMessage{},
		Aad:     buildAad(topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3})),
	})
	require.NoError(t, err)
	require.True(t, clientEnvelope.TopicMatchesPayload())

	clientEnvelope, err = NewClientEnvelope(&envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_UploadKeyPackage{},
		Aad:     buildAad(topic.NewTopic(topic.TopicKindKeyPackagesV1, []byte{1, 2, 3})),
	})
	require.NoError(t, err)
	require.True(t, clientEnvelope.TopicMatchesPayload())

	// Mismatched topic and payload
	clientEnvelope, err = NewClientEnvelope(&envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_GroupMessage{},
		Aad:     buildAad(topic.NewTopic(topic.TopicKindKeyPackagesV1, []byte{1, 2, 3})),
	})
	require.NoError(t, err)
	require.False(t, clientEnvelope.TopicMatchesPayload())

	// Payer Report envelope with wrong topic
	clientEnvelope, err = NewClientEnvelope(&envelopesProto.ClientEnvelope{
		Payload: &envelopesProto.ClientEnvelope_PayerReport{},
		Aad:     buildAad(topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte{1, 2, 3})),
	})
	require.NoError(t, err)
	require.False(t, clientEnvelope.TopicMatchesPayload())
}

func TestRecoverSigner(t *testing.T) {
	nodeID := envelopeTestUtils.DefaultClientEnvelopeNodeID
	payerPrivateKey := testutils.RandomPrivateKey(t)
	rawPayerEnv := envelopeTestUtils.CreatePayerEnvelope(t, nodeID)

	payerSignature, err := utils.SignClientEnvelope(
		nodeID,
		rawPayerEnv.UnsignedClientEnvelope,
		payerPrivateKey,
	)
	require.NoError(t, err)
	rawPayerEnv.PayerSignature = &associations.RecoverableEcdsaSignature{
		Bytes: payerSignature,
	}

	payerEnv, err := NewPayerEnvelope(rawPayerEnv)
	require.NoError(t, err)

	signer, err := payerEnv.RecoverSigner()
	require.NoError(t, err)
	require.Equal(t, ethcrypto.PubkeyToAddress(payerPrivateKey.PublicKey).Hex(), signer.Hex())

	// Now test with an incorrect signature
	wrongPayerSignature, err := utils.SignClientEnvelope(
		nodeID,
		testutils.RandomBytes(128),
		payerPrivateKey,
	)
	require.NoError(t, err)
	rawPayerEnv.PayerSignature = &associations.RecoverableEcdsaSignature{
		Bytes: wrongPayerSignature,
	}
	payerEnv, err = NewPayerEnvelope(rawPayerEnv)
	require.NoError(t, err)

	// This will recover an incorrect signer address because the inputs to the signature
	// do not match the unsigned client envelope
	newSigner, err := payerEnv.RecoverSigner()
	require.NoError(t, err)
	require.NotEqual(t, signer.Hex(), newSigner.Hex())
}
