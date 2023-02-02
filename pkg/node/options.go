package node

import (
	"github.com/xmtp/xmtpd/pkg/api"
)

type LogOptions struct {
	Level    string `long:"level" description:"Log level. Support values: error, warn, info, debug" default:"info"`
	Encoding string `long:"encoding" description:"Log encoding format. Support values: console, json" default:"console"`
}

type Options struct {
	Log LogOptions  `group:"Log options" namespace:"log"`
	API api.Options `group:"API options" namespace:"api"`
}
