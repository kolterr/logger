package logger_test

import (
	"testing"
	"git.ycyz.org/kolter/logger"
	"time"
)

func TestSimpleLogger(t *testing.T) {
	//std := logger.NewStdHandler(os.Stdout)
	//l := logger.NewSimpleLogger(std)

	//l.Debug("debug")
	//l.Warn("warn")
	//l.Error("Error")
	//l.Fatal("fatal")
	//l.Info("info")
	ch :=make(chan bool,1)
	f, _ := logger.NewFileHandler("")
	ll := logger.NewSimpleLogger(f)
	for i:=0;i < 6;i++{
		ll.Debug("debug")
		ll.Warn("warn")
	}
	go func() {
		f.Close()
		time.Sleep(time.Second)
		ch<-true
	}()
	<-ch
}
