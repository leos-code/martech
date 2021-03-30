package attribution_log

import (
	"time"
)

type AttributionStatisticResult struct {
	AdId            string     `json:"ad_id"`
	Platform        string     `json:"platform"`
	ClickCount      int     `json:"click_count"`
	ConversionCount int     `json:"conversion_count"`
	ConversionRate  float64 `json:"conversion_rate"`
}

type AttributionLogRequest struct {
	DataSourceStr   string
	TimeInterval *TimeInterval
	DimensionsList   []string
}

type TimeInterval struct {
	BeginTime time.Time
	EndTime   time.Time
}

type AttributionStatisticItem struct {
	Platform string `json:"platform"`
	AdId string `json:"adId"`
	SumCount float64 `json:"sum_count"`
}
