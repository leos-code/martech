package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

var (
	RtaExpDataUrl     = "https://api.rta.qq.com/api/v1/exp/data"
	ExpDataHttpClient = &http.Client{}
)

// ExpDataGranularity 查询时间粒度，5分钟粒度支持3天内的数据查询，小时粒度的支持7天内的数据查询，天数据支持90天内数据查询
type ExpDataGranularity string

const (
	ExpDataGranularityFiveMinute ExpDataGranularity = "min_5"
	ExpDataGranularityHour       ExpDataGranularity = "hour"
	ExpDataGranularityDay        ExpDataGranularity = "day"
)

// ExpDataGroupBy 数据聚合维度
type ExpDataGroupBy string

const (
	ExpDataGroupByAdvertiserID ExpDataGroupBy = "advertiser_id"
	ExpDataGroupByCampaignID   ExpDataGroupBy = "campaign_id"
	ExpDataGroupByAppID        ExpDataGroupBy = "app_id"
	ExpDataGroupByAdID         ExpDataGroupBy = "ad_id"
)

// RtaExpDataRequest RTA报表数据请求
type RtaExpDataRequest struct {
	*RtaExpRequestBase `json:"-"`
	Data               *RtaExpDataRequestData `json:"data"`
}

// RtaExpDataRequestData RTA报表数据请求中的data结构体
type RtaExpDataRequestData struct {
	Granularity ExpDataGranularity    `json:"Granularity"` // 否 时间粒度，仅支持：min_5,hour,day，分别是5分钟、小时、天
	BeginTime   RtaExpDataRequestTime `json:"BeginTime"`   // 是 开始时间，格式202011011500
	EndTime     RtaExpDataRequestTime `json:"EndTime"`     // 是 结束时间，格式202011012300
	ExpID       []int64               `json:"ExpID"`       // 是 查询指定的实验ID
	UID         string                `json:"UId"`         // 否 广告主ID
	AppID       string                `json:"AppId"`       // 否 APP ID
	CID         string                `json:"CId"`         // 否 推广计划ID
	AID         string                `json:"AId"`         // 否 广告ID
	TotalFlag   int                   `json:"TotalFlag"`   // 否 是否汇总，0：不汇总  1：汇总 ，不传默认不汇总
	GroupBy     []ExpDataGroupBy      `json:"GroupBy"`     // 否 展开字段，指定GroupBy后，UID/AppID/CID/AID不能为空
}

type RtaExpDataRequestTime time.Time

// MarshalJSON 覆盖json marshal
func (t RtaExpDataRequestTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("200601021504"))
}

// RtaExpDataRequest RTA报表响应结构
type RtaExpDataResponse struct {
	Status string          `json:"status"`
	Code   int             `json:"code"`
	Data   RtaExpDataSlice `json:"data"`
}

// RtaExpDataSlice ams返回的结构体可能不是list，需要自定义unmarshal去兼容
type RtaExpDataSlice []*ExpData

// UnmarshalJSON 覆盖json unmarshal
func (s *RtaExpDataSlice) UnmarshalJSON(d []byte) error {
	if d[0] == '[' {
		tmp := make([]*ExpData, 0)
		if err := json.Unmarshal(d, &tmp); err != nil {
			return err
		}
		*s = tmp
	}

	return nil
}

// SortRtaExpDataSlice 按照时间对实验报表结果进行排序
func SortRtaExpDataSlice(s RtaExpDataSlice) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Time < s[j].Time
	})
}

// ExpData 实验报表数据
type ExpData struct {
	Time             ExpDataTime `json:"time"`              // 时间
	ExpId            int64       `json:"exp_id"`            // 实验ID
	GroupInfo        *GroupInfo  `json:"group_info"`        // 分组信息
	Cost             float64     `json:"cost"`              // 消耗
	Exposure         int         `json:"exposure"`          // 曝光
	Click            int         `json:"click"`             // 点击
	Cpm              float64     `json:"cpm"`               // cpm
	Cpc              float64     `json:"cpc"`               // cpc
	Ctr              float64     `json:"ctr"`               // ctr
	Conversion       int         `json:"conversion"`        // 浅层转化量
	ConversionSecond int         `json:"conversion_second"` // 深层转化量（第二目标）
	Cvr              float64     `json:"cvr"`               // 浅层cvr
	CvrSecond        float64     `json:"cvr_second"`        // 深层cvr（第二目标）
}

type ExpDataTime int64

// UnmarshalJSON 覆盖json unmarshal
func (t *ExpDataTime) UnmarshalJSON(bytes []byte) error {
	var tm time.Time
	var err error
	length := len(bytes)

	if length == 8 {
		tm, err = time.ParseInLocation("20060102", string(bytes), time.Local)
	} else if length == 10 {
		tm, err = time.ParseInLocation("2006010215", string(bytes), time.Local)
	} else if length == 12 {
		tm, err = time.ParseInLocation("200601021504", string(bytes), time.Local)
	}

	if err != nil {
		return err
	}

	*t = ExpDataTime(tm.Unix())
	return nil
}

// GroupInfo 聚合结果
type GroupInfo struct {
	AdId         string `json:"ad_id"`         // GroupBy选择了ad_id就会出现
	AdvertiserId string `json:"advertiser_id"` // GroupBy选择了advertiser_id就会出现
	AppId        string `json:"app_id"`        // GroupBy选择了app_id就会出现
	CampaignId   string `json:"campaign_id"`   // GroupBy选择了campaign_id就会出现
}

// GetRtaExpData 获取实验参数
func GetRtaExpData(req *RtaExpDataRequest) (*RtaExpDataResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", RtaExpDataUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	httpResp, err := ExpDataHttpClient.Do(httpReq)
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

	var resp RtaExpDataResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("%s, code[%d]: %s", resp.Status, resp.Code, RtaExpRespCodeMap[resp.Code])
	}

	return &resp, nil
}
