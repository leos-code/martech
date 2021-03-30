package api

import (
	"encoding/json"
	"fmt"
	"github.com/tencentad/martech/pkg/druid"
	"github.com/awatercolorpen/godruid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Skip()
	clickQuery := &godruid.QueryGroupBy{
		DataSource:   "click_sink_test",
		Intervals:    []string{"2021-02-24/2021-02-27"},
		Granularity:  godruid.GranAll,
		Dimensions:   []godruid.DimSpec{"adId", "platform"},
		Aggregations: []godruid.Aggregation{*godruid.AggLongSum("sum_count", "count")},
	}
	client := druid.GetDruidClient()
	assert.NoError(t, client.Query(clickQuery, ""))
	fmt.Println(client.LastRequest)
	fmt.Println(client.LastResponse)
	fmt.Println(clickQuery.QueryResult)
}

func TestGetAllData(t *testing.T) {
	t.Skip()
	client := druid.GetDruidClient()
	result, err := GetAttributionStaticData(client)
	assert.NoError(t, err)
	assert.NoError(t, err)
	bytes, err := json.Marshal(result)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}
