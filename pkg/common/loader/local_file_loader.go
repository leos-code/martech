package loader

import (
	"os"
	"time"

	"github.com/golang/glog"
)

// LocalFileLoader 本地文件加载器
type LocalFileLoader struct {
	filePath string
	fileLoad fileLoadFunc
	alloc    allocFunc

	lastModify   time.Time
	detectModify time.Time
}

// NewLocalFileLoader
func NewLocalFileLoader(filePath string, alloc allocFunc, loader fileLoadFunc) *LocalFileLoader {
	return &LocalFileLoader{
		filePath: filePath,
		fileLoad: loader,
		alloc:    alloc,
	}
}

// DetectNewFile implement FileLoader
func (l *LocalFileLoader) DetectNewFile() (string, bool) {
	statInfo, err := os.Stat(l.filePath)
	if err != nil {
		glog.Errorf("failed to stat[%s], err: %v", l.filePath, err)
		return "", false
	}

	l.detectModify = statInfo.ModTime()
	return l.filePath, l.detectModify.After(l.lastModify)
}

// Alloc implement FileLoader
func (l *LocalFileLoader) Alloc() interface{} {
	return l.alloc()
}

// Load implement FileLoader
func (l *LocalFileLoader) Load(filePath string, i interface{}) error {
	reader, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = reader.Close()
	}()

	err = l.fileLoad(reader, i)
	if err != nil {
		return err
	}

	l.lastModify = l.detectModify
	return nil
}
