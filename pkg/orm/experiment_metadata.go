package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllExperimentMetadata 从数据库中获取所有的ExperimentMetadata的数据
func GetAllExperimentMetadata(db *gorm.DB) ([]*types.ExperimentMetadata, error) {
	var experimentMetadata []*types.ExperimentMetadata
	if err := db.Find(&experimentMetadata).Error; err != nil {
		return nil, err
	}
	return experimentMetadata, nil
}

// GetExperimentMetadataById 根据id从数据库中获取所有的ExperimentMetadata的信息
func GetExperimentMetadataById(db *gorm.DB, id uint64) (*types.ExperimentMetadata, error) {
	experimentMetadata := &types.ExperimentMetadata{}
	if err := db.Take(experimentMetadata, id).Error; err != nil {
		return nil, err
	}
	return experimentMetadata, nil
}

// LoadExperimentMetadataExperimentParameter 加载ExperimentMetadata相关的实验参数
func LoadExperimentMetadataExperimentParameter(db *gorm.DB, metadata *types.ExperimentMetadata) error {
	var parameter types.ExperimentParameter
	if err := db.Model(metadata).Association("ExperimentParameter").Find(&parameter); err != nil {
		return err
	}
	metadata.ExperimentParameter = &parameter
	return nil
}

// UpsertExperimentMetadata 插入或者更新ExperimentMetadata的数据
func UpsertExperimentMetadata(db *gorm.DB, experimentMetadata *types.ExperimentMetadata) error {
	if experimentMetadata.ID == 0 {
		return db.Create(experimentMetadata).Error
	}
	return db.Updates(experimentMetadata).Error
}

// DeleteExperimentMetadataById 根据id从数据库中软删除对应的ExperimentMetadata
func DeleteExperimentMetadataById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.ExperimentMetadata{}, id).Error
}

// UnbindExperimentMetadata 删除当前metadata和ExperimentItem的关联关系
func UnbindExperimentMetadata(db *gorm.DB, metadata *types.ExperimentMetadata) error {
	return db.Model(metadata).UpdateColumn("experiment_item_id", nil).Error
}
