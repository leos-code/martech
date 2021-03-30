package logger

import (
	"io/ioutil"

	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"
)

var (
	DefaultLevel = "info"
)

// ElasticLogOption ElasticLog配置
type ElasticLogOption struct {
	Url       string `json:"url"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	IndexName string `json:"index"`
}

// LogOption 日志配置
type LogOption struct {
	Level      string            `json:"level"`
	DiscardStd bool              `json:"discard_std"`
	Host       string            `json:"-"`
	Elastic    *ElasticLogOption `json:"elastic"`
}

func init() {
	log.SetReportCaller(true)
}

func log4elasticSearch(option *ElasticLogOption, host string) error {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{option.Url},
		Username:  option.UserName,
		Password:  option.Password,
	})

	if err != nil {
		return err
	}

	hook, err := elogrus.NewAsyncElasticHook(client, host, log.InfoLevel, option.IndexName)

	if err != nil {
		return err
	}

	log.AddHook(hook)
	return nil
}

// Init 日志模块初始化
func Init(option *LogOption) error {
	if option == nil {
		return nil
	}

	if option.Level == "" {
		option.Level = DefaultLevel
	}

	if level, err := log.ParseLevel(option.Level); err != nil {
		return err
	} else {
		log.SetLevel(level)
	}

	if option.DiscardStd {
		log.SetOutput(ioutil.Discard)
	}

	if option.Elastic != nil {
		if err := log4elasticSearch(option.Elastic, option.Host); err != nil {
			return err
		}
	}

	return nil
}
