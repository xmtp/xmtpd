package crdttest

import (
	"bytes"
	"fmt"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

func Dump(ev *types.Event) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "CID: %x\n", ev.Cid)
	fmt.Fprintln(&b, "Links:")
	for _, l := range ev.Links {
		fmt.Fprintf(&b, "- %x\n", l)
	}
	fmt.Fprintf(&b, "Topic: %s\n", ev.ContentTopic)
	fmt.Fprintf(&b, "TimestampNs: %d\n", ev.TimestampNs)
	fmt.Fprintf(&b, "Message: %s\n", string(ev.Message))
	return b.String()
}
