package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllFeature 获取所有的前端功能
func GetAllFeature(db *gorm.DB) ([]*types.Feature, error) {
	var feature []*types.Feature
	if err := db.Find(&feature).Error; err != nil {
		return nil, err
	}

	return feature, nil
}

// ListFeatureById 获取一批前端功能
func ListFeatureById(db *gorm.DB, id []uint64) ([]*types.Feature, error) {
	var feature []*types.Feature
	if err := db.Find(&feature, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	return feature, nil
}

// TakeFeature 根据租户获取前端功能
func TakeFeature(db *gorm.DB, feature *types.Feature) error {
	return db.Where(feature).Take(feature).Error
}

// UpsertFeature 创建或者更新前端功能
func UpsertFeature(db *gorm.DB, feature *types.Feature) error {
	if feature.ID == 0 {
		return insertOrRecover(db, feature)
	}

	return db.Updates(feature).Error
}

// DeleteFeatureById 根据id从数据库中软删除对应的前端功能
func DeleteFeatureById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Feature{}, id).Error
}
