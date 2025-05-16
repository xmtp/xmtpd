package retryable_errors

type RetryableError interface {
	error
	ShouldRetry() bool
}

type NonRetryableError struct {
	err error
}

func (e NonRetryableError) Error() string {
	return e.err.Error()
}

func (e NonRetryableError) ShouldRetry() bool {
	return false
}

func NewUnrecoverableLogStorageError(err error) NonRetryableError {
	return NonRetryableError{err: err}
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
