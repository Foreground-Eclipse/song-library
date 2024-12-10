package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(logLevel string) (*Logger, error) {
	config := zap.NewProductionConfig()
	switch logLevel {
	case "DEBUG":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "ERROR":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.Encoding = "console"

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

func (l *Logger) LogInfo(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) LogError(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

func (l *Logger) LogDebug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}
