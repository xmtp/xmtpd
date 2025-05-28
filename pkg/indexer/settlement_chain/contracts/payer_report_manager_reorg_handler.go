package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

type PayerReportManagerReorgHandler struct {
	logger *zap.Logger
}

func NewPayerReportManagerReorgHandler(logger *zap.Logger) *PayerReportManagerReorgHandler {
	return &PayerReportManagerReorgHandler{logger: logger.Named("reorg-handler")}
}

func (h *PayerReportManagerReorgHandler) HandleLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", zap.Any("log", event))
	return nil
}
