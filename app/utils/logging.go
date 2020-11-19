package utils

import (
	"github.com/tunedmystic/commits.lol/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getZapLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}

// SetupLogging ...
func SetupLogging() {
	zapconfig := zap.NewDevelopmentConfig()
	zapLevel := getZapLogLevel(config.App.LogLevel)
	zapconfig.Level.SetLevel(zapLevel)

	logger, _ := zapconfig.Build()
	zap.ReplaceGlobals(logger)

	zap.L().Debug("Logger initialized with", zap.String("level", config.App.LogLevel))
}
