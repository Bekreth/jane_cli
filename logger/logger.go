package logger

type Logger interface {
	Infoln(args ...interface{})
	Infof(format string, args ...interface{})
	Debugln(args ...interface{})
	Debugf(format string, args ...interface{})
	AddContext(key string, value string) Logger
	EnableDebugger()
}
