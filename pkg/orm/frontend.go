package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllFrontend 获取所有的页面实体
func GetAllFrontend(db *gorm.DB) ([]*types.Frontend, error) {
	var frontend []*types.Frontend
	if err := db.Find(&frontend).Error; err != nil {
		return nil, err
	}

	return frontend, nil
}

// ListFrontendById 获取一批页面实体
func ListFrontendById(db *gorm.DB, id []uint64) ([]*types.Frontend, error) {
	var frontend []*types.Frontend
	if err := db.Find(&frontend, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	return frontend, nil
}

// TakeFrontend 根据租户获取页面实体
func TakeFrontend(db *gorm.DB, frontend *types.Frontend) error {
	return db.Where(frontend).Take(frontend).Error
}

// UpsertFrontend 创建或者更新页面实体
func UpsertFrontend(db *gorm.DB, frontend *types.Frontend) error {
	if frontend.ID == 0 {
		return insertOrRecover(db, frontend)
	}

	return db.Updates(frontend).Error
}

// DeleteFrontendById 根据id从数据库中软删除对应的页面实体
func DeleteFrontendById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Frontend{}, id).Error
}
