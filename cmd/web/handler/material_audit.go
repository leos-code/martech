package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MaterialAuditPostHandler 批量提交审核
func MaterialAuditPostHandler(c *gin.Context) {
	list := make([]*types.MaterialAudit, 0)
	dbEditHandler(c, &list, func(db *gorm.DB) error {
		user, err := userSessionGet(c)
		if err != nil {
			return err
		}
		for _, audit := range list {
			audit.UserID = user.ID
		}
		return orm.SubmitMaterialAudit(db, list)
	})
}

// MaterialAuditGetHandler 获取审核列表
func MaterialAuditGetHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.ListMaterialAuditDetail(db, option)
	})
}
