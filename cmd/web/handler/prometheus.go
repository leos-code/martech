package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const namespace = "experiment_system"

var (
	labelNames = []string{"status", "url", "method"}

	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "total number of HTTP requests",
		}, labelNames,
	)

	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "cost_time",
			Help:      "HTTP handle cost time",
		}, labelNames,
	)

	reqErrCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "req_err_count_total",
			Help:      "total number of HTTP requests that reports an error",
		}, labelNames,
	)

	excludeRexStatus = ""
	excludeRexUrl    = "/metrics$"
	excludeRexMethod = ""
)

func checkLabelValue(labelValue string, pattern string) bool {
	if pattern == "" {
		return true
	}
	matched, err := regexp.MatchString(pattern, labelValue)
	if err != nil {
		return true
	}
	return !matched
}

func init() {
	prometheus.MustRegister(reqCount, reqDuration, reqErrCount)
}

func getFormattedUrl(c *gin.Context) string {
	url := c.Request.URL.Path
	id := c.Param("id")
	if id == "" {
		return url
	}
	index := strings.LastIndex(url, "/"+id)
	if index == -1 {
		return url
	}
	return url[0:index]
}

// PrometheusMetricsMiddleware 添加prometheus监控的中间件
func PrometheusMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		url := getFormattedUrl(c)
		method := c.Request.Method

		if isOk := checkLabelValue(status, excludeRexStatus) &&
			checkLabelValue(url, excludeRexUrl) &&
			checkLabelValue(method, excludeRexMethod); !isOk {
			return
		}

		labelValues := []string{status, url, method}

		reqCount.WithLabelValues(labelValues...).Inc()
		reqDuration.WithLabelValues(labelValues...).Observe(time.Since(start).Seconds())
		if len(c.Errors) > 0 {
			reqErrCount.WithLabelValues(labelValues...).Inc()
		}
	}
}

// PromHandler 暴露提供prometheus拉取数据的接口
func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
