package types

import (
	"time"

	"gorm.io/gorm"
)

// Role 用户角色信息
type Role struct {
	ID        uint64         `gorm:"column:id;primaryKey"                   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"                 json:"-"`
	Name      string         `gorm:"column:name;index:idx_role,unique"      json:"name,omitempty"`
	Object    string         `gorm:"column:object"                          json:"object,omitempty"`
	TenantID  uint64         `gorm:"column:tenant_id;index:idx_role,unique" json:"tenant_id,omitempty"`
	ParentID  uint64         `gorm:"-"                                      json:"parent_id"`
	User      []*User        `gorm:"-"                                      json:"user,omitempty"`
	Policy    []*Policy      `gorm:"-"                                      json:"policy,omitempty"`
}

// Policy 权限表示信息
type Policy struct {
	Object *Object `json:"object,omitempty"`
	Read   bool    `json:"read"`
	Write  bool    `json:"write"`
}

func (r *Role) GetObject() string {
	return r.Object
}
