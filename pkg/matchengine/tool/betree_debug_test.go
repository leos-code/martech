package tool

import (
	_ "encoding/xml"
	"testing"

	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBETreeDebug(t *testing.T) {
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
		        		type: ID
						id: 1001
		      		>
		    	>
			>
		>
	>`

	var target targeting.Targeting
	assert.NoError(t, proto.UnmarshalText(text, &target))
	t.Log(BETreeDebugString(target.Betree))
}
