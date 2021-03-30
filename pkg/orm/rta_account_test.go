package orm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"
)

func TestRtaAccount(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.RtaAccount{
		Token: "rta_account",
	}
	assert.NoError(t, UpsertRtaAccount(db, item))

	list, err := GetAllRtaAccount(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	findMetadata, err := GetRtaAccountById(db, list[0].ID)
	assert.NoError(t, err)
	fmt.Println(findMetadata)

	assert.NoError(t, DeleteRtaAccountById(db, list[0].ID))

	list, err = GetAllRtaAccount(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestSyncRtaAccountExp(t *testing.T) {
	err := gohook.HookMethod(
		api.ExpListHttpClient, "Do",
		func(_ *http.Client, req *http.Request) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`
{
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
}
`))}, nil
		}, nil)

	assert.NoError(t, err)

	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.RtaAccount{
		Token: "rta_account",
		RtaID: "rta_id",
	}
	assert.NoError(t, UpsertRtaAccount(db, item))

	list, err := GetAllRtaAccount(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	account, err := GetRtaAccountById(db, list[0].ID)
	assert.NoError(t, err)
	fmt.Println(account)

	assert.NoError(t, SyncRtaAccountExp(db, account))
	account, err = GetRtaAccountById(db, list[0].ID)
	assert.NoError(t, err)
	assert.NoError(t, LoadAccountRtaExp(db, account))
	assert.Len(t, account.RtaExp, 3)

	assert.NoError(t, SyncRtaAccountExp(db, account))
}
