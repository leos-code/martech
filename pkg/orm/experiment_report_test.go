package orm

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
)

var (
	dataStr = `
{
  "status": "success",
  "code": 0,
  "data": [
    {
      "time": 2020112601,
      "exp_id": 11,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112601,
      "exp_id": 12,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112601,
      "exp_id": 22,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112601,
      "exp_id": 32,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112601,
      "exp_id": 31,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112601,
      "exp_id": 33,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 11,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 12,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 22,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 32,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 31,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    },
    {
      "time": 2020112600,
      "exp_id": 33,
      "group_info": {},
      "cost": 623.01596,
      "exposure": 25256,
      "click": 635,
      "cpm": 0.024668037694013302,
      "cpc": 0.9811274960629921,
      "ctr": 0.025142540386442827,
      "conversion": 233,
      "conversion_second": 0,
      "cvr": 0.3669291338582677,
      "cvr_second": 0
    }
  ]
}
`
)

func TestGetExperimentReport(t *testing.T) {
	{
		err := gohook.HookMethod(
			api.ExpDataHttpClient, "Do",
			func(_ *http.Client, req *http.Request) (*http.Response, error) {
				return &http.Response{
					Body: ioutil.NopCloser(bytes.NewBufferString(dataStr))}, nil
			}, nil)

		// 构造DB数据
		db, err := getTestDB(t)
		assert.NoError(t, Setup(db))
		assert.NoError(t, err)

		assert.NoError(t, UpsertRtaAccount(db, &types.RtaAccount{
			RtaID: "123",
			Token: "456",
		}))
		assert.NoError(t, UpsertExperimentGroup(db, &types.ExperimentGroup{
			RtaAccountID: 1,
		}))

		stage := &types.ExperimentStage{
			ExperimentGroupID: 1,
			ExperimentItem: []*types.ExperimentItem{
				{
					Name: "base",
					RtaExp: []*types.RtaExp{
						{
							ID: "ams::11",
						},
						{
							ID: "ams::12",
						},
					},
				},
				{
					Name: "lab",
					RtaExp: []*types.RtaExp{
						{
							ID: "ams::21",
						},
						{
							ID: "ams::22",
						},
					},
				},
				{
					Name: "aa",
					RtaExp: []*types.RtaExp{
						{
							ID: "ams::31",
						},
						{
							ID: "ams::32",
						},
					},
				},
			},
		}

		assert.NoError(t, UpsertExperimentStage(db, stage))
		for _, item := range stage.ExperimentItem {
			assert.NoError(t, UpsertExperimentItem(db, item))
		}
		{
			param := &types.ExperimentReportParameter{
				ExperimentStageId:    1,
				BaseExperimentItemId: 1,
				LabExperimentItemId:  []uint64{3, 2},
				Granularity:          api.ExpDataGranularityFiveMinute,
				TimeMergeType:        types.TimeMergeNo,
			}

			resp, err := GetExperimentReport(db, param)
			assert.NoError(t, err)
			assert.Len(t, resp.Records, 6)
		}
		{
			param := &types.ExperimentReportParameter{
				ExperimentStageId:    1,
				BaseExperimentItemId: 1,
				LabExperimentItemId:  []uint64{3, 2},
				Granularity:          api.ExpDataGranularityFiveMinute,
				TimeMergeType:        types.TimeMergeYes,
			}

			resp, err := GetExperimentReport(db, param)
			assert.NoError(t, err)
			assert.Len(t, resp.Records, 3)
		}
	}
}
