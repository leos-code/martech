package ginx

import (
	"strconv"

	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
)

const defaultPageSize = 10

// getPaginationOption 获取分页信息
func GetPaginationOption(c *gin.Context) (opt *orm.PaginationOption, err error) {
	pageStr, ok := c.GetQuery("page")
	if !ok {
		return
	}

	var page int
	page, err = strconv.Atoi(pageStr)
	if err != nil {
		return
	}

	var pageSize int
	if pageSizeStr, ok := c.GetQuery("page_size"); ok {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return
		}
	} else {
		pageSize = defaultPageSize
	}

	opt = &orm.PaginationOption{
		Page:     page,
		PageSize: pageSize,
	}
	return
}
