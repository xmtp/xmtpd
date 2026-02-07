package contracts

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/constants"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	// https://github.com/xmtp/libxmtp/blob/main/crates/xmtp_id/src/associations/association_log.rs#L10
	associationErrorGeneric                      = "Error creating association"
	associationErrorMultipleCreate               = "Multiple create operations detected"
	associationErrorNotCreated                   = "XID not yet created"
	associationErrorSignature                    = "Signature validation failed"
	associationErrorMemberNotAllowed             = "not allowed to add"
	associationErrorMissingExistingMember        = "Missing existing member"
	associationErrorLegacySignatureReuse         = "Legacy key is only allowed to be associated using a legacy signature with nonce 0"
	associationErrorNewMemberIDSignatureMismatch = "The new member identifier does not match the signer"
	associationErrorWrongInboxID                 = "Wrong inbox_id specified on association"
	associationErrorSignatureNotAllowed          = "Signature not allowed for role"
	associationErrorReplay                       = "Replay detected"
	associationErrorDeserialization              = "Deserialization error"
	associationErrorMissingIdentityUpdate        = "Missing identity update"
	associationErrorChainIDMismatch              = "Wrong chain id."
	associationErrorInvalidAccountAddress        = "Invalid account address: Must be 42 hex characters, starting with '0x'."
	associationErrorNotIdentifier                = "are not a public identifier"
	associationErrorConvert                      = "Conversion error"

	ErrAdvisoryLockSequence   = "advisory lock failed"
	ErrParseIdentityUpdate    = "error parsing identity update"
	ErrGetLatestSequenceID    = "get latest sequence id failed"
	ErrValidateIdentityUpdate = "validate identity update failed"
	ErrInsertAddressLog       = "insert address log failed"
	ErrRevokeAddressFromLog   = "revoke address from log failed"
)

var associationErrorPatterns = []string{
	associationErrorGeneric,
	associationErrorMultipleCreate,
	associationErrorNotCreated,
	associationErrorSignature,
	associationErrorMemberNotAllowed,
	associationErrorMissingExistingMember,
	associationErrorLegacySignatureReuse,
	associationErrorNewMemberIDSignatureMismatch,
	associationErrorWrongInboxID,
	associationErrorSignatureNotAllowed,
	associationErrorReplay,
	associationErrorDeserialization,
	associationErrorMissingIdentityUpdate,
	associationErrorChainIDMismatch,
	associationErrorInvalidAccountAddress,
	associationErrorNotIdentifier,
	associationErrorConvert,
}

type IdentityUpdateStorer struct {
	contract          *iu.IdentityUpdateBroadcaster
	db                *sql.DB
	logger            *zap.Logger
	validationService mlsvalidate.MLSValidationService
}

var _ c.ILogStorer = &IdentityUpdateStorer{}

func NewIdentityUpdateStorer(
	db *sql.DB,
	logger *zap.Logger,
	contract *iu.IdentityUpdateBroadcaster,
	validationService mlsvalidate.MLSValidationService,
) *IdentityUpdateStorer {
	return &IdentityUpdateStorer{
		db:                db,
		logger:            logger.Named(utils.StorerLoggerName),
		contract:          contract,
		validationService: validationService,
	}
}

// StoreLog validates and stores an identity update log event
func (s *IdentityUpdateStorer) StoreLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	msgSent, err := s.contract.ParseIdentityUpdateCreated(event)
	if err != nil {
		return re.NewNonRecoverableError(ErrParseIdentityUpdate, err)
	}

	// NOTE ON CONCURRENCY CONTROL:
	//
	// We intentionally run this transaction at READ COMMITTED and take a
	// transaction-scoped advisory lock (pg_advisory_xact_lock) keyed on
	// (originator_node_id).
	//
	// Why this works:
	//   • The advisory lock guarantees that only one HA worker can process the
	//     same event concurrently.
	//   • READ COMMITTED means each statement sees the latest committed state,
	//     so GetLatestSequenceId reflects any inserts that committed after we
	//     waited on the advisory lock. At REPEATABLE READ we would be stuck
	//     with the snapshot taken before the lock, which could cause duplicate
	//     processing or serialization errors.
	//   • Our INSERT/UPDATE statements on address_log are monotonic (only-if-newer
	//     guards), so stale writes become no-ops even if two different events
	//     touch the same (address,inbox_id).
	//
	// In short: the advisory lock prevents two workers from processing the same
	// event in parallel, and READ COMMITTED ensures that once we hold the lock
	// we see the latest database state. This combination avoids the
	// “could not serialize access due to concurrent update” errors that show
	// up under REPEATABLE READ with UPSERTs.

	err = db.RunInTx(
		ctx,
		s.db,
		&sql.TxOptions{Isolation: sql.LevelReadCommitted},
		func(ctx context.Context, querier *queries.Queries) error {
			err := db.NewAdvisoryLocker().
				LockIdentityUpdateInsert(ctx, querier, uint32(constants.IdentityUpdateOriginatorID))
			if err != nil {
				return re.NewNonRecoverableError(ErrAdvisoryLockSequence, err)
			}

			latestSequenceID, err := querier.GetLatestSequenceId(
				ctx,
				constants.IdentityUpdateOriginatorID,
			)
			if err != nil {
				return re.NewNonRecoverableError(ErrGetLatestSequenceID, err)
			}

			if uint64(latestSequenceID) >= msgSent.SequenceId {
				s.logger.Debug(
					"identity update already inserted, skipping",
					utils.LastSequenceIDField(latestSequenceID),
					utils.SequenceIDField(int64(msgSent.SequenceId)),
				)
				return nil
			}

			messageTopic := topic.NewTopic(topic.TopicKindIdentityUpdatesV1, msgSent.InboxId[:])

			if s.logger.Core().Enabled(zap.DebugLevel) {
				s.logger.Debug(
					"inserting identity update from contract",
					utils.TopicField(messageTopic.String()),
				)
			}

			clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Update)
			if err != nil {
				s.logger.Error(
					ErrParseClientEnvelope,
					utils.TopicField(messageTopic.String()),
					zap.Error(err),
				)
				return re.NewNonRecoverableError(ErrParseClientEnvelope, err)
			}

			associationState, validationError := s.validateIdentityUpdate(
				ctx,
				querier,
				msgSent.InboxId,
				clientEnvelope,
			)
			if validationError != nil {
				s.logger.Error(
					ErrValidateIdentityUpdate,
					utils.TopicField(messageTopic.String()),
					zap.Error(validationError),
				)

				return validationError
			}

			var (
				inboxID     = utils.HexEncode(msgSent.InboxId[:])
				sequenceID  = int64(msgSent.SequenceId)
				insertBatch = make([]string, 0)
				revokeBatch = make([]string, 0)
			)

			for _, newMember := range associationState.StateDiff.NewMembers {
				if s.logger.Core().Enabled(zap.DebugLevel) {
					s.logger.Debug("new member", utils.BodyField(newMember))
				}

				if address, ok := newMember.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
					insertBatch = append(insertBatch, address.EthereumAddress)
				}
			}

			for _, removedMember := range associationState.StateDiff.RemovedMembers {
				if s.logger.Core().Enabled(zap.DebugLevel) {
					s.logger.Debug("removed member", utils.BodyField(removedMember))
				}

				if address, ok := removedMember.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
					revokeBatch = append(revokeBatch, address.EthereumAddress)
				}
			}

			if len(insertBatch) > 0 {
				_, err := querier.InsertAddressLogsBatch(ctx, queries.InsertAddressLogsBatchParams{
					Addresses:             insertBatch,
					InboxID:               inboxID,
					AssociationSequenceID: sequenceID,
				})
				if err != nil {
					return re.NewRecoverableError(ErrInsertAddressLog, err)
				}
			}

			if len(revokeBatch) > 0 {
				_, err := querier.RevokeAddressFromLogBatch(
					ctx,
					queries.RevokeAddressFromLogBatchParams{
						Addresses:            revokeBatch,
						InboxID:              inboxID,
						RevocationSequenceID: sequenceID,
					},
				)
				if err != nil {
					return re.NewRecoverableError(ErrRevokeAddressFromLog, err)
				}
			}

			originatorEnvelope, err := buildOriginatorEnvelope(
				constants.IdentityUpdateOriginatorID,
				msgSent.SequenceId,
				msgSent.Update,
			)
			if err != nil {
				s.logger.Error(
					ErrBuildOriginatorEnvelope,
					utils.TopicField(messageTopic.String()),
					zap.Error(err),
				)
				return re.NewNonRecoverableError(ErrBuildOriginatorEnvelope, err)
			}

			signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
				originatorEnvelope,
				event.TxHash,
			)
			if err != nil {
				s.logger.Error(
					ErrBuildSignedOriginatorEnvelope,
					utils.TopicField(messageTopic.String()),
					zap.Error(err),
				)
				return re.NewNonRecoverableError(ErrBuildSignedOriginatorEnvelope, err)
			}

			originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
			if err != nil {
				s.logger.Error(
					ErrMarshallOriginatorEnvelope,
					utils.TopicField(messageTopic.String()),
					zap.Error(err),
				)
				return re.NewNonRecoverableError(ErrMarshallOriginatorEnvelope, err)
			}

			_, err = db.InsertGatewayEnvelopeWithChecksTransactional(
				ctx,
				querier,
				queries.InsertGatewayEnvelopeParams{
					OriginatorNodeID:     constants.IdentityUpdateOriginatorID,
					OriginatorSequenceID: int64(msgSent.SequenceId),
					Topic:                messageTopic.Bytes(),
					OriginatorEnvelope:   originatorEnvelopeBytes,
					Expiry:               math.MaxInt64,
				},
			)
			if err != nil {
				s.logger.Error(
					ErrInsertEnvelopeFromSmartContract,
					utils.TopicField(messageTopic.String()),
					zap.Error(err),
				)
				return re.NewRecoverableError(ErrInsertEnvelopeFromSmartContract, err)
			}

			return nil
		},
	)
	if err != nil {
		var logStorageErr re.RetryableError
		if errors.As(err, &logStorageErr) {
			return logStorageErr
		}
		// If the error was not a LogStorageError we can assume it's a DB error and it should be retried
		return re.NewRecoverableError(ErrInsertBlockchainMessage, err)
	}

	return nil
}

func (s *IdentityUpdateStorer) validateIdentityUpdate(
	ctx context.Context,
	querier *queries.Queries,
	inboxID [32]byte,
	clientEnvelope *envelopes.ClientEnvelope,
) (*mlsvalidate.AssociationStateResult, re.RetryableError) {
	gatewayEnvelopes, err := querier.SelectGatewayEnvelopesByTopics(
		ctx,
		queries.SelectGatewayEnvelopesByTopicsParams{
			Topics: []db.Topic{
				topic.NewTopic(topic.TopicKindIdentityUpdatesV1, inboxID[:]).Bytes(),
			},
			RowLimit: 256,
		},
	)
	// No rows returned means this is a new identity.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, re.NewRecoverableError(
			"could not retrieve envelopes to validate identity update",
			err,
		)
	}

	identityUpdate, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_IdentityUpdate)
	if !ok {
		return nil, re.NewNonRecoverableError(
			"client envelope payload is not an identity update",
			errors.New("client envelope payload is not an identity update"),
		)
	}

	result, err := s.validationService.GetAssociationStateFromEnvelopes(
		ctx,
		gatewayEnvelopes,
		identityUpdate.IdentityUpdate,
	)
	if err != nil {
		if shouldRetryValidationError(err) {
			return nil, re.NewRecoverableError(
				"could not get association state from envelopes",
				err,
			)
		}

		return nil, re.NewNonRecoverableError(
			"could not get association state from envelopes",
			err,
		)
	}

	return result, nil
}

func buildOriginatorEnvelope(
	originatorID uint32,
	sequenceID uint64,
	clientEnvelopeBytes []byte,
) (*envelopesProto.UnsignedOriginatorEnvelope, error) {
	payerEnvelope := &envelopesProto.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvelopeBytes,
	}
	payerEnvelopeBytes, err := proto.Marshal(payerEnvelope)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     originatorID,
		OriginatorSequenceId: sequenceID,
		OriginatorNs:         time.Now().UnixNano(),
		PayerEnvelopeBytes:   payerEnvelopeBytes,
	}, nil
}

func buildSignedOriginatorEnvelope(
	originatorEnvelope *envelopesProto.UnsignedOriginatorEnvelope,
	transactionHash common.Hash,
) (*envelopesProto.OriginatorEnvelope, error) {
	envelopeBytes, err := proto.Marshal(originatorEnvelope)
	if err != nil {
		return nil, err
	}

	return &envelopesProto.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: envelopeBytes,
		Proof: &envelopesProto.OriginatorEnvelope_BlockchainProof{
			BlockchainProof: &envelopesProto.BlockchainProof{
				TransactionHash: transactionHash[:],
			},
		},
	}, nil
}

// shouldRetryValidationError returns true if the error is a validation error that should be retried.
// Note: this approach is fragile as it depends on us creating new error messages for new validation errors.
// This function should rely on gRPC error codes instead, but it's not possible at the moment.
// Read https://github.com/xmtp/libxmtp/issues/3130
func shouldRetryValidationError(err error) bool {
	errMsg := err.Error()
	for _, pattern := range associationErrorPatterns {
		if strings.Contains(errMsg, pattern) {
			return false
		}
	}
	return true
}
