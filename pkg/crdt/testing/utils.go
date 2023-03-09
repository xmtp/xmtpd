package crdttest

import (
	"fmt"
	"io"

	v1 "github.com/xmtp/proto/v3/go/message_api/v1"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"go.uber.org/zap/buffer"
)

func DumpEvent(ev *types.Event) string {
	var w buffer.Buffer
	fmt.Fprintf(&w, "CID: %x\n", ev.Cid)
	fmt.Fprintln(&w, "Links:")
	for _, l := range ev.Links {
		fmt.Fprintf(&w, "- %x\n", l)
	}
	DumpEnvelope(&w, ev.Envelope)
	return w.String()
}

func DumpEnvelope(w io.Writer, env *v1.Envelope) {
	fmt.Fprintf(w, "Topic: %s\n", env.ContentTopic)
	fmt.Fprintf(w, "TimestampNs: %d\n", env.TimestampNs)
	fmt.Fprintf(w, "Message: %s\n", string(env.Message))
}
