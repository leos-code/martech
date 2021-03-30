package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RTATargetingListHandler 获取RTA Targeting 列表
func RTATargetingListHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.TargetingPage(db, option)
	})
}

// RTATargetingEditHandler 编辑RTA Targeting
func RTATargetingEditHandler(c *gin.Context) {
	item := &types.Targeting{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		return orm.UpsertTargeting(db, item)
	})
}

// RTATargetingDeleteHandler 删除RTA Targeting的接口
func RTATargetingDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, orm.DeleteTargetingByID)
}
