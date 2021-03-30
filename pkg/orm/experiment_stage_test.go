package orm

import (
	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExperimentStage(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	stage := &types.ExperimentStage{
		Version:        1,
		ExperimentItem: []*types.ExperimentItem{{ID: 1}, {ID: 2}},
	}
	stage2 := &types.ExperimentStage{
		Version:        2,
		ExperimentItem: []*types.ExperimentItem{{ID: 1}, {ID: 2}},
	}
	assert.NoError(t, UpsertExperimentStage(db, stage))
	assert.NoError(t, UpsertExperimentStage(db, stage2))

	list, err := GetAllExperimentStage(db)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	item, err := GetAllExperimentItem(db)
	assert.NoError(t, err)
	assert.Len(t, item, 2)
}
