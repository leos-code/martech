package orm

import (
	"fmt"

	"github.com/tencentad/martech/api/types"
	"github.com/ahmetb/go-linq/v3"
	"gorm.io/gorm"
)

// TakeUser 根据用户获取用户信息
func TakeUser(db *gorm.DB, user *types.User) error {
	return db.Where(user).Take(user).Error
}

// UpsertUser 创建或者更新用户记录
func UpsertUser(db *gorm.DB, user *types.User) error {
	if user.ID == 0 {
		return db.Create(user).Error
	}
	return db.Updates(user).Error
}

// ListUserByIdWithLoginUser 获取一批用户和关联的登录帐号信息
func ListUserByIdWithLoginUser(db *gorm.DB, id []uint64) ([]*types.User, error) {
	var user []*types.User
	if err := db.Preload("LoginUser").Find(&user, "id IN ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// LoadUserLoginUser 获取用户关联的登录帐号信息
func LoadUserLoginUser(db *gorm.DB, user *types.User) error {
	return db.Model(user).Association("LoginUser").Find(&user.LoginUser)
}

// SearchUser 根据info搜索相关用户
func SearchUser(db *gorm.DB, info *types.UserSearch) ([]*types.User, error) {
	var user []*types.User
	if info.IsEmpty() {
		return user, nil
	}

	op, format := "=", "%v"
	if info.Approximate {
		op, format = "LIKE", "%%%v%%"
	}

	d := db.Preload("LoginUser")
	if info.PhoneNumber != "" {
		query := fmt.Sprintf("phone_number %v ?", op)
		arg := fmt.Sprintf(format, info.PhoneNumber)
		d = d.Where(query, arg)
	}
	if info.Email != "" {
		query := fmt.Sprintf("email %v ?", op)
		arg := fmt.Sprintf(format, info.Email)
		d = d.Where(query, arg)
	}
	if info.NickName != "" {
		query := fmt.Sprintf("JOIN login_users ON login_users.user_id = users.id AND login_users.nick_name %v ?", op)
		arg := fmt.Sprintf(format, info.NickName)
		d = d.Joins(query, arg)
	}

	if err := d.Find(&user).Error; err != nil {
		return nil, err
	}

	linq.From(user).DistinctBy(func(v interface{}) interface{} {
		return v.(*types.User).ID
	}).ToSlice(&user)

	return user, nil
}
