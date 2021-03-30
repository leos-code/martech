package orm

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestExperimentParameter(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.ExperimentParameter{
		Name:               "experimentParameter",
		ExperimentMetadata: []*types.ExperimentMetadata{{Value: "metadata1"}, {Value: "metadata2"}},
	}
	assert.NoError(t, UpsertExperimentParameter(db, item))

	list, err := GetAllExperimentParameter(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	findParameter, err := GetExperimentParameterById(db, list[0].ID)
	assert.NoError(t, err)

	assert.NoError(t, LoadParameterExperimentMetadata(db, findParameter))
	b, _ := json.Marshal(findParameter)
	fmt.Println(string(b))

	assert.NoError(t, DeleteExperimentParameterById(db, list[0].ID))

	list, err = GetAllExperimentParameter(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}
