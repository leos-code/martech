package matchengine

import (
	"testing"

	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestIndexBuilder(t *testing.T) {
	var ad targeting.Targeting
	var text = `id:100 
	betree: < 
		not: false
		op: And
	  	predicate: <
			value: <
		    	type: RANGE
		      	range: <
		        	begin: 1
		        	end: 2
		      	>
		    >
		>
		betree: <
			op: And
			not: false
		    predicate: <
		    	value: <
		        	type: RANGE
		        	range: <
		          		begin: 100
		          		end: 200
		        	>
		      	>
		    >
		    betree: <
				op: And
				not: true
		    	predicate: <
		    		value: <
		        		type: RANGE
		        		range: <
		          			begin: 1000
		          			end: 2000
		        		>
		      		>
		    	>
			>
		>
	>`
	assert.NoError(t, proto.UnmarshalText(text, &ad))
	var ad2 targeting.Targeting
	var text2 = `id:200 
	betree: < 
		not: false
		op: And
	  	predicate: <
			value: <
		    	type: RANGE
		      	range: <
		        	begin: 1
		        	end: 2
		      	>
		    >
		>
		betree: <
			op: And
			not: false
		    predicate: <
		    	value: <
		        	type: RANGE
		        	range: <
		          		begin: 100
		          		end: 200
		        	>
		      	>
		    >
		    betree: <
				op: And
				not: true
		    	predicate: <
		    		value: <
		        		type: RANGE
		        		range: <
		          			begin: 1000
		          			end: 2000
		        		>
		      		>
		    	>
			>
		>
	>`
	assert.NoError(t, proto.UnmarshalText(text2, &ad2))

	var beList []*targeting.Targeting
	beList = append(beList, &ad)
	beList = append(beList, &ad2)
	var index targeting.TargetingIndex
	builder := NewTargetingIndexBuilder(&index)
	builder.Build(beList)
	t.Log(proto.MarshalTextString(&index))
}
