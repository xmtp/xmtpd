package migrator

import (
	"fmt"

	"github.com/xmtp/xmtpd/pkg/envelopes"
)

type transformer struct{}

func NewTransformer() *transformer {
	return &transformer{}
}

func (t *transformer) Transform(record ISourceRecord) (*envelopes.OriginatorEnvelope, error) {
	switch record.TableName() {
	case addressLogTableName:
		return t.TransformAddressLog(record.(*AddressLog))

	case groupMessagesTableName:
		return t.TransformGroupMessage(record.(*GroupMessage))

	case inboxLogTableName:
		return t.TransformInboxLog(record.(*InboxLog))

	case installationsTableName:
		return t.TransformInstallation(record.(*Installation))

	case welcomeMessagesTableName:
		return t.TransformWelcomeMessage(record.(*WelcomeMessage))

	default:
		return nil, fmt.Errorf(
			"Transform not implemented for table: %s",
			record.TableName(),
		)
	}
}

// TransformAddressLog converts AddressLog to appropriate XMTP envelope format.
func (t *transformer) TransformAddressLog(
	addressLog *AddressLog,
) (*envelopes.OriginatorEnvelope, error) {
	return nil, fmt.Errorf(
		"TransformAddressLog not implemented",
	)
}

// TransformGroupMessage converts GroupMessage to appropriate XMTP envelope format.
func (t *transformer) TransformGroupMessage(
	groupMessage *GroupMessage,
) (*envelopes.OriginatorEnvelope, error) {
	return nil, fmt.Errorf(
		"TransformGroupMessage not implemented",
	)
}

// TransformInboxLog converts InboxLog to appropriate XMTP envelope format.
func (t *transformer) TransformInboxLog(
	inboxLog *InboxLog,
) (*envelopes.OriginatorEnvelope, error) {
	return nil, fmt.Errorf(
		"TransformInboxLog not implemented",
	)
}

// TransformInstallation converts Installation to appropriate XMTP envelope format.
func (t *transformer) TransformInstallation(
	installation *Installation,
) (*envelopes.OriginatorEnvelope, error) {
	return nil, fmt.Errorf(
		"TransformInstallation not implemented",
	)
}

// TransformWelcomeMessage converts WelcomeMessage to appropriate XMTP envelope format.
func (t *transformer) TransformWelcomeMessage(
	welcomeMessage *WelcomeMessage,
) (*envelopes.OriginatorEnvelope, error) {
	return nil, fmt.Errorf(
		"TransformWelcomeMessage not implemented",
	)
}

// TODO: Helper method.
func (t *transformer) createBasicEnvelope(
	_ []byte,
	_ []byte,
) (*envelopes.OriginatorEnvelope, error) {
	// 1. Create a ClientEnvelope with the appropriate authenticated data and payload
	// 2. Create a PayerEnvelope containing the ClientEnvelope
	// 3. Create an UnsignedOriginatorEnvelope with the PayerEnvelope
	// 4. Create the final OriginatorEnvelope with appropriate signatures/proofs

	return nil, fmt.Errorf(
		"createBasicEnvelope not implemented",
	)
}
