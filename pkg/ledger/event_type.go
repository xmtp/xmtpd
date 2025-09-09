package ledger

type EventType int16

const (
	// Triggered by a Deposit smart contract event
	eventTypeDeposit EventType = 0
	// Triggered by a WithdrawalRequested smart contract event
	eventTypeWithdrawal EventType = 1
	// Triggered by a UsageSettled smart contract event
	eventTypeSettlement EventType = 2
	// Triggered by a WithdrawalCancelled smart contract event
	eventTypeCanceledWithdrawal EventType = 3
	// Triggered by our reorg handler
	eventTypeReorgReversal EventType = 4
)
