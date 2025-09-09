package migrator_test

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/migrator"
	"github.com/xmtp/xmtpd/pkg/migrator/testdata"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	protobuf "google.golang.org/protobuf/proto"
)

type transformerTest struct {
	ctx             context.Context
	cleanup         func()
	db              *sql.DB
	transformer     *migrator.Transformer
	payerPrivateKey *ecdsa.PrivateKey
	nodePrivateKey  *ecdsa.PrivateKey
	payerAddress    string
	nodeAddress     string
}

func newTransformerTest(t *testing.T) *transformerTest {
	var (
		ctx             = t.Context()
		db, _, cleanup  = testdata.NewMigratorTestDB(t, ctx)
		payerPrivateKey = testutils.RandomPrivateKey(t)
		nodePrivateKey  = testutils.RandomPrivateKey(t)
		payerAddress    = crypto.PubkeyToAddress(payerPrivateKey.PublicKey).Hex()
		nodeAddress     = crypto.PubkeyToAddress(nodePrivateKey.PublicKey).Hex()
	)

	transformer := migrator.NewTransformer(
		payerPrivateKey,
		nodePrivateKey,
	)

	return &transformerTest{
		ctx:             ctx,
		cleanup:         cleanup,
		db:              db,
		transformer:     transformer,
		payerPrivateKey: payerPrivateKey,
		nodePrivateKey:  nodePrivateKey,
		payerAddress:    payerAddress,
		nodeAddress:     nodeAddress,
	}
}

func TestTransformGroupMessage(t *testing.T) {
	var (
		test   = newTransformerTest(t)
		reader = migrator.NewGroupMessageReader(test.db)
	)

	defer test.cleanup()

	records, err := reader.Fetch(test.ctx, 0, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.IsType(t, &migrator.GroupMessage{}, records[0])

	migratedGroupMessage, ok := records[0].(*migrator.GroupMessage)
	require.True(t, ok)

	envelope, err := test.transformer.Transform(migratedGroupMessage)
	require.NoError(t, err)
	require.NotNil(t, envelope)

	// OriginatorEnvelope check: Target topic has to be equal to TOPIC_KIND_GROUP_MESSAGES_V1 and the groupID.
	checkTopic(
		t,
		envelope,
		topic.NewTopic(topic.TopicKindGroupMessagesV1, migratedGroupMessage.GroupID[:]),
	)

	// OriginatorEnvelope check: Originator ID has to be hardcoded with GroupMessageOriginatorID.
	require.Equal(t, migrator.GroupMessageOriginatorID, envelope.OriginatorNodeID())

	// OriginatorEnvelope check: Sequence ID has to be the ID of the record.
	require.Equal(t, uint64(migratedGroupMessage.ID), envelope.OriginatorSequenceID())

	// OriginatorEnvelope check: Payload checks.
	payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()
	require.NotNil(t, payload)
	require.IsType(t, &proto.ClientEnvelope_GroupMessage{}, payload)

	groupMessagePayload, ok := payload.(*proto.ClientEnvelope_GroupMessage)
	require.True(t, ok)
	require.NotNil(t, groupMessagePayload.GroupMessage)
	require.IsType(t, &mlsv1.GroupMessageInput{}, groupMessagePayload.GroupMessage)

	groupMessageV1 := groupMessagePayload.GroupMessage.GetV1()
	require.NotNil(t, groupMessageV1)
	require.Equal(t, migratedGroupMessage.Data, groupMessageV1.GetData())

	// Payer checks: expiration. Should not expire.
	require.Equal(
		t,
		uint32(math.MaxInt32),
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
	checkPayerSignature(t, envelope, test.payerAddress)
	checkOriginatorSignature(t, envelope, test.nodePrivateKey, test.nodeAddress)
}

func TestTransformInboxLog(t *testing.T) {
	var (
		test   = newTransformerTest(t)
		reader = migrator.NewInboxLogReader(test.db)
	)

	defer test.cleanup()

	records, err := reader.Fetch(test.ctx, 0, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.IsType(t, &migrator.InboxLog{}, records[0])

	migratedInboxLog, ok := records[0].(*migrator.InboxLog)
	require.True(t, ok)

	envelope, err := test.transformer.Transform(migratedInboxLog)
	require.NoError(t, err)
	require.NotNil(t, envelope)

	// OriginatorEnvelope check: Target topic has to be equal to TOPIC_KIND_GROUP_MESSAGES_V1 and the groupID.
	checkTopic(
		t,
		envelope,
		topic.NewTopic(topic.TopicKindIdentityUpdatesV1, migratedInboxLog.InboxID[:]),
	)

	// OriginatorEnvelope check: Originator ID has to be hardcoded with GroupMessageOriginatorID.
	require.Equal(t, migrator.InboxLogOriginatorID, envelope.OriginatorNodeID())

	// OriginatorEnvelope check: Sequence ID has to be the ID of the record.
	require.Equal(t, uint64(migratedInboxLog.SequenceID), envelope.OriginatorSequenceID())

	// OriginatorEnvelope check: Payload checks.
	payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()
	require.NotNil(t, payload)
	require.IsType(t, &proto.ClientEnvelope_IdentityUpdate{}, payload)

	identityUpdatePayload, ok := payload.(*proto.ClientEnvelope_IdentityUpdate)
	require.True(t, ok)
	require.NotNil(t, identityUpdatePayload.IdentityUpdate)
	require.IsType(t, &associations.IdentityUpdate{}, identityUpdatePayload.IdentityUpdate)

	require.Equal(
		t,
		hex.EncodeToString(migratedInboxLog.InboxID[:]),
		identityUpdatePayload.IdentityUpdate.InboxId,
	)

	migratedIdentityUpdateProto := &associations.IdentityUpdate{}
	err = protobuf.Unmarshal(migratedInboxLog.IdentityUpdateProto, migratedIdentityUpdateProto)
	require.NoError(t, err)

	require.True(
		t,
		protobuf.Equal(migratedIdentityUpdateProto, identityUpdatePayload.IdentityUpdate),
	)

	// Payer checks: expiration. Should not expire.
	require.Equal(
		t,
		uint32(math.MaxInt32),
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
	checkPayerSignature(t, envelope, test.payerAddress)
	checkOriginatorSignature(t, envelope, test.nodePrivateKey, test.nodeAddress)
}

func TestTransformKeyPackage(t *testing.T) {
	var (
		test   = newTransformerTest(t)
		reader = migrator.NewKeyPackageReader(test.db)
	)

	defer test.cleanup()

	records, err := reader.Fetch(test.ctx, 0, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.IsType(t, &migrator.KeyPackage{}, records[0])

	migratedInstallation, ok := records[0].(*migrator.KeyPackage)
	require.True(t, ok)

	envelope, err := test.transformer.Transform(migratedInstallation)
	require.NoError(t, err)
	require.NotNil(t, envelope)

	// OriginatorEnvelope check: Target topic has to be equal to TOPIC_KIND_GROUP_MESSAGES_V1 and the groupID.
	checkTopic(
		t,
		envelope,
		topic.NewTopic(topic.TopicKindKeyPackagesV1, migratedInstallation.InstallationID[:]),
	)

	// OriginatorEnvelope check: Originator ID has to be hardcoded with GroupMessageOriginatorID.
	require.Equal(t, migrator.KeyPackagesOriginatorID, envelope.OriginatorNodeID())

	// OriginatorEnvelope check: Sequence ID has to be the ID of the record.
	require.Equal(t, uint64(migratedInstallation.SequenceID), envelope.OriginatorSequenceID())

	// OriginatorEnvelope check: Payload checks.
	payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()
	require.NotNil(t, payload)
	require.IsType(t, &proto.ClientEnvelope_UploadKeyPackage{}, payload)

	uploadKeyPackagePayload, ok := payload.(*proto.ClientEnvelope_UploadKeyPackage)
	require.True(t, ok)
	require.NotNil(t, uploadKeyPackagePayload.UploadKeyPackage)
	require.IsType(t, &mlsv1.UploadKeyPackageRequest{}, uploadKeyPackagePayload.UploadKeyPackage)

	require.Equal(
		t,
		uploadKeyPackagePayload.UploadKeyPackage.KeyPackage.GetKeyPackageTlsSerialized(),
		migratedInstallation.KeyPackage,
	)

	// Payer checks: expiration.
	require.Equal(
		t,
		uint32(constants.DefaultStorageDurationDays),
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
	checkPayerSignature(t, envelope, test.payerAddress)
	checkOriginatorSignature(t, envelope, test.nodePrivateKey, test.nodeAddress)
}

func TestTransformWelcomeMessage(t *testing.T) {
	var (
		test   = newTransformerTest(t)
		reader = migrator.NewWelcomeMessageReader(test.db)
	)

	defer test.cleanup()

	records, err := reader.Fetch(test.ctx, 0, 1)
	require.NoError(t, err)
	require.Len(t, records, 1)
	require.IsType(t, &migrator.WelcomeMessage{}, records[0])

	migratedWelcomeMessage, ok := records[0].(*migrator.WelcomeMessage)
	require.True(t, ok)

	envelope, err := test.transformer.Transform(migratedWelcomeMessage)
	require.NoError(t, err)
	require.NotNil(t, envelope)

	// OriginatorEnvelope check: Target topic has to be equal to TOPIC_KIND_GROUP_MESSAGES_V1 and the groupID.
	checkTopic(
		t,
		envelope,
		topic.NewTopic(
			topic.TopicKindWelcomeMessagesV1,
			migratedWelcomeMessage.InstallationKey[:],
		),
	)

	// OriginatorEnvelope check: Originator ID has to be hardcoded with GroupMessageOriginatorID.
	require.Equal(t, migrator.WelcomeMessageOriginatorID, envelope.OriginatorNodeID())

	// OriginatorEnvelope check: Sequence ID has to be the ID of the record.
	require.Equal(t, uint64(migratedWelcomeMessage.ID), envelope.OriginatorSequenceID())

	// OriginatorEnvelope check: Payload checks.
	payload := envelope.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Payload()
	require.NotNil(t, payload)
	require.IsType(t, &proto.ClientEnvelope_WelcomeMessage{}, payload)

	welcomeMessagePayload, ok := payload.(*proto.ClientEnvelope_WelcomeMessage)
	require.True(t, ok)
	require.NotNil(t, welcomeMessagePayload.WelcomeMessage)
	require.IsType(t, &mlsv1.WelcomeMessageInput{}, welcomeMessagePayload.WelcomeMessage)

	welcomeMessageV1 := welcomeMessagePayload.WelcomeMessage.GetV1()
	require.NotNil(t, welcomeMessageV1)
	require.Equal(t, migratedWelcomeMessage.InstallationKey, welcomeMessageV1.InstallationKey)
	require.Equal(t, migratedWelcomeMessage.Data, welcomeMessageV1.Data)
	require.Equal(t, migratedWelcomeMessage.HpkePublicKey, welcomeMessageV1.HpkePublicKey)
	require.Equal(
		t,
		migratedWelcomeMessage.WrapperAlgorithm,
		int16(welcomeMessageV1.WrapperAlgorithm),
	)

	// Payer checks: expiration.
	require.Equal(
		t,
		uint32(constants.DefaultStorageDurationDays),
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
	checkPayerSignature(t, envelope, test.payerAddress)
	checkOriginatorSignature(t, envelope, test.nodePrivateKey, test.nodeAddress)
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

func checkPayerSignature(t *testing.T, env *envelopes.OriginatorEnvelope, payerAddress string) {
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

func checkOriginatorSignature(
	t *testing.T,
	env *envelopes.OriginatorEnvelope,
	nodeSigningKey *ecdsa.PrivateKey,
	nodeAddress string,
) {
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
	require.Equal(t, nodeAddress, crypto.PubkeyToAddress(*publicKey).Hex())
}
