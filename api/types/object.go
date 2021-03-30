package types

import (
	"time"

	"gorm.io/gorm"
)

// Object 数据实体信息（资源）
type Object struct {
	ID        uint64         `gorm:"column:id;primaryKey"                     json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                        json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                        json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"                   json:"-"`
	Name      string         `gorm:"column:name;index:idx_object,unique"      json:"name,omitempty"`
	Type      ObjectType     `gorm:"column:type"                              json:"type,omitempty"`
	Object    string         `gorm:"column:object"                            json:"object,omitempty"`
	TenantID  uint64         `gorm:"column:tenant_id;index:idx_object,unique" json:"tenant_id,omitempty"`
	ParentID  uint64         `gorm:"-"                                        json:"parent_id"`
}

type ObjectType string

const (
	ObjectTypeDefault    ObjectType = "default"
	ObjectTypeObject     ObjectType = "object"
	ObjectTypeRole       ObjectType = "role"
	ObjectTypeFrontend   ObjectType = "frontend"
	ObjectTypeBackend    ObjectType = "backend"
	ObjectTypeFeature    ObjectType = "feature"
	ObjectTypeMaterial   ObjectType = "material"
	ObjectTypeReport     ObjectType = "report"
	ObjectTypeAdvertiser ObjectType = "advertiser"
	ObjectTypeSchema     ObjectType = "schema"
	ObjectTypeTargeting  ObjectType = "targeting"
)

func (o *Object) GetObject() string {
	return o.Object
}
