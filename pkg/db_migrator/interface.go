package db_migrator

import (
	"github.com/xmtp/xmtpd/pkg/envelopes"
)

// DataTransformer defines the interface for transforming external data to xmtpd OriginatorEnvelope format.
type DataTransformer interface {
	TransformAddressLog(addressLog *AddressLog) (*envelopes.OriginatorEnvelope, error)
	TransformGroupMessage(groupMessage *GroupMessage) (*envelopes.OriginatorEnvelope, error)
	TransformInboxLog(inboxLog *InboxLog) (*envelopes.OriginatorEnvelope, error)
	TransformInstallation(installation *Installation) (*envelopes.OriginatorEnvelope, error)
	TransformWelcomeMessage(welcomeMessage *WelcomeMessage) (*envelopes.OriginatorEnvelope, error)
}
