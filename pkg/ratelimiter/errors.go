package ratelimiter

import "errors"

var (
	ErrCostMustBeGreaterThanZero = errors.New("cost must be > 0")
	ErrUnexpectedScriptResponse  = errors.New("unexpected script response")
	ErrNoLimitsProvided          = errors.New("no limits provided")
	ErrInvalidFailedLimit        = errors.New("invalid failed limit")
)
