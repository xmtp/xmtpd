package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	Level    string `long:"level" description:"Log level. Support values: error, warn, info, debug" default:"info"`
	Encoding string `long:"encoding" description:"Log encoding format. Support values: console, json" default:"console"`
}

// Logger embeds zap.Logger so that it inherits its public functions.
type Logger struct {
	*zap.Logger
}

// Named needs to rewrap the result so that zap.Logger doesn't leak out.
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

// With needs to rewrap the result so that zap.Logger doesn't leak out.
func (l *Logger) With(field zap.Field) *Logger {
	return &Logger{l.Logger.With(field)}
}

// NewLogger is the primary logger constructor.
func NewLogger(opts *Options) (*Logger, error) {
	atom := zap.NewAtomicLevel()
	level := zapcore.InfoLevel
	err := level.Set(opts.Level)
	if err != nil {
		return nil, err
	}
	atom.SetLevel(level)
	cfg := zap.Config{
		Encoding:         opts.Encoding,
		Level:            atom,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			NameKey:      "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{l}, nil
}

// NewDevelopmentLogger is used for testing.
func NewDevelopmentLogger(debug bool) (*Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	if !debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{l}, nil
}
