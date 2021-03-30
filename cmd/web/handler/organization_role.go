package handler

import (
	"fmt"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizationRoleGetHandler 拉取角色视角对应的用户数据的列表
func OrganizationRoleGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)

		e := GetEnforcer(c)

		// get all user's id and role's tree in current tenant
		var uid []uint64
		tree := map[uint64]uint64{}
		rules := e.GetFilteredGroupingPolicy(2, domain)
		for _, rule := range rules {
			if casbinUserIs(rule[0]) {
				x := casbinUserDecode(rule[0]).ID
				uid = append(uid, x)
			} else {
				x := casbinRoleDecode(rule[0]).ID
				y := casbinRoleDecode(rule[1]).ID
				tree[y] = x
			}
		}

		// get all user and role data
		linq.From(uid).Distinct().ToSlice(&uid)
		users, err := orm.ListUserByIdWithLoginUser(db, uid)
		if err != nil {
			return nil, err
		}
		um := map[uint64]*types.User{}
		for _, v := range users {
			um[v.ID] = v
		}

		roles, err := orm.ListRoleByTenantId(db, currentT.ID)
		if err != nil {
			return nil, err
		}

		// filter roles with permission
		linq.From(roles).Where(func(v interface{}) bool {
			ok, _ := e.Enforce(cu, domain, v.(*types.Role).Object, read)
			return ok
		}).ToSlice(&roles)

		// get users per role in current tenant
		for _, v := range roles {
			us := e.GetUsersForRoleInDomain(casbinRoleEncode(v), domain)
			for _, w := range us {
				if casbinUserIs(w) {
					v.User = append(v.User, um[casbinUserDecode(w).ID])
				}
			}

			if p, ok := tree[v.ID]; ok {
				v.ParentID = p
			}
		}

		return roles, nil
	})
}

// OrganizationRolePostHandler 更新一个角色的数据
func OrganizationRolePostHandler(c *gin.Context) {
	item := &types.Role{}
	dbGetHandler(c, item, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// check permission
		fn := func(id uint64) (pData, error) {
			v := &types.Role{ID: id}
			err := orm.TakeRole(db, v)
			return v, err
		}
		checker := &writePermissionChecker{
			e:       e,
			subject: cu,
			domain:  domain,
			item:    item,
			take:    fn,
		}
		if err := checker.toCheck(item.ID, item.ParentID); err != nil {
			return nil, err
		}

		// update role in db
		item.TenantID = currentT.ID
		if err := orm.UpsertRole(db, item); err != nil {
			return nil, err
		}

		// delete old parent
		role := casbinRoleEncode(item)
		rs := e.GetUsersForRoleInDomain(role, domain)
		for _, v := range rs {
			if !casbinRoleIs(v) {
				continue
			}
			if _, err := e.DeleteRoleForUserInDomain(v, role, domain); err != nil {
				return nil, err
			}
		}

		// add new parent
		if item.ParentID != 0 {
			parent := &types.Role{ID: item.ParentID}
			if _, err := e.AddRoleForUserInDomain(casbinRoleEncode(parent), role, domain); err != nil {
				return nil, err
			}
		}

		item.CreatedAt = time.Time{}
		item.UpdatedAt = time.Time{}
		return item, orm.TakeRole(db, item)
	})
}

// OrganizationRoleUserHandler 更新一个角色的关联用户数据
func OrganizationRoleUserHandler(c *gin.Context) {
	item := &types.Role{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		// verify input
		if item.ID == 0 {
			return ErrEmptyID
		}
		if err := orm.TakeRole(db, item); err != nil {
			return err
		}

		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		role := casbinRoleEncode(item)
		e := GetEnforcer(c)

		if ok, _ := e.Enforce(cu, domain, item.Object, write); !ok {
			return ErrNoWritePermission
		}

		// get all user's id include request's user in current tenant
		var uid, uid1, uid2 []uint64
		rs := e.GetUsersForRoleInDomain(role, domain)
		for _, v := range rs {
			if casbinUserIs(v) {
				uid1 = append(uid1, casbinUserDecode(v).ID)
			}
		}
		for _, v := range item.User {
			if v.ID != 0 {
				uid2 = append(uid2, v.ID)
			}
		}

		// get all user  data
		uid = append(uid, uid1...)
		uid = append(uid, uid2...)
		users, err := orm.ListUserByIdWithLoginUser(db, uid)
		if err != nil {
			return err
		}
		um := map[uint64]*types.User{}
		for _, v := range users {
			um[v.ID] = v
		}

		// make source and target user list
		var source, target []interface{}
		for _, v := range uid1 {
			if u, ok := um[v]; ok {
				source = append(source, casbinUserEncode(u))
			}
		}
		for _, v := range uid2 {
			if u, ok := um[v]; ok {
				target = append(target, casbinUserEncode(u))
			}
		}

		// get diff to add and remove
		add, remove := Diff(source, target)
		for _, v := range add {
			if _, err := e.AddRoleForUserInDomain(v.(string), role, domain); err != nil {
				return err
			}
		}

		for _, v := range remove {
			if _, err := e.DeleteRoleForUserInDomain(v.(string), role, domain); err != nil {
				return err
			}
		}

		return nil
	})
}

// OrganizationRoleDeleteHandler 删除一个角色对应的关联数据
func OrganizationRoleDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		// verify input
		item := &types.Role{ID: id}
		if err := orm.TakeRole(db, item); err != nil {
			return err
		}

		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		role := casbinRoleEncode(item)
		e := GetEnforcer(c)

		// filter roles with permission
		if ok, _ := e.Enforce(cu, domain, item.Object, write); !ok {
			return fmt.Errorf("current user has no wirte permission for role data in domain")
		}

		// delete all user's id for the role in current tenant
		us := e.GetUsersForRoleInDomain(role, domain)
		for _, v := range us {
			if _, err := e.DeleteRoleForUserInDomain(v, role, domain); err != nil {
				return err
			}
		}

		// delete all role's id for the role in current tenant
		rs := e.GetRolesForUserInDomain(role, domain)
		for _, v := range rs {
			if _, err := e.DeleteRoleForUserInDomain(role, v, domain); err != nil {
				return err
			}
		}

		// delete role policy
		if _, err := e.RemoveFilteredPolicy(0, role, domain); err != nil {
			return err
		}

		return orm.DeleteRoleById(db, id)
	})
}
