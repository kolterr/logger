package logger_test

import (
	"testing"
	"git.ycyz.org/kolter/logger"
)

func TestSimpleLogger(t *testing.T) {
	//std := logger.NewStdHandler(os.Stdout)
	//l := logger.NewSimpleLogger(std)

	//l.Debug("debug")
	//l.Warn("warn")
	//l.Error("Error")
	//l.Fatal("fatal")
	//l.Info("info")

	f, _ := logger.NewFileHandler("")
	ll := logger.NewSimpleLogger(f)
	ll.SetLevel(logger.LevelDebug)
	ll.Debug("debug")
	ll.Warn("warn")
}
