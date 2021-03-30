package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOption_value(t *testing.T) {
	opt := &FilterOption{
		Operation: FilterOperationLike,
		Value: "a",
	}

	assert.EqualValues(t, "%a%", opt.value())
}