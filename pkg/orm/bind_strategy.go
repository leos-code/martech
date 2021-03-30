package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllBindStrategy 从数据库中查询得到所有的BindStrategy
func GetAllBindStrategy(db *gorm.DB, option *StatementOption) ([]*types.BindStrategy, error) {
	var bindStrategy []*types.BindStrategy

	if err := getDBWithOption(db, option).Find(&bindStrategy).Error; err != nil {
		return nil, err
	}

	return bindStrategy, nil
}

// GetBindStrategyByID 根据id从数据库中查询得到对应的BindStrategy
func GetBindStrategyByID(db *gorm.DB, id uint64) (*types.BindStrategy, error) {
	bindStrategy := &types.BindStrategy{}
	if err := db.Take(bindStrategy, id).Error; err != nil {
		return nil, err
	}

	return bindStrategy, nil
}

// UpsertBindStrategy 更新或者插入一个BindStrategy的信息
func UpsertBindStrategy(db *gorm.DB, bindStrategy *types.BindStrategy) error {
	if bindStrategy.ID == 0 {
		return db.Create(bindStrategy).Error
	}

	return db.Updates(bindStrategy).Error
}

// DeleteBindStrategyById 根据id从数据库中软删除一个BindStrategy
func DeleteBindStrategyById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.BindStrategy{}, id).Error
}
