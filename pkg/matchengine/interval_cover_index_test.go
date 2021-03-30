package matchengine

import (
	"testing"

	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/stretchr/testify/assert"
)

func TestRangeIndex(t *testing.T) {
	rMapping := make(map[Interval]uint32)
	rMapping[Interval{2, 7}] = 0
	rMapping[Interval{3, 6}] = 7
	rMapping[Interval{5, 9}] = 3
	var ri targeting.IntervalCoverIndex
	BuildIntervalCoverIndex(rMapping, &ri)

	{
		tokenIds := SearchIntervalCoverIndex(&ri, 1)
		assert.Len(t, tokenIds, 0)
	}

	{
		tokenIds := SearchIntervalCoverIndex(&ri, 2)
		assert.Len(t, tokenIds, 1)
		assert.EqualValues(t, 0, tokenIds[0])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 3)
		assert.Len(t, tokenIds, 2)
		assert.EqualValues(t, 0, tokenIds[0])
		assert.EqualValues(t, 7, tokenIds[1])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 4)
		assert.Len(t, tokenIds, 2)
		assert.EqualValues(t, 0, tokenIds[0])
		assert.EqualValues(t, 7, tokenIds[1])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 5)
		assert.Len(t, tokenIds, 3)
		assert.EqualValues(t, 0, tokenIds[0])
		assert.EqualValues(t, 3, tokenIds[1])
		assert.EqualValues(t, 7, tokenIds[2])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 6)
		assert.Len(t, tokenIds, 2)
		assert.EqualValues(t, 0, tokenIds[0])
		assert.EqualValues(t, 3, tokenIds[1])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 7)
		assert.Len(t, tokenIds, 1)
		assert.EqualValues(t, 3, tokenIds[0])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 8)
		assert.Len(t, tokenIds, 1)
		assert.EqualValues(t, 3, tokenIds[0])
	}
	{
		tokenIds := SearchIntervalCoverIndex(&ri, 9)
		assert.Len(t, tokenIds, 0)
	}

	{
		tokenIds := SearchIntervalCoverIndex(&ri, 10)
		assert.Len(t, tokenIds, 0)
	}
}

func TestRangeIndex2(t *testing.T) {
	rMapping := make(map[Interval]uint32)
	rMapping[Interval{2, 3}] = 0
	rMapping[Interval{2, 3}] = 1
	rMapping[Interval{3, 4}] = 2
	rMapping[Interval{5, 9}] = 3
	var ri targeting.IntervalCoverIndex
	BuildIntervalCoverIndex(rMapping, &ri)

	t.Log(SearchIntervalCoverIndex(&ri, 1))
	t.Log(SearchIntervalCoverIndex(&ri, 2))
	t.Log(SearchIntervalCoverIndex(&ri, 3))
	t.Log(SearchIntervalCoverIndex(&ri, 4))
	t.Log(SearchIntervalCoverIndex(&ri, 5))
	t.Log(SearchIntervalCoverIndex(&ri, 6))
	t.Log(SearchIntervalCoverIndex(&ri, 7))
	t.Log(SearchIntervalCoverIndex(&ri, 8))
	t.Log(SearchIntervalCoverIndex(&ri, 9))
	t.Log(SearchIntervalCoverIndex(&ri, 10))
}
