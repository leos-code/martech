package handler

import (
	"fmt"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

// OrganizationPolicyGetHandler 拉取用户视角对应的角色数据的列表
func OrganizationPolicyGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// get role's tree in current tenant
		rtree := map[uint64]uint64{}
		rules1 := e.GetFilteredGroupingPolicy(2, domain)
		for _, rule := range rules1 {
			if !casbinUserIs(rule[0]) {
				continue
			}
			x := casbinRoleDecode(rule[0]).ID
			y := casbinRoleDecode(rule[1]).ID
			rtree[y] = x
		}

		// // get object's tree in current tenant
		// otree := map[uint64]uint64{}
		// rules2 := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
		// for _, rule := range rules2 {
		// 	if !casbinObjectIs(rule[1]) {
		// 		continue
		// 	}
		// 	// x := casbinObjectDecode(rule[0]).ID
		// 	// y := casbinObjectDecode(rule[1]).ID
		// 	// otree[x] = y
		// }

		// get all role and object data
		roles, err := orm.ListRoleByTenantId(db, currentT.ID)
		if err != nil {
			return nil, err
		}

		objects, err := orm.ListObjectByTenantId(db, currentT.ID)
		if err != nil {
			return nil, err
		}
		// get feature object
		in0 := &types.Object{Type: types.ObjectTypeFeature, TenantID: SuperAdminDomainID}
		objects0, err := orm.FindObject(db, in0)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objects0...)

		// filter roles and objects with permission
		linq.From(roles).Where(func(v interface{}) bool {
			ok, _ := e.Enforce(cu, domain, v.(*types.Role).Object, read)
			return ok
		}).ToSlice(&roles)

		linq.From(objects).Where(func(v interface{}) bool {
			ok, _ := e.Enforce(cu, domain, v.(*types.Object).Object, read)
			return ok
		}).ToSlice(&objects)

		rm := map[uint64]*types.Role{}
		for _, v := range roles {
			rm[v.ID] = v
		}

		om := map[uint64]*types.Object{}
		for _, v := range objects {
			om[v.ID] = v
		}

		// get all policies
		policies := e.GetFilteredPolicy(1, domain)
		for _, policy := range policies {
			r := casbinRoleDecode(policy[0])
			o := casbinObjectDecode(policy[2])
			rr, ok := rm[r.ID]
			action := policy[3]
			if !ok {
				continue
			}
			oo, ok := om[o.ID]
			if !ok {
				continue
			}

			p := &types.Policy{Object: oo}
			switch action {
			case read:
				p.Read = true
			case write:
				p.Write = true
			}

			rr.Policy = append(rr.Policy, p)
		}

		// merge policies
		for _, role := range roles {
			pm := map[uint64]*types.Policy{}
			for _, p := range role.Policy {
				if _, ok := pm[p.Object.ID]; !ok {
					pm[p.Object.ID] = p
				}
				v := pm[p.Object.ID]
				v.Read = v.Read || p.Read
				v.Write = v.Write || p.Write
			}

			role.Policy = nil
			for _, v := range pm {
				role.Policy = append(role.Policy, v)
			}
		}

		return roles, nil
	})
}

// OrganizationUserPostHandler 更新一个用户对应的角色数据
func OrganizationPolicyPostHandler(c *gin.Context) {
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

		// get all object's id include object's role in current tenant
		var oid, oid1, oid2 []uint64
		policy := e.GetPermissionsForUser(role, domain)
		for _, v := range policy {
			if casbinObjectIs(v[2]) {
				oid1 = append(oid1, casbinObjectDecode(v[2]).ID)
			}
		}
		for _, v := range item.Policy {
			if v.Object.ID != 0 {
				oid2 = append(oid2, v.Object.ID)
			}
		}

		// get all object data
		oid = append(oid, oid1...)
		oid = append(oid, oid2...)
		objects, err := orm.ListObjectById(db, oid)
		if err != nil {
			return err
		}

		// filter objects with permission
		om := map[uint64]*types.Object{}
		for _, v := range objects {
			if v.Type == types.ObjectTypeFeature {
				if ok, _ := e.Enforce(cu, domain, v.Object, read); ok {
					om[v.ID] = v
				}
			}
			if ok, _ := e.Enforce(cu, domain, v.Object, write); ok {
				om[v.ID] = v
			}
		}

		// make source and target role list
		var source, target []interface{}
		for _, v := range policy {
			if !casbinObjectIs(v[2]) {
				continue
			}
			if _, ok := om[casbinObjectDecode(v[2]).ID]; ok {
				source = append(source, fmt.Sprintf("%v,%v", v[2], v[3]))
			}
		}

		for _, v := range item.Policy {
			if _, ok := om[v.Object.ID]; !ok {
				continue
			}
			if v.Read {
				target = append(target, fmt.Sprintf("%v,%v", casbinObjectEncode(v.Object), read))
			}
			if v.Write && v.Object.Type != types.ObjectTypeFeature {
				target = append(target, fmt.Sprintf("%v,%v", casbinObjectEncode(v.Object), write))
			}
		}

		// get diff to add and remove
		add, remove := Diff(source, target)
		for _, v := range add {
			s := strings.Split(v.(string), ",")
			if _, err := e.AddPermissionForUser(role, domain, s[0], s[1]); err != nil {
				return err
			}
		}

		for _, v := range remove {
			s := strings.Split(v.(string), ",")
			if _, err := e.DeletePermissionForUser(role, domain, s[0], s[1]); err != nil {
				return err
			}
		}

		return nil
	})
}
