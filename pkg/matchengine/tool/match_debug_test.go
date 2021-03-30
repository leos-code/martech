package tool

import (
	"testing"

	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewMatchDebug(t *testing.T) {
	md := NewMatchDebug()

	var target targeting.Targeting
	var text = ` id:100 
	betree: < 
		not: false
		op: And
	  	predicate: <
			field: "income"
			value: <
		    	type: RANGE
		      	range: <
		        	begin: 8000
		        	end: 10000
		      	>
		    >
		>
		betree: <
			op: And
		    predicate: <
				field: "age"
		    	value: <
		        	type: RANGE
		        	range: <
		          		begin: 30
		          		end: 35
		        	>
		      	>
		    >
		    betree: <
				op: And
		    	predicate: <
					field: "interest"
		    		value: <
		        		type: String
						str: "moba" 
		      		>
		    	>
			>
		>
	>`
	assert.NoError(t, proto.UnmarshalText(text, &target))

	t.Log(md.MatchDebug(&target, &retrieval.RetrievalRequest{
		Feature: []*retrieval.Feature{
			{
				Field: "income",
				Value: []*retrieval.Feature_Value{
					{
						Type: retrieval.Feature_Value_ID,
						Id:   5000,
					},
				},
			},
			{
				Field: "age",
				Value: []*retrieval.Feature_Value{
					{
						Type: retrieval.Feature_Value_ID,
						Id:   15,
					},
				},
			},
			{
				Field: "interest",
				Value: []*retrieval.Feature_Value{
					{
						Type: retrieval.Feature_Value_String,
						Str:  "rpg",
					},
				},
			},
		},
	}))
}
