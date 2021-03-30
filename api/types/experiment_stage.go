package types

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ExperimentStage 实验阶段，进行实验流量调整会产生新的阶段
type ExperimentStage struct {
	ID                uint64                `gorm:"column:id;primaryKey"         json:"id,omitempty"`
	CreatedAt         time.Time             `gorm:"column:created_at"            json:"created_at,omitempty"`
	UpdatedAt         time.Time             `gorm:"column:updated_at"            json:"updated_at,omitempty"`
	DeletedAt         gorm.DeletedAt        `gorm:"column:delete_at;index"       json:"-"`
	StartTime         time.Time             `gorm:"column:start_time"            json:"start_time,omitempty"`
	EndTime           time.Time             `gorm:"column:end_time"              json:"end_time,omitempty"`
	Status            ExperimentStageStatus `gorm:"column:status"                json:"status,omitempty"`
	ExperimentGroupID uint64                `gorm:"column:experiment_group_id"   json:"experiment_group_id,omitempty"`
	ExperimentItem    []*ExperimentItem     `gorm:"foreignKey:ExperimentStageID" json:"experiment_item,omitempty"`
	Version           uint64                `gorm:"column:version"               json:"version,omitempty"`
}

// ExperimentStageStatus 实验阶段状态
type ExperimentStageStatus string

const (
	ExperimentStageDraft   ExperimentStageStatus = "Draft"
	ExperimentStageRunning ExperimentStageStatus = "Running"
	ExperimentStageStop    ExperimentStageStatus = "Stop"
)

// StageEqual 从stage的角度判断是否相等
func (stage *ExperimentStage) StageEqual(other *ExperimentStage) bool {
	if other == nil {
		return false
	}

	if len(stage.ExperimentItem) != len(other.ExperimentItem) {
		return false
	}

	SortExperimentItemSlice(stage.ExperimentItem)
	SortExperimentItemSlice(other.ExperimentItem)

	length := len(stage.ExperimentItem)

	for i := 0; i < length; i++ {
		if !other.ExperimentItem[i].StageEqual(stage.ExperimentItem[i]) {
			return false
		}
	}

	return true
}

// GetUpdateTime 获取更新时间
func (stage *ExperimentStage) GetUpdateTime() time.Time {
	return stage.UpdatedAt
}

// GetID 获取实验阶段的ID
func (stage *ExperimentStage) GetID() interface{} {
	return stage.ID
}

// Update 更新状态
func (stage *ExperimentStage) Update(db *gorm.DB) error {
	if stage.EndTime.Before(time.Now()) {
		if stage.Status == ExperimentStageRunning {
			if err := db.Model(stage).UpdateColumn("status", ExperimentStageStop).Error; err != nil {
				return fmt.Errorf("fail to update experiment stage, err: %v", err)
			}
			stage.Status = ExperimentStageStop
		}
	}
	return nil
}
