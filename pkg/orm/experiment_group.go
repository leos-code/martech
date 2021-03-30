package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllExperimentGroup 从数据库中获取所有的实验组数据
func GetAllExperimentGroup(db *gorm.DB, option *StatementOption) ([]*types.ExperimentGroup, error) {
	var group []*types.ExperimentGroup

	if err := getDBWithOption(db, option).Find(&group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// GetExperimentGroupById 根据id从数据库中获取对应的实验组信息
func GetExperimentGroupById(db *gorm.DB, id uint64) (*types.ExperimentGroup, error) {
	group := &types.ExperimentGroup{}
	if err := db.Take(group, id).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// LoadGroupExperimentStage 加载实验组下包含的ExperimentStage数据
func LoadGroupExperimentStage(db *gorm.DB, group *types.ExperimentGroup) error {
	return db.Model(group).Association("ExperimentStage").Find(&group.ExperimentStage)
}

// LoadGroupUser 加载管理该实验组的相关用户
func LoadGroupUser(db *gorm.DB, group *types.ExperimentGroup) error {
	return db.Model(group).Association("User").Find(&group.User)
}

// LoadGroupModifyRecord 加载实验组历史修改记录
func LoadGroupModifyRecord(db *gorm.DB, group *types.ExperimentGroup) error {
	return db.Model(group).Association("ModifyRecord").Find(&group.ModifyRecord)
}

// LoadGroupRtaAccount 加载实验组对应的账号
func LoadGroupRtaAccount(db *gorm.DB, group *types.ExperimentGroup) error {
	var rtaAccount types.RtaAccount
	if err := db.Model(group).Association("RtaAccount").Find(&rtaAccount); err != nil {
		return err
	}
	group.RtaAccount = &rtaAccount
	return nil
}

// UpsertExperimentGroup 插入或者更新对应实验组的数据
func UpsertExperimentGroup(db *gorm.DB, group *types.ExperimentGroup) error {
	if group.ID == 0 {
		return db.Create(group).Error
	}
	return db.Updates(group).Error
}

// DeleteExperimentGroupById 根据id从数据库中软删除对应的实验组
func DeleteExperimentGroupById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.ExperimentGroup{}, id).Error
}
