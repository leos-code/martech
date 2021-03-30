package orm

import (
	"fmt"
	"math"
	"reflect"

	"gorm.io/gorm"
)

// Upsert 插入或者更新记录，记录需要有ID字段
func Upsert(db *gorm.DB, values ...interface{}) (err error) {
	for _, value := range values {
		if err = upsert(db, value); err != nil {
			return err
		}
	}

	return nil
}

func upsert(db *gorm.DB, value interface{}) (err error) {
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("panic: %v", reason)
		}
	}()
	s := reflect.ValueOf(value).Elem()
	id := s.FieldByName("ID").Uint()

	if id == 0 {
		err = db.Create(value).Error
	} else {
		err = db.Updates(value).Error
	}
	return
}

func getPageCount(db *gorm.DB, model interface{}, option *StatementOption) (int64, error) {
	paginationOpt := getPaginationOption(option)
	if paginationOpt == nil {
		return 0, NoPageError
	}

	var count int64
	if err := getDBWithFilterOption(db, option).Model(model).Count(&count).Error; err != nil {
		return 0, err
	}

	return int64(math.Ceil(float64(count) / float64(paginationOpt.PageSize))), nil
}
