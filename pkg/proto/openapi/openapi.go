package openapi

import (
	"embed"
)

// Embed all swagger files
//
//go:embed xmtpv4/envelopes/envelopes.swagger.json
//go:embed xmtpv4/message_api/message_api.swagger.json
//go:embed xmtpv4/message_api/misbehavior_api.swagger.json
//go:embed xmtpv4/metadata_api/metadata_api.swagger.json
//go:embed xmtpv4/payer_api/payer_api.swagger.json
var SwaggerFS embed.FS
