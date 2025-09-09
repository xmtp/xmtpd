// Package deserializer implements the deserializer for the MLS messages.
package deserializer

import (
	"bytes"

	envelopesProto "github.com/xmtp/xmtpd/pkg/proto/xmtpv4/envelopes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DeserializeGroupMessage(
	payload *envelopesProto.ClientEnvelope_GroupMessage,
) (*MlsMessageIn, error) {
	if payload == nil || payload.GroupMessage == nil || payload.GroupMessage.GetV1() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload")
	}

	if payload.GroupMessage.GetV1().Data == nil {
		return nil, status.Errorf(codes.InvalidArgument, "can not process empty payload")
	}

	r := bytes.NewReader(payload.GroupMessage.GetV1().Data)

	msg := MlsMessageIn{}
	err := msg.TLSDeserialize(r)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"failed to deserialize MLS message: %v",
			err,
		)
	}

	return &msg, nil
}

func IsGroupMessageCommit(
	payload *envelopesProto.ClientEnvelope_GroupMessage,
) (bool, error) {
	msg, err := DeserializeGroupMessage(payload)
	if err != nil {
		return false, err
	}

	var ct ContentType

	switch msg := msg.Body.(type) {
	case *PublicMessageIn:
		ct = msg.Content.Body.ContentType()
	case *PrivateMessageIn:
		ct = msg.ContentType
	default:
		return false, status.Errorf(codes.InvalidArgument, "invalid message type")
	}

	return ct == ContentTypeCommit, nil
}
