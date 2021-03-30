package orm

import (
	"encoding/json"
	"fmt"
	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentItem(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.ExperimentItem{
		Name:               "experimentItem",
		ExperimentMetadata: []*types.ExperimentMetadata{{Value: "metadata1"}, {Value: "metadata2"}},
		RtaExp:             []*types.RtaExp{{RtaAccountID: 1}, {RtaAccountID: 1}},
	}
	assert.NoError(t, UpsertExperimentItem(db, item))

	list, err := GetAllExperimentItem(db)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	b, _ := json.Marshal(list)
	fmt.Println(string(b))

	findItem, err := GetExperimentItemById(db, list[0].ID)
	assert.NoError(t, err)
	fmt.Println(findItem)

	assert.NoError(t, LoadItemExperimentMetadata(db, findItem))
	assert.NoError(t, LoadItemRtaExp(db, findItem))

	b, _ = json.Marshal(findItem)
	fmt.Println(string(b))

	assert.NoError(t, DeleteExperimentItemById(db, list[0].ID))

	list, err = GetAllExperimentItem(db)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestLock(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	rtaExp := &types.RtaExp{ID: "rtaExp1", RtaAccountID: 1}
	rtaExp2 := &types.RtaExp{ID: "rtaExp2", RtaAccountID: 1}

	expItem := &types.ExperimentItem{
		Name:              "expItem1",
		ExperimentStageID: 0,
		RtaExp:            []*types.RtaExp{rtaExp},
	}

	expItem2 := &types.ExperimentItem{
		Name:              "expItem2",
		ExperimentStageID: 0,
	}

	assert.NoError(t, UpsertExperimentItem(db, expItem))
	assert.NoError(t, UpsertExperimentItem(db, expItem2))

	tx := db.Begin()
	expItem, err := GetExperimentItemById(db, 1)
	assert.NoError(t, err)

	rtaExp, err = GetRtaExpById(db, "rtaExp1")
	assert.NoError(t, err)

	assert.NoError(t, tx.Set("gorm:query_option", "for update").Error)

	if err = tx.Model(expItem).Association("RtaExp").Delete(rtaExp); err != nil {
		tx.Rollback()
		assert.Error(t, err)
	}
	if err = tx.Model(expItem).Association("RtaExp").Append(rtaExp2); err != nil {
		tx.Rollback()
		assert.NoError(t, err)
	}

	expItem2, err = GetExperimentItemById(db, 2)
	assert.NoError(t, err)
	if err := tx.Set("gorm:query_option", "for update").Model(expItem2).Association("RtaExp").Append(rtaExp); err != nil {
		tx.Rollback()
		assert.Error(t, err)
	}

	tx.Commit()

	itemList, err2 := GetAllExperimentItem(db)
	assert.NoError(t, err2)
	for _, item := range itemList {
		assert.NoError(t, LoadItemRtaExp(db, item))
		b1, err := json.Marshal(item)
		assert.NoError(t, err)
		fmt.Println(string(b1))
	}
}
