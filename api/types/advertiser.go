package types

import (
	"database/sql/driver"
	"time"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"gorm.io/gorm"
)

// Advertiser 平台广告主帐号
type Advertiser struct {
	ID          uint64                    `gorm:"column:id;primaryKey"                         json:"id,omitempty"`
	CreatedAt   time.Time                 `gorm:"column:created_at"                            json:"created_at,omitempty"`
	UpdatedAt   time.Time                 `gorm:"column:updated_at"                            json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt            `gorm:"column:delete_at;index"                       json:"-"`
	Platform    sdk.MarketingPlatformType `gorm:"column:platform"                              json:"platform,omitempty"`
	Name        string                    `gorm:"column:name;index:idx_advertiser,unique"      json:"name,omitempty"`
	TenantID    uint64                    `gorm:"column:tenant_id;index:idx_advertiser,unique" json:"tenant_id,omitempty"`
	ObjectID    uint64                    `gorm:"column:object_id"                             json:"object_id,omitempty"`
	Object      *Object                   `gorm:"foreignKey:ObjectID"                          json:"object,omitempty"`
	AuthAccount *AuthAccount              `gorm:"column:auth_account"                          json:"auth_account,omitempty"`
}

// AdvertiserAuthorization 单个跳转平台广告主帐号结构
type AdvertiserAuthorization struct {
	Platform sdk.MarketingPlatformType `json:"platform"`
	Url      string                    `json:"url"`
}

// AdvertiserAuthorizations 跳转平台广告主帐号结构
type AdvertiserAuthorizations []*AdvertiserAuthorization

type AuthAccount sdk.AuthAccount

// Scan scan value into Jsonb, implements sql.Scanner interface
func (a *AuthAccount) Scan(value interface{}) error {
	return scan(value, a)
}

// Value return json value, implement driver.Valuer interface
func (a AuthAccount) Value() (driver.Value, error) {
	return value(a)
}
