package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizationUserGetHandler 拉取用户视角对应的角色数据的列表
func OrganizationUserGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// get all user's id and role's id in current tenant
		var uid []uint64
		rules := e.GetFilteredGroupingPolicy(2, domain)
		for _, rule := range rules {
			if casbinUserIs(rule[0]) {
				uid = append(uid, casbinUserDecode(rule[0]).ID)
			}
		}

		// get all user and role data
		linq.From(uid).Distinct().ToSlice(&uid)
		users, err := orm.ListUserByIdWithLoginUser(db, uid)
		if err != nil {
			return nil, err
		}
		roles, err := orm.ListRoleByTenantId(db, currentT.ID)
		if err != nil {
			return nil, err
		}

		// filter roles with permission
		rm := map[uint64]*types.Role{}
		for _, v := range roles {
			if ok, _ := e.Enforce(cu, domain, v.Object, read); ok {
				rm[v.ID] = v
			}
		}

		// get roles per user in current tenant
		for _, v := range users {
			rs := e.GetRolesForUserInDomain(casbinUserEncode(v), domain)
			for _, r := range rs {
				if role, ok := rm[casbinRoleDecode(r).ID]; ok {
					v.Role = append(v.Role, role)
				}
			}
		}

		return users, nil
	})
}

// OrganizationUserPostHandler 更新一个用户对应的角色数据
func OrganizationUserPostHandler(c *gin.Context) {
	item := &types.User{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		// verify input
		if err := orm.TakeUser(db, item); err != nil {
			return err
		}

		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		user := casbinUserEncode(item)
		e := GetEnforcer(c)

		// get all role's id include request's role in current tenant
		var rid, rid1, rid2 []uint64
		rs := e.GetRolesForUserInDomain(user, domain)
		for _, v := range rs {
			if casbinRoleIs(v) {
				rid1 = append(rid1, casbinRoleDecode(v).ID)
			}
		}
		for _, v := range item.Role {
			if v.ID != 0 {
				rid2 = append(rid2, v.ID)
			}
		}

		// get all role data
		rid = append(rid, rid1...)
		rid = append(rid, rid2...)
		roles, err := orm.ListRoleById(db, rid)
		if err != nil {
			return err
		}

		// filter roles with permission
		rm := map[uint64]*types.Role{}
		for _, v := range roles {
			if ok, _ := e.Enforce(cu, domain, v.Object, write); ok {
				rm[v.ID] = v
			}
		}

		// make source and target role list
		var source, target []interface{}
		for _, v := range rid1 {
			if r, ok := rm[v]; ok {
				source = append(source, casbinRoleEncode(r))
			}
		}
		for _, v := range rid2 {
			if r, ok := rm[v]; ok {
				target = append(target, casbinRoleEncode(r))
			}
		}

		// get diff to add and remove
		add, remove := Diff(source, target)
		for _, v := range add {
			if _, err := e.AddRoleForUserInDomain(user, v.(string), domain); err != nil {
				return err
			}
		}

		for _, v := range remove {
			if _, err := e.DeleteRoleForUserInDomain(user, v.(string), domain); err != nil {
				return err
			}
		}

		return nil
	})
}
