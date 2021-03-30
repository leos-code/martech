package config

import (
	"fmt"

	"github.com/tencentad/union-marketing-go-sdk/api/manager"
	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	sdkConfig "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	sdkOrm "github.com/tencentad/union-marketing-go-sdk/pkg/sdk/orm"
	"github.com/tencentad/martech/pkg/common/s3x"
	"github.com/tencentad/martech/pkg/config"
	"github.com/tencentad/martech/pkg/logger"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
)

// CasbinOption casbin配置
type CasbinOption struct {
	ModelConfPath string `json:"model_conf_path"`
}

// LoginOption 单个登录配置
type LoginOption struct {
	AppID        string `json:"app_id"`
	Token        string `json:"token"`
	CallBackHost string `json:"callback_host"`
}

// LoginOptions 登录配置
type LoginOptions struct {
	Insecure bool         `json:"insecure"`
	Rio      *LoginOption `json:"rio"`
}

type SessionStoreType string

const (
	SessionStoreTypeMemStore SessionStoreType = "memstore"
	SessionStoreTypeRedis    SessionStoreType = "redis"
)

// SessionOption Session配置
type SessionOption struct {
	Type     SessionStoreType `json:"type"`
	Address  string           `json:"address"`
	Password string           `json:"password"`
}

// SDKOption SDK配置
type SDKOption struct {
	DB  *sdkOrm.Option    `json:"db"`
	AMS *sdkConfig.Config `json:"ams"`
	OceanEngine *sdkConfig.Config `json:"ocean_engine"`
}

// WebConfiguration web服务的配置
type WebConfiguration struct {
	ServerAddress  string         `json:"server_address"`
	MetricsAddress string         `json:"metrics_address"`
	WebFileRoot    string         `json:"web_file_root"`
	DocsFileRoot   string         `json:"docs_file_root"`
	Casbin         *CasbinOption  `json:"casbin"`
	Login          *LoginOptions  `json:"login"`
	Session        *SessionOption `json:"session"`
}

// APIOption api配置
type APIOption struct {
	RtaExpListUrl string `json:"rta_exp_list_url"`
	RtaExpDataUrl string `json:"rta_exp_data_url"`
}

// MaterialFileOption 素材文件配置
type MaterialFileOption struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}

// Configuration 所有配置的集合
type Configuration struct {
	Host         string              `json:"-"`
	DB           *orm.Option         `json:"db"`
	Log          *logger.LogOption   `json:"log"`
	Web          *WebConfiguration   `json:"web"`
	API          *APIOption          `json:"api"`
	S3           *s3x.Option         `json:"s3"`
	MaterialFile *MaterialFileOption `json:"material_file"`
	SDK          *SDKOption          `json:"sdk"`
}

// Load 加载配置
func Load(configFile ...string) (*Configuration, error) {
	if err := config.Init(configFile...); err != nil {
		return nil, err
	}

	// host env
	host := config.GetIP()
	configuration := &Configuration{Host: host}

	if err := config.Scan(configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}

// Setup 使用配置
func Setup(configuration *Configuration) error {
	// log configuration
	if configuration.Log != nil {
		configuration.Log.Host = configuration.Host
	}
	if err := logger.Init(configuration.Log); err != nil {
		return err
	}

	// db configuration
	if db := orm.GetDB(configuration.DB); db == nil {
		return fmt.Errorf("db is nil")
	}

	// api configuration
	if configuration.API != nil {
		apiOption := configuration.API
		config.AssignStringIfNotEmpty(apiOption.RtaExpListUrl, &api.RtaExpListUrl)
		config.AssignStringIfNotEmpty(apiOption.RtaExpDataUrl, &api.RtaExpDataUrl)
	}

	// s3 configuration
	if configuration.S3 != nil {
		if s3 := s3x.GetS3(configuration.S3); s3 == nil {
			return fmt.Errorf("s3 is nil")
		}
	}

	// sdk configuration
	if configuration.SDK != nil {
		if db := sdkOrm.GetDB(configuration.SDK.DB); db == nil {
			return fmt.Errorf("sdk db not init ok")
		}

		if configuration.SDK.AMS != nil {
			if err := manager.Register(sdk.AMS, configuration.SDK.AMS); err != nil {
				return err
			}
		}
		if configuration.SDK.OceanEngine != nil {
			if err := manager.Register(sdk.OceanEngine, configuration.SDK.OceanEngine); err != nil {
				return err
			}
		}
	}

	return nil
}
