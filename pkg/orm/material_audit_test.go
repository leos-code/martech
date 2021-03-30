package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestListMaterialAuditDetail(t *testing.T) {
	material1 := &types.Material{
		Name: "m1",
	}

	material2 := &types.Material{
		Name: "m2",
	}

	user := &types.User{
		Email: "test@qq.com",
	}

	db, err := getTestDB(t)
	assert.NoError(t, err)
	assert.NoError(t, Setup(db))

	assert.NoError(t, Upsert(db, user, material1, material2))

	materialAudit1 := &types.MaterialAudit{
		MaterialID:  1,
		UserID:      1,
		AuditStatus: types.MaterialAuditPass,
	}

	materialAudit2 := &types.MaterialAudit{
		MaterialID:   2,
		UserID:       1,
		AuditStatus:  types.MaterialAuditReject,
		RejectReason: "not valid",
	}
	assert.NoError(t, SubmitMaterialAudit(db, []*types.MaterialAudit{materialAudit1, materialAudit2}))

	var listMaterialAudit []*types.MaterialAudit
	listMaterialAudit, err = ListMaterialAuditDetail(db, &StatementOption{
		Filter: []*FilterOption{
			{
				Field:     "User__id",
				Operation: FilterOperationIn,
				Value:     []float64{1},
			},
		},
	})
	assert.NoError(t, err)

	assert.Len(t, listMaterialAudit, 2)
	assert.EqualValues(t, 1, listMaterialAudit[0].User.ID)
	assert.EqualValues(t, "m2", listMaterialAudit[1].Material.Name)

	var materials []*types.Material
	var page *Page
	page, err = MaterialPage(db, &StatementOption{Pagination: &PaginationOption{PageSize: 1, Page: 1}})
	assert.NoError(t, err)
	assert.EqualValues(t, 2, page.Total)
	materials = page.List.([]*types.Material)
	assert.Len(t, materials, 1)
	assert.EqualValues(t, "m1", materials[0].Name)
	assert.EqualValues(t, types.MaterialAuditPass, materials[0].Audit[0].AuditStatus)
}
