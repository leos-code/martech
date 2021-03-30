package types

import (
	"time"

	"gorm.io/gorm"
)

// ExperimentParameter 实验参数
type ExperimentParameter struct {
	ID                 uint64                `gorm:"column:id;primaryKey"             json:"id,omitempty"`
	CreatedAt          time.Time             `gorm:"column:created_at"                json:"created_at,omitempty"`
	UpdatedAt          time.Time             `gorm:"column:updated_at"                json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt        `gorm:"column:delete_at;index"           json:"-"`
	Name               string                `gorm:"column:name;unique"               json:"name,omitempty"`
	Type               string                `gorm:"column:type"                      json:"type,omitempty"`
	Description        string                `gorm:"column:description"               json:"description,omitempty"`
	ExperimentMetadata []*ExperimentMetadata `gorm:"foreignKey:ExperimentParameterID" json:"experiment_metadata,omitempty"`
}
