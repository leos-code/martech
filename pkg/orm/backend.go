package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllBackend 获取所有的后台API
func GetAllBackend(db *gorm.DB) ([]*types.Backend, error) {
	var backend []*types.Backend
	if err := db.Find(&backend).Error; err != nil {
		return nil, err
	}

	return backend, nil
}

// ListBackendById 获取一批后台API
func ListBackendById(db *gorm.DB, id []uint64) ([]*types.Backend, error) {
	var backend []*types.Backend
	if err := db.Find(&backend, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	return backend, nil
}

// TakeBackend 根据租户获取后台API
func TakeBackend(db *gorm.DB, backend *types.Backend) error {
	return db.Where(backend).Take(backend).Error
}

// UpsertBackend 创建或者更新后台API
func UpsertBackend(db *gorm.DB, backend *types.Backend) error {
	if backend.ID == 0 {
		return insertOrRecover(db, backend)
	}

	return db.Updates(backend).Error
}

// DeleteBackendById 根据id从数据库中软删除对应的后台API
func DeleteBackendById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Backend{}, id).Error
}
