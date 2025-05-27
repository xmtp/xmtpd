package testutils

import (
	"path/filepath"
	"runtime"
)

func GetScriptPath(paths ...string) string {
	// Use runtime.Caller to get the absolute path of the file that calls this function.
	_, filename, _, _ := runtime.Caller(1) // 1 to get the caller of the function
	baseDir := filepath.Dir(filename)
	return filepath.Join(append([]string{baseDir}, paths...)...)
}
