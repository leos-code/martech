package orm

import (
	"testing"
)

func TestExperimentManagement(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	//item := &types.ExperimentItem{
	//	Name: "exp_item1",
	//	RtaExp: []*types.RtaExp{
	//		{
	//			ID: "rta_exp1",
	//		},
	//		{
	//			ID: "rta_exp2",
	//		},
	//	},
	//}
	//
	//item2 := &types.ExperimentItem{
	//	Name: "exp_item2",
	//	RtaExp: []*types.RtaExp{
	//		{
	//			ID: "rta_exp1",
	//		},
	//	},
	//}
	//
	//assert.NoError(t, UpsertExperimentItem(db, item))
	//assert.NoError(t, UpsertExperimentItem(db, item2))
	//
	//item, err := GetExperimentItemById(db, 1)
	//assert.NoError(t, err)
	//
	//rtaExpList, err := GetAllRtaExp(db)
	//
	//tx := db.Begin()
	//
	//for _, rtaExp := range rtaExpList {
	//	assert.NoError(t, unbindExperimentItemAndRtaExp(tx, item, rtaExp))
	//}
	//
	//tx.Commit()
	//
	//assert.NoError(t, LoadItemRtaExp(db, item))
	//bytes, _ := json.Marshal(item)
	//fmt.Println(string(bytes))
	//
	//assert.NoError(t, UnbindStageExpItem(db, []*types.ExperimentItem{item}, []*types.ExperimentItem{item2}))

	//duration, _ := time.ParseDuration("24h")
	//group := &types.ExperimentGroup{
	//	ID: 0,
	//Name: "用来测试update的group第三次修改了名字",
	//ExperimentStage: []*types.ExperimentStage{
	//	{
	//		Status: types.ExperimentStageDraft,
	//		ExperimentItem: []*types.ExperimentItem{
	//			{
	//				Name: "testExpItem1",
	//				RtaExp: []*types.RtaExp{
	//					{
	//						ID: "rta_exp1",
	//					},
	//					{
	//						ID: "rta_exp2",
	//					},
	//				},
	//				ExperimentMetadata: []*types.ExperimentMetadata{
	//					{
	//						Value: "metadata1",
	//					},
	//					{
	//						Value: "metadata2",
	//					},
	//				},
	//			},
	//		},
	//	},
	//{
	//	StartTime: time.Now(),
	//	EndTime:   time.Now().Add(duration),
	//	Status:    types.ExperimentStageRunning,
	//	ExperimentItem: []*types.ExperimentItem{
	//		{
	//			Name: "testExpItem2",
	//			RtaExp: []*types.RtaExp{
	//				{
	//					ID: "rta_exp3",
	//				},
	//				{
	//					ID: "rta_exp4",
	//				},
	//			},
	//			ExperimentMetadata: []*types.ExperimentMetadata{
	//				{
	//					Value: "metadata3",
	//				},
	//				{
	//					Value: "metadata4",
	//				},
	//			},
	//		},
	//	},
	//},
	//},
	//}

	//draft := &types.ExperimentStage{
	//	ID:        1,
	//	Status:    types.ExperimentStageDraft,
	//	StartTime: time.Now(),
	//	EndTime:   time.Now().Add(duration),
	//	ExperimentItem: []*types.ExperimentItem{
	//		{
	//			ID:   1,
	//			Name: "firstItem第三次修改了名字",
	//			RtaExp: []*types.RtaExp{
	//				{
	//					ID: "rta_exp1",
	//				},
	//				{
	//					ID: "rta_exp2",
	//				},
	//			},
	//		},
	//	},
	//}

	//group.Draft = draft
	//
	//assert.NoError(t, EditExperimentGroup(db, group))

	//assert.NoError(t, PromoteExperimentGroup(db, group))
	//assert.NoError(t, StopExperimentGroup(db, group))
}
