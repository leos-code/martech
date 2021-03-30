package action

import "github.com/tencentad/martech/cmd/rta/server/data"

// FLResultFetchAction 联邦学习结果拉取
type FLResultFetchAction struct {
}

// NewFLResultFetchAction
func NewFLResultFetchAction() *FLResultFetchAction {
	return &FLResultFetchAction{}
}

// Run implement workflow.Runnable
func (action *FLResultFetchAction) Run(i interface{}) {
	c := i.(*data.RTAContext)

	// TODO 从模型缓存中获取对应的分数

	c.FLModelScore = map[string]float64{
		"new_user": 0.8,
	}
}





