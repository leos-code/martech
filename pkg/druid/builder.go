package druid

import "github.com/awatercolorpen/godruid"

type queryBuilder interface {
	DataSource() string
	Intervals() []string
	Granularity() godruid.Granlarity
	Dimensions() []godruid.DimSpec
	Aggregations() []godruid.Aggregation
	PostAggregations() []godruid.PostAggregation
	Filter() *godruid.Filter
	LimitSpec() *godruid.Limit
}

type builder struct {
	qb queryBuilder
}

func (b *builder) Build() interface{} {
	return &godruid.QueryGroupBy{
		DataSource:       b.qb.DataSource(),
		Dimensions:       b.qb.Dimensions(),
		Granularity:      b.qb.Granularity(),
		LimitSpec:        b.qb.LimitSpec(),
		Filter:           b.qb.Filter(),
		Aggregations:     b.qb.Aggregations(),
		PostAggregations: b.qb.PostAggregations(),
		Intervals:        b.qb.Intervals(),
		Context: map[string]interface{}{
			"timeout": 180000,
		},
	}
}

func NewBuilder(qb queryBuilder) *builder {
	return &builder{qb: qb}
}
