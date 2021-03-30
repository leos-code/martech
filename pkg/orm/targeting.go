package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllTargeting 从数据库中获取所有的Rta定向信息
func GetAllTargeting(db *gorm.DB, option *StatementOption) ([]*types.Targeting, error) {
	var targeting []*types.Targeting
	if err := getDBWithOption(db, option).Find(&targeting).Error; err != nil {
		return nil, err
	}

	return targeting, nil
}

func TargetingPage(db *gorm.DB, option *StatementOption) (*Page, error) {
	pageCount, err := getPageCount(db, types.TargetingTemp, option)
	if err != nil {
		return nil, err
	}

	var targeting []*types.Targeting
	if err := getDBWithOption(db, option).Preload("BindStrategy").Find(&targeting).Error; err != nil {
		return nil, err
	}

	return &Page{
		Total: pageCount,
		List: targeting,
	}, nil
}

// GetTargetingByID 根据id从数据库中获取对应Rta定向的信息
func GetTargetingByID(db *gorm.DB, id uint) (*types.Targeting, error) {
	targeting := &types.Targeting{}
	if err := db.Take(targeting, id).Error; err != nil {
		return nil, err
	}

	return targeting, nil
}

// UpsertTargeting 插入或者更新一个新的Rta定向信息
func UpsertTargeting(db *gorm.DB, targeting *types.Targeting) error {
	if targeting.ID == 0 {
		return db.Create(targeting).Error
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(targeting).Association("BindStrategy").Replace(targeting.BindStrategy); err != nil {
			return err
		}
		targeting.BindStrategy = nil

		return tx.Updates(targeting).Error
	})
}

// DeleteTargetingByID 根据id从数据库中软删除对应的Rta定向
func DeleteTargetingByID(db *gorm.DB, id uint64) error {
	return db.Delete(&types.Targeting{}, id).Error
}

// LoadTargetingBindStrategy 加载定向绑定的策略
func LoadTargetingBindStrategy(db *gorm.DB, targeting *types.Targeting) error {
	return db.Model(targeting).Association("BindStrategy").Find(&targeting.BindStrategy)
}
