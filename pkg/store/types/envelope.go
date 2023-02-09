package types

import (
	"bytes"

	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/protobuf/proto"
)

type Envelope struct {
	*messagev1.Envelope
	Cid multihash.Multihash
}

func WrapEnvelope(env *messagev1.Envelope) (*Envelope, error) {
	envB, err := proto.Marshal(env)
	if err != nil {
		return nil, err
	}
	cid, err := multihash.Sum(envB, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Envelope: env,
		Cid:      cid,
	}, nil
}

func (e *Envelope) Compare(env *Envelope) int {
	res := e.TimestampNs - env.TimestampNs
	if res != 0 {
		return int(res)
	}
	return bytes.Compare(e.Cid, env.Cid)
}
