package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseWithStatus200 以状态200返回响应数据
func ResponseWithStatus200(c *gin.Context, msg string, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  msg,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  msg,
			"data": data[0],
		})
	}
}

// ResponseWithError 以状态code返回响应数据和错误
func ResponseWithError(c *gin.Context, code int, err error) {
	_ = c.Error(err)
	c.AbortWithStatusJSON(code, gin.H{
		"code": -1,
		"msg":  err.Error(),
	})
}
