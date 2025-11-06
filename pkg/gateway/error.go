package gateway

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
)

type GatewayServiceError interface {
	error
	Code() connect.Code
	ClientMessage() string
	RetryAfter() *time.Duration
}

type PermissionDeniedError struct {
	msg string
	err error
}

func (e PermissionDeniedError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e PermissionDeniedError) ClientMessage() string {
	return e.msg
}

func (e PermissionDeniedError) Code() connect.Code {
	return connect.CodePermissionDenied
}

func (e PermissionDeniedError) RetryAfter() *time.Duration {
	return nil
}

func NewPermissionDeniedError(msg string, err error) *PermissionDeniedError {
	return &PermissionDeniedError{msg: msg, err: err}
}

type UnauthenticatedError struct {
	msg string
	err error
}

func (e UnauthenticatedError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e UnauthenticatedError) ClientMessage() string {
	return e.msg
}

func (e UnauthenticatedError) Code() connect.Code {
	return connect.CodeUnauthenticated
}

func (e UnauthenticatedError) RetryAfter() *time.Duration {
	return nil
}

func NewUnauthenticatedError(msg string, err error) *UnauthenticatedError {
	return &UnauthenticatedError{msg: msg, err: err}
}

type RateLimitExceededError struct {
	err        error
	retryAfter time.Duration
}

func (e RateLimitExceededError) Error() string {
	if e.err == nil {
		return "rate limit exceeded"
	}

	return fmt.Sprintf("rate limit exceeded: %s", e.err.Error())
}

func (e RateLimitExceededError) ClientMessage() string {
	return "rate limit exceeded"
}

func (e RateLimitExceededError) Code() connect.Code {
	return connect.CodeResourceExhausted
}

func (e RateLimitExceededError) RetryAfter() *time.Duration {
	if e.retryAfter == 0 {
		return nil
	}
	return &e.retryAfter
}

func NewRateLimitExceededError(err error, retryAfter time.Duration) *RateLimitExceededError {
	return &RateLimitExceededError{
		err:        err,
		retryAfter: retryAfter,
	}
}
