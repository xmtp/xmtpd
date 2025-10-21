package misbehavior

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/envelopes"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	testEnvelopes "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMisbehaviorService(t *testing.T) {
	t.Skip("skipping misbehavior test")
	env1, err := envelopes.NewOriginatorEnvelope(testEnvelopes.CreateOriginatorEnvelope(t, 1, 1))
	require.NoError(t, err)
	env2, err := envelopes.NewOriginatorEnvelope(testEnvelopes.CreateOriginatorEnvelope(t, 1, 2))
	require.NoError(t, err)

	report, err := NewSafetyFailureReport(
		1,
		proto.Misbehavior_MISBEHAVIOR_DUPLICATE_SEQUENCE_ID,
		true,
		[]*envelopes.OriginatorEnvelope{env1, env2},
	)
	require.NoError(t, err)

	core, observedLogs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	service := NewLoggingMisbehaviorService(logger)

	err = service.SafetyFailure(report)
	require.NoError(t, err)

	logs := observedLogs.All()
	require.Len(t, logs, 1)
	logEntry := logs[0]
	require.Equal(t, "misbehavior detected", logEntry.Message)
	require.Equal(t, "MISBEHAVIOR_DUPLICATE_SEQUENCE_ID", logEntry.ContextMap()["misbehavior_type"])
	require.Equal(t, uint32(1), logEntry.ContextMap()["misbehaving_node_id"])
	require.Equal(t, true, logEntry.ContextMap()["submitted_by_node"])
	require.Len(t, logEntry.ContextMap()["envelopes"], 2)
}

func TestNewSafetyFailureReportValidations(t *testing.T) {
	t.Skip("skipping misbehavior test")
	// Test case: No envelopes provided
	_, err := NewSafetyFailureReport(
		1,
		proto.Misbehavior_MISBEHAVIOR_DUPLICATE_SEQUENCE_ID,
		true,
		nil,
	)
	require.Error(t, err)
	require.Equal(t, "no envelopes provided", err.Error())

	// Test case: Misbehaving node ID is zero
	env, _ := envelopes.NewOriginatorEnvelope(testEnvelopes.CreateOriginatorEnvelope(t, 1, 1))
	_, err = NewSafetyFailureReport(
		0,
		proto.Misbehavior_MISBEHAVIOR_DUPLICATE_SEQUENCE_ID,
		true,
		[]*envelopes.OriginatorEnvelope{env},
	)
	require.Error(t, err)
	require.Equal(t, "misbehaving node id is required", err.Error())

	// Test case: Misbehavior type is unspecified
	_, err = NewSafetyFailureReport(
		1,
		proto.Misbehavior_MISBEHAVIOR_UNSPECIFIED,
		true,
		[]*envelopes.OriginatorEnvelope{env},
	)
	require.Error(t, err)
	require.Equal(t, "misbehavior type is required", err.Error())
}
