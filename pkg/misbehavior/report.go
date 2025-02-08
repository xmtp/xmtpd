package misbehavior

import (
	"errors"

	"github.com/xmtp/xmtpd/pkg/envelopes"
	proto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
)

type SafetyFailureReport struct {
	misbehavingNodeId uint32
	misbehaviorType   proto.Misbehavior
	submittedByNode   bool
	envelopes         []*envelopes.OriginatorEnvelope
}

func NewSafetyFailureReport(
	misbehavingNodeId uint32,
	misbehaviorType proto.Misbehavior,
	submittedByNode bool,
	envs []*envelopes.OriginatorEnvelope,
) (*SafetyFailureReport, error) {
	if len(envs) == 0 {
		return nil, errors.New("no envelopes provided")
	}

	if misbehavingNodeId == 0 {
		return nil, errors.New("misbehaving node id is required")
	}

	if misbehaviorType == proto.Misbehavior_MISBEHAVIOR_UNSPECIFIED {
		return nil, errors.New("misbehavior type is required")
	}

	return &SafetyFailureReport{
		misbehavingNodeId: misbehavingNodeId,
		misbehaviorType:   misbehaviorType,
		submittedByNode:   submittedByNode,
		envelopes:         envs,
	}, nil
}
