package handler

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/tencentad/martech/pkg/druid"
	"github.com/tencentad/martech/pkg/druid/api"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExperimentReportGetHandler 获取实验报表的接口
func ExperimentReportGetHandler(c *gin.Context) {
	parameter := &types.ExperimentReportParameter{}
	dbGetHandler(c, parameter, func(db *gorm.DB) (interface{}, error) {
		return orm.GetExperimentReport(db, parameter)
	})
}

func ExperimentReportGetAttributionHandler(c *gin.Context) {
	client := druid.GetDruidClient()
	queryResult, err := api.GetAttributionStaticData(client)
	if err != nil {
		ginx.ResponseWithError(c, 500, err)
		return
	}
	handleWithData(c, "query druid success", queryResult, err)

}
