package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

type GroupMessageReorgHandler struct {
	logger *zap.Logger
}

func NewGroupMessageReorgHandler(logger *zap.Logger) *GroupMessageReorgHandler {
	return &GroupMessageReorgHandler{logger: logger.Named("reorg-handler")}
}

func (h *GroupMessageReorgHandler) HandleLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", zap.Any("log", event))
	return nil
}
