package storer

type LogStorageError interface {
	error
	ShouldRetry() bool
}

type logStorageError struct {
	err         error
	shouldRetry bool
}

func NewLogStorageError(err error, shouldRetry bool) logStorageError {
	return logStorageError{
		err:         err,
		shouldRetry: shouldRetry,
	}
}

func (e logStorageError) Error() string {
	return e.err.Error()
}

func (e logStorageError) ShouldRetry() bool {
	return e.shouldRetry
}
