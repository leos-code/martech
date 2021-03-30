package handler

import (
	"fmt"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SuperAdminGetHandler 拉取超级管理员列表
func SuperAdminGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		e := GetEnforcer(c)
		subjects := e.GetUsersForRoleInDomain(superAdminRole, superAdminDomain)

		var id []uint64
		for _, v := range subjects {
			id = append(id, casbinUserDecode(v).ID)
		}

		return orm.ListUserByIdWithLoginUser(db, id)
	})
}

// SuperAdminPostHandler 新增超级管理员信息
func SuperAdminPostHandler(c *gin.Context) {
	item := &types.User{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		e := GetEnforcer(c)
		return superAdminOneHandler(db, item.ID, e.AddGroupingPolicy)
	})
}

// SuperAdminDeleteHandler 删除某个超级管理员
func SuperAdminDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		e := GetEnforcer(c)
		return superAdminOneHandler(db, id, e.RemoveGroupingPolicy)
	})
}

func superAdminOneHandler(db *gorm.DB, id uint64, f func(...interface{}) (bool, error)) error {
	if id == 0 {
		return fmt.Errorf("invaild user id")
	}

	user := &types.User{ID: id}
	if err := orm.TakeUser(db, user); err != nil {
		return err
	}

	_, err := f(casbinUserEncode(user), superAdminRole, superAdminDomain)
	return err
}
