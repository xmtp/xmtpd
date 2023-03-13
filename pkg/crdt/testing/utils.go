package crdttest

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
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

func RequireEventsEqual(t *testing.T, exp []*types.Event, act []*types.Event) {
	require.Equal(t, len(exp), len(act))
	for i := range exp {
		RequireEventEqual(t, exp[i], act[i], "event[%d]", i)
	}
}

func RequireEventEqual(t *testing.T, exp, act *types.Event, msgAndArgs ...any) {
	require.Equal(t, DumpEvent(exp), DumpEvent(act), msgAndArgs)
}
