package orm

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/tencentad/martech/api/types"
)

func TestExperimentGroup(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	expGroup := &types.ExperimentGroup{
		Name: "experiment_group",
		ExperimentStage: []*types.ExperimentStage{
			{
				ExperimentItem: []*types.ExperimentItem{
					{
						Name: "item1",
						ExperimentMetadata: []*types.ExperimentMetadata{
							{Value: "100"}, {Value: "200"},
						},
						RtaExp: []*types.RtaExp{
							{ID: "ams::1"}, {ID: "ams::2"},
						},
					},
					{
						Name: "item2",
					}},
			},
			{
				ExperimentItem: []*types.ExperimentItem{{Name: "item3"}, {Name: "item4"}},
			},
		},
		User: []*types.User{{PhoneNumber: "1"}, {PhoneNumber: "2"}},
		RtaAccount: &types.RtaAccount{
			RtaExp: []*types.RtaExp{
				{ID: "ams::1"}, {ID: "ams::2"},
			},
		},
	}
	assert.NoError(t, UpsertExperimentGroup(db, expGroup))

	list, err := GetAllExperimentGroup(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	foundExpGroup, err := GetExperimentGroupById(db, list[0].ID)
	assert.NoError(t, err)

	assert.NoError(t, db.
		Preload("User").
		Preload("ExperimentStage.ExperimentItem.ExperimentMetadata").
		Preload("ExperimentStage.ExperimentItem.RtaExp").
		Find(foundExpGroup).Error)

	strByte, _ := json.Marshal(foundExpGroup)
	fmt.Println(string(strByte))

	assert.NoError(t, DeleteExperimentGroupById(db, list[0].ID))

	list, err = GetAllExperimentGroup(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestUpdate(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	group := &types.ExperimentGroup{
		ID:   1,
		Name: "修改过的group2",
		ExperimentStage: []*types.ExperimentStage{
			{
				ID:     1,
				Status: types.ExperimentStageRunning,
				ExperimentItem: []*types.ExperimentItem{
					{
						ID:   1,
						Name: "修改过的firstItem",
					},
					{
						Name: "新增加的secondItem",
					},
				},
			},
		},
	}
	// 这里测试结果显示能够插入新的数据并且获得新插入的id
	assert.NoError(t, db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(group).Error)
	bytes, _ := json.Marshal(group)
	fmt.Println(string(bytes))
}

func TestLoadGroupAccount(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	group := &types.ExperimentGroup{
		ID:           1,
		Name:         "修改过的group2",
		RtaAccountID: 1,
	}

	account := &types.RtaAccount{
		RtaID: "rr",
	}
	assert.NoError(t, UpsertRtaAccount(db, account))
	assert.NoError(t, LoadGroupRtaAccount(db, group))
	assert.EqualValues(t, "rr", group.RtaAccount.RtaID)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Format("2006/01/02 Monday 03:04:05"))
}
