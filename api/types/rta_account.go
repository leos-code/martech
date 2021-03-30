package types

import (
	"time"

	"gorm.io/gorm"
)

// RtaAccount RTA账号信息
type RtaAccount struct {
	ID              uint64             `gorm:"column:id;primaryKey"    json:"id,omitempty"`
	CreatedAt       time.Time          `gorm:"column:created_at"       json:"created_at,omitempty"`
	UpdatedAt       time.Time          `gorm:"column:updated_at"       json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt     `gorm:"column:delete_at;index"  json:"-"`
	RtaID           string             `gorm:"column:rta_id;unique"    json:"rta_id,omitempty"`
	Token           string             `gorm:"column:token"            json:"token,omitempty"`
	Description     string             `gorm:"column:description"      json:"description,omitempty"`
	RtaExp          []*RtaExp          `gorm:"foreignKey:RtaAccountID" json:"rta_exp,omitempty"`
	ExperimentGroup []*ExperimentGroup `gorm:"foreignKey:RtaAccountID" json:"experiment_group,omitempty"`
}
