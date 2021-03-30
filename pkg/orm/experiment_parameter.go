package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllExperimentParameter从数据库中获取所有的实验参数信息
func GetAllExperimentParameter(db *gorm.DB, option *StatementOption) ([]*types.ExperimentParameter, error) {
	var experimentParameters []*types.ExperimentParameter
	if err := getDBWithOption(db, option).Find(&experimentParameters).Error; err != nil {
		return nil, err
	}
	return experimentParameters, nil
}

// GetExperimentParameterById 根据id从数据库中获取对应的试验参数
func GetExperimentParameterById(db *gorm.DB, id uint64) (*types.ExperimentParameter, error) {
	experimentParameter := &types.ExperimentParameter{}
	if err := db.Take(experimentParameter, id).Error; err != nil {
		return nil, err
	}
	return experimentParameter, nil
}

// UpsertExperimentParameter 插入或者更新实验参数的信息
func UpsertExperimentParameter(db *gorm.DB, experimentParameter *types.ExperimentParameter) error {
	if experimentParameter.ID == 0 {
		return insertOrRecover(db, experimentParameter)
	}
	return db.Updates(experimentParameter).Error
}

// DeleteExperimentParameterById 根据id从数据库中删除对应的实验参数
func DeleteExperimentParameterById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.ExperimentParameter{}, id).Error
}

// LoadParameterExperimentMetadata 加载实验参数相关的ExperimentMetadata
func LoadParameterExperimentMetadata(db *gorm.DB, parameter *types.ExperimentParameter) error {
	return db.Model(parameter).Association("ExperimentMetadata").Find(&parameter.ExperimentMetadata)
}
