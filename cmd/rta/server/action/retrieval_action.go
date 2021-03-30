package action

import (
	"reflect"

	pb "github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/pkg/matchengine/helper"
	"github.com/tencentad/martech/pkg/retrieval"
	log "github.com/sirupsen/logrus"
)

// RetrievalAction 检索满足条件的rta策略
type RetrievalAction struct {
	schemaHelper *helper.Helper
	serviceImpl  *retrieval.ServiceImpl
}

// NewRetrievalAction
func NewRetrievalAction() *RetrievalAction {
	return &RetrievalAction{
		serviceImpl: retrieval.GetServiceImpl(),
	}
}

// Init 初始化
func (action *RetrievalAction) Init(schemaPath string) error {
	h, err := helper.NewSchemaHelper(schemaPath)
	if err != nil {
		return err
	}

	action.schemaHelper = h
	return nil
}

// Run
func (action *RetrievalAction) Run(i interface{}) {
	c := i.(*data.RTAContext)
	if err := action.retrieval(c); err != nil {
		log.Errorf("failed to execute retrieval action")
	}
}

func (action *RetrievalAction) retrieval(c *data.RTAContext) error {
	if action.serviceImpl != nil {
		log.Warnf("retrieval impl is nil, skip retrieval")
		return nil
	}

	req, err := buildRetrievalReq(action.schemaHelper, c.UserProfile)
	if err != nil {
		return err
	}

	resp, _ := action.serviceImpl.Retrieve(c.Ctx, req)
	for _, info := range resp.HitRtaInfo {
		for _, strategy := range info.BindStrategy {
			c.HitRtaStrategy = append(c.HitRtaStrategy, strategy)
		}
	}

	return nil
}

func buildRetrievalReq(h *helper.Helper, prof *data.UserProfile) (*pb.RetrievalRequest, error) {
	req := &pb.RetrievalRequest{}
	for name, value := range prof.Feature {
		featureValues, err := h.GetFeatureValues(name, convertValue(value)...)
		if err != nil {
			log.Errorf("failed to GetFeatureValues, name[%s], value[%v], err: %v", name, value, err)
		}
		req.Feature = append(req.Feature, &pb.Feature{
			Field: name,
			Value: featureValues,
		})
	}
	return req, nil
}

// convertValue 假如是slice，显示的转化成slice，
func convertValue(i interface{}) []interface{} {
	obj := reflect.ValueOf(i)

	switch obj.Kind() {
	case reflect.Slice, reflect.Array:
		length := obj.Len()
		ret := make([]interface{}, 0, length)
		for i := 0; i < length; i++ {
			ret = append(ret, obj.Index(i).Interface())
		}
		return ret
	default:
		return []interface{}{i}
	}
}
