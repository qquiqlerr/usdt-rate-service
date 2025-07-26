package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func MustLoadLogger(logLevel string) *zap.Logger {
	// Set the log level based on the provided string
	var level zap.AtomicLevel
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		panic(fmt.Sprintf("unknown log level: %s", logLevel))
	}
	// Create a new zap logger with the specified log level
	logger, err := zap.NewProductionConfig().Build(zap.IncreaseLevel(level))
	if err != nil {
		panic(err)
	}

	// Return the created logger
	return logger
}
