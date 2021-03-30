package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExperimentGroupListHandler 获取实验组列表数据的接口
func ExperimentGroupListHandler(c *gin.Context) {
	dbPageHandler(c, func(db *gorm.DB, option *orm.StatementOption) (interface{}, error) {
		return orm.GetAllExperimentGroup(db, option)
	})
}

// ExperimentGroupGetHandler 获取单个实验组详细信息的接口
func ExperimentGroupGetHandler(c *gin.Context) {
	dbGetByIdHandler(c, func(db *gorm.DB, id uint64) (interface{}, error) {
		group, err := orm.GetExperimentGroupById(db, id)
		if err != nil {
			return nil, err
		}
		if err := orm.FillExperimentGroup(db, group); err != nil {
			return nil, err
		}
		return group, nil
	})
}

// ExperimentGroupEditHandler 编辑实验组信息的接口
func ExperimentGroupEditHandler(c *gin.Context) {
	item := &types.ExperimentGroup{}
	dbGetHandler(c, item, func(db *gorm.DB) (interface{}, error) {
		return item, EditExperimentGroup(db, item)
	})
}

// ExperimentGroupPromptHandler 提升实验组，将草稿提升为运行状态的接口
func ExperimentGroupPromptHandler(c *gin.Context) {
	dbGetByIdHandler(c, func(db *gorm.DB, id uint64) (interface{}, error) {
		group, err := orm.GetExperimentGroupById(db, id)
		if err != nil {
			return nil, err
		}
		if err := PromoteExperimentGroup(db, group, c); err != nil {
			return nil, err
		}
		return group, nil
	})
}

// ExperimentGroupStopHandler 停止实验的接口
func ExperimentGroupStopHandler(c *gin.Context) {
	dbGetByIdHandler(c, func(db *gorm.DB, id uint64) (interface{}, error) {
		group, err := orm.GetExperimentGroupById(db, id)
		if err != nil {
			return nil, err
		}
		if err := StopExperimentGroup(db, group, c); err != nil {
			return nil, err
		}
		return group, nil
	})
}
