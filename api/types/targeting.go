package types

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

var (
	TargetingTemp = &Targeting{}
)

// Targeting RTA定向
type Targeting struct {
	ID            uint64          `gorm:"column:id;primaryKey"                json:"id,omitempty"`
	CreatedAt     time.Time       `gorm:"column:created_at"                   json:"created_at,omitempty"`
	UpdatedAt     time.Time       `gorm:"column:updated_at"                   json:"updated_at,omitempty"`
	DeletedAt     gorm.DeletedAt  `gorm:"column:delete_at;index"              json:"-"`
	Name          string          `gorm:"column:name"                         json:"name,omitempty"`
	TargetingInfo TargetingInfos  `gorm:"column:targeting_info"               json:"targeting_info,omitempty"`
	BindStrategy  []*BindStrategy `gorm:"foreignKey:TargetingID"              json:"bind_strategy,omitempty"`
	Status        TargetingStatus `gorm:"column:status;default:Valid"         json:"status,omitempty"`
	FreqControl   *FreqControl    `gorm:"column:freq_control"                 json:"freq_control,omitempty"`
}

// TargetingStatus 定向状态
type TargetingStatus string

const (
	TargetingStatusValid  TargetingStatus = "Valid"
	TargetingStatusFrozen TargetingStatus = "Frozen"
)

// TargetingInfos 定向信息列表
type TargetingInfos []*TargetingInfo

// TargetingInfo 一个定向条件
type TargetingInfo struct {
	Name   string            `json:"name,omitempty"`
	Not    bool              `json:"not,omitempty"`
	Values *TargetingValue `json:"values,omitempty"`
}

// TargetingValueType 定向值类型
type TargetingValueType string

// 定向值类型枚举
const (
	TargetingValueTypeString TargetingValueType = "string"
	TargetingValueTypeRange  TargetingValueType = "range"
)

// TargetingValue 定向值
type TargetingValue struct {
	Type   TargetingValueType `json:"type,omitempty"`
	String []string           `json:"string,omitempty"`
	Range  []*Range           `json:"range,omitempty"`
}

// Range 区间 [Begin, End), 前闭后开
type Range struct {
	Begin uint64 `json:"start,omitempty"`
	End   uint64 `json:"end,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *TargetingInfos) Scan(value interface{}) error {
	return scan(value, t)
}

// Value return json value, implement driver.Valuer interface
func (t TargetingInfos) Value() (driver.Value, error) {
	return value(t)
}

// FreqControl 频次控制
type FreqControl struct {
	Rules []*FreqControlRule `json:"rules"`
}

type FreqControlForType string

const (
	FreqControlForClick FreqControlForType = "click"
)

type FreqControlRule struct {
	Platform        ADPlatformType     `json:"platform"`         // 如果没有填，默认为全渠道一起控制
	For             FreqControlForType `json:"for"`              // 为什么行为，控制频次
	DayLimit        int                `json:"limit"`            // 天级限制次数
	AccumulateLimit int                `json:"accumulate_limit"` // 总共限制次数
	Interval        int                `json:"interval"`         // 间隔时间, 多久内只能出现一次，单位秒
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (t *FreqControl) Scan(value interface{}) error {
	return scan(value, t)
}

// Value return json value, implement driver.Valuer interface
func (t FreqControl) Value() (driver.Value, error) {
	return value(t)
}
