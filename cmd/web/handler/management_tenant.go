package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TenantGetHandler 拉取所有租户信息列表
func TenantGetHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		return orm.GetAllTenant(db)
	})
}

// TenantPostHandler 新增或者编辑租户信息
func TenantPostHandler(c *gin.Context) {
	item := &types.Tenant{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		// if exist and online, it will only update
		// else it will be reentrant to initialize tenant
		if err := orm.TakeTenant(db, item); err == nil {
			return orm.UpsertTenant(db, item)
		}

		return tenantInitialization(c, db, item)
	})
}

// TenantDeleteHandler 删除某个租户
func TenantDeleteHandler(c *gin.Context) {
	dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
		tenant := &types.Tenant{ID: id}
		if err := orm.TakeTenant(db, tenant); err != nil {
			return err
		}

		domain := casbinTenantEncode(tenant)
		e := GetEnforcer(c)
		gp := e.GetFilteredGroupingPolicy(2, domain)
		var rules [][]string
		for _, v := range gp {
			if casbinUserIs(v[0]) {
				rules = append(rules, v)
			}
		}

		if _, err := e.RemoveGroupingPolicies(rules); err != nil {
			return err
		}

		return orm.DeleteTenantById(db, id)
	})
}

func tenantInitialization(c *gin.Context, db *gorm.DB, tenant *types.Tenant) error {
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := orm.UpsertTenant(tx, tenant); err != nil {
			return err
		}

		object1 := &types.Object{Name: "数据实体", Type: types.ObjectTypeObject, TenantID: tenant.ID}
		object2 := &types.Object{Name: "角色", Type: types.ObjectTypeRole, TenantID: tenant.ID}
		object3 := &types.Object{Name: "默认", Type: types.ObjectTypeDefault, TenantID: tenant.ID}
		object4 := &types.Object{Name: "素材库", Type: types.ObjectTypeMaterial, TenantID: tenant.ID}
		object5 := &types.Object{Name: "数据报表", Type: types.ObjectTypeReport, TenantID: tenant.ID}
		object6 := &types.Object{Name: "帐号管家", Type: types.ObjectTypeAdvertiser, TenantID: tenant.ID}
		object7 := &types.Object{Name: "Schema", Type: types.ObjectTypeSchema, TenantID: tenant.ID}
		object8 := &types.Object{Name: "RTA策略管理", Type: types.ObjectTypeTargeting, TenantID: tenant.ID}
		for _, v := range []*types.Object{object1, object2, object3, object4, object5, object6, object7, object8} {
			if err := orm.UpsertObject(tx, v); err != nil {
				return err
			}
		}

		oo1 := casbinObjectEncode(object1)
		oo2 := casbinObjectEncode(object2)
		object1.Object = oo1
		object2.Object = oo1
		object3.Object = oo1
		object4.Object = oo1
		object5.Object = oo1
		object6.Object = oo1
		object7.Object = oo1
		object8.Object = oo1
		for _, v := range []*types.Object{object1, object2, object3, object4, object5, object6, object7, object8} {
			if err := orm.UpsertObject(tx, v); err != nil {
				return err
			}
		}

		role1 := &types.Role{Name: "管理员", Object: oo2, TenantID: tenant.ID}
		role2 := &types.Role{Name: "成员", Object: oo2, TenantID: tenant.ID}
		for _, v := range []*types.Role{role1, role2} {
			if err := orm.UpsertRole(tx, v); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	t := casbinTenantEncode(tenant)

	object0 := getFeatureRootObject()
	object1 := &types.Object{Name: "数据实体", Type: types.ObjectTypeObject, TenantID: tenant.ID}
	object2 := &types.Object{Name: "角色", Type: types.ObjectTypeRole, TenantID: tenant.ID}
	object3 := &types.Object{Name: "默认", Type: types.ObjectTypeDefault, TenantID: tenant.ID}
	object4 := &types.Object{Name: "素材库", Type: types.ObjectTypeMaterial, TenantID: tenant.ID}
	object5 := &types.Object{Name: "数据报表", Type: types.ObjectTypeReport, TenantID: tenant.ID}
	object6 := &types.Object{Name: "帐号管家", Type: types.ObjectTypeAdvertiser, TenantID: tenant.ID}
	object7 := &types.Object{Name: "Schema", Type: types.ObjectTypeSchema, TenantID: tenant.ID}
	object8 := &types.Object{Name: "RTA策略管理", Type: types.ObjectTypeTargeting, TenantID: tenant.ID}
	for _, v := range []*types.Object{object1, object2, object3, object4, object5, object6, object7, object8} {
		if err := orm.TakeObject(db, v); err != nil {
			return err
		}
	}
	o0 := casbinObjectEncode(object0)
	o1 := casbinObjectEncode(object1)
	o2 := casbinObjectEncode(object2)
	o3 := casbinObjectEncode(object3)
	o4 := casbinObjectEncode(object4)
	o5 := casbinObjectEncode(object5)
	o6 := casbinObjectEncode(object6)
	o7 := casbinObjectEncode(object7)
	o8 := casbinObjectEncode(object8)

	role1 := &types.Role{Name: "管理员", TenantID: tenant.ID}
	role2 := &types.Role{Name: "成员", TenantID: tenant.ID}
	for _, v := range []*types.Role{role1, role2} {
		if err := orm.TakeRole(db, v); err != nil {
			return err
		}
	}
	r1 := casbinRoleEncode(role1)
	r2 := casbinRoleEncode(role2)

	e := GetEnforcer(c)

	for _, v := range [][]string{
		{r1, t, o0, read},
		{r1, t, o1, read},
		{r1, t, o1, write},
		{r1, t, o2, read},
		{r1, t, o2, write},
		{r1, t, o3, read},
		{r1, t, o3, write},
		{r1, t, o4, read},
		{r1, t, o4, write},
		{r1, t, o5, read},
		{r1, t, o5, write},
		{r1, t, o6, read},
		{r1, t, o6, write},
		{r1, t, o7, read},
		{r1, t, o7, write},
		{r1, t, o8, read},
		{r1, t, o8, write},
		{r2, t, o3, read},
		{r2, t, o3, write},
	} {
		if _, err := e.AddPolicy(v); err != nil {
			return err
		}
	}

	return nil
}
