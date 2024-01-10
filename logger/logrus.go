package logger

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	*logrus.Logger
}

type LogrusEntry struct {
	*logrus.Entry
}

func NewLogrusLogger(config Config) (Logger, error) {
	output := LogrusLogger{
		logrus.New(),
	}
	path := filepath.FromSlash(config.Output)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return output, err
	}
	output.SetOutput(file)
	if config.Debugger {
		output.EnableDebugger()
	}
	return output, err
}

func (logger LogrusLogger) AddContext(key string, value string) Logger {
	level := logger.GetLevel()
	entry := logger.WithField(key, value)
	entry.Logger.SetLevel(level)
	return LogrusEntry{entry}
}

func (logger LogrusEntry) AddContext(key string, value string) Logger {
	level := logger.Level
	entry := logger.WithField(key, value)
	entry.Logger.SetLevel(level)
	return LogrusEntry{entry}
}

func (logger LogrusLogger) EnableDebugger() {
	logger.Logger.SetLevel(logrus.DebugLevel)
}

func (logger LogrusEntry) EnableDebugger() {
	logger.Logger.SetLevel(logrus.DebugLevel)
}
