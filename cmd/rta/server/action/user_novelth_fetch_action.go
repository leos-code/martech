package action

import (
	noveltyPB "github.com/tencentad/martech/api/proto/novelty"
	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/pkg/novelty"

	//"github.com/tencentad/martech/pkg/novelty"

)

// UserNoveltyFetchAction 获取用户新鲜度
type UserNoveltyFetchAction struct {
	// TODO implement
	store novelty.UserNoveltyStore
}

// NewUserNoveltyFetchAction
func NewUserNoveltyFetchAction() *UserNoveltyFetchAction {
	return &UserNoveltyFetchAction{}
}

// Run implement workflow.Runnable
func (action *UserNoveltyFetchAction) Run(i interface{}) {
	c := i.(*data.RTAContext)

	c.Novelty = &noveltyPB.Novelty{
		Click: &noveltyPB.Click{

		},
	}
}
