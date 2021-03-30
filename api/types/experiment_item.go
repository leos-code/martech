package types

import (
	"sort"
	"time"

	"gorm.io/gorm"
)

// ExperimentItem 实验组下的实验
type ExperimentItem struct {
	ID                 uint64                `gorm:"column:id;primaryKey"               json:"id,omitempty"`
	CreatedAt          time.Time             `gorm:"column:created_at"                  json:"created_at,omitempty"`
	UpdatedAt          time.Time             `gorm:"column:updated_at"                  json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt        `gorm:"column:delete_at;index"             json:"-"`
	OuterID            uint64                `gorm:"column:outer_id"                    json:"outer_id,omitempty"`
	Name               string                `gorm:"column:name"                        json:"name,omitempty"`
	ExperimentStageID  uint64                `gorm:"column:experiment_stage_id"         json:"experiment_stage_id,omitempty"`
	RtaExp             []*RtaExp             `gorm:"many2many:experiment_item_rta_exps" json:"rta_exp,omitempty"`
	ExperimentMetadata []*ExperimentMetadata `gorm:"foreignKey:ExperimentItemID"        json:"experiment_metadata,omitempty"`
}

// SortExperimentItemSlice 对实验列表进行排序
func SortExperimentItemSlice(s []*ExperimentItem) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].ID < s[j].ID
	})
}

// StageEqual 从Stage角度判断是否相等
func (item *ExperimentItem) StageEqual(other *ExperimentItem) bool {
	if item.OuterID != other.OuterID {
		return false
	}

	exp := item.RtaExp
	otherExp := other.RtaExp

	if len(exp) != len(otherExp) {
		return false
	}

	SortRtaExpSlice(exp)
	SortRtaExpSlice(otherExp)

	for i := 0; i < len(exp); i++ {
		if exp[i].ID != otherExp[i].ID {
			return false
		}
	}

	return true
}
