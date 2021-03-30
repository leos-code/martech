package router

import (
	"fmt"

	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/tencentad/martech/cmd/web/handler/login/rio"
	"github.com/gin-gonic/gin"
)

func setupPortal(router gin.IRouter, config *config.LoginOptions) error {
	if config == nil {
		return fmt.Errorf("there is no login config")
	}

	login := router.Group("/login")
	login.GET("/oauth", handler.LoginOauthHandler(config))

	cb := login.Group("/callback")
	if config.Rio != nil {
		cb.GET(rio.RelativePath, rio.CallBackHandler(config.Rio), handler.LoginUserHandler)
	}

	router.POST("/register", handler.RegisterPostHandler)
	router.GET("/register", handler.RegisterGetHandler)

	router.GET("/logout", handler.LogoutHandler)

	return nil
}
