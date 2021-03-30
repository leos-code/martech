package media

import (
	"net/http"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/cmd/rta/server/media/ams"
)

var Path2Media = map[types.ADPlatformType]Handler{
	types.ADPlatformTypeAMS: &ams.Handler{},
}

// GetMediaHandler get handler by media name
func GetMediaHandler(mediaName types.ADPlatformType) Handler {
	return Path2Media[mediaName]
}

// Handler handler interface definition
type Handler interface {
	// Parse 解析请求
	Parse(req *http.Request) (*data.RTAContext, error)

	// SendRsp 发送回包, 当err不为nil时，ctx可能为nil
	SendRsp(w http.ResponseWriter, ctx *data.RTAContext, err error)
}
