package ams

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tencentad/martech/api/proto/rta/ams"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/pkg/deviceid"
	"github.com/golang/protobuf/proto"
)

// Handler ams handler implements
type Handler struct {
}

// Parse 解析请求
func (handler *Handler) Parse(req *http.Request) (*data.RTAContext, error) {
	body, err:= ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	rtaReq := &ams.RtaRequest{}
	if err = proto.Unmarshal(body, rtaReq); err != nil {
		return nil, err
	}

	s := rtaReq.String()
	fmt.Println(s)

	ctx := data.NewRTAContext()
	ctx.Platform = types.ADPlatformTypeAMS

	ctx.RequestId = rtaReq.GetId()
	ctx.SiteSetId = rtaReq.GetSitesetId()

	device := rtaReq.Device

	ctx.ID = &deviceid.ID{
		IDFAMD5:      device.IdfaMd5Sum,
		IMEIMD5:      device.ImeiMd5Sum,
		OAID:         device.Oaid,
		OAIDMD5:      device.OaidMd5Sum,
		AndroidIDMD5: device.AndroidIdMd5Sum,
		MACMD5:       device.MacMd5Sum,
		QAID:         device.Oaid,
	}

	if requestInfo := rtaReq.GetRequestInfo(); requestInfo != nil {
		if requestInfo.RequestType == ams.RtaRequest_CLICK_REQUEST {
			ctx.RequestType = data.RequestTypeClick
		}
	}

	return ctx, nil
}

// SendRsp 返回数据
func (handler *Handler) SendRsp(w http.ResponseWriter, ctx *data.RTAContext, err error) {
	w.Header().Add("Connection", "Keep-Alive")
	var rsp *ams.RtaResponse
	if err != nil {
		rsp = &ams.RtaResponse{
			RequestId:   ctx.RequestId,
			Code:        0,
			OutTargetId: []string{},
		}
	} else {
		rsp = &ams.RtaResponse{
			RequestId:   ctx.RequestId,
			Code:        0,
			OutTargetId: []string{},
		}
	}

	b, _ := proto.Marshal(rsp)
	_, _ = w.Write(b)
}
