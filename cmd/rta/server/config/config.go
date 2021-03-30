package config

import (
	"fmt"

	"github.com/tencentad/martech/cmd/rta/server/logic"
	"github.com/tencentad/martech/pkg/common/s3x"
	"github.com/tencentad/martech/pkg/config"
	"github.com/tencentad/martech/pkg/logger"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/tencentad/martech/pkg/retrieval"
)

type configuration struct {
	DB         *orm.Option                  `json:"db"`
	Log        *logger.LogOption            `json:"log"`
	S3         *s3x.Option                  `json:"s3"`
	Retrieval  *retrieval.ServiceImplOption `json:"retrieval"`
	RTAService *logic.RTAServiceOption      `json:"rta_service"`
}

var Configuration configuration

// Load 加载配置
func Load(configFile ...string) error {
	if err := config.Init(configFile...); err != nil {
		return err
	}

	if err := config.Scan(&Configuration); err != nil {
		return err
	}

	return setUp()
}

func setUp() error {
	if err := logger.Init(Configuration.Log); err != nil {
		return err
	}

	// db configuration
	if db := orm.GetDB(Configuration.DB); db == nil {
		return fmt.Errorf("db is nil")
	}

	if Configuration.S3 != nil {
		if s3 := s3x.GetS3(Configuration.S3); s3 == nil {
			return fmt.Errorf("s3 is nil")
		}
	}

	if Configuration.Retrieval != nil {
		if err := retrieval.InitServiceImpl(Configuration.Retrieval); err != nil {
			return err
		}
	}

	return nil
}
