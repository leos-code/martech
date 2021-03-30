package types

import (
	"time"

	"gorm.io/gorm"
)

// ExperimentMetadata 实验参数值
type ExperimentMetadata struct {
	ID                    uint64               `gorm:"column:id;primaryKey"             json:"id,omitempty"`
	CreatedAt             time.Time            `gorm:"column:created_at"                json:"created_at,omitempty"`
	UpdatedAt             time.Time            `gorm:"column:updated_at"                json:"updated_at,omitempty"`
	DeletedAt             gorm.DeletedAt       `gorm:"column:delete_at;index"           json:"-"`
	Value                 string               `gorm:"column:value"                     json:"value,omitempty"`
	ExperimentParameterID uint64               `gorm:"column:experiment_parameter_id"   json:"experiment_parameter_id,omitempty"`
	ExperimentParameter   *ExperimentParameter `gorm:"foreignKey:ExperimentParameterID" json:"experiment_parameter"`
	ExperimentItemID      uint64               `gorm:"column:experiment_item_id"        json:"experiment_item_id,omitempty"`
}
