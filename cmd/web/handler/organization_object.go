package handler

import (
	"fmt"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrganizationObjectGetHandler 拉取实体数据
func OrganizationObjectGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		ty := c.DefaultQuery("type", "")
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// get object's tree in current tenant
		tree := map[uint64]uint64{}
		rules := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
		for _, rule := range rules {
			if !casbinObjectIs(rule[1]) {
				continue
			}
			x := casbinObjectDecode(rule[0]).ID
			y := casbinObjectDecode(rule[1]).ID
			tree[x] = y
		}

		// get all object data in current tenant
		in := &types.Object{Type: types.ObjectType(ty), TenantID: currentT.ID}
		objects, err := orm.FindObject(db, in)
		if err != nil {
			return nil, err
		}

		// get feature object
		if types.ObjectType(ty) == types.ObjectTypeFeature || ty == "" {
			in0 := &types.Object{Type: types.ObjectTypeFeature, TenantID: SuperAdminDomainID}
			objects0, err := orm.FindObject(db, in0)
			if err != nil {
				return nil, err
			}

			objects = append(objects, objects0...)
		}

		// filter objects read permission
		linq.From(objects).Where(func(v interface{}) bool {
			ok, _ := e.Enforce(cu, domain, v.(*types.Object).Object, read)
			return ok
		}).ToSlice(&objects)

		// build parent id
		for _, v := range objects {
			if p, ok := tree[v.ID]; ok {
				v.ParentID = p
			}
		}

		return objects, nil
	})
}

// OrganizationObjectPostHandler 更新实体数据
func OrganizationObjectPostHandler(c *gin.Context) {
	item := &types.Object{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		// check permission
		fn := func(id uint64) (pData, error) {
			v := &types.Object{ID: id}
			err := orm.TakeObject(db, v)
			if err == nil && v.Type != item.Type {
				return nil, ErrNoWritePermission
			}
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
			return err
		}
		if item.ParentID == 0 {
			return fmt.Errorf("can't add a new root data")
		}

		// update role in db
		item.TenantID = currentT.ID
		if err := orm.UpsertObject(db, item); err != nil {
			return err
		}

		// delete old parent
		object := casbinObjectEncode(item)
		os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
		for _, v := range os {
			if _, err := e.RemoveNamedGroupingPolicy("g2", []string{object, v, domain}); err != nil {
				return err
			}
		}

		// add new parent
		if item.ParentID != 0 {
			parent := &types.Object{ID: item.ParentID}
			if _, err := e.AddNamedGroupingPolicy("g2", object, casbinObjectEncode(parent), domain); err != nil {
				return err
			}
		}

		return nil
	})
}

// OrganizationObjectDeleteHandler 删除实体数据
func OrganizationObjectDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		// verify input
		item := &types.Object{ID: id}
		if err := orm.TakeObject(db, item); err != nil {
			return err
		}

		currentU, _ := userSessionGet(c)
		currentT, _ := currentTenantSessionGet(c)
		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		object := casbinObjectEncode(item)
		e := GetEnforcer(c)

		// check permission
		if ok, _ := e.Enforce(cu, domain, item.Object, write); !ok {
			return ErrNoWritePermission
		}

		// delete object as item in current tenant
		os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
		for _, v := range os {
			if _, err := e.RemoveNamedGroupingPolicy("g2", object, v, domain); err != nil {
				return err
			}
		}

		// delete object as group in current tenant
		if _, err := e.RemoveFilteredNamedGroupingPolicy("g2", 1, object, domain); err != nil {
			return err
		}

		// delete object policy
		if _, err := e.RemoveFilteredPolicy(1, domain, object); err != nil {
			return err
		}

		return orm.DeleteObjectById(db, id)
	})
}
