package sync

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/registrant"
	"github.com/xmtp/xmtpd/pkg/registry"
	"go.uber.org/zap"
)

type SyncServer struct {
	ctx        context.Context
	log        *zap.Logger
	registrant *registrant.Registrant
	store      *sql.DB
	worker     *syncWorker
}

func NewSyncServer(
	ctx context.Context,
	log *zap.Logger,
	nodeRegistry registry.NodeRegistry,
	registrant *registrant.Registrant,
	store *sql.DB,
) (*SyncServer, error) {
	worker, err := startSyncWorker(ctx, log, nodeRegistry, registrant, store)
	if err != nil {
		return nil, err
	}

	s := &SyncServer{
		ctx:        ctx,
		log:        log.Named("sync"),
		registrant: registrant,
		store:      store,
		worker:     worker,
	}

	// TODO(rich): Add healthcheck
	return s, nil
}

func (s *SyncServer) Close() {
	s.log.Info("closing")
	s.worker.close()
	s.log.Info("closed")
}
