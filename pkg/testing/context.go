package testing

import (
	"testing"

	"github.com/xmtp/xmtpd/pkg/context"
)

func NewContext(t testing.TB) context.Context {
	return context.New(context.Background(), NewLogger(t))
}
