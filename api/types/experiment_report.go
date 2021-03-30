package types

import (
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
)

// ExperimentReportParameter 实验报表请求参数
type ExperimentReportParameter struct {
	ExperimentStageId    uint64                           `json:"experiment_stage_id"`
	BaseExperimentItemId uint64                           `json:"base_experiment_item_id"`
	LabExperimentItemId  []uint64                         `json:"lab_experiment_item_id"`
	BeginTime            int64                            `json:"begin_time"`
	EndTime              int64                            `json:"end_time"`
	Granularity          api.ExpDataGranularity           `json:"granularity"`
	TimeMergeType        TimeMergeType                    `json:"time_merge_type"`
	Filter               *ExperimentReportFilterParameter `json:"filter"`
}

type TimeMergeType int

const (
	TimeMergeNo  TimeMergeType = 0
	TimeMergeYes TimeMergeType = 1
)

// ShouldTimeMerge 需要把时间合并到一起
func (p *ExperimentReportParameter) ShouldTimeMerge() bool {
	return p.TimeMergeType == TimeMergeYes
}

// ExperimentReportFilterParameter 实验报表过滤条件
type ExperimentReportFilterParameter struct {
	UID   []string `json:"uid"`    // 广告主ID
	AppID []string `json:"app_id"` // App ID
	CID   []string `json:"cid"`    // 推广计划ID
	AID   []string `json:"aid"`    // 广告ID
}

// ExperimentReportResponse 实验报表响应结果
type ExperimentReportResponse struct {
	Records []*ExperimentReportRecord `json:"records"`
}

// ExperimentReportRecord 实验报表响应记录
type ExperimentReportRecord struct {
	Time             [2]int64                `json:"time"`              // 时间戳，如果两个元素都不为0，表示区间
	ExperimentItem   *ExperimentItem         `json:"experiment_item"`   // 实验
	Cost             *ExperimentReportMetric `json:"cost"`              // 消耗
	Exposure         *ExperimentReportMetric `json:"exposure"`          // 曝光
	Click            *ExperimentReportMetric `json:"click"`             // 点击
	Cpm              *ExperimentReportMetric `json:"cpm"`               // cpm
	Cpc              *ExperimentReportMetric `json:"cpc"`               // cpc
	Ctr              *ExperimentReportMetric `json:"ctr"`               // ctr
	Conversion       *ExperimentReportMetric `json:"conversion"`        // 浅层转化量
	ConversionSecond *ExperimentReportMetric `json:"conversion_second"` // 深层转化量（第二目标）
	Cvr              *ExperimentReportMetric `json:"cvr"`               // 浅层cvr
	CvrSecond        *ExperimentReportMetric `json:"cvr_second"`        // 深层cvr（第二目标）
}

// CalculateExperimentReportRecordDelta 计算实验组相对于对照组的delta值
func CalculateExperimentReportRecordDelta(records []*ExperimentReportRecord) {
	if len(records) <= 1 {
		return
	}

	base := records[0]
	for _, lab := range records[1:] {
		calculateExperimentReportRecordDelta(base, lab)
	}
}

func calculateExperimentReportRecordDelta(base, lab *ExperimentReportRecord) {
	calcExperimentReportMetricDelta(base.Cost, lab.Cost)
	calcExperimentReportMetricDelta(base.Exposure, lab.Exposure)
	calcExperimentReportMetricDelta(base.Click, lab.Click)
	calcExperimentReportMetricDelta(base.Cpm, lab.Cpm)
	calcExperimentReportMetricDelta(base.Cpc, lab.Cpc)
	calcExperimentReportMetricDelta(base.Ctr, lab.Ctr)
	calcExperimentReportMetricDelta(base.Conversion, lab.Conversion)
	calcExperimentReportMetricDelta(base.ConversionSecond, lab.ConversionSecond)
	calcExperimentReportMetricDelta(base.Cvr, lab.Cvr)
	calcExperimentReportMetricDelta(base.CvrSecond, lab.CvrSecond)
}

// CalculateDerivedMetric 计算派生指标
func (r *ExperimentReportRecord) CalculateDerivedMetric() {
	var cpm float64 = 0
	var ctr float64 = 0
	if r.Exposure.Value != 0 {
		cpm = r.Cost.Value / r.Exposure.Value * 1000
		ctr = r.Click.Value / r.Exposure.Value
	}

	r.Cpm = NewExperimentReportMetric(cpm)
	r.Ctr = NewExperimentReportMetric(ctr)

	var cpc float64 = 0
	var cvr float64 = 0
	var cvrSecond float64 = 0
	if r.Click.Value != 0 {
		cpc = r.Cost.Value / r.Click.Value
		cvr = r.Conversion.Value / r.Click.Value
		cvrSecond = r.ConversionSecond.Value / r.Click.Value
	}
	r.Cpc = NewExperimentReportMetric(cpc)
	r.Cvr = NewExperimentReportMetric(cvr)
	r.CvrSecond = NewExperimentReportMetric(cvrSecond)
}

// ExperimentReportMetric 指标数据
type ExperimentReportMetric struct {
	Value float64  `json:"value"`
	Delta *float64 `json:"delta,omitempty"`
}

// NewExperimentReportMetric 创建实验报表新的指标
func NewExperimentReportMetric(value float64) *ExperimentReportMetric {
	return &ExperimentReportMetric{
		Value: value,
	}
}

// AddExperimentReportMetric 加法
func AddExperimentReportMetric(metric *ExperimentReportMetric, value float64) *ExperimentReportMetric {
	if metric == nil {
		return NewExperimentReportMetric(value)
	} else {
		metric.Value += value
		return metric
	}
}

var zeroDelta float64 = 0

func calcExperimentReportMetricDelta(base, lab *ExperimentReportMetric) {
	if base == nil || lab == nil {
		return
	}

	if base.Value != 0 {
		delta := (lab.Value - base.Value) / base.Value
		lab.Delta = &delta
	} else {
		lab.Delta = &zeroDelta
	}
}
