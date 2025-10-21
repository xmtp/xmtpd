package utils

// Unimplemented panics with a message indicating that the function is not implemented.
func Unimplemented(message string, unusedVariables ...interface{}) {
	panic("unimplemented: " + message)
}

// Unused is a no-op function that takes any number of arguments and does nothing with them.
// Useful for temporarily silencing "unused variable" warnings in development.
func Unused(unusedVariables ...interface{}) {
	// Do nothing
}
