package router

import (
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/gin-gonic/gin"
)

func setupDocs(router gin.IRouter, config *config.WebConfiguration) error {
	router.Static("/", config.DocsFileRoot)
	return nil
}
