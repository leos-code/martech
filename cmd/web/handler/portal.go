package handler

import (
	"net/http"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/cmd/web/handler/login/rio"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginOauthHandler 拉取登录列表接口
func LoginOauthHandler(config *config.LoginOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		list := types.RedirectLogins{}
		if config.Rio != nil {
			list = append(list, rio.GetRedirect(config.Rio))
		}

		ginx.ResponseWithStatus200(c, "get success", list)
	}
}

// LogoutHandler 登出接口
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	ginx.ResponseWithStatus200(c, "logout success")
}

// LoginUserHandler 用户登录处理接口
func LoginUserHandler(c *gin.Context) {
	loginUser := c.MustGet(gin.AuthUserKey).(*types.LoginUser)
	if err := loginUserSessionSet(c, loginUser); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	db := orm.GetDB()
	if err := orm.TakeLoginUser(db, loginUser); err != nil {
		c.Redirect(http.StatusFound, "/page/proxy#login/callback?redirect_type=register")
		return
	}

	user := &types.User{ID: loginUser.UserID}
	if err := orm.TakeUser(db, user); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	if err := userSessionSet(c, user); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/page/proxy#login/callback?redirect_type=dashboard")
}

// RegisterPostHandler 提交用户信息接口
func RegisterPostHandler(c *gin.Context) {
	loginUser, err := loginUserSessionGet(c)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusUnauthorized, err)
		return
	}

	info := &types.RegisterUserInfo{}
	if err := c.ShouldBindJSON(info); err != nil {
		ginx.ResponseWithError(c, http.StatusBadRequest, err)
		return
	}

	user := &types.User{
		PhoneNumber: info.PhoneNumber,
		Email:       info.Email,
		LoginUser:   []*types.LoginUser{loginUser},
	}

	db := orm.GetDB()
	_ = orm.TakeUser(db, user)
	if err := orm.UpsertUser(db, user); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	if err := userSessionSet(c, user); err != nil {
		ginx.ResponseWithError(c, http.StatusInternalServerError, err)
		return
	}

	ginx.ResponseWithStatus200(c, "register success")
}

// RegisterGetHandler 获取用户登录信息
func RegisterGetHandler(c *gin.Context) {
	loginUser, err := loginUserSessionGet(c)
	if err != nil {
		ginx.ResponseWithError(c, http.StatusUnauthorized, err)
		return
	}

	ginx.ResponseWithStatus200(c, "get success", loginUser)
}
