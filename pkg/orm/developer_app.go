package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

func GetAllDeveloperApp(db *gorm.DB) ([]*types.DeveloperApp, error) {
	var developerApps []*types.DeveloperApp
	if err := db.Find(&developerApps).Error; err != nil {
		return nil, err
	}
	return developerApps, nil
}

func TakeDeveloperApp(db *gorm.DB, developerApp *types.DeveloperApp) error {
	return db.Where(developerApp).Take(developerApp).Error
}

func UpsertDeveloperApp(db *gorm.DB, developerApp *types.DeveloperApp) error {
	if developerApp.ID == 0 {
		return insertOrRecover(db, developerApp)
	}
	return db.Updates(developerApp).Error
}

func DeleteDeveloperApp(db *gorm.DB, id uint64) error {
	return db.Delete(&types.DeveloperApp{ID: id}).Error
}
