package orm

import (
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
)

// GetExperimentReport 获取实验报表
func GetExperimentReport(db *gorm.DB, param *types.ExperimentReportParameter) (*types.ExperimentReportResponse, error) {
	c := &experimentReportContext{
		db:       db,
		param:    param,
		response: &types.ExperimentReportResponse{},
	}
	var err error
	if err = loadExperimentReportDBData(c); err != nil {
		return nil, err
	}

	filterStageExperimentItem(c)
	sortStageExperimentItem(c)
	buildExperimentReportExpIDIndex(c)

	if err = getRtaDataFromAMS(c); err != nil {
		return nil, err
	}

	if err = buildExperimentReportResponse(c); err != nil {
		return nil, err
	}

	return c.response, nil
}

type experimentReportContext struct {
	db          *gorm.DB
	param       *types.ExperimentReportParameter
	expIDIdx    map[int64]int
	account     *types.RtaAccount
	stage       *types.ExperimentStage
	amsResponse *api.RtaExpDataResponse
	response    *types.ExperimentReportResponse
}

func loadExperimentReportDBData(c *experimentReportContext) error {
	db := c.db
	param := c.param

	stage, err := GetExperimentStageById(db, param.ExperimentStageId)
	if err != nil {
		return err
	}
	if err = LoadStageExperimentItem(db, stage); err != nil {
		return err
	}
	for _, item := range stage.ExperimentItem {
		if err := LoadItemRtaExp(db, item); err != nil {
			return err
		}
	}
	c.stage = stage

	var group *types.ExperimentGroup
	group, err = GetExperimentGroupById(db, stage.ExperimentGroupID)
	if err != nil {
		return err
	}
	if err = LoadGroupRtaAccount(db, group); err != nil {
		return err
	}
	c.account = group.RtaAccount
	return nil
}

func getRtaDataFromAMS(c *experimentReportContext) error {
	req, err := buildAMSRtaDataRequest(c)
	if err != nil {
		return err
	}

	resp, err := api.GetRtaExpData(req)
	if err != nil {
		return err
	}

	c.amsResponse = resp
	return nil
}

// 如果item不在请求的列表中，不需要处理
func filterStageExperimentItem(c *experimentReportContext) {
	param := c.param
	stage := c.stage
	reqItemSet := make(map[uint64]struct{})
	reqItemSet[param.BaseExperimentItemId] = struct{}{}
	for _, itemId := range param.LabExperimentItemId {
		reqItemSet[itemId] = struct{}{}
	}
	validItem := make([]*types.ExperimentItem, 0)
	for _, item := range stage.ExperimentItem {
		if _, ok := reqItemSet[item.ID]; ok {
			validItem = append(validItem, item)
		}
	}
	stage.ExperimentItem = validItem
}

func sortStageExperimentItem(c *experimentReportContext) {
	baseItemId := c.param.BaseExperimentItemId
	labOrder := make(map[uint64]int)
	for i, itemId := range c.param.LabExperimentItemId {
		labOrder[itemId] = i + 1
	}
	sort.Slice(c.stage.ExperimentItem, func(i, j int) bool {
		if c.stage.ExperimentItem[i].ID == baseItemId {
			return true
		}
		if c.stage.ExperimentItem[j].ID == baseItemId {
			return false
		}
		return labOrder[c.stage.ExperimentItem[i].ID] < labOrder[c.stage.ExperimentItem[j].ID]
	})
}

func buildExperimentReportExpIDIndex(c *experimentReportContext) {
	expIDIdx := make(map[int64]int)

	for i, item := range c.stage.ExperimentItem {
		for _, exp := range item.RtaExp {
			expId, _ := api.ParseAmsExpId(exp.ID)
			expIDIdx[expId] = i
		}
	}
	c.expIDIdx = expIDIdx
}

func buildAMSRtaDataRequest(c *experimentReportContext) (*api.RtaExpDataRequest, error) {
	param := c.param

	expIDIdx := c.expIDIdx
	expIDs := make([]int64, 0, len(expIDIdx))
	for expID := range expIDIdx {
		expIDs = append(expIDs, expID)
	}

	req := &api.RtaExpDataRequest{
		RtaExpRequestBase: &api.RtaExpRequestBase{
			RtaID: c.account.RtaID,
			Token: c.account.Token,
		},
		Data: &api.RtaExpDataRequestData{
			Granularity: param.Granularity,
			BeginTime:   api.RtaExpDataRequestTime(time.Unix(param.BeginTime, 0)),
			EndTime:     api.RtaExpDataRequestTime(time.Unix(param.EndTime, 0)),
			ExpID:       expIDs,
			TotalFlag:   0,
		},
	}
	if param.Filter != nil {
		req.Data.UID = strings.Join(param.Filter.UID, ",")
		req.Data.AppID = strings.Join(param.Filter.AppID, ",")
		req.Data.CID = strings.Join(param.Filter.CID, ",")
		req.Data.AID = strings.Join(param.Filter.AID, ",")
	}
	return req, nil
}

func buildExperimentReportResponse(c *experimentReportContext) error {
	api.SortRtaExpDataSlice(c.amsResponse.Data)
	buf := make(api.RtaExpDataSlice, 0, 0)
	var groupTime api.ExpDataTime = 0
	for _, rtaData := range c.amsResponse.Data {
		if !c.param.ShouldTimeMerge() && rtaData.Time != groupTime {
			if len(buf) != 0 {
				c.response.Records = append(c.response.Records, mergeByExperimentItem(c, buf)...)
			}
			buf = make(api.RtaExpDataSlice, 0, 0)
			buf = append(buf, rtaData)
			groupTime = rtaData.Time
		} else {
			buf = append(buf, rtaData)
		}
	}
	if len(buf) != 0 {
		c.response.Records = append(c.response.Records, mergeByExperimentItem(c, buf)...)
	}

	return nil
}

func mergeByExperimentItem(
	c *experimentReportContext, rtaData api.RtaExpDataSlice) []*types.ExperimentReportRecord {

	expIDIdx := c.expIDIdx
	ret := constructExperimentReportRecords(c)

	for _, data := range rtaData {
		expId := data.ExpId
		if i, ok := expIDIdx[expId]; !ok {
			continue
		} else {
			mergeToExperimentReportRecord(ret[i], data)
		}
	}

	types.CalculateExperimentReportRecordDelta(ret)

	return ret
}

func constructExperimentReportRecords(c *experimentReportContext) []*types.ExperimentReportRecord {
	ret := make([]*types.ExperimentReportRecord, len(c.stage.ExperimentItem))
	for i, item := range c.stage.ExperimentItem {
		ret[i] = &types.ExperimentReportRecord{
			ExperimentItem: item,
		}
	}
	return ret
}

func mergeToExperimentReportRecord(record *types.ExperimentReportRecord, rtaData *api.ExpData) {
	mergeExperimentReportRtaDataTime(&record.Time, int64(rtaData.Time))
	record.Cost = types.AddExperimentReportMetric(record.Cost, rtaData.Cost)
	record.Exposure = types.AddExperimentReportMetric(record.Exposure, float64(rtaData.Exposure))
	record.Click = types.AddExperimentReportMetric(record.Click, float64(rtaData.Click))
	record.Conversion = types.AddExperimentReportMetric(record.Conversion, float64(rtaData.Conversion))
	record.ConversionSecond = types.AddExperimentReportMetric(
		record.ConversionSecond, float64(rtaData.ConversionSecond))
	record.CalculateDerivedMetric()
}

func mergeExperimentReportRtaDataTime(r *[2]int64, n int64) {
	if r[0] == 0 {
		r[0] = n
	} else if r[0] == n {
		return
	} else {
		r[1] = n
	}
}
