package storer

type LogStorageError interface {
	error
	ShouldRetry() bool
}

type UnrecoverableLogStorageError struct {
	err error
}

func (e UnrecoverableLogStorageError) Error() string {
	return e.err.Error()
}

func (e UnrecoverableLogStorageError) ShouldRetry() bool {
	return false
}

func NewUnrecoverableLogStorageError(err error) UnrecoverableLogStorageError {
	return UnrecoverableLogStorageError{err: err}
}

type RetryableLogStorageError struct {
	err error
}

func (e RetryableLogStorageError) Error() string {
	return e.err.Error()
}

func (e RetryableLogStorageError) ShouldRetry() bool {
	return true
}

func NewRetryableLogStorageError(err error) RetryableLogStorageError {
	return RetryableLogStorageError{err: err}
}
