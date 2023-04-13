package zap

import (
	"go.uber.org/zap"
)

// Re-export relevant zap functions
var String = zap.String
var Strings = zap.Strings
var Error = zap.Error
var Int = zap.Int
var Duration = zap.Duration
