package deserializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PrivateMessageIn struct {
	GroupID     []byte
	Epoch       uint64
	ContentType ContentType
	// and more
}

func (s *Sender) IsMember() bool {
	return s.Type == SenderTypeMember
}

func (s *Sender) AsMember() *uint32 {
	if s.Type == SenderTypeMember {
		return s.LeafIndex
	}
	return nil
}

func (p *PrivateMessageIn) TLSDeserialize(r *bytes.Reader) error {
	groupID, err := readVariableOpaqueVec(r)
	if err != nil {
		return fmt.Errorf("failed to read group_id: %w", err)
	}

	var epoch uint64
	if err := binary.Read(r, binary.BigEndian, &epoch); err != nil {
		return fmt.Errorf("failed to read epoch: %w", err)
	}

	contentTypeByte, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read content_type: %w", err)
	}

	// we don't care about anything after this

	p.GroupID = groupID
	p.Epoch = epoch
	p.ContentType = ContentType(contentTypeByte)
	return nil
}
