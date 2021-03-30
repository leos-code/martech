package druid

import (
	"github.com/tencentad/martech/pkg/druid/common"
	"github.com/awatercolorpen/godruid"
)

type queryParser interface {
	ParseQuery(interface{}) interface{}
}

type parser struct {
	qp queryParser
}

func (p *parser) Parse(resp interface{}) interface{} {
	result := resp.(*godruid.QueryGroupBy).QueryResult
	ch := common.Parallel(result, func(v interface{}) interface{} {
		return p.qp.ParseQuery(v)
	}, 3)
	return ch
}

func NewParser(qb queryParser) *parser {
	return &parser{qp: qb}
}
