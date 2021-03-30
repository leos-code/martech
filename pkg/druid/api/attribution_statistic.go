package api

import (
	"github.com/tencentad/martech/pkg/druid"
	"github.com/tencentad/martech/pkg/druid/service/attribution_log"
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/godruid"
	"time"
)

type adIdPlatformPair struct {
	adId     string
	platform string
}

func GetAttributionStaticData(client *godruid.Client) ([]*attribution_log.AttributionStatisticResult, error) {
	clickRequest := &attribution_log.AttributionLogRequest{
		DataSourceStr:  "click_sink_test",
		TimeInterval:   &attribution_log.TimeInterval{
			BeginTime: time.Date(2021, 2, 20, 0,0,0,0,time.Local),
			EndTime:   time.Now(),
		},
		DimensionsList: []string{"adId", "platform"},
	}
	builder := &attribution_log.QueryBuilder{AttributionLogRequest: clickRequest}
	clickQuery := druid.NewBuilder(builder).Build().(godruid.Query)
	if err := client.Query(clickQuery, ""); err != nil {
		return nil, err
	}
	parser := druid.NewParser(&attribution_log.QueryParser{})
	clickResult := parser.Parse(clickQuery)
	clickMap := setMap(clickResult)

	conversionRequest := &attribution_log.AttributionLogRequest{
		DataSourceStr:  "conversion_sink_test",
		TimeInterval:   &attribution_log.TimeInterval{
			BeginTime: time.Date(2021, 2, 20, 0,0,0,0,time.Local),
			EndTime:   time.Now(),
		},
		DimensionsList: []string{"adId", "platform"},
	}
	builder = &attribution_log.QueryBuilder{AttributionLogRequest: conversionRequest}
	conversionQuery := druid.NewBuilder(builder).Build().(godruid.Query)
	if err := client.Query(conversionQuery, ""); err != nil {
		return nil, err
	}
	conversionResult := parser.Parse(conversionQuery)
	conversionMap := setMap(conversionResult)

	result := make([]*attribution_log.AttributionStatisticResult, 0)
	for pair, clickCount := range clickMap {
		conversionCount, found := conversionMap[pair]
		if found {
			attributionLog := &attribution_log.AttributionStatisticResult{
				AdId:            pair.adId,
				Platform:        pair.platform,
				ClickCount:      int(clickCount),
				ConversionCount: int(conversionCount),
				ConversionRate:  conversionCount / clickCount,
			}
			result = append(result, attributionLog)
		}
	}

	return result, nil
}

func setMap(v interface{}) map[adIdPlatformPair]float64 {
	var d []*attribution_log.AttributionStatisticItem
	linq.From(v).ToSlice(&d)
	resultMap := make(map[adIdPlatformPair]float64)
	for _, item := range d {
		pair := adIdPlatformPair{
			adId:     item.AdId,
			platform: item.Platform,
		}
		resultMap[pair] = item.SumCount
	}
	return resultMap
}
