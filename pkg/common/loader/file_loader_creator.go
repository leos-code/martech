package loader

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	hdfsPathPrefix = "hdfs://"
	s3PathPrefix   = "s3://"
)

// CreateFileLoader 创建文件加载器
func CreateFileLoader(path string, alloc allocFunc, loadFunc fileLoadFunc) FileLoader {
	if strings.HasPrefix(path, hdfsPathPrefix) {
		return NewHdfsHotLoaderWithDefaultClient(path, alloc, loadFunc)
	} else if strings.HasPrefix(path, s3PathPrefix) {
		return NewS3FileLoader(path, alloc, loadFunc)
	} else {
		return NewLocalFileLoader(path, alloc, loadFunc)
	}
}

type WaitDone func() error

// CreateFileLoaderWithWaitDone, 创建文件加载器，并返回WaitDone方法
func CreateFileLoaderWithWaitDone(
	path string, alloc allocFunc, loadFunc fileLoadFunc, timeout time.Duration) (FileLoader, WaitDone) {

	loadDone := make(chan struct{}, 1)
	firstRound := false

	newLoadFunc :=  func(reader io.Reader, i interface{}) error {
		if err := loadFunc(reader, i); err != nil {
			return err
		}

		if firstRound {
			loadDone <- struct{}{}
			firstRound = false
		}

		return nil
	}

	wait := func() error {
		tm := time.NewTimer(timeout * time.Second)

		select {
		case <- loadDone:
			return nil
		case <-tm.C:
			return fmt.Errorf("load data timeout")
		}
	}

	return CreateFileLoader(path, alloc, newLoadFunc), wait
}



// CreatePBFileLoader 创建pb格式的文件加载器
func CreatePBFileLoader(path string, alloc allocFunc) FileLoader {
	return CreateFileLoader(path, alloc, LoadPB)
}

// LoadPB 从reader中读取pb数据，发序列化成对应的结构体
func LoadPB(reader io.Reader, i interface{}) error {
	message := i.(proto.Message)

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return proto.Unmarshal(bytes, message)
}
