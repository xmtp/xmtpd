package ledger

type EventType int16

const (
	// Triggered by a Deposit smart contract event
	EVENT_TYPE_DEPOSIT EventType = 0
	// Triggered by a WithdrawalRequested smart contract event
	EVENT_TYPE_WITHDRAWAL EventType = 1
	// Triggered by a UsageSettled smart contract event
	EVENT_TYPE_SETTLEMENT EventType = 2
	// Triggered by a WithdrawalCancelled smart contract event
	EVENT_TYPE_CANCELED_WITHDRAWAL EventType = 3
	// Triggered by our reorg handler
	EVENT_TYPE_REORG_REVERSAL EventType = 4
)
