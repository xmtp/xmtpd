package node

import (
	apigateway "github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Options struct {
	Log zap.Options        `group:"Log options" namespace:"log"`
	API apigateway.Options `group:"API options" namespace:"api"`
}
