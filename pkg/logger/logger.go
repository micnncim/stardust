package logger

import (
	"fmt"
	"strings"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	l, err := parseLogLevel(level)
	if err != nil {
		return nil, err
	}

	c := zapdriver.NewProductionConfig()
	c.Level = zap.NewAtomicLevelAt(l)
	return c.Build()
}

func parseLogLevel(levelStr string) (zapcore.Level, error) {
	switch strings.ToUpper(levelStr) {
	case zapcore.DebugLevel.CapitalString():
		return zapcore.DebugLevel, nil
	case zapcore.InfoLevel.CapitalString():
		return zapcore.InfoLevel, nil
	case zapcore.WarnLevel.CapitalString():
		return zapcore.WarnLevel, nil
	case zapcore.ErrorLevel.CapitalString():
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("undefined log level: %s", levelStr)
	}
}
