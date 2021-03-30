package retrieval

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
)

type incUpdater struct {
	offlineCache   map[uint64]int64
	lastUpdateTime int64
}

func newIncUpdater() *incUpdater {
	return &incUpdater{
		offlineCache:   make(map[uint64]int64),
		lastUpdateTime: 0,
	}
}

func (up *incUpdater) updateFromDB(searcher *rtaSearcher) error {
	if searcher.dumpTime == 0 {
		return nil
	}

	if up.lastUpdateTime == 0 {
		up.lastUpdateTime = searcher.dumpTime
	}

	db := orm.GetDB()

	dbT := make([]*types.Targeting, 0)
	if err := db.Select("id, updated_at").
		Where("unix_timestamp(updated_at) >= ? and status != ?", up.lastUpdateTime, types.TargetingStatusValid).
		Order("unix_timestamp(updated_at)").
		Find(&dbT).Error; err != nil {

		return err
	}

	for _, record := range dbT {
		up.updateRecord(searcher, record)
	}

	if len(dbT) != 0 {
		latestUpdateTime := dbT[len(dbT)-1].UpdatedAt.Unix()
		if latestUpdateTime > up.lastUpdateTime {
			up.lastUpdateTime = latestUpdateTime
		} else if latestUpdateTime == up.lastUpdateTime {
			up.lastUpdateTime++
		}
	}

	return nil
}

func (up *incUpdater) updateRecord(searcher *rtaSearcher, dbT *types.Targeting) {
	upTime := dbT.UpdatedAt.Unix()
	if upTime < searcher.dumpTime {
		return
	}

	tid := dbT.ID
	localId, ok := searcher.GetLocalId(tid)
	if ok {
		targeting := searcher.GetTargeting(localId)
		targeting.UpdateTime = upTime
		targeting.Offline = true
	}
	up.offlineCache[tid] = upTime
}

func (up *incUpdater) updateFromCache(searcher *rtaSearcher) {
	for tid, upTime := range up.offlineCache {
		if upTime < searcher.dumpTime {
			delete(up.offlineCache, tid)
		}
		localId, ok := searcher.GetLocalId(tid)
		if ok {
			targeting := searcher.GetTargeting(localId)
			if targeting.UpdateTime < upTime {
				targeting.UpdateTime = upTime
				targeting.Offline = true
			} else if searcher.targetingIndex.Targeting[localId].UpdateTime > upTime {
				delete(up.offlineCache, tid)
			}
		}

	}
}
