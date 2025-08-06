package ledger

import "errors"

var (
	ErrInvalidAmount      = errors.New("amount must be greater than 0")
	ErrInvalidEventID     = errors.New("event ID must be greater than 0")
	ErrWithdrawalNotFound = errors.New(
		"trying to cancel a withdrawal that has not been stored",
	)
	ErrWithdrawalAlreadyCanceled = errors.New("withdrawal already canceled")
)
