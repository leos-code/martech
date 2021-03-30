package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content string = `
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

func TestSchema(t *testing.T) {
	var s Schema
	assert.NoError(t, json.Unmarshal([]byte(content), &s))
	assert.NoError(t, s.Check())

	enums := s.Fields[1].GetAllEnum()
	assert.Len(t, enums, 6)
}

func TestField_check(t *testing.T) {
	{
		f := &Field{}
		assert.Error(t, f.check())
	}

	{
		f := &Field{
			Name: "age",
			Type: "double",
		}
		assert.Error(t, f.check())
	}

	{
		f := &Field{
			Name: "age",
			Type: FieldTypeInteger,
			Range: &Range{
				Min: 2,
				Max: 1,
			},
		}
		assert.Error(t, f.check())
	}

	{
		f := &Field{
			Name: "interest",
			Type: FieldTypeEnum,
		}
		assert.Error(t, f.check())
	}

	{
		f := &Field{
			Name: "interest",
			Type: FieldTypeEnum,
			Enum: []*Enum{
				{
					Value: "",
				},
			},
		}
		assert.Error(t, f.check())
	}

	{
		f := &Field{
			Name: "interest",
			Type: FieldTypeEnum,
			Enum: []*Enum{
				{
					Value: "music",
					Children: []*Enum{
						{
							Value: "",
						},
					},
				},
			},
		}
		assert.Error(t, f.check())
		t.Log(f.check())
	}
}
