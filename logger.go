package logger

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var (
	prefix []string = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	colors          = [6]string{"\033[41;37m", "\033[36m", "\033[32m", "\033[33m", "\033[5;41;33m", "\033[41;37m"}
)

const (
	timeFormatDefault = "2006-01-02 15:04:05"
	timeFormatShort   = "2006-01-02"
	call3Depth        = 2
)

type Logger interface {
	SetLevel(level int)
	outPut(call3Depth, level int, format string, arg ...interface{})
	run()
	SetFlags(int)
	Debug(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Trace(format string, v ...interface{})
	NoColor()
}

func SetColors(c [6]string) {
	colors = c
}
