package retrieval

import (
	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/pkg/matchengine"
)

type rtaSearcher struct {
	se             *matchengine.TargetingIndexSearcher
	targetingIndex *targeting.TargetingIndex // 索引数据
	rtaInfo        []*rta.RTAInfo            // 正排数据
	version        int64
	dumpTime       int64
}

func (searcher *rtaSearcher) init(rtaIndex *rta.RTAIndex) {
	searcher.version = rtaIndex.Version
	searcher.dumpTime = rtaIndex.DumpTime
	searcher.targetingIndex = rtaIndex.Index
	searcher.rtaInfo = rtaIndex.RtaTargeting
	searcher.se = matchengine.NewTargetingIndexSearcher(rtaIndex.Index)
}

func (searcher *rtaSearcher) search(
	sc *matchengine.SearchContext, req *retrieval.RetrievalRequest) *retrieval.RetrievalResponse {
	localIds := searcher.se.Search(sc, req)
	hitRtaInfo := make([]*rta.RTAInfo, 0)
	for _, localId := range localIds {
		if searcher.targetingIndex.Targeting[localId].Offline {
			continue
		}
		hitRtaInfo = append(hitRtaInfo, searcher.rtaInfo[localId])
	}
	return &retrieval.RetrievalResponse{
		RetCode:    0,
		HitRtaInfo: hitRtaInfo,
	}
}

// GetLocalId 根据定向ID获取下标
func (searcher *rtaSearcher) GetLocalId(targetingId uint64) (uint32, bool) {
	return searcher.se.GetLocalId(targetingId)
}

// GetTargeting 根据下标获取定向
func (searcher *rtaSearcher) GetTargeting(localId uint32) *targeting.Targeting {
	return searcher.targetingIndex.Targeting[localId]
}

// GetRTAInfo 根据下标获取RTA定向以及对应的策略
func (searcher *rtaSearcher) GetRTAInfo(localId uint32) *rta.RTAInfo {
	return searcher.rtaInfo[localId]
}
