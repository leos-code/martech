package handler

import (
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizationInfoHandler 拉取组织信息
func OrganizationInfoHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		current, _ := currentTenantSessionGet(c)
		_ = orm.TakeTenant(db, current)
		return current, nil
	})
}
