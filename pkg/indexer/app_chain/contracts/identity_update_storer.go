package contracts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pingcap/log"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
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
	// We may not want to hardcode this to 1 and have an originator ID for each smart contract?
	IDENTITY_UPDATE_ORIGINATOR_ID = 1
)

var (
	ErrAdvisoryLockSequence   = "advisory lock failed"
	ErrParseIdentityUpdate    = "error parsing identity update"
	ErrGetLatestSequenceId    = "get latest sequence id failed"
	ErrValidateIdentityUpdate = "validate identity update failed"
	ErrInsertAddressLog       = "insert address log failed"
	ErrRevokeAddressFromLog   = "revoke address from log failed"
)

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
		logger:            logger.Named("storer"),
		contract:          contract,
		validationService: validationService,
	}
}

// Validate and store an identity update log event
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
				LockIdentityUpdateInsert(ctx, querier, uint32(IDENTITY_UPDATE_ORIGINATOR_ID))
			if err != nil {
				return re.NewNonRecoverableError(ErrAdvisoryLockSequence, err)
			}

			latestSequenceId, err := querier.GetLatestSequenceId(ctx, IDENTITY_UPDATE_ORIGINATOR_ID)
			if err != nil {
				return re.NewNonRecoverableError(ErrGetLatestSequenceId, err)
			}

			if uint64(latestSequenceId) >= msgSent.SequenceId {
				s.logger.Debug(
					"Identity update already inserted. Skipping... ",
					zap.Uint64("latest_sequence_id", uint64(latestSequenceId)),
					zap.Uint64("msg_sequence_id", msgSent.SequenceId),
				)
				return nil
			}

			messageTopic := topic.NewTopic(topic.TopicKindIdentityUpdatesV1, msgSent.InboxId[:])

			s.logger.Info(
				"Inserting identity update from contract",
				zap.String("topic", messageTopic.String()),
			)

			clientEnvelope, err := envelopes.NewClientEnvelopeFromBytes(msgSent.Update)
			if err != nil {
				s.logger.Error(ErrParseClientEnvelope, zap.Error(err))
				return re.NewNonRecoverableError(ErrParseClientEnvelope, err)
			}

			associationState, err := s.validateIdentityUpdate(
				ctx,
				querier,
				msgSent.InboxId,
				clientEnvelope,
			)
			if err != nil {
				s.logger.Error(ErrValidateIdentityUpdate, zap.Error(err))
				return re.NewNonRecoverableError(ErrValidateIdentityUpdate, err)
			}

			inboxId := utils.HexEncode(msgSent.InboxId[:])

			for _, newMember := range associationState.StateDiff.NewMembers {
				s.logger.Info("New member", zap.Any("member", newMember))
				if address, ok := newMember.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
					numRows, err := querier.InsertAddressLog(ctx, queries.InsertAddressLogParams{
						Address: address.EthereumAddress,
						InboxID: inboxId,
						AssociationSequenceID: sql.NullInt64{
							Valid: true,
							Int64: int64(msgSent.SequenceId),
						},
					})
					if err != nil {
						return re.NewRecoverableError(ErrInsertAddressLog, err)
					}
					if numRows == 0 {
						s.logger.Warn(
							"Could not insert address log",
							zap.String("address", address.EthereumAddress),
							zap.String("inbox_id", inboxId),
							zap.Int("sequence_id", int(msgSent.SequenceId)),
						)
					}
				}
			}

			for _, removed_member := range associationState.StateDiff.RemovedMembers {
				log.Info("Removed member", zap.Any("member", removed_member))
				if address, ok := removed_member.Kind.(*associations.MemberIdentifier_EthereumAddress); ok {
					rows, err := querier.RevokeAddressFromLog(
						ctx,
						queries.RevokeAddressFromLogParams{
							Address: address.EthereumAddress,
							InboxID: inboxId,
							RevocationSequenceID: sql.NullInt64{
								Valid: true,
								Int64: int64(msgSent.SequenceId),
							},
						},
					)
					if err != nil {
						return re.NewRecoverableError(ErrRevokeAddressFromLog, err)
					}
					if rows == 0 {
						s.logger.Warn(
							"Could not find address log entry to revoke",
							zap.String("address", address.EthereumAddress),
							zap.String("inbox_id", inboxId),
						)
					}
				}
			}

			originatorEnvelope, err := buildOriginatorEnvelope(
				IDENTITY_UPDATE_ORIGINATOR_ID,
				msgSent.SequenceId,
				msgSent.Update,
			)
			if err != nil {
				s.logger.Error(ErrBuildOriginatorEnvelope, zap.Error(err))
				return re.NewNonRecoverableError(ErrBuildOriginatorEnvelope, err)
			}

			signedOriginatorEnvelope, err := buildSignedOriginatorEnvelope(
				originatorEnvelope,
				event.TxHash,
			)
			if err != nil {
				s.logger.Error(ErrBuildSignedOriginatorEnvelope, zap.Error(err))
				return re.NewNonRecoverableError(ErrBuildSignedOriginatorEnvelope, err)
			}

			originatorEnvelopeBytes, err := proto.Marshal(signedOriginatorEnvelope)
			if err != nil {
				s.logger.Error(ErrMarshallOriginatorEnvelope, zap.Error(err))
				return re.NewNonRecoverableError(ErrMarshallOriginatorEnvelope, err)
			}

			if _, err = querier.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     IDENTITY_UPDATE_ORIGINATOR_ID,
				OriginatorSequenceID: int64(msgSent.SequenceId),
				Topic:                messageTopic.Bytes(),
				OriginatorEnvelope:   originatorEnvelopeBytes,
				Expiry:               sql.NullInt64{Int64: math.MaxInt64, Valid: true},
			}); err != nil {
				s.logger.Error(ErrInsertEnvelopeFromSmartContract, zap.Error(err))
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
	inboxId [32]byte,
	clientEnvelope *envelopes.ClientEnvelope,
) (*mlsvalidate.AssociationStateResult, error) {
	gatewayEnvelopes, err := querier.SelectGatewayEnvelopes(
		ctx,
		queries.SelectGatewayEnvelopesParams{
			Topics: []db.Topic{
				db.Topic(topic.NewTopic(topic.TopicKindIdentityUpdatesV1, inboxId[:]).Bytes()),
			},
			OriginatorNodeIds: []int32{IDENTITY_UPDATE_ORIGINATOR_ID},
			RowLimit:          256,
		},
	)
	if err != nil {
		return nil, err
	}

	identityUpdate, ok := clientEnvelope.Payload().(*envelopesProto.ClientEnvelope_IdentityUpdate)
	if !ok {
		return nil, fmt.Errorf("client envelope payload is not an identity update")
	}

	return s.validationService.GetAssociationStateFromEnvelopes(
		ctx,
		gatewayEnvelopes,
		identityUpdate.IdentityUpdate,
	)
}

func buildOriginatorEnvelope(
	originatorId uint32,
	sequenceId uint64,
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
		OriginatorNodeId:     originatorId,
		OriginatorSequenceId: sequenceId,
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
