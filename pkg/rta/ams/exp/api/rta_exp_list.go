package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	RtaExpListUrl     = "https://api.rta.qq.com/api/v1/exp/list"
	ExpListHttpClient = &http.Client{}
)

// RtaExpListRequest 请求RtaExp列表的请求数据格式
type RtaExpListRequest struct {
	*RtaExpRequestBase
}

// RtaExpListResponse 请求RtaExp列表数据的请求对应的响应格式
type RtaExpListResponse struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   RtaExpSlice `json:"data"`
}

// RtaExpSlice ams返回的结构体可能不是list，需要自定义unmarshal去兼容
type RtaExpSlice []*RtaExp

// UnmarshalJSON 覆盖json.Unmarshal
func (s *RtaExpSlice) UnmarshalJSON(d []byte) error {
	if d[0] == '[' {
		tmp := make([]*RtaExp, 0)
		if err := json.Unmarshal(d, &tmp); err != nil {
			return err
		}
		*s = tmp
	}

	return nil
}

// RtaExp rta实验信息
type RtaExp struct {
	ExpId    int64         // 实验ID
	ExpName  string        // 实验名称
	FlowRate int           // 流量百分比
	EndTime  RtaExpEndTime // 结束时间
	SiteSet  []int         // 站点集，未配置默认为全部
	Status   int           // 实验状态，定义见下方
}

// ExpStatusMap rta实验状态
var ExpStatusMap = map[int]string{
	1: "已启用",
	2: "已暂停",
	3: "已结束",
}

type RtaExpEndTime time.Time

// UnmarshalJSON 覆盖json.Unmarshal
func (t *RtaExpEndTime) UnmarshalJSON(d []byte) error {
	s := strings.Trim(string(d), `"`)

	tmp, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	if err != nil {
		return err
	}

	*t = RtaExpEndTime(tmp)
	return nil
}

// GetRtaExpList 获取rta实验列表
func GetRtaExpList(req *RtaExpListRequest) (*RtaExpListResponse, error) {
	httpReq, err := http.NewRequest("POST", RtaExpListUrl, nil)
	if err != nil {
		return nil, err
	}

	FillRTAExpRequestHeader(httpReq.Header, req.RtaID, req.Token)

	httpResp, err := ExpListHttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return nil, err
	}

	var resp RtaExpListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if strings.ToLower(resp.Status) != "success" {
		return nil, fmt.Errorf("%s, code[%d]: %s", resp.Status, resp.Code, RtaExpRespCodeMap[resp.Code])
	}

	return &resp, nil
}

const AmsExpIdPrefix = "ams::"

// FormatAmsExpId 生成包含ams前缀的实验ID
func FormatAmsExpId(amsExpId int64) string {
	return AmsExpIdPrefix + strconv.FormatInt(amsExpId, 10)
}

// ParseAmsExpId 解析出原始的实验ID
func ParseAmsExpId(expId string) (int64, error) {
	if !strings.HasPrefix(expId, AmsExpIdPrefix) {
		return 0, fmt.Errorf("not valid ams exp id[%s]", expId)
	} else {
		return strconv.ParseInt(expId[len(AmsExpIdPrefix):], 10, 64)
	}
}
