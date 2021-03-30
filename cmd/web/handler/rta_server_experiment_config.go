package handler

import (
	"net/http"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/gin-gonic/gin"
)

// GetRtaServerExperimentConfig 获取Rta服务器实验配置的接口
func GetRtaServerExperimentConfig(c *gin.Context) {
	db := orm.GetDB()
	rtaExp, err := orm.GetRtaServerExperimentConfig(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newExperimentConfigErrResponse(err))
	} else {
		c.JSON(http.StatusOK, convertRtaExp(rtaExp))
	}
}

func newExperimentConfigErrResponse(err error) *experimentConfigResponse {
	return &experimentConfigResponse{
		Code:    -1,
		Message: err.Error(),
	}
}

type experimentConfigResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    []*experimentRecord `json:"data"`
}

type experimentRecord struct {
	ExpId     string                 `json:"exp_id"`
	Parameter []*experimentParameter `json:"parameter"`
}

type experimentParameter struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func convertRtaExp(rtaExp []*types.RtaExp) *experimentConfigResponse {
	resp := &experimentConfigResponse{
		Code:    0,
		Message: "success",
	}

	for _, exp := range rtaExp {
		expId := exp.ID

		if exp.BindExperimentItem == nil {
			continue
		}

		parameters := make([]*experimentParameter, 0, len(exp.BindExperimentItem.ExperimentMetadata))
		for _, metadata := range exp.BindExperimentItem.ExperimentMetadata {
			parameters = append(parameters, &experimentParameter{
				Name:  metadata.ExperimentParameter.Name,
				Type:  metadata.ExperimentParameter.Type,
				Value: metadata.Value,
			})
		}

		resp.Data = append(resp.Data, &experimentRecord{
			ExpId:     expId,
			Parameter: parameters,
		})
	}
	return resp
}
