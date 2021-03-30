package handler

import (
	"fmt"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

const (
	defaultUser = "yigeshadiao"
)
// UserInfoGetHandler 拉取该用户数据
func UserInfoGetHandler(c *gin.Context) {
	user, _ := userSessionGet(c)
	loginUser, _ := loginUserSessionGet(c)
	current, _ := currentTenantSessionGet(c)
	e := GetEnforcer(c)

	tenant, err := userTenant(c, user)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
	}

	if isSuperAdmin(e, user) {
		tenant = append(tenant, &types.Tenant{
			ID:   SuperAdminDomainID,
			Name: SuperAdminDomainName,
		})
	}

	for _, v := range tenant {
		if current != nil && current.ID == v.ID {
			current = v
		}
	}

	info := &types.UserInfo{
		User:          user,
		LoginUser:     loginUser,
		CurrentTenant: current,
		Tenant:        tenant,
	}

	ginx.ResponseWithStatus200(c, "get success", info)
}

// UserInfoPostHandler 编辑用户信息
func UserInfoPostHandler(c *gin.Context) {
	item := &types.User{}
	dbEditHandler(c, item, func(db *gorm.DB) error {
		user, _ := userSessionGet(c)
		if item.ID != user.ID {
			return fmt.Errorf("it is not current login user")
		}

		if err := orm.UpsertUser(db, item); err != nil {
			return err
		}

		return userSessionSet(c, item)
	})

}

// UserTenantHandler 选择该用户当前组织
func UserTenantHandler(c *gin.Context) {
	tenant := &types.Tenant{}
	dbEditHandler(c, tenant, func(db *gorm.DB) error {
		user, _ := userSessionGet(c)
		u := casbinUserEncode(user)
		t := casbinTenantEncode(tenant)

		e := GetEnforcer(c)
		if !isSuperAdmin(e, user) {
			roles := e.GetRolesForUserInDomain(u, t)
			if len(roles) == 0 {
				return fmt.Errorf("current user has no permission in domain: %v", tenant)
			}
		}

		return sessionSet(c, currentTenant, tenant)
	})
}

// UserAuthorityHandler 拉取该用户权限数据
func UserAuthorityHandler(c *gin.Context) {
	dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
		subject := c.DefaultQuery("subject", defaultUser)

		data := gin.H{
			"m": GetCasbinModel(c),
		}

		currentU, _ := userSessionGet(c)
		currentT, err := currentTenantSessionGet(c)
		if err != nil {
			return data, nil
		}

		cu := casbinUserEncode(currentU)
		domain := casbinTenantEncode(currentT)
		e := GetEnforcer(c)

		var p [][]string
		var g [][]string
		var g2 [][]string

		if isSuperAdmin(e, cu) && currentT.ID == SuperAdminDomainID {
			g = append(g, []string{"g", subject, superAdminRole, superAdminDomain})
		} else {
			frontend, _ := orm.GetAllFrontend(db)
			for _, v := range frontend {
				o := casbinObjectEncode(&types.Object{ID: v.ObjectID})
				if ok, _ := e.Enforce(cu, domain, o, read); !ok {
					continue
				}

				key := fmt.Sprint(v.Key, "#", v.Type)
				p = append(p, []string{"p", subject, domain, key, read})
			}
		}

		data["p"] = p
		data["g"] = g
		data["g2"] = g2

		return data, nil
	})
}

// UserSearchHandler 搜索用户信息
func UserSearchHandler(c *gin.Context) {
	info := &types.UserSearch{}
	dbGetHandler(c, info, func(db *gorm.DB) (interface{}, error) {
		return orm.SearchUser(db, info)
	})
}

func userTenant(c *gin.Context, user *types.User) ([]*types.Tenant, error) {
	db := orm.GetDB()
	e := GetEnforcer(c)
	if isSuperAdmin(e, user) {
		return orm.GetAllTenant(db)
	}

	rules := e.GetFilteredGroupingPolicy(0, casbinUserEncode(user))
	var id []uint64
	for _, rule := range rules {
		v := rule[2]
		id = append(id, casbinTenantDecode(v).ID)
	}

	linq.From(id).Distinct().ToSlice(&id)
	return orm.ListTenantById(db, id)
}
