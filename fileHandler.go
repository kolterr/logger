package logger

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
	"sync"
	"path"
	"strings"
	"os/exec"
	"sync/atomic"
)

type FileHandler struct {
	file     *os.File
	lock     sync.Mutex
	fileName string
	prefix   string
	maxLine  int
	maxSize  int
	line     *int64
}

func (f *FileHandler) Write(data []byte) (int, error) {
	if err := f.nextFile(); err != nil {
		return 0, err
	}
	n, err := f.file.Write(data)
	atomic.AddInt64(f.line, 1)
	return n, err
}

func (f *FileHandler) Close() error {
	return f.file.Close()
}

func (f *FileHandler) SetLine(line int) {
	f.maxLine = line
}

func (f *FileHandler) SetMaxSize(maxSize int) {
	f.maxSize = maxSize
}

// nextFile 检查是否跨日
func (f *FileHandler) nextFile() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	fileName := strings.TrimSuffix(f.fileName, path.Ext(f.fileName))
	if time.Now().Format(timeFormatShort) != fileName { // 跨日
		f.file.Close()
		fileName := fmt.Sprintf(f.prefix+"%v.log", time.Now().Format(timeFormatShort))
		file, err := createFile(fileName)
		if err != nil {
			return err
		}
		f.file = file
		f.fileName = fileName
	}
	f.packAge() //  打包
	return nil
}

// 日志打包 按月、年
func (f *FileHandler) packAge() {
	oldName := f.fileName
	fileName := strings.TrimSuffix(oldName, path.Ext(oldName))
	//exec.Command("tar", "czf", fileName+".tar.gz", f.fileName).Run()
	os.Remove(f.fileName)
	last := strings.LastIndexAny(f.fileName, "-")
	exec.Command("tar", "rf", f.fileName[:last]+".tar", f.fileName+".tar.gz").Run()
	fmt.Println(f.fileName[:last], fileName)
}

func NewFileHandler(prefix string) (*FileHandler, error) {
	fileName := fmt.Sprintf(prefix+"%v.log", time.Now().Format(timeFormatShort))
	file, err := createFile(fileName)
	if err != nil {
		return nil, err
	}
	return &FileHandler{
		file:     file,
		fileName: fileName,
		prefix:   prefix,
	}, nil
}

func createFile(fileName string) (*os.File, error) {
	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, 0755)
	return os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
}
