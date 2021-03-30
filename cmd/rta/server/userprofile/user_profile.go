package userprofile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tencentad/martech/cmd/rta/server/data"
)

// Getter 获取用户画像接口
type Getter interface {
	GetUserProfile(c *data.RTAContext) error
}

type httpGetter struct {
	url string
}

// NewHttpGetter
func NewHttpGetter(url string) *httpGetter {
	return &httpGetter{
		url: url,
	}

}

// GetUserProfile implement Getter
func (g *httpGetter) GetUserProfile(c *data.RTAContext) error {
	resp, err := http.Get(g.url)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(content, &c.UserProfile.Feature); err != nil {
		return err
	}

	return nil
}
