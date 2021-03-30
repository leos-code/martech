package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MaterialGetHandler 获取素材列表
func MaterialGetHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.MaterialPage(db, option)
	})
}

// MaterialPostHandler 编辑素材
func MaterialPostHandler(c *gin.Context) {
	material := &types.Material{}
	dbEditHandler(c, material, func(db *gorm.DB) error {
		return orm.Upsert(db, material)
	})
}

// MaterialDeleteManyHandler 删除素材
func MaterialDeleteManyHandler(c *gin.Context) {
	deleteMany := &types.DeleteMany{}
	dbDeleteMultiple(c, deleteMany, func(db *gorm.DB) error {
		return orm.MaterialDelete(db, deleteMany.ID...)
	})
}
