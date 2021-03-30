package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"
)

func TestGetRtaExpData(t *testing.T) {
	err := gohook.HookMethod(
		ExpDataHttpClient, "Do",
		func(_ *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`
{
  "status": "success",
  "code": 0,
  "data": [
    {
      "time": 2020112600,
      "exp_id": 124996,
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
`))}, nil
		}, nil)

	assert.NoError(t, err)
	resp, err := GetRtaExpData(&RtaExpDataRequest{
		RtaExpRequestBase: &RtaExpRequestBase{
			RtaID: "",
			Token: "",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.EqualValues(t, 1606320000, resp.Data[0].Time)
	assert.Len(t, resp.Data, 1)
}

func TestRtaExpDataRequestTime_MarshalJSON(t *testing.T) {
	var timestamp int64 = 1606320000
	tm := time.Unix(timestamp, 0)
	data, err := (RtaExpDataRequestTime(tm)).MarshalJSON()
	assert.NoError(t, err)
	assert.EqualValues(t, `"202011260000"`, string(data))
}
