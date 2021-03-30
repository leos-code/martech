package loader

import (
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/tencentad/martech/pkg/common/hdfsx"
	"github.com/colinmarc/hdfs"
	"github.com/golang/glog"
)

// HdfsFileLoader HDFS文件加载器
type HdfsFileLoader struct {
	client   *hdfs.Client
	filePath string
	fileLoad fileLoadFunc
	alloc    allocFunc

	detectModify time.Time
	lastModify   time.Time
}

// fileLoadFunc 从reader中读取数据，生成对应的结构体
type fileLoadFunc func(reader io.Reader, i interface{}) error

// NewHdfsFileLoader
func NewHdfsFileLoader(client *hdfs.Client, filePath string, loader fileLoadFunc, alloc allocFunc) *HdfsFileLoader {
	return &HdfsFileLoader{
		client:   client,
		filePath: filePath,
		fileLoad: loader,
		alloc:    alloc,
	}
}

// NewHdfsHotLoaderWithDefaultClient 使用默认的HDFS客户端，创建HDFS文件加载器
func NewHdfsHotLoaderWithDefaultClient(filePath string, alloc allocFunc, loader fileLoadFunc) *HdfsFileLoader {
	return NewHdfsFileLoader(hdfsx.GetHDFSClient(), filePath, loader, alloc)
}

// Alloc
func (l *HdfsFileLoader) Alloc() interface{} {
	return l.alloc()
}

// DetectNewFile
func (l *HdfsFileLoader) DetectNewFile() (string, bool) {
	statInfo, err := l.client.Stat(l.filePath)
	if err != nil {
		glog.Errorf("failed to stat[%s], err: %v", l.filePath, err)
		return "", false
	}

	l.detectModify = statInfo.ModTime()
	return l.filePath, l.detectModify.After(l.lastModify)
}

// Load
func (l *HdfsFileLoader) Load(filePath string, i interface{}) error {
	reader, err := l.client.Open(filePath)
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

// ParseHDFSPath 解析hdfs地址，支持
// 1. hdfs://master_url:9200/path
// 2. /path
func ParseHDFSPath(fullPath string) (masterUrl string, path string, err error) {
	if !strings.HasPrefix(fullPath, hdfsPathPrefix) {
		path = fullPath
		return
	}

	var u *url.URL
	u, err = url.Parse(fullPath)
	if err != nil {
		return
	}
	masterUrl = u.Host
	path = u.RequestURI()
	return
}
