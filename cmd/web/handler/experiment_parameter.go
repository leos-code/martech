package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExperimentParameterListHandler 获取实验参数列表的接口
func ExperimentParameterListHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.GetAllExperimentParameter(db, option)
	})
}

// ExperimentParameterEditHandler 编译实验参数的接口
func ExperimentParameterEditHandler(c *gin.Context) {
	item := &types.ExperimentParameter{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		return orm.UpsertExperimentParameter(db, item)
	})
}

// ExperimentParameterDeleteHandler 删除实验参数的接口
func ExperimentParameterDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, orm.DeleteExperimentParameterById)
}
