package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllExperimentItem 从数据库中获取到所有的ExperimentItem数据
func GetAllExperimentItem(db *gorm.DB) ([]*types.ExperimentItem, error) {
	var experimentItems []*types.ExperimentItem
	if err := db.Find(&experimentItems).Error; err != nil {
		return nil, err
	}
	return experimentItems, nil
}

// GetExperimentItemById 根据id从数据库中获取到相应的ExperimentItem数据
func GetExperimentItemById(db *gorm.DB, id uint64) (*types.ExperimentItem, error) {
	experimentItem := &types.ExperimentItem{}
	if err := db.Take(experimentItem, id).Error; err != nil {
		return nil, err
	}
	return experimentItem, nil
}

// UpsertExperimentItem 更新或者插入ExperimentItem的数据
func UpsertExperimentItem(db *gorm.DB, experimentItem *types.ExperimentItem) error {
	if experimentItem.ID == 0 {
		if err := db.Create(experimentItem).Error; err != nil {
			return err
		}
		return UpdateOuterId(db, experimentItem)
	}
	return db.Updates(experimentItem).Error
}

// DeleteExperimentItemById 根据id从数据库中软删除对应的ExperimentItem
func DeleteExperimentItemById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.ExperimentItem{}, id).Error
}

// LoadItemExperimentMetadata 加载ExperimentItem使用的ExperimentMetadata
func LoadItemExperimentMetadata(db *gorm.DB, item *types.ExperimentItem) error {
	return db.Model(item).Association("ExperimentMetadata").Find(&item.ExperimentMetadata)
}

// LoadItemRtaExp 加载ExperimentItem绑定的Rta实验
func LoadItemRtaExp(db *gorm.DB, item *types.ExperimentItem) error {
	return db.Model(item).Association("RtaExp").Find(&item.RtaExp)
}

// UpdateOuterId 更新item的outerId
func UpdateOuterId(db *gorm.DB, item *types.ExperimentItem) error {
	item.OuterID = item.ID
	return db.Updates(item).Error
}
