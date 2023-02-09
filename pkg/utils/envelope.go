package utils

import (
	"github.com/multiformats/go-multihash"
	messagev1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"google.golang.org/protobuf/proto"
)

func BuildEnvelopeCid(env *messagev1.Envelope) (multihash.Multihash, error) {
	envB, err := proto.Marshal(env)
	if err != nil {
		return nil, err
	}
	cid, err := multihash.Sum(envB, multihash.SHA2_256, -1)
	if err != nil {
		return nil, err
	}
	return cid, nil
}
