package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	re "github.com/xmtp/xmtpd/pkg/errors"
	"go.uber.org/zap"
)

type PayerRegistryReorgHandler struct {
	logger *zap.Logger
}

func NewPayerRegistryReorgHandler(logger *zap.Logger) *PayerRegistryReorgHandler {
	return &PayerRegistryReorgHandler{logger: logger.Named("reorg-handler")}
}

func (h *PayerRegistryReorgHandler) HandleLog(
	ctx context.Context,
	event types.Log,
) re.RetryableError {
	h.logger.Info("handling reorged event", zap.Any("log", event))
	return nil
}
