package logger

import "github.com/sirupsen/logrus"

type LogrusLogger struct {
	*logrus.Logger
}

type LogrusEntry struct {
	*logrus.Entry
}

func NewLogrusLogger() Logger {
	return LogrusLogger{
		logrus.New(),
	}
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
