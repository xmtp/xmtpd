package node

import (
	"github.com/xmtp/xmtpd/pkg/apigateway"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log zap.Options        `group:"Log options" namespace:"log"`
	API apigateway.Options `group:"API options" namespace:"api"`
}
