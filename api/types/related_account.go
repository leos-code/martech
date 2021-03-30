package types

import (
	"gorm.io/gorm"
	"time"
)

// RelatedAccount 渠道帐号信息
type RelatedAccount struct {
	ID        uint64         `gorm:"column:id;primaryKey"             json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at"                json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at"                json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_at;index"           json:"-"`
}
