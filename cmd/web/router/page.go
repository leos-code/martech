package router

import (
	"io/ioutil"
	"strings"

	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/gin-gonic/gin"
)

func setupPage(router gin.IRouter, config *config.WebConfiguration) error {
	files, err := ioutil.ReadDir(config.WebFileRoot)
	if err != nil {
		return err
	}

	for _, v := range files {
		relativePath := "/" + v.Name()
		root := config.WebFileRoot + "/" + v.Name()
		switch v.IsDir() {
		case true:
			router.Static(relativePath, root)
		case false:
			// skip dashboard.html and portal.html
			// dashboard.html must use auth middleware
			if strings.HasSuffix(v.Name(), ".html") {
				continue
			}

			router.StaticFile(relativePath, root)
		}
	}

	dashboard := router.Group("/dashboard")
	dashboard.Use(handler.AuthLoginMiddleware())
	dashboard.StaticFile("", config.WebFileRoot+"/"+"dashboard.html")

	portal := router.Group("/portal")
	portal.StaticFile("", config.WebFileRoot+"/"+"portal.html")

	proxy := router.Group("/proxy")
	proxy.StaticFile("", config.WebFileRoot+"/"+"proxy.html")
	return nil
}
