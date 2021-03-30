package orm

import (
	"testing"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestExperimentModifyRecord(t *testing.T) {
	db, _ := getTestDB(t)
	assert.NoError(t, Setup(db))

	account := &types.RtaAccount{}
	assert.NoError(t, UpsertRtaAccount(db, account))
	group := &types.ExperimentGroup{
		RtaAccountID: account.ID,
	}
	assert.NoError(t, UpsertExperimentGroup(db, group))

	stage := &types.ExperimentStage{
		StartTime:         time.Now(),
		EndTime:           time.Now(),
		ExperimentGroupID: group.ID,
	}
	assert.NoError(t, UpsertExperimentStage(db, stage))

	item := &types.ExperimentModifyRecord{
		Data: &types.ExperimentGroupWrapper{
			Name: "1",
		},
		User: &types.User{
			PhoneNumber: "1",
		},
		ExperimentGroupID: group.ID,
		ExperimentStageID: stage.ID,
	}

	assert.NoError(t, UpsertExperimentModifyRecord(db, item))

	var err error
	item, err = GetExperimentModifyRecordById(db, item.ID)
	t.Log(item)
	assert.NoError(t, err)

	err = DeleteExperimentModifyRecordById(db, item.ID)
	assert.NoError(t, err)

	item, err = GetExperimentModifyRecordById(db, item.ID)
	assert.Error(t, err)
}
