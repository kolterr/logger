package logger

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"io"
	"bytes"
)

const (
	bufferSize = 1024 * 1024
)

type FileHandler struct {
	file          *os.File
	lock          sync.Mutex
	fileName      string
	prefix        string
	maxBufferSize int
	transportChan chan []byte
	bufferPool    []byte
	close         bool
	errorChan     chan error
}

func (f *FileHandler) SetBufferSize(size int) {
	f.maxBufferSize = size
}

func (f *FileHandler) run() {
	for {
		select {
		case data := <-f.transportChan:
			b:=bytes.NewBuffer(f.bufferPool)
			b.Write(data)
			f.bufferPool = b.Bytes()
			if len(f.bufferPool)>=f.maxBufferSize {
				f.doWrite()
			}
		case err := <-f.errorChan:
			f.file.Close()
			io.WriteString(os.Stdout,fmt.Sprintf("write error %+v",err))
		}
	}
}

func (f *FileHandler) Write(data []byte) error {
	if f.close {
		return nil
	}
	if err := f.nextFile(); err != nil {
		f.errorChan<-err
		return err
	}
	f.transportChan <- data
	return nil
}

func (f *FileHandler) doWrite() {
	if len(f.bufferPool)==0{
		return
	}
	if _,err:=f.file.Write(f.bufferPool); err !=nil{
		f.errorChan <- err
	}
	f.bufferPool = make([]byte,0)
}

func (f *FileHandler) Close() error {
	f.close = true
	f.doWrite() // 防止数据丢失
	return f.file.Close()
}


// nextFile 检查是否跨时间
func (f *FileHandler) nextFile() error {
	fileName := strings.TrimSuffix(f.fileName, path.Ext(f.fileName))
	if time.Now().Format(timeFormatShort) != fileName { // 跨日
		f.file.Close()
		f.packAge() //  打包
		fileName := fmt.Sprintf(f.prefix+"%v.log", time.Now().Format(timeFormatShort))
		file, err := createFile(fileName)
		if err != nil {
			return err
		}
		f.file = file
		f.fileName = fileName
	}
	return nil
}

// 日志打包 按月、年
func (f *FileHandler) packAge() {
	oldName := f.fileName
	fileName := strings.TrimSuffix(oldName, path.Ext(oldName))
	last := strings.LastIndexAny(f.fileName, "-")
	exec.Command("tar", "czf", fileName[:last]+".tar.gz", oldName).Run()
	os.Remove(f.fileName)
	for i := 1; i <= 31; i++ { // 不支持正则表达式
		fileName := fmt.Sprintf("%s-%s.log", fileName[:last], i)
		exec.Command("tar", "rf", fileName[:last]+".tar.gz", fileName).Run()
		os.Remove(fileName)
	}
}

func NewFileHandler(prefix string) (*FileHandler, error) {
	fileName := fmt.Sprintf(prefix+"%v.log", time.Now().Format(timeFormatShort))
	file, err := createFile(fileName)
	if err != nil {
		return nil, err
	}
	o := &FileHandler{
		file:          file,
		fileName:      fileName,
		prefix:        prefix,
		transportChan: make(chan []byte,bufferSize),
		errorChan:     make(chan error, 1),
		bufferPool:    make([]byte, 0),
		maxBufferSize:1024*1024,
	}
	go o.run()
	return o, nil
}

func createFile(fileName string) (*os.File, error) {
	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, 0755)
	return os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
}
