package orm

import (
	"fmt"
	"testing"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRtaExp(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.RtaExp{
		ID:             "id",
		ExpirationTime: time.Now(),
	}
	assert.NoError(t, InsertRtaExp(db, item))

	list, err := GetAllRtaExp(db)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	findMetadata, err := GetRtaExpById(db, list[0].ID)
	assert.NoError(t, err)
	fmt.Println(findMetadata)

	assert.NoError(t, DeleteRtaExpById(db, list[0].ID))

	list, err = GetAllRtaExp(db)
	assert.NoError(t, err)
	assert.Len(t, list, 0)

	assert.NoError(t, InsertRtaExp(db, &types.RtaExp{
		ID: "2",
	}))
}

func TestGetRtaServerExperimentConfig(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)
	experimentItem := &types.ExperimentItem{
		ExperimentMetadata: []*types.ExperimentMetadata{
			{
				Value: "1",
			},
		},
	}
	assert.NoError(t, UpsertExperimentItem(db, experimentItem))

	rtaExp := &types.RtaExp{
		BindExperimentItemID: 1,
	}
	assert.NoError(t, InsertRtaExp(db, rtaExp))

	rtaExps, err := GetRtaServerExperimentConfig(db)
	assert.NoError(t, err)
	assert.Len(t, rtaExps, 1)
	assert.Len(t, rtaExps[0].BindExperimentItem.ExperimentMetadata, 1)
	assert.EqualValues(t, "1", rtaExps[0].BindExperimentItem.ExperimentMetadata[0].Value)
}

func TestUniqueRTAID(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	account := &types.RtaAccount{
		RtaID: "123",
	}

	assert.NoError(t, UpsertRtaAccount(db, account))

	assert.NoError(t, InsertRtaExp(db, &types.RtaExp{
		ID:             "2",
		RtaAccountID:   account.ID,
		ExpirationTime: time.Now(),
	}))

	assert.NoError(t, InsertRtaExp(db, &types.RtaExp{
		ID:             "3",
		RtaAccountID:   account.ID,
		ExpirationTime: time.Now(),
	}))

}
