package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestBindStrategy(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	targeting := &types.Targeting{
		Name: "targeting0",
	}
	assert.NoError(t, UpsertTargeting(db, targeting))

	bindStrategy := &types.BindStrategy{
		Name:        "bindStrategy",
		Platform:    "unknown",
		TargetingID: 1,
		Strategy: &types.Strategy{
			StrategyID:   []string{"1"},
			AdvertiserID: []string{"2"},
		},
	}
	assert.NoError(t, UpsertBindStrategy(db, bindStrategy))

	var related1 []*types.BindStrategy
	assert.NoError(t, db.Model(targeting).Association("BindStrategy").Find(&related1))
	assert.Len(t, related1, 1)
	var related2 []*types.BindStrategy
	assert.NoError(t, db.Model(&types.Targeting{ID: 2}).Association("BindStrategy").Find(&related2))
	assert.Len(t, related2, 0)

	bs, err := GetBindStrategyByID(db, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), bs.ID)
	assert.Equal(t, "bindStrategy", bs.Name)
	assert.NotNil(t, bs.Strategy)
	assert.Len(t, bs.Strategy.StrategyID, 1)
	assert.Len(t, bs.Strategy.AdvertiserID, 1)
	assert.Len(t, bs.Strategy.CampaignID, 0)
	assert.Equal(t, "1", bs.Strategy.StrategyID[0])
	assert.Equal(t, "2", bs.Strategy.AdvertiserID[0])

	bs.Strategy.StrategyID = nil
	bs.Strategy.AdvertiserID = append(bs.Strategy.AdvertiserID, "3")

	assert.NoError(t, UpsertBindStrategy(db, bs))
	bs2, err := GetAllBindStrategy(db, nil)
	assert.NoError(t, err)
	assert.Len(t, bs2, 1)
	assert.Len(t, bs2[0].Strategy.StrategyID, 0)
	assert.Len(t, bs2[0].Strategy.AdvertiserID, 2)
	assert.Equal(t, "3", bs2[0].Strategy.AdvertiserID[1])

	assert.NoError(t, DeleteBindStrategyById(db, 1))
	bs3, err := GetAllBindStrategy(db, nil)
	assert.NoError(t, err)
	assert.Len(t, bs3, 0)
}

func TestStrategyExpression(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	bindStrategy := &types.BindStrategy{
		Name:        "bindStrategy",
		Platform:    "unknown",
		TargetingID: 1,
		Strategy: &types.Strategy{
			StrategyID:   []string{"1"},
			AdvertiserID: []string{"2"},
		},
	}
	assert.NoError(t, UpsertBindStrategy(db, bindStrategy))

	var bs []*types.BindStrategy
	assert.NoError(t, db.Find(&bs, datatypes.JSONQuery("strategy").HasKey("strategy_id")).Error)
}
