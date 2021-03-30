package router

import (
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/gin-gonic/gin"
)

func setupAPI(router gin.IRouter) error {
	router.Use(handler.DeveloperAppMiddleware())
	v1 := router.Group("/v1")
	rtaServer := v1.Group("/rta_server")
	{
		rtaServer.GET("/experiment_config", handler.GetRtaServerExperimentConfig)
	}

	return nil
}
