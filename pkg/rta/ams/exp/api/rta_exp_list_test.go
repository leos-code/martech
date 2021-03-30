package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"
)

func TestRTAExpListResponse(t *testing.T) {
	{
		var resp RtaExpListResponse
		data := `{
  "status": "success",
  "data": [
    {
      "ExpId": 123478,
      "ExpName": "基础实验",
      "FlowRate": 20,
      "EndTime": "2021-02-09 15:56:23",
      "SiteSet": [
        21,
        15,
        25
      ],
      "Status": 1
    },
    {
      "ExpId": 123480,
      "ExpName": "控制实验",
      "FlowRate": 80,
      "EndTime": "2021-02-09 15:56:23",
      "SiteSet": [
        21,
        15,
        25
      ],
      "Status": 1
    },
    {
      "ExpId": 123508,
      "ExpName": "基础实验",
      "FlowRate": 21,
      "EndTime": "2020-11-12 14:02:00",
      "SiteSet": [],
      "Status": 3
    }
  ],
 "code": 1
}`
		err := json.Unmarshal([]byte(data), &resp)
		assert.NoError(t, err)

		assert.EqualValues(t, 1, resp.Code)
		assert.EqualValues(t, "success", resp.Status)
		assert.Len(t, resp.Data, 3)
		assert.EqualValues(t, 123480, resp.Data[1].ExpId)
		assert.EqualValues(t, "控制实验", resp.Data[1].ExpName)
		assert.EqualValues(t, 80, resp.Data[1].FlowRate)
		assert.EqualValues(t, time.Date(2021, 2, 9, 15, 56, 23, 0, time.Local), resp.Data[1].EndTime)
		assert.Len(t, resp.Data[1].SiteSet, 3)
		assert.EqualValues(t, 1, resp.Data[1].Status)
	}

	{
		var resp RtaExpListResponse
		data := `{
  "status": "success",
  "data": [""]
}`

		err := json.Unmarshal([]byte(data), &resp)
		assert.Error(t, err)
	}

	{
		var resp RtaExpListResponse
		data := `{
  "status": "success",
  "data": ""
}`

		err := json.Unmarshal([]byte(data), &resp)
		assert.NoError(t, err)
	}
}

func TestGetRtaExpList(t *testing.T) {
	err := gohook.HookMethod(
		ExpListHttpClient, "Do",
		func(_ *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{
  "status": "success",
  "code": 0
}
`))}, nil
		}, nil)

	assert.NoError(t, err)
	resp, err := GetRtaExpList(&RtaExpListRequest{
		RtaExpRequestBase: &RtaExpRequestBase{
			RtaID: "",
			Token: "",
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFormatAmsExpId(t *testing.T) {
	{
		var amsExpId int64 = 123
		expId := FormatAmsExpId(amsExpId)

		parsedAmsExpId, err := ParseAmsExpId(expId)
		assert.NoError(t, err)
		assert.EqualValues(t, amsExpId, parsedAmsExpId)
	}
}
