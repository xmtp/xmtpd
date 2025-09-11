package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type IdentityUpdateReorgHandler struct {
	logger *zap.Logger
}

func NewIdentityUpdateReorgHandler(logger *zap.Logger) *IdentityUpdateReorgHandler {
	return &IdentityUpdateReorgHandler{logger: logger.Named("reorg-handler")}
}

func (h *IdentityUpdateReorgHandler) HandleLog(
	_ context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", zap.Any("blockNumber", event.BlockNumber))
	return nil
}
