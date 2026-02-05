package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapWriter struct {
	logger *zap.Logger
	level  zapcore.Level
}

func Init(mode string) (*zap.Logger, error) {
	normalized := strings.ToLower(strings.TrimSpace(mode))
	var cfg zap.Config
	if normalized == "release" || normalized == "prod" || normalized == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	instance, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(instance)
	return instance, nil
}

func Writer(logger *zap.Logger, level zapcore.Level) *zapWriter {
	return &zapWriter{logger: logger, level: level}
}

func (w *zapWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	if msg == "" {
		return len(p), nil
	}
	if w.level >= zapcore.ErrorLevel {
		w.logger.Error(msg)
	} else {
		w.logger.Info(msg)
	}
	return len(p), nil
}
