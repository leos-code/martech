package types

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租户信息
type Tenant struct {
	ID        uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name      string         `gorm:"column:name;unique"     json:"name,omitempty"`
}
