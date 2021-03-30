package matchengine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildConjunctionHit() [][]*conjunctionHit {
	return [][]*conjunctionHit{
		{
			{ConjunctionId: 1, PredicateId: 0},
			{ConjunctionId: 3, PredicateId: 0},
			{ConjunctionId: 7, PredicateId: 0},
			{ConjunctionId: 15, PredicateId: 0},
		},
		{
			{ConjunctionId: 2, PredicateId: 0},
			{ConjunctionId: 3, PredicateId: 2},
			{ConjunctionId: 7, PredicateId: 1},
			{ConjunctionId: 9, PredicateId: 0},
		},
		{
			{ConjunctionId: 0, PredicateId: 0},
			{ConjunctionId: 3, PredicateId: 1},
			{ConjunctionId: 18, PredicateId: 1},
		},
	}
}

func TestPostingList(t *testing.T) {
	{
		slice := buildConjunctionHit()
		postingList := NewConjunctionPostingList(slice)

		assert.True(t, postingList.IsValid(), "not valid")

		assert.EqualValues(t, 0, postingList.CurEntry().ConjunctionId)
		assert.EqualValues(t, 0, postingList.CurEntry().PredicateId)

		count := 0
		t.Log("------- from 0 ---------")
		for postingList.IsValid() {
			t.Log(postingList.CurEntry().ConjunctionId, ",", postingList.CurEntry().PredicateId)
			count++
			postingList.Next()
		}
		assert.EqualValues(t, 11, count)
	}

	{
		slice := buildConjunctionHit()
		postingList := NewConjunctionPostingList(slice)

		t.Log("------- from 15 ---------")
		postingList.SeekTo(15)
		count := 0
		for postingList.IsValid() {
			t.Log(postingList.CurEntry().ConjunctionId, ",", postingList.CurEntry().PredicateId)
			count++
			postingList.Next()
		}
		assert.EqualValues(t, 2, count)
	}
}
