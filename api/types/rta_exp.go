package types

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RtaExp RTA实验信息，同步媒体实验平台获取
type RtaExp struct {
	ID                   string            `gorm:"column:id;primaryKey"               json:"id,omitempty"`
	CreatedAt            time.Time         `gorm:"column:created_at"                  json:"created_at,omitempty"`
	UpdatedAt            time.Time         `gorm:"column:updated_at"                  json:"updated_at,omitempty"`
	DeletedAt            gorm.DeletedAt    `gorm:"column:delete_at;index"             json:"-"`
	RtaAccountID         uint64            `gorm:"column:rta_account_id"              json:"rta_account_id,omitempty"`
	FlowRate             float64           `gorm:"column:flow_rate"                   json:"flow_rate,omitempty"` // 单位%
	ExpirationTime       time.Time         `gorm:"column:expiration_time"             json:"expiration_time,omitempty"`
	JsonProperty         datatypes.JSON    `gorm:"column:json_property"               json:"json_property,omitempty"`
	ExperimentItem       []*ExperimentItem `gorm:"many2many:experiment_item_rta_exps" json:"experiment_item,omitempty"`
	Status               RtaExpStatus      `gorm:"column:status;default:Valid"        json:"status,omitempty"`
	BindExperimentItemID uint64            `gorm:"column:bind_experiment_item_id"     json:"bind_experiment_item_id,omitempty"`
	BindExperimentItem   *ExperimentItem   `gorm:"-"                                  json:"bind_experiment_item,omitempty"`
	BindStatus           RtaExpBindStatus  `gorm:"column:bind_status;default:Idle"    json:"bind_status,omitempty"`
}

// RtaExpBindStatus 判断媒体侧RTA实验是否绑定
type RtaExpBindStatus string

const (
	RtaExpIdle RtaExpBindStatus = "Idle"
	RtaExpBusy RtaExpBindStatus = "Busy"
)

// RtaExpBindStatus 判断媒体侧RTA实验的状态
type RtaExpStatus string

const (
	RtaExpValid  RtaExpStatus = "Valid"
	RtaExpPause  RtaExpStatus = "Pause"
	RtaExpExpire RtaExpStatus = "Expire"
)

// CanDelete 判断媒体侧RTA实验是否能被删除
func (exp *RtaExp) CanDelete() bool {
	return exp.BindStatus == RtaExpIdle
}

// IDLess 判断媒体侧RTA的ID小于关系
func (exp *RtaExp) IDLess(other *RtaExp) bool {
	return exp.ID < other.ID
}

// SortRtaExpSlice 对媒体侧RTA实验列表进行排序
func SortRtaExpSlice(s []*RtaExp) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].IDLess(s[j])
	})
}

// GetUpdateTime 获取更新时间
func (exp *RtaExp) GetUpdateTime() time.Time {
	return exp.UpdatedAt
}

// GetID 获取ID
func (exp *RtaExp) GetID() interface{} {
	return exp.ID
}

// Update 更新状态
func (exp *RtaExp) Update(db *gorm.DB) error {
	if exp.ExpirationTime.Before(time.Now()) {
		if exp.Status == RtaExpValid {
			if err := db.Model(exp).Update("Status", RtaExpExpire).Error; err != nil {
				return fmt.Errorf("failed to update rta exp status, err: %v", err)
			} else {
				exp.Status = RtaExpExpire
			}
		}
	}
	return nil
}
