package tool

import (
	"strconv"

	"github.com/tencentad/martech/api/proto/targeting"
)

func BETreeDebugString(be *targeting.BETree) string {
	return "\n" + beTreeDebugString(be, "")
}

func beTreeDebugString(be *targeting.BETree, indent string) string {
	res := ""
	opLine := indent
	if be.Not {
		opLine += "!"
	}
	if be.Op == targeting.LogicalOp_And {
		opLine += "&&"
	} else {
		opLine += "||"
	}
	res += opLine + "\n"
	nextIndent := indent + "\t"
	if be.Predicate != nil {
		preLine := nextIndent
		preLine += predicateDebugString(be.Predicate) + "\n"
		res += preLine
	}
	for _, subBe := range be.Betree {
		res += beTreeDebugString(subBe, nextIndent)
	}
	return res
}

func predicateDebugString(p *targeting.Predicate) string {
	res := p.Field + ":"
	if p.Not {
		res += "!"
	}
	res += "{"
	for i := range p.Value {
		v := p.Value[i]
		if i != 0 {
			res += ";"
		}
		switch v.Type {
		case targeting.Predicate_Value_ID:
			res += strconv.FormatUint(v.Id, 10)
		case targeting.Predicate_Value_RANGE:
			res += "[" + strconv.FormatUint(v.Range.Begin, 10)
			res += "," + strconv.FormatUint(v.Range.End, 10) + ")"
		case targeting.Predicate_Value_String:
			res += v.Str
		}
	}
	res += "}"
	return res
}
