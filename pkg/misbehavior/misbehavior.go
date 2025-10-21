// Package misbehavior implements the misbehavior reports service.
package misbehavior

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// LoggingMisbehaviorService provides an implementation of the MisbehaviorService interface that
// logs misbehavior reports without storing them or forwarding to the network.
type LoggingMisbehaviorService struct {
	logger *zap.Logger
}

func NewLoggingMisbehaviorService(logger *zap.Logger) *LoggingMisbehaviorService {
	return &LoggingMisbehaviorService{
		logger: logger.Named(utils.MisbehaviorLoggerName),
	}
}

func (m *LoggingMisbehaviorService) SafetyFailure(report *SafetyFailureReport) error {
	if report == nil {
		return errors.New("report is nil")
	}
	m.logger.Warn(
		"misbehavior detected",
		zap.String("misbehavior_type", report.misbehaviorType.String()),
		zap.Uint32("misbehaving_node_id", report.misbehavingNodeID),
		zap.Bool("submitted_by_node", report.submittedByNode),
		zap.Any("envelopes", report.envelopes),
	)

	return nil
}
