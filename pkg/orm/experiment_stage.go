package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllExperimentStage 从数据库中获取所有的实验阶段数据
func GetAllExperimentStage(db *gorm.DB) ([]*types.ExperimentStage, error) {
	var stage []*types.ExperimentStage
	if err := db.Find(&stage).Error; err != nil {
		return nil, err
	}
	return stage, nil
}

// GetExperimentStageById 根据id从数据库中获取到对应的实验阶段的数据
func GetExperimentStageById(db *gorm.DB, id uint64) (*types.ExperimentStage, error) {
	stage := &types.ExperimentStage{}
	if err := db.Take(stage, id).Error; err != nil {
		return nil, err
	}
	return stage, nil
}

// LoadStageExperimentItem 加载实验阶段下对应的ExperimentItem
func LoadStageExperimentItem(db *gorm.DB, stage *types.ExperimentStage) error {
	return db.Model(stage).Association("ExperimentItem").Find(&stage.ExperimentItem)
}

// UpsertExperimentStage 插入或更新实验阶段的数据
func UpsertExperimentStage(db *gorm.DB, stage *types.ExperimentStage) error {
	if stage.ID == 0 {
		return db.Create(stage).Error
	}
	return db.Updates(stage).Error
}

// DeleteExperimentStageById 根据id从数据库中软删除对应的实验阶段
func DeleteExperimentStageById(db *gorm.DB, id uint) error {
	return db.Delete(&types.ExperimentStage{}, id).Error
}

// 获取最新的版本
func GetLatestVersionOfGroup(db *gorm.DB, groupId uint64) (uint64, error) {
	type maxStruct struct {
		Max uint64 `gorm:"column:max"`
	}
	t := &maxStruct{}
	if err := db.
		Raw("SELECT MAX(version) AS max FROM experiment_stages WHERE experiment_group_id = ?", groupId).
		Scan(t).Error; err != nil {
		return 0, err
	}
	return t.Max, nil
}
