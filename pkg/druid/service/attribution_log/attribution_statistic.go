package attribution_log

import (
	"encoding/json"
	"github.com/tencentad/martech/pkg/druid/common"
	"github.com/awatercolorpen/godruid"
)

type QueryBuilder struct {
	*AttributionLogRequest
}

func (qb *QueryBuilder) DataSource() string {
	return qb.DataSourceStr
}

func (qb *QueryBuilder) Intervals() []string {
	return common.BuildIntervals(qb.TimeInterval.BeginTime, qb.TimeInterval.EndTime)
}

func (qb *QueryBuilder) Granularity() godruid.Granlarity {
	return godruid.GranAll
}

func (qb *QueryBuilder) Dimensions() []godruid.DimSpec {
	var ds []interface{}
	for _, d := range qb.DimensionsList {
		ds = append(ds, d)
	}
	return common.BuildDimensions(ds)
}

func (qb *QueryBuilder) Aggregations() []godruid.Aggregation {
	return []godruid.Aggregation{*godruid.AggLongSum("sum_count", "count")}
}

func (qb *QueryBuilder) PostAggregations() []godruid.PostAggregation {
	return nil
}

func (qb *QueryBuilder) Filter() *godruid.Filter {
	return nil
}

func (qb *QueryBuilder) LimitSpec() *godruid.Limit {
	return nil
}

type QueryParser struct{}

func (qp *QueryParser) ParseQuery(v interface{}) interface{} {
	u := &AttributionStatisticItem{}
	w := v.(godruid.GroupbyItem)

	b, _ := json.Marshal(w.Event)
	_ = json.Unmarshal(b, u)

	return u
}
