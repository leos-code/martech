package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// TakeObject 根据实体获取实体信息
func TakeObject(db *gorm.DB, object *types.Object) error {
	return db.Where(object).Take(object).Error
}

// UpsertObject 创建或者更新实体的信息
func UpsertObject(db *gorm.DB, object *types.Object) error {
	if object.ID == 0 {
		return insertOrRecover(db, object)
	}

	return db.Updates(object).Error
}

// ListObjectById 获取一批实体的信息
func ListObjectById(db *gorm.DB, id []uint64) ([]*types.Object, error) {
	var object []*types.Object
	if err := db.Find(&object, "id IN ?", id).Error; err != nil {
		return nil, err
	}
	return object, nil
}

// ListObjectByTenantId 获取一批角色的信息
func ListObjectByTenantId(db *gorm.DB, tenantID uint64) ([]*types.Object, error) {
	var object []*types.Object
	if err := db.Where(&types.Object{TenantID: tenantID}).Find(&object).Error; err != nil {
		return nil, err
	}
	return object, nil
}

// FindObject 获取一批角色的信息
func FindObject(db *gorm.DB, in *types.Object) ([]*types.Object, error) {
	var object []*types.Object
	if err := db.Where(in).Find(&object).Error; err != nil {
		return nil, err
	}
	return object, nil
}

// GetObjectInDomain new API to replace ListObjectByTenantId API
func GetObjectInDomain(db *gorm.DB, tenantID uint64, objectType ...types.ObjectType) ([]*types.Object, error) {
	var object []*types.Object
	d := db.Where("tenant_id = ?", tenantID)
	if len(objectType) > 0 {
		d = d.Where("type IN ?", objectType)
	}

	if err := d.Find(&object).Error; err != nil {
		return nil, err
	}
	return object, nil
}

// DeleteObjectById 根据id从数据库中软删除对应的实体
func DeleteObjectById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Object{}, id).Error
}
