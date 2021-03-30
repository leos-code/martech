package orm

import (
	"fmt"
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestExperimentMetadata(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.ExperimentMetadata{
		Value: "experimentMetadata",
	}
	assert.NoError(t, UpsertExperimentMetadata(db, item))

	list, err := GetAllExperimentMetadata(db)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	findMetadata, err := GetExperimentMetadataById(db, list[0].ID)
	assert.NoError(t, err)
	fmt.Println(findMetadata)

	assert.NoError(t, DeleteExperimentMetadataById(db, list[0].ID))

	list, err = GetAllExperimentMetadata(db)
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestLoadExperimentMetadataExperimentParameter(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	parameter := &types.ExperimentParameter{
		Name: "a",
	}
	assert.NoError(t, UpsertExperimentParameter(db, parameter))

	metadata := &types.ExperimentMetadata{
		Value:                 "experimentMetadata",
		ExperimentParameterID: 1,
	}
	assert.NoError(t, UpsertExperimentMetadata(db, metadata))

	dbMetadata, err := GetExperimentMetadataById(db, 1)
	assert.NoError(t, err)

	assert.NoError(t, LoadExperimentMetadataExperimentParameter(db, dbMetadata))
	assert.EqualValues(t, "a", dbMetadata.ExperimentParameter.Name)
}
