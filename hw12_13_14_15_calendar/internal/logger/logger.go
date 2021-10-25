package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// panics in case of unknow log level
func New(level, file string) *zap.Logger {
	loggerCfg := zap.NewProductionConfig()

	l := strings.ToLower(level)
	switch l {
	case "debug":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "panic":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		panic("Unknown log level")
	}
	if file == "" {
		loggerCfg.OutputPaths = []string{"stderr"}
	} else {
		loggerCfg.OutputPaths = []string{"stderr", file}
	}
	logger, _ := loggerCfg.Build()
	return logger
}
