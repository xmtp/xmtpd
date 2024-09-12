package storer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

type IdentityUpdateStorer struct {
	contract *abis.IdentityUpdates
	queries  *queries.Queries
	logger   *zap.Logger
}

func NewIdentityUpdateStorer(
	queries *queries.Queries,
	logger *zap.Logger,
	contract *abis.IdentityUpdates,
) *IdentityUpdateStorer {
	return &IdentityUpdateStorer{
		queries:  queries,
		logger:   logger.Named("IdentityUpdateStorer"),
		contract: contract,
	}
}

// Validate and store an identity update log event
func (s *IdentityUpdateStorer) StoreLog(ctx context.Context, event types.Log) LogStorageError {
	msgSent, err := s.contract.ParseIdentityUpdateCreated(event)
	if err != nil {
		return NewLogStorageError(err, false)
	}

	// TODO:nm figure out topic structure
	topic := BuildInboxTopic(msgSent.InboxId)

	s.logger.Debug("Inserting identity update from contract", zap.String("topic", topic))

	/**
	TODO:nm validate the identity update
	**/

	if _, err = s.queries.InsertGatewayEnvelope(ctx, queries.InsertGatewayEnvelopeParams{
		// We may not want to hardcode this to 0 and have an originator ID for each smart contract?
		OriginatorNodeID:     0,
		OriginatorSequenceID: int64(msgSent.SequenceId),
		Topic:                []byte(topic),
		OriginatorEnvelope:   msgSent.Update, // TODO:nm parse originator envelope and do some validation
	}); err != nil {
		s.logger.Error("Error inserting envelope from smart contract", zap.Error(err))
		return NewLogStorageError(err, true)
	}

	return nil
}

func BuildInboxTopic(inboxId [32]byte) string {
	return fmt.Sprintf("1/i/%x", inboxId)
}
