package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// TakeRole 根据角色获取角色信息
func TakeRole(db *gorm.DB, role *types.Role) error {
	return db.Where(role).Take(role).Error
}

// UpsertRole 创建或者更新角色的信息
func UpsertRole(db *gorm.DB, role *types.Role) error {
	if role.ID == 0 {
		return insertOrRecover(db, role)
	}

	return db.Updates(role).Error
}

// ListRoleById 获取一批角色的信息
func ListRoleById(db *gorm.DB, id []uint64) ([]*types.Role, error) {
	var role []*types.Role
	if err := db.Find(&role, "id IN ?", id).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// ListRoleByTenantId 获取一批角色的信息
func ListRoleByTenantId(db *gorm.DB, tenantID uint64) ([]*types.Role, error) {
	var role []*types.Role
	if err := db.Where(&types.Role{TenantID: tenantID}).Find(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// DeleteRoleById 根据id从数据库中软删除对应的角色
func DeleteRoleById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Role{}, id).Error
}
