package utils

import (
	"github.com/tunedmystic/commits.lol/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SetupLogging ...
func SetupLogging() {
	zapconfig := zap.NewDevelopmentConfig()

	zapLevel := zapcore.InfoLevel

	switch config.App.LogLevel {
	case "DEBUG":
		zapLevel = zapcore.DebugLevel
	case "INFO":
		zapLevel = zapcore.InfoLevel
	case "WARN":
		zapLevel = zapcore.WarnLevel
	case "ERROR":
		zapLevel = zapcore.ErrorLevel
	}

	zapconfig.Level.SetLevel(zapLevel)

	logger, err := zapconfig.Build()
	if err != nil {
		panic("Failed to setup logger: " + err.Error())
	}

	zap.ReplaceGlobals(logger)
	zap.L().Debug("Logger initialized with", zap.String("level", config.App.LogLevel))
}
