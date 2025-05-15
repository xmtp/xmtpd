package storer

import "fmt"

type LogStorageError interface {
	error
	ShouldRetry() bool
}

type UnrecoverableLogStorageError struct {
	msg string
	err error
}

func (e UnrecoverableLogStorageError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e UnrecoverableLogStorageError) ShouldRetry() bool {
	return false
}

func NewUnrecoverableLogStorageError(msg string, err error) UnrecoverableLogStorageError {
	return UnrecoverableLogStorageError{err: err, msg: msg}
}

type RetryableLogStorageError struct {
	msg string
	err error
}

func (e RetryableLogStorageError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func (e RetryableLogStorageError) ShouldRetry() bool {
	return true
}

func NewRetryableLogStorageError(msg string, err error) RetryableLogStorageError {
	return RetryableLogStorageError{err: err, msg: msg}
}
