package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.LoginUser{
		OpenID:    "1",
		LoginType: types.LoginTypeRio,
	}
	assert.NoError(t, UpsertLoginUser(db, item))
	{
		list, err := GetAllLoginUser(db)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	}

	item.Avatar = "avatar"
	assert.NoError(t, UpsertLoginUser(db, item))
	{
		list, err := GetAllLoginUser(db)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	}

	user0 := &types.LoginUser{
		OpenID:    "1",
		LoginType: types.LoginTypeRio,
	}
	assert.NoError(t, TakeLoginUser(db, user0))
	assert.Equal(t, "avatar", user0.Avatar)

	user1 := &types.LoginUser{
		OpenID:    "0",
		LoginType: types.LoginTypeRio,
	}
	assert.Error(t, TakeLoginUser(db, user1))
}
