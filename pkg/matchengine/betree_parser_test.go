package matchengine

import (
	"testing"

	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBEtreeParser(t *testing.T) {
	var ad targeting.Targeting
	var text string = ` id:100 
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
			op: Or
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
				not: false
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
	t.Log(proto.MarshalTextString(&ad))
	var dnf targeting.TargetingDNF
	ConvertToTargetingDNF(&ad, &dnf)
	t.Log(proto.MarshalTextString(&dnf))
}
