package updater

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"gorm.io/gorm"
	"time"
)

// RtaExpTableSnapshotSync 同步RtaExp表
type RtaExpTableSnapshotSync struct{}

// NewRtaExpTableSnapshotSync 构造函数
func NewRtaExpTableSnapshotSync() *RtaExpTableSnapshotSync {
	return &RtaExpTableSnapshotSync{}
}

// GetAll 获取所有的RtaExp数据
func (r RtaExpTableSnapshotSync) GetAll(db *gorm.DB) ([]Record, error) {
	list, err := orm.GetAllRtaExp(db)
	if err != nil {
		return nil, err
	}
	return convertRtaExpListToRecordList(list), nil

}

// GetInc 增量获取新修改的RtaExp数据
func (r RtaExpTableSnapshotSync) GetInc(db *gorm.DB, updateFrom time.Time) ([]Record, error) {
	var inc []*types.RtaExp
	if err := db.Find(&inc, "updated_at >= ?", updateFrom).Error; err != nil {
		return nil, err
	}
	return convertRtaExpListToRecordList(inc), nil
}

func convertRtaExpListToRecordList(list []*types.RtaExp) []Record {
	result := make([]Record, len(list))
	for i, v := range list {
		result[i] = v
	}
	return result
}
