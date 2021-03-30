/*
 * Copyright (c) 2018-2118
 * Author: linceyou
 * LastModified: 19-4-9 下午4:51
 */

package tool

import (
	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/pkg/matchengine"
)

type MatchDebug struct {
}

func NewMatchDebug() *MatchDebug {
	return &MatchDebug{}
}

func (md *MatchDebug) MatchDebug(target *targeting.Targeting, req *retrieval.RetrievalRequest) string {
	var dnf targeting.TargetingDNF
	matchengine.ConvertToTargetingDNF(target, &dnf)
	return md.DNFMatchDebug(&dnf, req)
}

func (md *MatchDebug) DNFMatchDebug(dnf *targeting.TargetingDNF, req *retrieval.RetrievalRequest) string {
	result := "\n"
	for _, conj := range dnf.Conjunction {
		result += "||" + "\n"
		result += md.conjunctionMatchDebug(conj, req, "\t")
	}
	return result
}

func (md *MatchDebug) conjunctionMatchDebug(
	conj *targeting.TargetingDNF_Conjunction, req *retrieval.RetrievalRequest, indent string) string {

	result := ""
	reqHelper := NewRetrievalReqHelper(req)

	result += indent + "&&" + "\n"
	for _, predicate := range conj.Predicate {
		result += indent + "\t" + predicateDebugString(predicate)
		if match := md.matchPredicate(predicate, reqHelper.GetFeature(predicate.Field)); !match {
			result += "[NO_MATCH]"
		}
		result += "\n"
	}

	return result
}

func (md *MatchDebug) matchPredicate(predicate *targeting.Predicate, feature *retrieval.Feature) (result bool) {
	if feature == nil {
		return predicate.Not
	}

	for _, val := range predicate.Value {
		if matchValue(val, feature.Value) {
			return !predicate.Not
		}
	}

	return predicate.Not
}

func matchValue(val *targeting.Predicate_Value, values []*retrieval.Feature_Value) bool {
	switch val.Type {

	case targeting.Predicate_Value_ID:
		for _, v := range values {
			if v.Type == retrieval.Feature_Value_ID && val.Id == v.Id {
				return true
			}
		}
	case targeting.Predicate_Value_RANGE:
		for _, v := range values {
			if v.Type == retrieval.Feature_Value_ID && val.Range.Begin <= v.Id && v.Id < val.Range.End {
				return true
			}
		}
	case targeting.Predicate_Value_String:
		for _, v := range values {
			if v.Type == retrieval.Feature_Value_String && val.Str == v.Str {
				return true
			}
		}

	}

	return false
}
