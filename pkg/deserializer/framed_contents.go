package deserializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type FramedContentIn struct {
	GroupID           GroupID
	Epoch             uint64
	Sender            Sender
	AuthenticatedData []byte
	Body              FramedContentBodyIn
}

func (f *FramedContentIn) TLSDeserialize(r *bytes.Reader) error {
	groupID, err := readVariableOpaqueVec(r)
	if err != nil {
		return err
	}

	var epoch uint64
	if err := binary.Read(r, binary.BigEndian, &epoch); err != nil {
		return err
	}

	var sender Sender
	if err := sender.TLSDeserialize(r); err != nil {
		return err
	}

	authData, err := readVariableOpaqueVec(r)
	if err != nil {
		return err
	}

	contentTypeByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	contentType := ContentType(contentTypeByte)

	var body FramedContentBodyIn
	switch contentType {
	case ContentTypeApplication:
		body = &ApplicationContent{}
	case ContentTypeProposal:
		body = &ProposalIn{}
	case ContentTypeCommit:
		body = &CommitIn{}
	default:
		return fmt.Errorf("unknown content type: %d", contentType)
	}

	if err := body.TLSDeserialize(r); err != nil {
		return err
	}

	f.GroupID = groupID
	f.Epoch = epoch
	f.Sender = sender
	f.AuthenticatedData = authData
	f.Body = body
	return nil
}

type FramedContentBodyIn interface {
	TLSDeserialize(r *bytes.Reader) error
	ContentType() ContentType
}

// ApplicationContent (stub)
type ApplicationContent struct{}

func (a *ApplicationContent) TLSDeserialize(r *bytes.Reader) error {
	return nil
}
func (*ApplicationContent) ContentType() ContentType { return ContentTypeApplication }

// ProposalIn (stub)
type ProposalIn struct{}

func (p *ProposalIn) TLSDeserialize(r *bytes.Reader) error { return nil }
func (*ProposalIn) ContentType() ContentType               { return ContentTypeProposal }

// CommitIn (stub)
type CommitIn struct{}

func (c *CommitIn) TLSDeserialize(r *bytes.Reader) error { return nil }
func (*CommitIn) ContentType() ContentType               { return ContentTypeCommit }
