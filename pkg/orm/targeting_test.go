package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestTargeting(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	targetingInfo := &types.TargetingInfo{
		Name: "name1",
		Values: &types.TargetingValue{
			Type: types.TargetingValueTypeRange,
			Range: []*types.Range{
				{
					Begin: 1,
					End:   2,
				},
			},
		},
	}

	targeting0 := &types.Targeting{
		Name:          "targeting0",
		TargetingInfo: types.TargetingInfos{targetingInfo},
	}

	assert.NoError(t, UpsertTargeting(db, targeting0))

	list, err := GetAllTargeting(db, nil)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	targeting1 := &types.Targeting{
		Name:          "targeting1",
		TargetingInfo: types.TargetingInfos{targetingInfo},
	}
	assert.NoError(t, UpsertTargeting(db, targeting1))

	tg, err := GetTargetingByID(db, 2)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), tg.ID)
	assert.Equal(t, "targeting1", tg.Name)
	assert.Len(t, tg.TargetingInfo, 1)
	assert.Equal(t, "name1", tg.TargetingInfo[0].Name)
	assert.Len(t, tg.TargetingInfo[0].Values.Range, 1)

	tg.TargetingInfo = append(tg.TargetingInfo, &types.TargetingInfo{
		Name: "name2",
		Values: &types.TargetingValue{
			Type: types.TargetingValueTypeRange,
			Range: []*types.Range{
				{
					Begin: 1,
					End:   2,
				},
			},
		},
	})

	assert.NoError(t, UpsertTargeting(db, tg))
	tg, err = GetTargetingByID(db, 2)
	assert.NoError(t, err)
	assert.Len(t, tg.TargetingInfo, 2)
	assert.Len(t, tg.TargetingInfo[1].Values.Range, 1)
	assert.EqualValues(t, 1, tg.TargetingInfo[1].Values.Range[0].Begin)
	assert.EqualValues(t, 2, tg.TargetingInfo[1].Values.Range[0].End)

	assert.NoError(t, DeleteTargetingByID(db, 1))
	tg1, err := GetAllTargeting(db, nil)
	assert.NoError(t, err)
	assert.Len(t, tg1, 1)
}
