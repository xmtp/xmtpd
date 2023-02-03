package node

import (
	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log zap.Options `group:"Log options" namespace:"log"`
	API api.Options `group:"API options" namespace:"api"`
}
