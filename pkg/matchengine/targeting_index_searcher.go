package matchengine

import (
	"sort"

	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/glog"
)

// TargetingIndexSearcher 定向检索器
type TargetingIndexSearcher struct {
	targetingIndex *targeting.TargetingIndex

	intervalVec      []*targeting.Interval
	conjIndexVec     []*targeting.ConjunctionIndex
	tokenIndexVec    []*targeting.TokenIndex
	idIndexVec       []*targeting.IdIndex
	searchTokenIndex map[string]*targeting.SearchTokenIndex
	conjEvaluator    *conjunctionEvaluator
}

// NewTargetingIndexSearcher
func NewTargetingIndexSearcher(index *targeting.TargetingIndex) *TargetingIndexSearcher {
	return &TargetingIndexSearcher{
		targetingIndex:   index,
		conjIndexVec:     index.ConjunctionIndex,
		tokenIndexVec:    index.TokenIndex,
		idIndexVec:       index.IdIndex,
		searchTokenIndex: index.SearchTokenIndex,
		conjEvaluator:    NewConjunctionEvaluator(index.ConjunctionIndex),
	}
}

// SearchContext 当次检索生效的数据
type SearchContext struct {
	tokenResMap     map[string]map[uint32]struct{} // feature field => tokenId => {}
	targetingResMap map[uint32]bool
	localIdBuf      []uint32
	postingListVec  []*conjunctionPostingList
	conjResult      []uint32
}

// NewSearchContext
func NewSearchContext() *SearchContext {
	return &SearchContext{
		tokenResMap:     make(map[string]map[uint32]struct{}),
		targetingResMap: make(map[uint32]bool),
		localIdBuf:      make([]uint32, 0, 2048),
		postingListVec:  make([]*conjunctionPostingList, 0, 4096),
		conjResult:      make([]uint32, 0, 1024),
	}
}

// Clear 清空当次检索的数据，复用结构
func (sc *SearchContext) Clear() {
	sc.tokenResMap = make(map[string]map[uint32]struct{})
	sc.targetingResMap = make(map[uint32]bool)

	sc.localIdBuf = sc.localIdBuf[:0]
	sc.postingListVec = sc.postingListVec[:0]
	sc.conjResult = sc.conjResult[:0]
}

// Search 检索符合条件的定向
func (se *TargetingIndexSearcher) Search(sc *SearchContext, req *retrieval.RetrievalRequest) []uint32 {
	for _, feature := range req.Feature {
		if index, ok := se.searchTokenIndex[feature.Field]; ok {
			tokenMap := make(map[uint32]struct{})
			for _, val := range feature.Value {
				if val.Type == retrieval.Feature_Value_ID {
					for _, tokenId := range SearchIntervalCoverIndex(index.CoverIndex, val.Id) {
						tokenMap[tokenId] = struct{}{}
					}
				} else if val.Type == retrieval.Feature_Value_String {
					if tokenId, ok := index.StringIndex[val.Str]; ok {
						tokenMap[tokenId] = struct{}{}
					}
				}
			}
			if len(tokenMap) != 0 {
				sc.tokenResMap[feature.Field] = tokenMap
			}
		}
	}

	for _, token := range sc.tokenResMap {
		hitVec := make([][]*conjunctionHit, 0, len(token))
		for tokenId := range token {
			hitVec = append(hitVec, se.tokenIndexVec[tokenId].ConjunctionHit)
		}
		pl := NewConjunctionPostingList(hitVec)
		sc.postingListVec = append(sc.postingListVec, pl)
	}
	conjList := se.conjEvaluator.Evaluate(sc)

	for _, conjId := range conjList {
		localIds := se.conjIndexVec[conjId].LocalIds
		for _, localId := range localIds {
			sc.targetingResMap[localId] = true
		}
	}
	for id := range sc.targetingResMap {
		sc.localIdBuf = append(sc.localIdBuf, id)
	}
	return sc.localIdBuf
}

// GetLocalId 根据定向的ID标识，找到检索使用的localID
func (se *TargetingIndexSearcher) GetLocalId(targetingId uint64) (uint32, bool) {
	length := len(se.idIndexVec)
	f := func(i int) bool {
		return se.idIndexVec[i].TargetingId >= targetingId
	}
	pos := sort.Search(length, f)
	if pos < length && se.idIndexVec[pos].TargetingId == targetingId {
		return se.idIndexVec[pos].LocalId, true
	}
	return 0, false
}

// GetLocalIds
func (se *TargetingIndexSearcher) GetLocalIds(targetingIds []uint64) []uint32 {
	var result []uint32
	for _, id := range targetingIds {
		if localId, exist := se.GetLocalId(id); exist {
			result = append(result, localId)
		} else {
			glog.Errorf("aid[%d] not exist", id)
		}
	}
	return result
}

// GetTargeting 根据localID，找到对应的定向
func (se *TargetingIndexSearcher) GetTargeting(localID uint32) *targeting.Targeting {
	return se.targetingIndex.Targeting[localID]
}

// GetTargetingID
func (se *TargetingIndexSearcher) GetTargetingID(localID uint32) uint64 {
	return se.targetingIndex.IdIndex[localID].TargetingId
}

func (se *TargetingIndexSearcher) IdIndexVec() []*targeting.IdIndex {
	return se.idIndexVec
}
