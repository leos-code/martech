package handler

import (
	"fmt"

	"github.com/tencentad/martech/api/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	loginUserKey  = "login_user"
	currentTenant = "current_tenant"
)

func sessionSet(c *gin.Context, k string, v interface{}) error {
	session := sessions.Default(c)
	session.Set(k, v)
	return session.Save()
}

func sessionGet(c *gin.Context, k string) interface{} {
	session := sessions.Default(c)
	return session.Get(k)
}

func userSessionSet(c *gin.Context, user *types.User) error {
	session := sessions.Default(c)
	session.Set(gin.AuthUserKey, user)
	session.Set(currentTenant, nil)
	return session.Save()
}

func userSessionGet(c *gin.Context) (*types.User, error) {
	v, ok := sessionGet(c, gin.AuthUserKey).(*types.User)
	if !ok {
		return nil, fmt.Errorf("session is invalid")
	}
	return v, nil
}

func loginUserSessionSet(c *gin.Context, user *types.LoginUser) error {
	return sessionSet(c, loginUserKey, user)
}

func loginUserSessionGet(c *gin.Context) (*types.LoginUser, error) {
	v, ok := sessionGet(c, loginUserKey).(*types.LoginUser)
	if !ok {
		return nil, fmt.Errorf("session is invalid")
	}
	return v, nil
}

func currentTenantSessionGet(c *gin.Context) (*types.Tenant, error) {
	v, ok := sessionGet(c, currentTenant).(*types.Tenant)
	if !ok {
		return nil, fmt.Errorf("session is invalid")
	}
	return v, nil
}
