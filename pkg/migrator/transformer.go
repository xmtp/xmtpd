package migrator

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/currency"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/fees"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	mlsv1 "github.com/xmtp/xmtpd/pkg/proto/mls/api/v1"
	messageContents "github.com/xmtp/xmtpd/pkg/proto/mls/message_contents"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	protobuf "google.golang.org/protobuf/proto"
)

type Transformer struct {
	feeCalculator   fees.IFeeCalculator
	payerPrivateKey *ecdsa.PrivateKey
	nodeSigningKey  *ecdsa.PrivateKey
}

func NewTransformer(
	feeCalculator fees.IFeeCalculator,
	payerPrivateKey *ecdsa.PrivateKey,
	nodeSigningKey *ecdsa.PrivateKey,
) *Transformer {
	return &Transformer{
		feeCalculator:   feeCalculator,
		payerPrivateKey: payerPrivateKey,
		nodeSigningKey:  nodeSigningKey,
	}
}

func (t *Transformer) Transform(record ISourceRecord) (*envelopes.OriginatorEnvelope, error) {
	switch record.TableName() {
	case groupMessagesTableName:
		data, ok := record.(*GroupMessage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformGroupMessage(data)

	case commitMessagesTableName:
		data, ok := record.(*CommitMessage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformCommitMessage(data)

	case inboxLogTableName:
		data, ok := record.(*InboxLog)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformInboxLog(data)

	case keyPackagesTableName:
		data, ok := record.(*KeyPackage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformKeyPackage(data)

	case welcomeMessagesTableName:
		data, ok := record.(*WelcomeMessage)
		if !ok {
			return nil, fmt.Errorf("invalid record type: %T", record)
		}

		return t.TransformWelcomeMessage(data)

	default:
		return nil, fmt.Errorf(
			"transform not implemented for table '%s'",
			record.TableName(),
		)
	}
}

// TransformGroupMessage converts GroupMessage to appropriate XMTPD envelope format.
func (t *Transformer) TransformGroupMessage(
	groupMessage *GroupMessage,
) (*envelopes.OriginatorEnvelope, error) {
	protoClientEnvelope, err := transformGroupMessage(groupMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to transform group message: %w", err)
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		GroupMessageOriginatorID,
		uint64(groupMessage.ID),
		groupMessage.CreatedAt,
	)
}

// TransformCommitMessage converts CommitMessage to appropriate XMTPD envelope format.
func (t *Transformer) TransformCommitMessage(
	commitMessage *CommitMessage,
) (*envelopes.OriginatorEnvelope, error) {
	protoClientEnvelope, err := transformGroupMessage(&commitMessage.GroupMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to transform group message: %w", err)
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		CommitMessageOriginatorID,
		uint64(commitMessage.ID),
		commitMessage.CreatedAt,
	)
}

// TransformInboxLog converts InboxLog to appropriate XMTPD IdentityUpdate envelope format.
func (t *Transformer) TransformInboxLog(
	inboxLog *InboxLog,
) (*envelopes.OriginatorEnvelope, error) {
	if inboxLog == nil {
		return nil, fmt.Errorf("inboxLog is nil")
	}

	if inboxLog.InboxID == nil {
		return nil, fmt.Errorf("inboxID is nil")
	}

	_, err := utils.ParseInboxID(inboxLog.InboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inbox ID: %w", err)
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
			TargetTopic: topic.NewTopic(topic.TopicKindIdentityUpdatesV1, inboxLog.InboxID[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		InboxLogOriginatorID,
		uint64(inboxLog.SequenceID),
		time.Unix(0, inboxLog.ServerTimestampNs),
	)
}

// TransformKeyPackage converts a V3 KeyPackage to appropriate XMTPD KeyPackage envelope format.
func (t *Transformer) TransformKeyPackage(
	keyPackage *KeyPackage,
) (*envelopes.OriginatorEnvelope, error) {
	if keyPackage == nil {
		return nil, fmt.Errorf("keyPackage is nil")
	}

	if len(keyPackage.KeyPackage) <= 0 {
		return nil, fmt.Errorf("keyPackage is empty")
	}

	protoClientEnvelope := &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_UploadKeyPackage{
			UploadKeyPackage: &mlsv1.UploadKeyPackageRequest{
				KeyPackage: &mlsv1.KeyPackageUpload{
					KeyPackageTlsSerialized: keyPackage.KeyPackage,
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindKeyPackagesV1, keyPackage.InstallationID[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		KeyPackagesOriginatorID,
		uint64(keyPackage.SequenceID),
		time.Unix(0, keyPackage.CreatedAt),
	)
}

// TransformWelcomeMessage converts WelcomeMessage to appropriate XMTPD envelope format.
func (t *Transformer) TransformWelcomeMessage(
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
						WelcomeMetadata: welcomeMessage.WelcomeMetadata,
					},
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindWelcomeMessagesV1, welcomeMessage.InstallationKey[:]).
				Bytes(),
		},
	}

	return t.originatorEnvelope(
		protoClientEnvelope,
		WelcomeMessageOriginatorID,
		uint64(welcomeMessage.ID),
		welcomeMessage.CreatedAt,
	)
}

// originatorEnvelope builds and signs an originator envelope from a client envelope.
func (t *Transformer) originatorEnvelope(
	protoClientEnvelope *proto.ClientEnvelope,
	originatorID uint32,
	sequenceID uint64,
	creationTime time.Time,
) (*envelopes.OriginatorEnvelope, error) {
	if protoClientEnvelope == nil {
		return nil, fmt.Errorf("protoClientEnvelope is nil")
	}

	payerEnvelope, err := t.buildAndSignPayerEnvelope(protoClientEnvelope, originatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign payer envelope: %w", err)
	}

	originatorEnvelope, err := t.buildAndSignOriginatorEnvelope(
		payerEnvelope,
		sequenceID,
		creationTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build and sign originator envelope: %w", err)
	}

	return originatorEnvelope, nil
}

func (t *Transformer) buildAndSignPayerEnvelope(
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

	// Lower than MaxUint32 to avoid overflow.
	retentionDays := uint32(math.MaxInt32)

	if isDatabaseDestination(originatorID) {
		retentionDays = constants.DefaultStorageDurationDays
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

func (t *Transformer) buildAndSignOriginatorEnvelope(
	payerEnvelope *envelopes.PayerEnvelope,
	sequenceID uint64,
	creationTime time.Time,
) (*envelopes.OriginatorEnvelope, error) {
	if payerEnvelope == nil {
		return nil, fmt.Errorf("payerEnvelope is nil")
	}

	payerEnvelopeBytes, err := payerEnvelope.Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get payer envelope bytes: %w", err)
	}

	var (
		now     = time.Now()
		baseFee currency.PicoDollar
	)

	// WARNING: we are doing some time hackery here
	// the expiration is calculated from the original creation date of the V3 payload
	// but fees are calculated based on the migration date

	if isDatabaseDestination(payerEnvelope.TargetOriginator) {
		baseFee, err = t.calculateFees(
			now,
			int64(len(payerEnvelopeBytes)),
			payerEnvelope.Proto().GetMessageRetentionDays(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate fees: %w", err)
		}
	}

	unsignedEnv := proto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:         payerEnvelope.TargetOriginator,
		OriginatorSequenceId:     sequenceID,
		OriginatorNs:             now.UnixNano(),
		PayerEnvelopeBytes:       payerEnvelopeBytes,
		BaseFeePicodollars:       uint64(baseFee),
		CongestionFeePicodollars: 0, // Migrator does not pay congestion fees.
		ExpiryUnixtime: uint64(
			creationTime.UTC().
				Add(time.Hour * 24 * time.Duration(payerEnvelope.Proto().GetMessageRetentionDays())).
				Unix(),
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

// transformGroupMessage transforms a GroupMessage (commit or not) to a ClientEnvelope.
func transformGroupMessage(groupMessage *GroupMessage) (*proto.ClientEnvelope, error) {
	if groupMessage == nil {
		return nil, fmt.Errorf("groupMessage is nil")
	}

	if groupMessage.GroupID == nil {
		return nil, fmt.Errorf("groupID is nil")
	}

	_, err := utils.ParseGroupID(groupMessage.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse group ID: %w", err)
	}

	if len(groupMessage.Data) <= 0 {
		return nil, fmt.Errorf("data is empty")
	}

	return &proto.ClientEnvelope{
		Payload: &proto.ClientEnvelope_GroupMessage{
			GroupMessage: &mlsv1.GroupMessageInput{
				Version: &mlsv1.GroupMessageInput_V1_{
					V1: &mlsv1.GroupMessageInput_V1{
						Data:       groupMessage.Data,
						SenderHmac: groupMessage.SenderHmac,
						ShouldPush: groupMessage.ShouldPush.Bool,
					},
				},
			},
		},
		Aad: &proto.AuthenticatedData{
			TargetTopic: topic.NewTopic(topic.TopicKindGroupMessagesV1, groupMessage.GroupID[:]).
				Bytes(),
		},
	}, nil
}

func (t *Transformer) calculateFees(
	originatorTime time.Time,
	envelopeLength int64,
	retentionDays uint32,
) (currency.PicoDollar, error) {
	baseFee, err := t.feeCalculator.CalculateBaseFee(
		originatorTime,
		envelopeLength,
		retentionDays,
	)
	if err != nil {
		return 0, err
	}

	// Migrator does not pay congestion fees.
	return baseFee, nil
}
