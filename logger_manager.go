package logger

type LogManager struct {
	loggers []Logger
}

func (l *LogManager) Register(loggers ...Logger) {
	l.loggers = append(l.loggers, loggers...)
}

func (l *LogManager) Debug(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Debug(format, arg)
	}
}

func (l *LogManager) Error(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Error(format, arg)
	}
}
func (l *LogManager) Fatal(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Fatal(format, arg)
	}
}
func (l *LogManager) Info(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Info(format, arg)
	}
}
func (l *LogManager) Warn(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Warn(format, arg)
	}
}
func (l *LogManager) Trace(format string, arg ...interface{}) {
	for _, v := range l.loggers {
		v.Trace(format, arg)
	}
}
