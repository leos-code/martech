package orm

import (
	"testing"

	"github.com/tencentad/martech/api/types"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	item := &types.User{
		PhoneNumber: "1",
		Email:       "1",
	}
	assert.NoError(t, UpsertUser(db, item))

	u := &types.User{PhoneNumber: "1"}
	assert.NoError(t, TakeUser(db, u))
	assert.Equal(t, uint64(1), u.ID)

	u2 := &types.User{PhoneNumber: "3"}
	assert.Error(t, TakeUser(db, u2))
}

func TestSearchUser(t *testing.T) {
	db, _ := getTestDB(t)
	_ = Setup(db)

	user1 := &types.User{
		PhoneNumber: "123",
		Email:       "abc",
		LoginUser: []*types.LoginUser{
			{
				OpenID:    "1",
				LoginType: types.LoginTypeRio,
				NickName:  "alice",
			},
			{
				OpenID:    "2",
				LoginType: types.LoginTypeQQ,
				NickName:  "alic",
			},
		},
	}
	assert.NoError(t, UpsertUser(db, user1))

	user2 := &types.User{
		PhoneNumber: "124",
		Email:       "abd",
		LoginUser: []*types.LoginUser{
			{
				OpenID:    "3",
				LoginType: types.LoginTypeRio,
				NickName:  "ali",
			},
		},
	}
	assert.NoError(t, UpsertUser(db, user2))

	{
		r, err := SearchUser(db, &types.UserSearch{
			PhoneNumber: "12",
			Approximate: true,
		})
		assert.NoError(t, err)
		assert.Len(t, r, 2)
	}

	{
		r, err := SearchUser(db, &types.UserSearch{
			PhoneNumber: "12",
		})
		assert.NoError(t, err)
		assert.Len(t, r, 0)
	}

	{
		r, err := SearchUser(db, &types.UserSearch{
			PhoneNumber: "12",
			Email:       "abd",
			Approximate: true,
		})
		assert.NoError(t, err)
		assert.Len(t, r, 1)
	}

	{
		r, err := SearchUser(db, &types.UserSearch{
			NickName:    "alic",
			Approximate: true,
		})
		assert.NoError(t, err)
		assert.Len(t, r, 1)
	}

	{
		r, err := SearchUser(db, &types.UserSearch{
			NickName: "ali",
		})
		assert.NoError(t, err)
		assert.Len(t, r, 1)
	}

	{
		r, err := SearchUser(db, &types.UserSearch{})
		assert.NoError(t, err)
		assert.Len(t, r, 0)
	}

}
