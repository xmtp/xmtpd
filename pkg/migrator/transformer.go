package migrator

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	messageContents "github.com/xmtp/xmtpd/pkg/proto/mls/message_contents"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	protobuf "google.golang.org/protobuf/proto"
)

// TODO: Validate payer envelope.
// TODO: Validate key package.

type transformer struct {
	payerPrivateKey *ecdsa.PrivateKey
	nodeSigningKey  *ecdsa.PrivateKey
}

func NewTransformer(
	payerPrivateKey *ecdsa.PrivateKey,
	nodeSigningKey *ecdsa.PrivateKey,
) *transformer {
	return &transformer{
		payerPrivateKey: payerPrivateKey,
		nodeSigningKey:  nodeSigningKey,
	}
}

func (t *transformer) Transform(record ISourceRecord) (*envelopes.OriginatorEnvelope, error) {
	switch record.TableName() {
	case groupMessagesTableName:
		data, ok := record.(*GroupMessage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformGroupMessage(data)

	case inboxLogTableName:
		data, ok := record.(*InboxLog)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformInboxLog(data)

	case installationsTableName:
		data, ok := record.(*Installation)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformInstallation(data)

	case welcomeMessagesTableName:
		data, ok := record.(*WelcomeMessage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformWelcomeMessage(data)

	default:
		return nil, fmt.Errorf(
			"Transform not implemented for table: %s",
			record.TableName(),
		)
	}
}

// TransformGroupMessage converts GroupMessage to appropriate XMTPD envelope format.
func (t *transformer) TransformGroupMessage(
	groupMessage *GroupMessage,
) (*envelopes.OriginatorEnvelope, error) {
	if groupMessage == nil {
		return nil, fmt.Errorf("groupMessage is nil")
	}

	if groupMessage.GroupID == nil {
		return nil, fmt.Errorf("groupID is nil")
	}

	if len(groupMessage.Data) <= 0 {
		return nil, fmt.Errorf("data is empty")
	}

	protoClientEnvelope := &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data: groupMessage.Data,
					},
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, groupMessage.GroupID[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		groupMessageOriginatorID,
		uint64(groupMessage.ID),
	)
}

// TransformInboxLog converts InboxLog to appropriate XMTPD IdentityUpdate envelope format.
func (t *transformer) TransformInboxLog(
	inboxLog *InboxLog,
) (*envelopes.OriginatorEnvelope, error) {
	if inboxLog == nil {
		return nil, fmt.Errorf("inboxLog is nil")
	}

	if inboxLog.InboxID == nil {
		return nil, fmt.Errorf("inboxID is nil")
	}

	if len(inboxLog.IdentityUpdateProto) <= 0 {
		return nil, fmt.Errorf("identityUpdateProto is empty")
	}

	var identityUpdateProto associations.IdentityUpdate

	if err := protobuf.Unmarshal(inboxLog.IdentityUpdateProto, &identityUpdateProto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IdentityUpdateProto: %w", err)
	}

	// Is identityUpdateProto everything we need?
	protoClientEnvelope := &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_IdentityUpdate{
			IdentityUpdate: &identityUpdateProto,
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_IDENTITY_UPDATES_V1, inboxLog.InboxID[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		inboxLogOriginatorID,
		uint64(inboxLog.SequenceID),
	)
}

// TransformInstallation converts Installation to appropriate XMTPD KeyPackage envelope format.
func (t *transformer) TransformInstallation(
	installation *Installation,
) (*envelopes.OriginatorEnvelope, error) {
	if installation == nil {
		return nil, fmt.Errorf("installation is nil")
	}

	if len(installation.KeyPackage) <= 0 {
		return nil, fmt.Errorf("keyPackage is empty")
	}

	protoClientEnvelope := &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_UploadKeyPackage{
			UploadKeyPackage: &mlsv1.UploadKeyPackageRequest{
				KeyPackage: &mlsv1.KeyPackageUpload{
					KeyPackageTlsSerialized: installation.KeyPackage,
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_KEY_PACKAGES_V1, installation.ID[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		installationOriginatorID,
		uint64(installation.CreatedAt),
	)
}

// TransformWelcomeMessage converts WelcomeMessage to appropriate XMTPD envelope format.
func (t *transformer) TransformWelcomeMessage(
	welcomeMessage *WelcomeMessage,
) (*envelopes.OriginatorEnvelope, error) {
	if welcomeMessage == nil {
		return nil, fmt.Errorf("welcomeMessage is nil")
	}

	protoClientEnvelope := &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_WelcomeMessage{
			WelcomeMessage: &mlsv1.WelcomeMessageInput{
				Version: &mlsv1.WelcomeMessageInput_V1_{
					V1: &mlsv1.WelcomeMessageInput_V1{
						InstallationKey: welcomeMessage.InstallationKey,
						Data:            welcomeMessage.Data,
						HpkePublicKey:   welcomeMessage.HpkePublicKey,
						WrapperAlgorithm: messageContents.WelcomeWrapperAlgorithm(
							welcomeMessage.WrapperAlgorithm,
						),
					},
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TOPIC_KIND_WELCOME_MESSAGES_V1, welcomeMessage.InstallationKey[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		welcomeMessageOriginatorID,
		uint64(welcomeMessage.ID),
	)
}

// originatorEnvelope builds and signs an originator envelope from a client envelope.
func (t *transformer) originatorEnvelope(
	protoClientEnvelope *proto.ClientEnvelope,
	originatorID uint32,
	sequenceID uint64,
) (*envelopes.OriginatorEnvelope, error) {
	if protoClientEnvelope == nil {
		return nil, fmt.Errorf("protoClientEnvelope is nil")
	}

	payerEnvelope, err := t.buildAndSignPayerEnvelope(protoClientEnvelope, originatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign payer envelope: %w", err)
	}

	originatorEnvelope, err := t.buildAndSignOriginatorEnvelope(payerEnvelope, sequenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign originator envelope: %w", err)
	}

	return originatorEnvelope, nil
}

func (t *transformer) buildAndSignPayerEnvelope(
	protoClientEnvelope *proto.ClientEnvelope,
	originatorID uint32,
) (*envelopes.PayerEnvelope, error) {
	if !isValidOriginatorID(originatorID) {
		return nil, fmt.Errorf("invalid originatorID: %d", originatorID)
	}

	if protoClientEnvelope == nil {
		return nil, fmt.Errorf("protoClientEnvelope is nil")
	}

	clientEnvelope, err := envelopes.NewClientEnvelope(protoClientEnvelope)
	if err != nil {
		return nil, fmt.Errorf("failed to build client envelope: %w", err)
	}

	clientEnvelopeBytes, err := clientEnvelope.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get client envelope bytes: %w", err)
	}

	payerSignature, err := utils.SignClientEnvelope(
		originatorID,
		clientEnvelopeBytes,
		t.payerPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign client envelope: %w", err)
	}

	retentionDays := uint32(math.MaxUint32)

	// group messages and identity updates are stored forever (MaxUint32).
	// welcome messages and installations are stored for the default duration.
	if originatorID == installationOriginatorID || originatorID == welcomeMessageOriginatorID {
		retentionDays = constants.DEFAULT_STORAGE_DURATION_DAYS
	}

	protoPayerEnvelope := &proto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvelopeBytes,
		PayerSignature: &associations.RecoverableEcdsaSignature{
			Bytes: payerSignature,
		},
		TargetOriginator:     originatorID,
		MessageRetentionDays: retentionDays,
	}

	return envelopes.NewPayerEnvelope(protoPayerEnvelope)
}

// TODO: Set congestion fee.
// TODO: Set base fee.

func (t *transformer) buildAndSignOriginatorEnvelope(
	payerEnvelope *envelopes.PayerEnvelope,
	sequenceID uint64,
) (*envelopes.OriginatorEnvelope, error) {
	if payerEnvelope == nil {
		return nil, fmt.Errorf("payerEnvelope is nil")
	}

	payerEnvelopeBytes, err := payerEnvelope.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get payer envelope bytes: %w", err)
	}

	unsignedEnv := proto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:         payerEnvelope.TargetOriginator,
		OriginatorSequenceId:     sequenceID,
		OriginatorNs:             time.Now().UnixNano(),
		PayerEnvelopeBytes:       payerEnvelopeBytes,
		BaseFeePicodollars:       0,
		CongestionFeePicodollars: 0,
		ExpiryUnixtime: uint64(
			time.Now().AddDate(0, 0, int(payerEnvelope.Proto().GetMessageRetentionDays())).Unix(),
		),
	}

	unsignedBytes, err := protobuf.Marshal(&unsignedEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal unsigned originator envelope: %w", err)
	}

	sig, err := crypto.Sign(utils.HashOriginatorSignatureInput(unsignedBytes), t.nodeSigningKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign originator envelope: %w", err)
	}

	protoOriginatorEnvelope := &proto.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof: &proto.OriginatorEnvelope_OriginatorSignature{
			OriginatorSignature: &associations.RecoverableEcdsaSignature{
				Bytes: sig,
			},
		},
	}

	return envelopes.NewOriginatorEnvelope(protoOriginatorEnvelope)
}
