package updater

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"gorm.io/gorm"
	"time"
)

// ExperimentStageTableSnapshotSync 同步ExperimentStage表
type ExperimentStageTableSnapshotSync struct{}

// NewExperimentStageTableSnapshotSync 构造函数
func NewExperimentStageTableSnapshotSync() *ExperimentStageTableSnapshotSync {
	return &ExperimentStageTableSnapshotSync{}
}

// GetAll 获取所有的数据
func (u *ExperimentStageTableSnapshotSync) GetAll(db *gorm.DB) ([]Record, error) {
	list, err := orm.GetAllExperimentStage(db)
	if err != nil {
		return nil, err
	}
	return convertExperimentStageListToRecordList(list), nil
}

// GetInc 增量获取新修改的数据
func (u *ExperimentStageTableSnapshotSync) GetInc(db *gorm.DB, updateFrom time.Time) ([]Record, error) {
	var inc []*types.ExperimentStage

	if err := db.
		Find(&inc, "updated_at >= ? and status == ? ", updateFrom, types.ExperimentStageRunning).Error; err != nil {
		return nil, err
	}
	return convertExperimentStageListToRecordList(inc), nil
}

func convertExperimentStageListToRecordList(list []*types.ExperimentStage) []Record {
	result := make([]Record, len(list))
	for i, v := range list {
		result[i] = v
	}
	return result
}
