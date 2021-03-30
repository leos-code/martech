package matchengine

import (
	"container/heap"

	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/willf/bitset"
)

type postingListHeap []*conjunctionPostingList
type conjunctionIndex = targeting.ConjunctionIndex

// Len
func (h *postingListHeap) Len() int {
	return len(*h)
}

// Less 优先判断ConjunctionID, 对于相同的ConjunctionID, PredicateID更大的排在前面，做Not过滤
func (h *postingListHeap) Less(i, j int) bool {
	cid1, cid2 := (*h)[i].CurEntry().ConjunctionId, (*h)[j].CurEntry().ConjunctionId
	if cid1 != cid2 {
		return cid1 < cid2
	}
	pid1, pid2 := (*h)[i].CurEntry().PredicateId, (*h)[j].CurEntry().PredicateId
	return pid1 > pid2
}

// Swap
func (h *postingListHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Pop
func (h *postingListHeap) Pop() interface{} {
	v := (*h)[h.Len()-1]
	*h = (*h)[0 : h.Len()-1]
	return v
}

// Push
func (h *postingListHeap) Push(v interface{}) {
	*h = append(*h, v.(*conjunctionPostingList))
}

type conjunctionEvaluator struct {
	conjIndex []*conjunctionIndex
}

// NewConjunctionEvaluator 构造函数
func NewConjunctionEvaluator(conIndex []*conjunctionIndex) *conjunctionEvaluator {
	ret := &conjunctionEvaluator{
		conjIndex: conIndex,
	}
	return ret
}

// Evaluate 计算满足条件的ConjunctionID
func (ce *conjunctionEvaluator) Evaluate(sc *SearchContext) []uint32 {
	h := postingListHeap(sc.postingListVec)
	heap.Init(&h)

	maxSize := ce.conjIndex[len(ce.conjIndex)-1].HitCount
	predicateBitSet := bitset.New(uint(maxSize))
	topkPlist := make([]*conjunctionPostingList, 0, maxSize)

	for h.Len() > 0 {
		pl := h[0]
		conID := pl.CurEntry().ConjunctionId
		curSize := ce.conjIndex[conID].HitCount
		if h.Len() < int(curSize) {
			return sc.conjResult
		}
		topkPlist = topkPlist[:0]
		for i := 0; i < int(curSize)-1; i++ {
			pl = heap.Pop(&h).(*conjunctionPostingList)
			topkPlist = append(topkPlist, pl)
		}
		if h[0].CurEntry().ConjunctionId == conID {
			for h.Len() > 0 && h[0].CurEntry().ConjunctionId == conID {
				pl = heap.Pop(&h).(*conjunctionPostingList)
				topkPlist = append(topkPlist, pl)
			}
			// 没有命中not， 开始check
			if topkPlist[0].CurEntry().PredicateId < curSize {
				count := 0
				predicateBitSet.ClearAll()
				for _, pl := range topkPlist {
					pid := pl.CurEntry().PredicateId
					if !predicateBitSet.Test(uint(pid)) {
						count++
						predicateBitSet.Set(uint(pid))
					}
				}
				if count == int(curSize) {
					sc.conjResult = append(sc.conjResult, conID)
				}
			}
			for _, pl := range topkPlist {
				pl.Next()
				if pl.IsValid() {
					heap.Push(&h, pl)
				}
			}
		} else {
			nextID := h[0].CurEntry().ConjunctionId
			for _, pl := range topkPlist {
				if pl.CurEntry().ConjunctionId != nextID {
					pl.SeekTo(nextID)
				}
				if pl.IsValid() {
					heap.Push(&h, pl)
				}
			}
		}
	}
	return sc.conjResult
}
