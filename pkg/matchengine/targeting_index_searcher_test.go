package matchengine

import (
	"testing"

	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/targeting"
	pb "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestIndexSearcher(t *testing.T) {
	var err error
	var ad targeting.Targeting
	var text = ` id:300 
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
		        		type: ID
						id: 1001
		      		>
		    	>
			>
		>
	>`
	err = pb.UnmarshalText(text, &ad)
	assert.NoError(t, err)
	var ad2 targeting.Targeting
	var text2 = `id:200 
	betree: < 
		not: false
		op: And
	  	predicate: <
			field: "income"
			value: <
		    	type: RANGE
		      	range: <
		        	begin: 5000
		        	end: 8000
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
		          		begin: 10
		          		end: 20
		        	>
		      	>
		    >
		    betree: <
				op: And
				not: false
		    	predicate: <
					field: "interest"
		    		value: <
		        		type: String
		        		str: "王者荣耀"
		      		>
		    	>
			>
		>
	>`
	err = pb.UnmarshalText(text2, &ad2)
	assert.NoError(t, err)

	var beList []*targeting.Targeting
	beList = append(beList, &ad)
	beList = append(beList, &ad2)
	var index targeting.TargetingIndex
	builder := NewTargetingIndexBuilder(&index)
	builder.Build(beList)

	se := NewTargetingIndexSearcher(&index)
	t.Log(pb.MarshalTextString(&index))
	localIds := se.Search(NewSearchContext(), &retrieval.RetrievalRequest{
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
						Str:  "王者荣耀",
					},
				},
			},
		},
	})
	assert.Len(t, localIds, 1)
	assert.EqualValues(t, 1, localIds[0])
	t.Log(se.GetTargeting(localIds[0]))
}
