package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// MaterialPage 获取素材
func MaterialPage(db *gorm.DB, option *StatementOption) (*Page, error) {
	pageCount, err := getPageCount(db, types.MaterialTemp, option)
	if err != nil {
		return nil, err
	}

	var materials []*types.Material
	if err := getDBWithOption(db, option).Preload("Audit").Find(&materials).Error; err != nil {
		return nil, err
	}

	return &Page{
		Total: pageCount,
		List:  materials,
	}, nil
}

// MaterialDelete 删除素材
func MaterialDelete(db *gorm.DB, ids ...uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, id := range ids {
			if err := materialDelete(db, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func materialDelete(db *gorm.DB, id uint64) error {
	return db.Delete(types.MaterialTemp, id).Error
}
