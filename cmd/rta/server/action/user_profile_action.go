package action

import (
	"github.com/tencentad/martech/cmd/rta/server/data"
	log "github.com/sirupsen/logrus"
)

// UserProfileAction 获取用户画像
type UserProfileAction struct {
	Getter UserProfileGetter
}

// NewUserProfileAction
func NewUserProfileAction() *UserProfileAction {
	return &UserProfileAction{}
}

// Run
func (action *UserProfileAction) Run(i interface{}) {
	if action.Getter == nil {
		return
	}

	c := i.(*data.RTAContext)

	if err := action.Getter.GetUserProfile(c); err != nil {
		log.Errorf("failed to get user profile, err: %v", err)
	}
}

// UserProfileGetter 获取用户画像接口
type UserProfileGetter interface {
	GetUserProfile(c *data.RTAContext) error
}
