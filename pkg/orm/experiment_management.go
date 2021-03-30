package orm

import (
	"errors"
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// 填充实验组自身及其相关子结构的所有数据
func FillExperimentGroup(db *gorm.DB, group *types.ExperimentGroup) error {
	if err := db.Session(&gorm.Session{PrepareStmt: true}).
		Preload("RtaAccount.RtaExp").
		Preload("User").
		Preload("ExperimentStage.ExperimentItem.ExperimentMetadata.ExperimentParameter").
		Preload("ExperimentStage.ExperimentItem.RtaExp").
		Preload("ModifyRecord.User").
		Find(group).
		Error; err != nil {
		return err
	}
	group.Draft = GetDraftFromExperimentGroup(group)
	group.Current = GetCurrentFromExperimentGroup(group)
	return nil
}

func FillExperimentStage(db *gorm.DB, stage *types.ExperimentStage) error {
	if err := db.Preload("ExperimentItem.ExperimentMetadata.ExperimentParameter").
		Preload("ExperimentItem.RtaExp").Find(stage).
		Error; err != nil {
		return err
	}
	return nil
}

func fillExperimentGroupCurrent(db *gorm.DB, group *types.ExperimentGroup) error {
	if err := LoadGroupExperimentStage(db, group); err != nil {
		return err
	}
	group.Current = GetCurrentFromExperimentGroup(group)
	if err := FillExperimentStage(db, group.Current); err != nil {
		return err
	}

	return nil
}

func getStageFromExperimentGroup(group *types.ExperimentGroup,
	status types.ExperimentStageStatus) *types.ExperimentStage {
	for _, v := range group.ExperimentStage {
		if v.Status == status {
			return v
		}
	}
	return nil
}

func GetDraftFromExperimentGroup(group *types.ExperimentGroup) *types.ExperimentStage {
	return getStageFromExperimentGroup(group, types.ExperimentStageDraft)
}

func GetCurrentFromExperimentGroup(group *types.ExperimentGroup) *types.ExperimentStage {
	return getStageFromExperimentGroup(group, types.ExperimentStageRunning)
}

func getStageWithAllDate(db *gorm.DB, group *types.ExperimentGroup,
	status types.ExperimentStageStatus) (*types.ExperimentStage, error) {
	stage := &types.ExperimentStage{}
	if err := db.Where("status = ? AND experiment_group_id = ?", status, group.ID).Take(stage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := FillExperimentStage(db, stage); err != nil {
		return nil, err
	}
	return stage, nil
}

func GetDraftWithAllData(db *gorm.DB, group *types.ExperimentGroup) (*types.ExperimentStage, error) {
	return getStageWithAllDate(db, group, types.ExperimentStageDraft)
}

func GetCurrentWithAllData(db *gorm.DB, group *types.ExperimentGroup) (*types.ExperimentStage, error) {
	return getStageWithAllDate(db, group, types.ExperimentStageRunning)
}
