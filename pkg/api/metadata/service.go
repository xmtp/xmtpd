package metadata

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	metadata_api.UnimplementedMetadataApiServer
	ctx   context.Context
	log   *zap.Logger
	store *sql.DB
}

func NewMetadataApiService(
	ctx context.Context,
	log *zap.Logger,
	store *sql.DB,
) (*Service, error) {
	return &Service{
		ctx:   ctx,
		log:   log,
		store: store,
	}, nil
}

func (s *Service) GetSyncCursor(
	ctx context.Context,
	_ *metadata_api.GetSyncCursorRequest,
) (*metadata_api.GetSyncCursorResponse, error) {

	rows, err := queries.New(s.store).GetLatestCursor(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not select latest cursor: %v", err)
	}

	return convertToGetSyncCursorResponse(rows), nil
}

func convertToGetSyncCursorResponse(
	rows []queries.GetLatestCursorRow,
) *metadata_api.GetSyncCursorResponse {
	nodeIdToSequenceId := make(map[uint32]uint64)
	for _, row := range rows {
		fmt.Println("row", row)
		nodeIdToSequenceId[uint32(row.OriginatorNodeID)] = uint64(row.MaxSequenceID)
	}

	cursor := &envelopes.Cursor{
		NodeIdToSequenceId: nodeIdToSequenceId,
	}

	return &metadata_api.GetSyncCursorResponse{
		LatestSync: cursor,
	}
}
