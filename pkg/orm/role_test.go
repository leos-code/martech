package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item1 := &types.Role{
		Name:     "admin",
		TenantID: 1,
	}
	assert.NoError(t, UpsertRole(db, item1))

	item2 := &types.Role{
		Name:     "normal",
		TenantID: 1,
	}
	assert.NoError(t, UpsertRole(db, item2))

	item3 := &types.Role{
		Name:     "admin",
		TenantID: 2,
	}
	assert.NoError(t, UpsertRole(db, item3))

	item4 := &types.Role{
		Name:     "admin",
		TenantID: 1,
	}
	assert.NoError(t, UpsertRole(db, item4))
	assert.Equal(t, item1.ID, item4.ID)
}
