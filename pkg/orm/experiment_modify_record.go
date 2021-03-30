package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetExperimentModifyRecordById 根据ID获取实验修改记录
func GetExperimentModifyRecordById(db *gorm.DB, id uint64) (*types.ExperimentModifyRecord, error) {
	modifyRecord := &types.ExperimentModifyRecord{}
	if err := db.Take(modifyRecord, id).Error; err != nil {
		return nil, err
	}
	return modifyRecord, nil
}

// UpsertExperimentModifyRecord 创建或者插入实验修改记录
func UpsertExperimentModifyRecord(db *gorm.DB, modifyRecord *types.ExperimentModifyRecord) error {
	if modifyRecord.ID == 0 {
		return db.Create(modifyRecord).Error
	}
	return db.Updates(modifyRecord).Error
}

// DeleteExperimentModifyRecordById 删除实验修改记录
func DeleteExperimentModifyRecordById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.ExperimentModifyRecord{}, id).Error
}

// LoadExperimentModifyRecordUser 加载修改记录中的操作用户信息
func LoadExperimentModifyRecordUser(db *gorm.DB, record *types.ExperimentModifyRecord) error {
	user := &types.User{}
	if err := db.Model(record).Association("User").Find(user); err != nil {
		return err
	}
	record.User = user
	return nil
}

// insertExperimentModifyRecordAfterPromote 实验修改后插入修改记录
func InsertExperimentModifyRecordAfterPromote(db *gorm.DB, groupId uint64, userId uint64) error {
	group, err := GetExperimentGroupById(db, groupId)
	if err != nil {
		return err
	}
	if err = fillExperimentGroupCurrent(db, group); err != nil {
		return err
	}
	if err = LoadGroupUser(db, group); err != nil {
		return err
	}

	group.Draft = nil
	group.ExperimentStage = nil

	modifyRecord := &types.ExperimentModifyRecord{
		UserID:            userId,
		ExperimentGroupID: group.ID,
		ExperimentStageID: group.Current.ID,
		Operation:         types.ExperimentModifyOperationModify,
		Data:              (*types.ExperimentGroupWrapper)(group),
	}

	return UpsertExperimentModifyRecord(db, modifyRecord)
}

// insertExperimentModifyRecordAfterStop 停止实验后插入修改记录
func InsertExperimentModifyRecordAfterStop(db *gorm.DB, groupId uint64, stageId uint64, userId uint64) error {
	modifyRecord := &types.ExperimentModifyRecord{
		UserID:            userId,
		ExperimentGroupID: groupId,
		ExperimentStageID: stageId,
		Operation:         types.ExperimentModifyOperationStop,
	}

	return UpsertExperimentModifyRecord(db, modifyRecord)
}
