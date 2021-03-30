package matchengine

import (
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/golang/protobuf/proto"
	"github.com/willf/bitset"
)

const (
	ParseTypeDNF     = 1
	ParseTypeCNF     = 2
	maxPredicateSize = 1024
)

// 一个bit代表一个predicate, 一个bitset代表一个conjunction, []bitset代表一个dnf
type parseResult []*bitset.BitSet

// BETreeParser 解析BETree，生成范式
type BETreeParser struct {
	result       parseResult
	predicateVec []*targeting.Predicate
}

func needMultiply(parseType int, op targeting.LogicalOp) bool {
	if parseType == ParseTypeDNF && op == targeting.LogicalOp_And {
		return true
	}
	if parseType == ParseTypeCNF && op == targeting.LogicalOp_Or {
		return true
	}
	return false
}

func needAdd(parseType int, op targeting.LogicalOp) bool {
	if parseType == ParseTypeDNF && op == targeting.LogicalOp_Or {
		return true
	}
	if parseType == ParseTypeCNF && op == targeting.LogicalOp_And {
		return true
	}
	return false
}

func (p *BETreeParser) parse(betree *targeting.BETree, parseType int) {
	p.result = p.parseInternal(betree, parseType, false)
}

func (p *BETreeParser) add(resA, resB parseResult) parseResult {
	return append(resA, resB...)
}

func (p *BETreeParser) multiply(resA, resB parseResult) parseResult {
	var tmp parseResult
	for _, a := range resA {
		for _, b := range resB {
			tmp = append(tmp, a.Union(b))
		}
	}
	return tmp
}

func (p *BETreeParser) parseInternal(betree *targeting.BETree, parseType int, not bool) parseResult {
	localNot := betree.Not
	if not {
		localNot = !localNot
	}
	var result parseResult
	if betree.Predicate != nil {
		predicate := proto.Clone(betree.Predicate).(*targeting.Predicate)
		predicateNot := betree.Predicate.Not
		if localNot {
			predicateNot = !predicateNot
		}
		predicate.Not = predicateNot
		index := uint(len(p.predicateVec))
		p.predicateVec = append(p.predicateVec, predicate)
		bset := bitset.New(maxPredicateSize)
		bset.Set(index)
		result = append(result, bset)
	}
	logicalOp := betree.Op
	if localNot {
		if logicalOp == targeting.LogicalOp_And {
			logicalOp = targeting.LogicalOp_Or
		} else {
			logicalOp = targeting.LogicalOp_And
		}
	}
	for _, subtree := range betree.Betree {
		subRes := p.parseInternal(subtree, parseType, localNot)
		if len(result) == 0 {
			result = subRes
			continue
		}
		if needMultiply(parseType, logicalOp) {
			result = p.multiply(result, subRes)
		} else if needAdd(parseType, logicalOp) {
			result = p.add(result, subRes)
		}
	}
	return result
}

// 把定向转换成析取范式
func ConvertToTargetingDNF(target *targeting.Targeting, dnf *targeting.TargetingDNF) {
	var parser BETreeParser
	parser.parse(target.Betree, ParseTypeDNF)
	dnf.Id = target.Id
	for _, bSet := range parser.result {
		conj := &targeting.TargetingDNF_Conjunction{}
		for i, e := bSet.NextSet(0); e; i, e = bSet.NextSet(i + 1) {
			p := parser.predicateVec[i]
			predicate := proto.Clone(p).(*targeting.Predicate)
			conj.Predicate = append(conj.Predicate, predicate)
		}
		dnf.Conjunction = append(dnf.Conjunction, conj)
	}
}
