package misbehavior

import (
	"errors"

	"go.uber.org/zap"
)

// LoggingMisbehaviorService provides an implementation of the MisbehaviorService interface that
// logs misbehavior reports without storing them or forwarding to the network.
type LoggingMisbehaviorService struct {
	log *zap.Logger
}

func NewLoggingMisbehaviorService(log *zap.Logger) *LoggingMisbehaviorService {
	return &LoggingMisbehaviorService{
		log: log.Named("misbehavior"),
	}
}

func (m *LoggingMisbehaviorService) SafetyFailure(report *SafetyFailureReport) error {
	if report == nil {
		return errors.New("report is nil")
	}
	m.log.Warn(
		"misbehavior detected",
		zap.String("misbehavior_type", report.misbehaviorType.String()),
		zap.Uint32("misbehaving_node_id", report.misbehavingNodeId),
		zap.Bool("submitted_by_node", report.submittedByNode),
		zap.Any("envelopes", report.envelopes),
	)

	return nil
}
