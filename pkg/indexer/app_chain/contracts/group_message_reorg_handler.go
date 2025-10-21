package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type GroupMessageReorgHandler struct {
	logger *zap.Logger
}

func NewGroupMessageReorgHandler(logger *zap.Logger) *GroupMessageReorgHandler {
	return &GroupMessageReorgHandler{logger: logger.Named(utils.ReorgHandlerLoggerName)}
}

func (h *GroupMessageReorgHandler) HandleLog(
	_ context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", utils.BlockNumberField(event.BlockNumber))
	return nil
}
