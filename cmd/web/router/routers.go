package router

import (
	"net/http"
	"time"

	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// SetupRouter 配合router
func SetupRouter(config *config.WebConfiguration) (*gin.Engine, error) {
	r := gin.Default()

	r.Use(requestid.New())
	r.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	r.Use(handler.PrometheusMetricsMiddleware())

	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusFound, "/page")
	})
	r.GET("/page", func(context *gin.Context) {
		context.Redirect(http.StatusFound, "/page/dashboard")
	})

	if err := setupSessionStore(r, config.Session); err != nil {
		return nil, err
	}

	if err := setupEnforcer(r, config.Casbin); err != nil {
		return nil, err
	}

	page := r.Group("/page")
	if err := setupPage(page, config); err != nil {
		return nil, err
	}

	docs := r.Group("/docs")
	if err := setupDocs(docs, config); err != nil {
		return nil, err
	}

	dashboard := r.Group("/dashboard")
	if err := setupDashboard(dashboard); err != nil {
		return nil, err
	}

	login := r.Group("/portal")
	if err := setupPortal(login, config.Login); err != nil {
		return nil, err
	}

	api := r.Group("/api")
	if err := setupAPI(api); err != nil {
		return nil, err
	}

	r.GET("/metrics", handler.PromHandler(promhttp.Handler()))

	return r, nil
}
