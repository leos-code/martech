package matchengine

import (
	"sort"

	"github.com/tencentad/martech/api/proto/targeting"
)

type pointVec []*targeting.IntervalCoverIndex_PointInfo

// Len
func (p pointVec) Len() int { return len(p) }

// Less
func (p pointVec) Less(i, j int) bool { return p[i].Point < p[j].Point }

// Swap
func (p pointVec) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type uint32Slice []uint32

// Len
func (p uint32Slice) Len() int { return len(p) }

// Less
func (p uint32Slice) Less(i, j int) bool { return p[i] < p[j] }

// Swap
func (p uint32Slice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// BuildTokenCoverIndex 为所有的token构造范围索引，token包含维度的信息
func BuildTokenCoverIndex(allTokenMap map[string]*tokenMap,
	searchTokenIndex map[string]*targeting.SearchTokenIndex) {
	for field, tokenMap := range allTokenMap {
		var ri targeting.IntervalCoverIndex
		BuildIntervalCoverIndex(tokenMap.intervalMap, &ri)
		searchTokenIndex[field] = &targeting.SearchTokenIndex{
			CoverIndex:  &ri,
			StringIndex: tokenMap.stringMap,
		}
	}
}

// BuildIntervalCoverIndex interval和一个id对应
func BuildIntervalCoverIndex(rangeMapping map[Interval]uint32, ri *targeting.IntervalCoverIndex) {
	existMap := make(map[uint64]bool)
	existMap[0] = true
	existMap[^uint64(0)] = true
	for r := range rangeMapping {
		existMap[r.Begin] = true
		existMap[r.End] = true
	}

	for p := range existMap {
		pp := &targeting.IntervalCoverIndex_PointInfo{}
		pp.Point = p
		ri.Points = append(ri.Points, pp)
	}
	sort.Sort(pointVec(ri.Points))
	pMapIndex := make(map[uint64]int)
	for i, p := range ri.Points {
		pMapIndex[p.Point] = i
	}
	for r, id := range rangeMapping {
		b, e := pMapIndex[r.Begin], pMapIndex[r.End]
		for i := b; i < e; i++ {
			ri.Points[i].Ids = append(ri.Points[i].Ids, id)
		}
	}
	for i := range ri.Points {
		sort.Sort(uint32Slice(ri.Points[i].Ids))
	}
}

// SearchIntervalCoverIndex 求输入一个点，被哪些id（哪些range）覆盖
func SearchIntervalCoverIndex(ri *targeting.IntervalCoverIndex, point uint64) []uint32 {
	low, high := 0, len(ri.Points)
	for low < high {
		mid := (low + high) >> 1
		if ri.Points[mid].Point < point {
			low = mid + 1
		} else {
			high = mid
		}
	}
	// low是第一个 >= point的
	index := low
	if ri.Points[low].Point != point {
		index--
	}
	return ri.Points[index].Ids
}
