package tool

import "github.com/tencentad/martech/api/proto/retrieval"

type RetrievalReqHelper struct {
	featureMap map[string]*retrieval.Feature
}

func NewRetrievalReqHelper(req *retrieval.RetrievalRequest) *RetrievalReqHelper {
	h := &RetrievalReqHelper{
		featureMap: make(map[string]*retrieval.Feature),
	}

	for _, f := range req.Feature {
		h.featureMap[f.Field] = f
	}

	return h
}

func (h *RetrievalReqHelper) GetFeature(field string) *retrieval.Feature {
	return h.featureMap[field]
}
