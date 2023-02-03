package zap

import (
	"go.uber.org/zap"
)

// Re-export relevant zap functions
var String = zap.String
var Error = zap.Error
var Int = zap.Int
