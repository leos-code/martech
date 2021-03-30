package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllAdvertiser 获取所有的后台API
func GetAllAdvertiser(db *gorm.DB) ([]*types.Advertiser, error) {
	var advertisers []*types.Advertiser
	if err := db.Find(&advertisers).Error; err != nil {
		return nil, err
	}
	return advertisers, nil
}

// TakeAdvertiser 获取一批后台API
func TakeAdvertiser(db *gorm.DB, advertiser *types.Advertiser) error {
	return db.Where(advertiser).Take(advertiser).Error
}

// UpsertAdvertiser 根据账户获取数据的后台API
func UpsertAdvertiser(db *gorm.DB, advertiser *types.Advertiser) error {
	if advertiser.ID == 0 {
		return insertOrRecover(db, advertiser)
	}
	return db.Updates(advertiser).Error
}

// DeleteAdvertiser 根据id从数据库中软删除对应数据的后台API
func DeleteAdvertiser(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Advertiser{}, id).Error
}

func ListAdvertiserByTenantId(db *gorm.DB, tenantId uint64) ([]*types.Advertiser, error) {
	var advertisers []*types.Advertiser
	if err := db.Where(&types.Advertiser{TenantID: tenantId}).Find(&advertisers).Error; err!=nil {
		return nil, err
	}
	return advertisers, nil
}
