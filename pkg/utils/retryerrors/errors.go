// Package retryerrors implements the retryable error interface.
package retryerrors

import "fmt"

type RetryableError interface {
	error
	ShouldRetry() bool
}

type NonRecoverableError struct {
	msg string
	err error
}

func NewNonRecoverableError(msg string, err error) *NonRecoverableError {
	return &NonRecoverableError{msg: msg, err: err}
}

func (e NonRecoverableError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e NonRecoverableError) ShouldRetry() bool {
	return false
}

func (e NonRecoverableError) Unwrap() error {
	return e.err
}

type RecoverableError struct {
	msg string
	err error
}

func NewRecoverableError(msg string, err error) *RecoverableError {
	return &RecoverableError{msg: msg, err: err}
}

func (e RecoverableError) Error() string {
	if e.err == nil {
		return e.msg
	}

	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e RecoverableError) ShouldRetry() bool {
	return true
}

func (e RecoverableError) Unwrap() error {
	return e.err
}
