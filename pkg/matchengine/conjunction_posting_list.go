package matchengine

import (
	"container/heap"

	"github.com/tencentad/martech/api/proto/targeting"
)

type conjunctionHit = targeting.TokenIndex_ConjunctionHit

// conjunctionPostingList 将多个同一维度的命中的多个倒排合并到一起
type conjunctionPostingList struct {
	heap innerHeap
}

// NewConjunctionPostingList 构造函数
func NewConjunctionPostingList(hitVec [][]*conjunctionHit) *conjunctionPostingList {
	heapItem := make([]*innerPostingList, 0, len(hitVec))
	for _, hit := range hitVec {
		heapItem = append(heapItem, newInnerPostingList(hit))
	}
	return &conjunctionPostingList{
		heap: buildInnerHeap(heapItem),
	}
}

// isValid 是否有效
func (cList *conjunctionPostingList) IsValid() bool {
	return cList.heap[0].isValid()
}

// CurEntry 当前的倒排元素
func (cList *conjunctionPostingList) CurEntry() *conjunctionHit {
	return cList.heap[0].curEntry()
}

// Next
func (cList *conjunctionPostingList) Next() {
	cList.heap[0].next()
	heap.Fix(&cList.heap, 0)
}

// SeekTo 跳到conjunctionID所在的位置，假如不存在，跳到第一个比ConjunctionID
// 大的位置
func (cList *conjunctionPostingList) SeekTo(conID uint32) {
	for _, inner := range cList.heap {
		inner.seekTo(conID)
	}
	heap.Init(&cList.heap)
}

type innerHeap []*innerPostingList

func buildInnerHeap(list []*innerPostingList) innerHeap {
	h := innerHeap(list)
	heap.Init(&h)
	return h
}

// Len
func (h *innerHeap) Len() int {
	return len(*h)
}

// Less 无效的postingList不能放在堆顶
func (h *innerHeap) Less(i, j int) bool {
	a := (*h)[i]
	b := (*h)[j]
	if !a.isValid() {
		return false
	}
	if !b.isValid() {
		return true
	}

	cid1, cid2 := a.curEntry().ConjunctionId, b.curEntry().ConjunctionId
	if cid1 != cid2 {
		return cid1 < cid2
	}
	pid1, pid2 := a.curEntry().PredicateId, b.curEntry().PredicateId
	return pid1 > pid2
}

// Swap
func (h *innerHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

// Pop
func (h *innerHeap) Pop() interface{} {
	v := (*h)[h.Len()-1]
	*h = (*h)[0 : h.Len()-1]
	return v
}

// Push
func (h *innerHeap) Push(v interface{}) {
	pl := v.(*innerPostingList)
	if pl.isValid() {
		*h = append(*h, v.(*innerPostingList))
	}
}

type innerPostingList struct {
	list     []*conjunctionHit
	curIndex int
	length   int
}

func newInnerPostingList(hit []*conjunctionHit) *innerPostingList {
	return &innerPostingList{
		list:     hit,
		curIndex: 0,
		length:   len(hit),
	}
}

func (cList *innerPostingList) isValid() bool {
	return cList.curIndex < cList.length
}

func (cList *innerPostingList) curEntry() *conjunctionHit {
	return cList.list[cList.curIndex]
}

func (cList *innerPostingList) next() {
	if cList.curIndex < cList.length {
		cList.curIndex++
	}
}

func (cList *innerPostingList) seekTo(conID uint32) {
	low, high := cList.curIndex, cList.length
	for low < high {
		mid := (low + high) >> 1
		if cList.list[mid].ConjunctionId < conID {
			low = mid + 1
		} else {
			high = mid
		}
	}
	cList.curIndex = low
}
