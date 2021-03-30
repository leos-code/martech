package router

import (
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

var (
	DefaultCasbinModelConfPath = "casbin_model.conf"
)

func createEnforcer(option *config.CasbinOption) (casbin.IEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(orm.GetDB())
	if err != nil {
		return nil, err
	}
	return casbin.NewSyncedEnforcer(option.ModelConfPath, adapter)
}

func setupEnforcer(router gin.IRouter, option *config.CasbinOption) error {
	if option == nil {
		option = &config.CasbinOption{
			ModelConfPath: DefaultCasbinModelConfPath,
		}
	}

	model, err := ioutil.ReadFile(option.ModelConfPath)
	if err != nil {
		return err
	}

	enforcer, err := createEnforcer(option)
	if err != nil {
		return err
	}

	router.Use(handler.EnforcerMiddleware(enforcer, string(model)))
	return nil
}
