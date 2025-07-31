package deserializer

import "bytes"

type TLSDeserializable interface {
	TLSDeserialize(r *bytes.Reader) error
}
