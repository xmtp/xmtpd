package utils

import (
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildLogger(options config.LogOptions) (*zap.Logger, *zap.Config, error) {
	atom := zap.NewAtomicLevel()
	level := zapcore.InfoLevel
	err := level.Set(options.LogLevel)
	if err != nil {
		return nil, nil, err
	}
	atom.SetLevel(level)

	cfg := zap.Config{
		Encoding:         options.LogEncoding,
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
	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	logger = logger.Named("replication")

	return logger, &cfg, nil
}
