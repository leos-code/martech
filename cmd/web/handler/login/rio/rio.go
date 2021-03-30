package rio

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/config"
	"github.com/tencentad/martech/cmd/web/ginx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	cRioPaasID    = "x-rio-paasid"
	cRioSignature = "x-rio-signature"
	cRioTimestamp = "x-rio-timestamp"
	cRioNonce     = "x-rio-nonce"

	redirectUrl       = "https://passport.woa.com/modules/passport/signin.ashx?oauth=true&url=%v&appkey=%v"
	accessTokenUrl    = "http://idc.rio.tencent.com/ebus/tof4/api/v1/passport/AccessToken?code=%v"
	devAccessTokenUrl = "http://rio.tencent.com/ebus/tof4/api/v1/passport/AccessToken?code=%v"
)

var (
	RelativePath = "/rio"
)

type data struct {
	ChineseName  string    `json:"ChineseName"`
	LoginName    string    `json:"LoginName"`
	DeptId       int       `json:"DeptId"`
	DeptName     string    `json:"DeptName"`
	IssueDate    time.Time `json:"IssueDate"`
	Expiration   time.Time `json:"Expiration"`
	IsPersistent bool      `json:"IsPersistent"`
	Key          string    `json:"Key"`
	StaffId      int       `json:"StaffId"`
	Ticket       string    `json:"Ticket"`
	Token        string    `json:"Token"`
	Version      int       `json:"Version"`
}

// LoginUser 获取登录用户接口
func (d *data) LoginUser() *types.LoginUser {
	return &types.LoginUser{
		OpenID:    fmt.Sprint(d.StaffId),
		LoginType: types.LoginTypeRio,
		Avatar:    fmt.Sprintf("http://r.hrc.oa.com/photo/150/%v.png", d.LoginName),
		NickName:  d.LoginName,
	}
}

type response struct {
	Data       *data  `json:"Data"`
	ErrCode    int    `json:"ErrCode"`
	ErrMsg     string `json:"ErrMsg"`
	Ret        int    `json:"Ret"`
	StackTrace string `json:"StackTrace"`
}

func calcSign(token string, timestamp string, nonce string) string {
	d := fmt.Sprintf("%s%s%s%s", timestamp, token, nonce, timestamp)
	return strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(d))))
}

// GetRedirect 获取回调地址
func GetRedirect(config *config.LoginOption) *types.RedirectLogin {
	host, _ := url.Parse(config.CallBackHost)
	host.Path = path.Join(host.Path, "/portal/login/callback", RelativePath)
	u := url.QueryEscape(host.String())
	id := config.AppID
	return &types.RedirectLogin{
		Type: types.LoginTypeRio,
		Url:  fmt.Sprintf(redirectUrl, u, id),
	}
}

// CallBackHandler 登录接口
func CallBackHandler(config *config.LoginOption) gin.HandlerFunc {
	uu := accessTokenUrl
	if gin.IsDebugging() {
		uu = devAccessTokenUrl
	}

	return func(c *gin.Context) {
		code := c.Query("code")
		u := fmt.Sprintf(uu, code)
		req, _ := http.NewRequest("GET", u, nil)

		timestamp := fmt.Sprint(time.Now().Unix())
		nonce := uuid.New().String()
		signature := calcSign(config.Token, timestamp, nonce)
		req.Header.Set(cRioPaasID, config.AppID)
		req.Header.Set(cRioSignature, signature)
		req.Header.Set(cRioTimestamp, timestamp)
		req.Header.Set(cRioNonce, nonce)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ginx.ResponseWithError(c, http.StatusInternalServerError, err)
			return
		}

		defer func() { _ = resp.Body.Close() }()

		result := &response{}
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			ginx.ResponseWithError(c, http.StatusInternalServerError, err)
			return
		}

		if result.ErrCode != 0 {
			ginx.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("%v", result))
			return
		}

		c.Set(gin.AuthUserKey, result.Data.LoginUser())
	}
}
