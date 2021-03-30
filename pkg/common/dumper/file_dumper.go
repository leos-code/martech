package dumper

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
)

const (
	hdfsPathPrefix = "hdfs://"
	s3PathPrefix   = "s3://"
)

type FileDumper interface {
	Dump(obj interface{}, marshal Marshal) error
}

type Marshal func(interface{}) ([]byte, error)

func DumpToFile(obj interface{}, marshal Marshal, filepath string) error {
	if strings.HasPrefix(filepath, hdfsPathPrefix) {
		return fmt.Errorf("hdfs not support")
	} else if strings.HasPrefix(filepath, s3PathPrefix) {
		return NewS3FileDumper(filepath).Dump(obj, marshal)
	} else {
		d := &LocalFileDumper{filepath: filepath}
		return d.Dump(obj, marshal)
	}
}

func DumpMessageToFile(message proto.Message, filepath string) error {
	return DumpToFile(message, protoMarshal, filepath)
}

var protoMarshal = func(i interface{}) ([]byte, error) {
	msg := i.(proto.Message)
	return proto.Marshal(msg)
}
