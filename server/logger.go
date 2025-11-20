package server

import (
	"os"
	"strings"

	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(cfg *props.LoggingConfig, serviceName string) (*zap.Logger, error) {

	// 1. Select encoder
	var encoder zapcore.Encoder

	encoding := strings.ToLower(cfg.Format)
	if encoding == "json" || encoding == "" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// 2. Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level,
	)

	logger := zap.New(core,
		zap.AddCaller(),
		zap.Fields(zap.String("service", serviceName)),
	)

	return logger, nil
}
