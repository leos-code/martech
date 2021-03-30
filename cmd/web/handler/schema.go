package handler

import (
	"encoding/json"

	"github.com/tencentad/martech/pkg/schema"
	"github.com/gin-gonic/gin"
)

// SchemaGetHandler
func SchemaGetHandler(c *gin.Context) {
	mockSchema := `
{
	"version": 1,
	"fields": [{
			"name": "age",
            "display_name": "年龄",
			"type": "integer"
		},
		{
			"name": "last_active_time",
			"display_name": "上次活跃时间",
			"type": "enum",
			"enum": [{
					"value": "1周内"
				},
				{
					"value": "1周以外，1月以内"
				},
				{
					"value": "1月以外"
				}
			]
		},
		{
			"name": "self_active_rate",
			"display_name": "自启概率",
			"type": "integer"
		}
	]
}
`
	var s schema.Schema
	_ = json.Unmarshal([]byte(mockSchema), &s)
	handleWithData(c, "get success", &s, nil)
}
