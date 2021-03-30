package types

import (
	"time"

	"gorm.io/gorm"
)

// ExperimentGroup 实验组
type ExperimentGroup struct {
	ID              uint64                    `gorm:"column:id;primaryKey"             json:"id,omitempty"`
	CreatedAt       time.Time                 `gorm:"column:created_at"                json:"created_at,omitempty"`
	UpdatedAt       time.Time                 `gorm:"column:updated_at"                json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt            `gorm:"column:delete_at;index"           json:"-"`
	Name            string                    `gorm:"column:name"                      json:"name,omitempty"`
	Description     string                    `gorm:"column:description"               json:"description,omitempty"`
	RtaAccountID    uint64                    `gorm:"column:rta_account_id"            json:"rta_account_id,omitempty"`
	RtaAccount      *RtaAccount               `gorm:"foreignKey:RtaAccountID"          json:"rta_account,omitempty"`
	User            []*User                   `gorm:"many2many:user_experiment_groups" json:"user,omitempty"`
	ExperimentStage []*ExperimentStage        `gorm:"foreignKey:ExperimentGroupID"     json:"experiment_stage,omitempty"`
	Draft           *ExperimentStage          `gorm:"-"                                json:"draft,omitempty"`
	Current         *ExperimentStage          `gorm:"-"                                json:"current,omitempty"`
	ModifyRecord    []*ExperimentModifyRecord `gorm:"foreignKey:ExperimentGroupID"     json:"modify_record,omitempty"`
}
