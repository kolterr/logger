package logger

import (
	"sync"
	"fmt"
	"strings"
	"time"
	"runtime"
)

type SimpleLogger struct {
	prefix  string
	level   int
	handler Handler
	msg     chan []byte
	flag    int
	quit    bool
	locker  sync.Locker
}

func (l *SimpleLogger) SetLevel(level int) {
	l.level = level
}

func (l *SimpleLogger) outPut(call3Depth, level int, format string, arg ...interface{}) {
	if level < l.level {
		return
	}
	if strings.Contains(format, "%") && len(arg) > 0 {
		format = fmt.Sprintf(format, arg...)
	}

	t := time.Now().Format(timeFormatDefault)
	_, file, line, ok := runtime.Caller(call3Depth)
	if !ok {
		line = 0
		file = ""
	} else {
		files := strings.Split(file, "/")
		file = files[len(files)-1]
	}
	outPut := fmt.Sprintf("%s %s %s:%d  [%s] %s\n", colors[level], t, file, line, prefix[level], format)
	l.handler.Write([]byte(outPut))
}

func (l *SimpleLogger) SetFlags(f int) {
	l.flag = f
}

func (l *SimpleLogger) Debug(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelDebug, format, arg)
}

func (l *SimpleLogger) Error(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelError, format, arg)
}

func (l *SimpleLogger) Fatal(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelFatal, format, arg)
}
func (l *SimpleLogger) Info(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelInfo, format, arg)
}
func (l *SimpleLogger) Warn(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelWarn, format, arg)
}
func (l *SimpleLogger) Trace(format string, arg ...interface{}) {
	l.outPut(call3Depth, LevelTrace, format, arg)
}

func NewSimpleLogger(handler Handler) Logger {
	l := &SimpleLogger{
		level:   LevelInfo,
		handler: handler,
	}
	return l
}
