package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllTenant 获取所有的租户
func GetAllTenant(db *gorm.DB) ([]*types.Tenant, error) {
	var tenant []*types.Tenant
	if err := db.Find(&tenant).Error; err != nil {
		return nil, err
	}

	return tenant, nil
}

// ListTenantById 获取一批租户信息
func ListTenantById(db *gorm.DB, id []uint64) ([]*types.Tenant, error) {
	var tenant []*types.Tenant
	if err := db.Find(&tenant, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	return tenant, nil
}

// TakeTenant 根据租户获取租户信息
func TakeTenant(db *gorm.DB, tenant *types.Tenant) error {
	return db.Where(tenant).Take(tenant).Error
}

// UpsertTenant 创建或者更新租户的信息
func UpsertTenant(db *gorm.DB, tenant *types.Tenant) error {
	if tenant.ID == 0 {
		return insertOrRecover(db, tenant)
	}

	return db.Updates(tenant).Error
}

// DeleteTenantById 根据id从数据库中软删除对应的租户
func DeleteTenantById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Tenant{}, id).Error
}
