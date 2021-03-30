package types

import (
	"time"

	"gorm.io/gorm"
)

// DeveloperApp 开发者第三方app数据结构
type DeveloperApp struct {
	ID               uint64         `gorm:"column:id;primaryKey"                            json:"id,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at"                               json:"created_at,omitempty"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"                               json:"updated_at,omitempty"`
	DeletedAt        gorm.DeletedAt `gorm:"column:delete_at;index"                          json:"-"`
	Name             string         `gorm:"column:name;index:idx_developer_app,unique"      json:"name,omitempty"`
	TenantID         uint64         `gorm:"column:tenant_id;index:idx_developer_app,unique" json:"tenant_id,omitempty"`
	OAuthAccessToken string         `gorm:"column:oauth_access_token;unique"                json:"oauth_access_token,omitempty"`
}
