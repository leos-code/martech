package types

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

// BindStrategy 定向绑定的RTA策略信息
type BindStrategy struct {
	ID          uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name        string         `gorm:"column:name"            json:"name,omitempty"`
	Platform    ADPlatformType `gorm:"column:platform"        json:"platform,omitempty"`
	TargetingID uint64         `gorm:"column:targeting_id"    json:"targeting_id,omitempty"`
	Strategy    *Strategy      `gorm:"column:strategy"        json:"strategy,omitempty"`
}

// Strategy RTA策略相关信息
type Strategy struct {
	StrategyID          []string `json:"strategy_id,omitempty"`
	AdvertiserID        []string `json:"advertiser_id,omitempty"`
	CampaignID          []string `json:"campaign_id,omitempty"`
	ExcludeStrategyID   []string `json:"exclude_strategy_id,omitempty"`
	ExcludeAdvertiserID []string `json:"exclude_advertiser_id,omitempty"`
	ExcludeCampaignID   []string `json:"exclude_campaign_id,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (s *Strategy) Scan(value interface{}) error {
	return scan(value, s)
}

// Value return json value, implement driver.Valuer interface
func (s Strategy) Value() (driver.Value, error) {
	return value(s)
}
