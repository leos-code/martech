package hdfsx

import (
	"sync"

	"github.com/colinmarc/hdfs"
	log "github.com/sirupsen/logrus"
)

var (
	once      sync.Once
	singleton *hdfs.Client
)

// Option HDFS配置信息
type Option struct {
	Address []string `json:"address"`
	User    string   `json:"user"`
}

// GetHDFSClient 获取HDFS客户端
func GetHDFSClient(options ...*Option) *hdfs.Client {
	once.Do(func() {
		var opt *Option = nil
		if len(options) > 0 {
			opt = options[0]
		}
		if opt == nil {
			log.Errorf("call GetHDFSClient without opt for the first time")
			return
		}
		var err error
		singleton, err = hdfs.NewClient(
			hdfs.ClientOptions{
				Addresses: opt.Address,
				User:      opt.User,
			})
		if err != nil {
			log.Errorf("failed to init hdfs client, err: %v", err)
		}
	})

	return singleton
}
