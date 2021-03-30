package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RtaAccountListHandler 获取RtaAccount列表的接口
func RtaAccountListHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.GetAllRtaAccount(db, option)
	})
}

// RtaAccountEditHandler 编辑RtaAccount信息的接口
func RtaAccountEditHandler(c *gin.Context) {
	item := &types.RtaAccount{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		return orm.UpsertRtaAccount(db, item)
	})
}

// RtaAccountDeleteHandler 删除RtaAccount的接口
func RtaAccountDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, orm.DeleteRtaAccountById)
}

// RtaAccountSyncHandler 同步RtaAccount下绑定的实验的接口
func RtaAccountSyncHandler(c *gin.Context) {
	item := &types.RtaAccount{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		return orm.SyncRtaAccountExp(db, item)
	})
}
