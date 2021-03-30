package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMaterial(t *testing.T) {
	material1 := &types.Material{
		Name: "m1",
	}

	material2 := &types.Material{
		Name: "m2",
	}
	db, err := getTestDB(t)
	assert.NoError(t, err)
	assert.NoError(t, Setup(db))

	assert.NoError(t, Upsert(db, material1))
	assert.NoError(t, Upsert(db, material2))

	var materials []*types.Material
	var page *Page
	page, err = MaterialPage(db, &StatementOption{Pagination: &PaginationOption{PageSize: 1, Page: 2}})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, page.Total)
	materials = page.List.([]*types.Material)
	assert.Len(t, materials, 1)
	assert.EqualValues(t, "m2", materials[0].Name)

	assert.NoError(t, MaterialDelete(db, 1))
	page, err = MaterialPage(db, &StatementOption{Pagination: &PaginationOption{PageSize: 1, Page: 1}})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, page.Total)
	materials = page.List.([]*types.Material)
	assert.Len(t, materials, 1)
	assert.EqualValues(t, "m2", materials[0].Name)
}
