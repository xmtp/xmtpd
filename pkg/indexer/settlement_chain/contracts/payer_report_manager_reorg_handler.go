package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type PayerReportManagerReorgHandler struct {
	logger *zap.Logger
}

func NewPayerReportManagerReorgHandler(logger *zap.Logger) *PayerReportManagerReorgHandler {
	return &PayerReportManagerReorgHandler{logger: logger.Named(utils.ReorgHandlerLoggerName)}
}

func (h *PayerReportManagerReorgHandler) HandleLog(
	_ context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", utils.BlockNumberField(event.BlockNumber))
	return nil
}
