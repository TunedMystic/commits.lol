package utils

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func Test_SetupLogging(t *testing.T) {
	SetupLogging()
}

func Test_GetLogLevel(t *testing.T) {
	type test struct {
		level    string
		expected zapcore.Level
	}

	tests := []test{
		{"DEBUG", zapcore.DebugLevel},
		{"INFO", zapcore.InfoLevel},
		{"WARN", zapcore.WarnLevel},
		{"ERROR", zapcore.ErrorLevel},
		{"BADLEVEL", zapcore.InfoLevel},
	}

	for _, testItem := range tests {
		t.Run(testItem.level, func(t *testing.T) {
			AssertEqual(t, getZapLogLevel(testItem.level), testItem.expected)
		})
	}
}
