package orm

import (
	"github.com/tencentad/martech/api/types"
	"gorm.io/gorm"
)

// GetAllLoginUser 获取所有用户的登录
func GetAllLoginUser(db *gorm.DB) ([]*types.LoginUser, error) {
	var user []*types.LoginUser
	if err := db.Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// TakeLoginUser 根据用户登录获取用户的登录信息
func TakeLoginUser(db *gorm.DB, user *types.LoginUser) error {
	return db.Where(user).Take(user).Error
}

// UpsertLoginUser 创建或者更新用户的登录记录
func UpsertLoginUser(db *gorm.DB, user *types.LoginUser) error {
	if err := db.Create(user).Error; err != nil {
		return db.Updates(user).Error
	}
	return nil
}
