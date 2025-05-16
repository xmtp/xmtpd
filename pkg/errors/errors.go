package retryable_errors

type RetryableError interface {
	error
	ShouldRetry() bool
}

type NonRecoverableError struct {
	err error
}

func (e NonRecoverableError) Error() string {
	return e.err.Error()
}

func (e NonRecoverableError) ShouldRetry() bool {
	return false
}

func NewNonRecoverableError(err error) NonRecoverableError {
	return NonRecoverableError{err: err}
}

type RecoverableError struct {
	err error
}

func (e RecoverableError) Error() string {
	return e.err.Error()
}

func (e RecoverableError) ShouldRetry() bool {
	return true
}

func NewRecoverableError(err error) RecoverableError {
	return RecoverableError{err: err}
}
