package types

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

// ExperimentModifyRecord 实验修改记录
type ExperimentModifyRecord struct {
	ID                uint64                    `gorm:"column:id;primaryKey"         json:"id,omitempty"`
	CreatedAt         time.Time                 `gorm:"column:created_at"            json:"created_at,omitempty"`
	UpdatedAt         time.Time                 `gorm:"column:updated_at"            json:"updated_at,omitempty"`
	DeletedAt         gorm.DeletedAt            `gorm:"column:delete_at;index"       json:"-"`
	UserID            uint64                    `gorm:"column:user_id"               json:"user_id"`
	User              *User                     `gorm:"foreignKey:UserID"            json:"user"`
	ExperimentGroupID uint64                    `gorm:"column:experiment_group_id"   json:"experiment_group_id"`
	ExperimentStageID uint64                    `gorm:"column:experiment_stage_id"   json:"experiment_stage_id"`
	BelongStage       *ExperimentStage          `gorm:"foreignKey:ExperimentStageID" json:"belong_stage,omitempty"`
	Operation         ExperimentModifyOperation `gorm:"column:operation"             json:"operation"`
	Data              *ExperimentGroupWrapper   `gorm:"column:data"                  json:"data"`
}

type ExperimentModifyOperation string

const (
	ExperimentModifyOperationModify ExperimentModifyOperation = "Modify"
	ExperimentModifyOperationStop   ExperimentModifyOperation = "Stop"
)

// ExperimentGroupWrapper 方便DB保存实验组json结构
type ExperimentGroupWrapper ExperimentGroup

// Scan sql.Scanner
func (group *ExperimentGroupWrapper) Scan(value interface{}) error {
	return scan(value, group)
}

// Value sql.Valuer
func (group *ExperimentGroupWrapper) Value() (driver.Value, error) {
	return value(group)
}
