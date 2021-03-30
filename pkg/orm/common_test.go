package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsert(t *testing.T) {
	db, err := getTestDB(t)
	assert.NoError(t, err)
	assert.NoError(t, Setup(db))

	{
		notPointer := struct {
		}{}

		err := Upsert(db, notPointer)
		assert.Error(t, err)
	}

	{
		noID := &struct {
		}{}
		err := Upsert(db, noID)
		assert.Error(t, err)
	}

}
