package logger

import (
	"fmt"
	"testing"
)

type TestLogger struct {
	*testing.T
}

func NewTestLogger(t *testing.T) Logger {
	return TestLogger{t}
}

func (logger TestLogger) Infoln(args ...interface{}) {
	logger.T.Log(args...)
}

func (logger TestLogger) Infof(format string, args ...interface{}) {
	logger.T.Logf(format, args...)
}

func (logger TestLogger) Debugln(args ...interface{}) {
	logger.T.Log("DEBUG: ", args)
}

func (logger TestLogger) Debugf(format string, args ...interface{}) {
	debugFormat := fmt.Sprintf("DEBUG: %v", format)
	logger.T.Logf(debugFormat, args...)
}

func (logger TestLogger) AddContext(key string, value string) Logger {
	return logger
}

func (logger TestLogger) EnableDebugger() {
}
