package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StrategyListHandler 获取strategy列表的接口
func StrategyListHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.GetAllBindStrategy(db, option)
	})
}

// StrategyEditHandler 编辑strategy的接口
func StrategyEditHandler(c *gin.Context) {
	item := &types.BindStrategy{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		return orm.UpsertBindStrategy(db, item)
	})
}

// StrategyDeleteHandler 删除strategy的接口
func StrategyDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, orm.DeleteBindStrategyById)
}
