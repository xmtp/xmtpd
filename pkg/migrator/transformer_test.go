package migrator_test

import (
	"crypto/ecdsa"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
)

var (
	payerPrivateKey    *ecdsa.PrivateKey
	payerAddress       string
	nodeSigningKey     *ecdsa.PrivateKey
	nodeSigningAddress string
)

func newTestTransformer(t *testing.T) *migrator.Transformer {
	payerPrivateKey = testutils.RandomPrivateKey(t)
	nodeSigningKey = testutils.RandomPrivateKey(t)
	payerAddress = crypto.PubkeyToAddress(payerPrivateKey.PublicKey).Hex()
	nodeSigningAddress = crypto.PubkeyToAddress(nodeSigningKey.PublicKey).Hex()

	transformer := migrator.NewTransformer(
		payerPrivateKey,
		nodeSigningKey,
	)

	return transformer
}

func TestTransformGroupMessage(t *testing.T) {
	var (
		ctx         = t.Context()
		db, cleanup = testdata.NewTestDB(t, ctx)
		reader      = migrator.NewGroupMessageReader(db)
		transformer = newTestTransformer(t)
	)

	defer cleanup()

	records, _, err := reader.Fetch(ctx, 0, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.IsType(t, &migrator.GroupMessage{}, records[0])

	migratedGroupMessage, ok := records[0].(*migrator.GroupMessage)
	require.True(t, ok)

	envelope, err := transformer.Transform(migratedGroupMessage)
	require.NoError(t, err)
	require.NotNil(t, envelope)

	// OriginatorEnvelope check: Target topic has to be equal to TOPIC_KIND_GROUP_MESSAGES_V1 and the groupID.
	checkTopic(
		t,
		envelope,
		topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, migratedGroupMessage.GroupID[:]),
	)

	// OriginatorEnvelope check: Originator ID has to be hardcoded with GroupMessageOriginatorID.
	require.Equal(t, migrator.GroupMessageOriginatorID, envelope.OriginatorNodeID())

	// OriginatorEnvelope check: Sequence ID has to be the ID of the record.
	require.Equal(t, uint64(migratedGroupMessage.ID), envelope.OriginatorSequenceID())

	// OriginatorEnvelope check: Payload checks.
	payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()
	require.NotNil(t, payload)
	require.IsType(t, &proto.ClientEnvelope_GroupMessage{}, payload)

	groupMessagePayload := payload.(*proto.ClientEnvelope_GroupMessage)
	require.NotNil(t, groupMessagePayload.GroupMessage)
	require.IsType(t, &mlsv1.GroupMessageInput{
		Version: &mlsv1.GroupMessageInput_V1_{
			V1: &mlsv1.GroupMessageInput_V1{
				Data: migratedGroupMessage.Data,
			},
		},
	}, groupMessagePayload.GroupMessage)
	require.Equal(t, migratedGroupMessage.Data, groupMessagePayload.GroupMessage.GetV1().GetData())

	// Payer checks: expiration. Should not expire.
	require.Equal(
		t,
		uint32(math.MaxUint32),
		envelope.UnsignedOriginatorEnvelope.PayerEnvelope.RetentionDays(),
	)

	// Originator node checks: fees.
	require.Equal(t, uint64(0), envelope.UnsignedOriginatorEnvelope.Proto().BaseFeePicodollars)
	require.Equal(
		t,
		uint64(0),
		envelope.UnsignedOriginatorEnvelope.Proto().CongestionFeePicodollars,
	)

	// Signature checks.
	checkPayerSignature(t, envelope)
	checkOriginatorSignature(t, envelope)
}

func checkTopic(
	t *testing.T,
	envelope *envelopes.OriginatorEnvelope,
	expected *topic.Topic,
) {
	require.Equal(
		t,
		expected.Identifier(),
		envelope.TargetTopic().Identifier(),
	)

	require.Equal(t, expected.Kind(), envelope.TargetTopic().Kind())

	require.True(
		t,
		envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.TopicMatchesPayload(),
	)
}

func checkPayerSignature(t *testing.T, env *envelopes.OriginatorEnvelope) {
	// Can recover the payer signature.
	payerSignature := env.UnsignedOriginatorEnvelope.PayerEnvelope.Proto().GetPayerSignature()
	require.NotNil(t, payerSignature)

	// Can recover the payer signer.
	payerSigner, err := env.UnsignedOriginatorEnvelope.PayerEnvelope.RecoverSigner()
	require.NoError(t, err)
	require.Equal(
		t,
		payerAddress,
		payerSigner.Hex(),
	)
}

func checkOriginatorSignature(t *testing.T, env *envelopes.OriginatorEnvelope) {
	// Can recover the originator signature.
	recoveredSignature := env.Proto().GetOriginatorSignature()
	require.NotNil(t, recoveredSignature)

	// Can recover the unsigned envelope and sign it with the same node signing key.
	unsignedOriginatorEnvelopeBytes := env.Proto().GetUnsignedOriginatorEnvelope()
	require.NotNil(t, unsignedOriginatorEnvelopeBytes)

	hash := utils.HashOriginatorSignatureInput(unsignedOriginatorEnvelopeBytes)
	generatedSignature, err := crypto.Sign(
		hash,
		nodeSigningKey,
	)
	require.NoError(t, err)

	// Both signatures (recovered and generated) are the same.
	require.Equal(t, recoveredSignature.Bytes, generatedSignature)

	// Both addresses are the same.
	publicKey, err := crypto.SigToPub(hash, recoveredSignature.Bytes)
	require.NoError(t, err)
	require.Equal(t, nodeSigningAddress, crypto.PubkeyToAddress(*publicKey).Hex())
}
