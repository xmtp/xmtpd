package retryable_errors

import "fmt"

type RetryableError interface {
	error
	ShouldRetry() bool
}

type NonRecoverableError struct {
	msg string
	err error
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

func NewNonRecoverableError(msg string, err error) *NonRecoverableError {
	return &NonRecoverableError{msg: msg, err: err}
}

type RecoverableError struct {
	msg string
	err error
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

func NewRecoverableError(msg string, err error) *RecoverableError {
	return &RecoverableError{msg: msg, err: err}
}
