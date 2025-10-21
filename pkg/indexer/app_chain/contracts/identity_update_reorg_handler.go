package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type IdentityUpdateReorgHandler struct {
	logger *zap.Logger
}

func NewIdentityUpdateReorgHandler(logger *zap.Logger) *IdentityUpdateReorgHandler {
	return &IdentityUpdateReorgHandler{logger: logger.Named(utils.ReorgHandlerLoggerName)}
}

func (h *IdentityUpdateReorgHandler) HandleLog(
	_ context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", utils.BlockNumberField(event.BlockNumber))
	return nil
}
