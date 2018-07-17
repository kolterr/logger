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
	"io/ioutil"
	"bytes"
)

type FileHandler struct {
	file     *os.File
	lock     sync.Mutex
	fileName string
	prefix   string
	maxLine  int64
	maxSize  int64
	line     *int64
	size     *int64
}

func (f *FileHandler) Write(data []byte) (int, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.maxLine != 0 && f.maxLine <= *f.line {
		return 0, nil
	}
	if f.maxSize != 0 && int64(len(data)) + *f.size > f.maxSize {
		return 0, nil
	}
	if err := f.nextFile(); err != nil {
		return 0, err
	}
	n, err := f.file.Write(data)
	if err == nil {
		atomic.AddInt64(f.line, 1)
		atomic.AddInt64(f.line, int64(len(data)))
	}
	return n, err
}

func (f *FileHandler) Close() error {
	return f.file.Close()
}

func (f *FileHandler) SetLine(line int64) {
	f.maxLine = line
}

func (f *FileHandler) SetMaxSize(maxSize int64) {
	f.maxSize = maxSize
}

// nextFile 检查是否跨日
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
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	count, err := getFileLine(fileName)
	if err != nil {
		return nil, err
	}
	n, s := new(int64), new(int64)
	*n, *s = count, info.Size()
	return &FileHandler{
		file:     file,
		fileName: fileName,
		prefix:   prefix,
		line:     n,
		size:     s,
	}, nil
}

func createFile(fileName string) (*os.File, error) {
	dir := filepath.Dir(fileName)
	os.MkdirAll(dir, 0755)
	return os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
}

func getFileLine(fileName string) (int64, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}
	n := bytes.Count(b, []byte("\n"))
	return int64(n), nil
}
