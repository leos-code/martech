package action

import (
	"encoding/json"
	"testing"

	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/pkg/matchengine/helper"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func Test_buildRetrievalReq(t *testing.T) {
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
                  "value": "王者荣耀" 
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
	h, err := helper.NewSchemaHelperFromContent([]byte(content))
	assert.NoError(t, err)

	userFeaturesText := `
{
  "age": 10,
  "interest": ["王者荣耀", "阴阳师"],
  "download_song": ["鱼"]
}
`
	userFeature := make(map[string]interface{})
	err = json.Unmarshal([]byte(userFeaturesText), &userFeature)
	assert.NoError(t, err)

	prof := &data.UserProfile{
		Feature: userFeature,
	}

	req, err := buildRetrievalReq(h, prof)
	assert.NoError(t, err)
	t.Log(proto.MarshalTextString(req))
}
