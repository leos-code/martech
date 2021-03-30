package ginx

import (
	"encoding/json"

	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
)

// GetFilterOption 获取过滤参数
// 过滤参数为 ?filter=URL编码([]*orm.FilterOption)
func GetFilterOption(c *gin.Context) ([]*orm.FilterOption, error) {
	filter, ok := c.GetQuery("filter")
	if !ok || filter == "" {
		return nil, nil
	}
	var err error
	options := make([]*orm.FilterOption, 0)
	if err = json.Unmarshal([]byte(filter), &options); err != nil {
		return nil, err
	}

	return options, nil
}
