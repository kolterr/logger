package logger_test

import (
	"testing"
	"git.ycyz.org/kolter/logger"
)

func TestSimpleLogger(t *testing.T) {
	//std := logger.NewStdHandler(os.Stdout)
	//l := logger.NewSimpleLogger(std)
	//l.SetLevel(logger.LevelDebug)
	//l.Debug("debug")
	//l.Warn("warn")
	//l.Error("Error")
	//l.Fatal("fatal")
	//l.Info("info")

	f, _ := logger.NewFileHandler("")
	ll := logger.NewSimpleLogger(f)
	//ll.Debug("debug")
	ll.Warn("warn")
}

func testLogger(t *testing.T, logger logger.Logger) {

}
