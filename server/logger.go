package server

import (
	"strings"

	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(cfg *props.LoggingConfig, serviceName string) (*zap.Logger, error) {

	level := zap.InfoLevel
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	zapCfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			MessageKey:     "msg",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stack",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
		},
	}

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	// global service field
	logger = logger.With(
		zap.String("service", serviceName),
	)

	return logger, nil
}
