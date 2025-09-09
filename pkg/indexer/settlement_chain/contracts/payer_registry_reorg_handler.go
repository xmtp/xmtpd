package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/utils/retryerrors"
	"go.uber.org/zap"
)

type PayerRegistryReorgHandler struct {
	logger *zap.Logger
}

func NewPayerRegistryReorgHandler(logger *zap.Logger) *PayerRegistryReorgHandler {
	return &PayerRegistryReorgHandler{logger: logger.Named("reorg-handler")}
}

func (h *PayerRegistryReorgHandler) HandleLog(
	_ context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", zap.Any("blockNumber", event.BlockNumber))
	return nil
}
