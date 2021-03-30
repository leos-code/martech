package data

import (
	"github.com/tencentad/martech/api/proto/novelty"
	rtaPb "github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/common/data"
	"github.com/tencentad/martech/pkg/deviceid"
)

// RTAContext rta请求所有的上下文信息
type RTAContext struct {
	*data.BaseContext

	// 广告平台名称
	Platform    types.ADPlatformType
	RequestId   string
	RequestType RequestType

	SiteSetId uint64 // 站点集

	ID *deviceid.ID

	UserProfile *UserProfile

	HitRtaStrategy []*rtaPb.BindStrategy

	FLModelScore map[string]float64 // 联邦学习分数
	Novelty      *novelty.Novelty
}

// NewRTAContext
func NewRTAContext() *RTAContext {
	return &RTAContext{
		BaseContext: data.NewBaseContext(),
		UserProfile: newUserProfile(),
	}
}

// RequestType 请求类型
type RequestType string

// 请求类型列表
const (
	RequestTypeRTA   RequestType = "rta"   // 正常RTA请求
	RequestTypeClick RequestType = "click" // 点击
)
