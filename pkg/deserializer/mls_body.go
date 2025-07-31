package deserializer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type MlsMessageIn struct {
	Version ProtocolVersion
	Body    MlsMessageBodyIn
}

func (m *MlsMessageIn) TLSDeserialize(r *bytes.Reader) error {
	var version ProtocolVersion
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return err
	}

	body, err := DeserializeMlsMessageBodyIn(r)
	if err != nil {
		return err
	}

	if version != 1 {
		return errors.New("MlsMessage protocol version is not 1")
	}

	m.Version = version
	m.Body = body
	return nil
}

// ---------- MlsMessageBodyIn enum ----------

type MlsMessageBodyIn interface {
	TLSDeserializable
}

const (
	discriminantPublicMessage  = 1
	discriminantPrivateMessage = 2
	discriminantWelcome        = 3
	discriminantGroupInfo      = 4
	discriminantKeyPackage     = 5
)

func DeserializeMlsMessageBodyIn(r *bytes.Reader) (MlsMessageBodyIn, error) {
	var discriminant uint16
	if err := binary.Read(r, binary.BigEndian, &discriminant); err != nil {
		return nil, err
	}

	var body MlsMessageBodyIn
	switch discriminant {
	case discriminantPublicMessage:
		body = &PublicMessageIn{}
	case discriminantPrivateMessage:
		body = &PrivateMessageIn{}
	case discriminantWelcome:
		return nil, fmt.Errorf("unsupported discriminant: %d", discriminant)
	case discriminantGroupInfo:

		return nil, fmt.Errorf("unsupported discriminant: %d", discriminant)
	case discriminantKeyPackage:

		return nil, fmt.Errorf("unsupported discriminant: %d", discriminant)
	default:
		return nil, fmt.Errorf("unknown discriminant: %d", discriminant)
	}

	if err := body.TLSDeserialize(r); err != nil {
		return nil, err
	}
	return body, nil
}
