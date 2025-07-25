package deserializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Sender struct {
	Type        SenderType
	LeafIndex   *uint32 // Only set if Member
	SenderIndex *uint32 // Only set if External
}

func (s *Sender) TLSDeserialize(r *bytes.Reader) error {
	var senderTypeByte uint8
	if err := binary.Read(r, binary.BigEndian, &senderTypeByte); err != nil {
		return fmt.Errorf("failed to read sender type: %w", err)
	}
	s.Type = SenderType(senderTypeByte)

	s.LeafIndex = nil
	s.SenderIndex = nil

	switch s.Type {
	case SenderTypeMember:
		var idx uint32
		if err := binary.Read(r, binary.BigEndian, &idx); err != nil {
			return fmt.Errorf("failed to read member leaf index: %w", err)
		}
		s.LeafIndex = &idx

	case SenderTypeExternal:
		var idx uint32
		if err := binary.Read(r, binary.BigEndian, &idx); err != nil {
			return fmt.Errorf("failed to read external sender index: %w", err)
		}
		s.SenderIndex = &idx

	case SenderTypeNewMemberProposal, SenderTypeNewMemberCommit:
		// No additional payload

	default:
		return fmt.Errorf("unknown sender type: %d", s.Type)
	}

	return nil
}
