package helper

import (
	"testing"

	"github.com/tencentad/martech/pkg/schema"
	"github.com/stretchr/testify/assert"
)

func Test_numberToValueID(t *testing.T) {
	{
		f := &schema.Field{}
		var v int8 = -1
		_, err := numberToValueID(f, v)
		t.Log(err)
		assert.Error(t, err)
	}
	{
		f := &schema.Field{}
		var v int8 = 10
		id, err := numberToValueID(f, v)
		assert.NoError(t, err)
		assert.EqualValues(t, v, id)
	}
	{
		f := &schema.Field{}
		var v uint8 = 10
		id, err := numberToValueID(f, v)
		assert.NoError(t, err)
		assert.EqualValues(t, v, id)
	}
	{
		f := &schema.Field{}
		var v float32 = -1
		_, err := numberToValueID(f, v)
		assert.Error(t, err)
		t.Log(err)
	}
	{
		f := &schema.Field{}
		var v float32 = 1.1
		_, err := numberToValueID(f, v)
		assert.Error(t, err)
		t.Log(err)
	}

	{
		f := &schema.Field{}
		var v float32 = 1
		id, err := numberToValueID(f, v)
		assert.NoError(t, err)
		assert.EqualValues(t, v, id)
	}

	{
		f := &schema.Field{}
		var v = "a"
		_, err := numberToValueID(f, v)
		assert.Error(t, err)
		t.Log(err)
	}
}

func TestHelper(t *testing.T) {
	var content = `
{
  "version": 1,
  "fields": [
    {
      "name": "age",
      "type": "integer",
      "range": {
        "min": 0,
        "max": 200
      }
    },
    {
      "name": "interest",
      "type": "enum",
      "enum": [
        {
          "value": "game",
          "children": [
            {
              "value": "moba",
              "children": [
                {
                  "value": "wangzhe"
                }
              ]
            },
            {
              "value": "和平精英"
            },
            {
              "value": "阴阳师"
            }
          ]
        },
        {
          "value": "音乐"
        }
      ]
    },
    {
      "name": "download_song",
      "type": "string"
    }
  ]
}
`
	h, err := NewSchemaHelperFromContent([]byte(content))
	assert.NoError(t, err)

	{
		fv, err := h.GetFeatureValues("interest", "王者荣耀")
		assert.NoError(t, err)
		assert.EqualValues(t, "王者荣耀", fv[0].Str)
	}
	{
		_, err := h.GetFeatureValues("interest", 1)
		assert.Error(t, err)
	}
	{
		fv, err := h.GetFeatureValues("interest", "wangzhe", "王者荣耀")
		assert.NoError(t, err)
		assert.Len(t, fv, 4)
	}
	{
		fv, err := h.GetFeatureValues("download_song", "someone like you")
		assert.NoError(t, err)
		assert.EqualValues(t, "someone like you", fv[0].Str)
	}
	{
		fv, err := h.GetFeatureValues("age", 1)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, fv[0].Id)
	}
}
