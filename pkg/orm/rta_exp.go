package orm

import (
	"fmt"
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllRtaExp 从数据库中获取所有的Rta实验信息
func GetAllRtaExp(db *gorm.DB) ([]*types.RtaExp, error) {
	var rtaExp []*types.RtaExp
	if err := db.Find(&rtaExp).Error; err != nil {
		return nil, err
	}
	return rtaExp, nil
}

// GetRtaExpById 根据id从数据库中获得对应的Rta实验信息
func GetRtaExpById(db *gorm.DB, id string) (*types.RtaExp, error) {
	rtaExp := &types.RtaExp{}
	if err := db.Where("ID = ?", id).Find(rtaExp).Error; err != nil {
		return nil, err
	}
	return rtaExp, nil
}

// InsertRtaExp 插入新的Rta实验
func InsertRtaExp(db *gorm.DB, rtaExp *types.RtaExp) error {
	return db.Create(rtaExp).Error
}

// UpdateRtaExp 更新RtaExp实验的数据
func UpdateRtaExp(db *gorm.DB, rtaExp *types.RtaExp) error {
	return db.Updates(rtaExp).Error
}

// DeleteRtaExpById 根据id从数据库中软删除Rta实验
func DeleteRtaExpById(db *gorm.DB, id string) error {
	rtaExp, err := GetRtaExpById(db, id)
	if err != nil {
		return err
	}
	if !rtaExp.CanDelete() {
		return fmt.Errorf("rta exp can not delete")
	}
	return db.Delete(rtaExp).Error
}

func LoadRtaExpBindExperimentItem(db *gorm.DB, exp *types.RtaExp) error {
	var err error
	exp.BindExperimentItem, err = GetExperimentItemById(db, exp.BindExperimentItemID)
	return err
}

// GetRtaServerExperimentConfig 获取Rta实验的配置
func GetRtaServerExperimentConfig(db *gorm.DB) ([]*types.RtaExp, error) {
	rtaExp := make([]*types.RtaExp, 0)
	var err error
	if err = db.Where("bind_experiment_item_id > ?", 0).
		Where("status", types.RtaExpValid).
		Find(&rtaExp).Error; err != nil {
		return nil, err
	}

	for _, exp := range rtaExp {
		if err := LoadRtaExpBindExperimentItem(db, exp); err != nil {
			return nil, err
		}

		if err = LoadItemExperimentMetadata(db, exp.BindExperimentItem); err != nil {
			return nil, err
		}

		for _, metadata := range exp.BindExperimentItem.ExperimentMetadata {
			if err = LoadExperimentMetadataExperimentParameter(db, metadata); err != nil {
				return nil, err
			}
		}
	}

	return rtaExp, nil
}
