package action

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/tencentad/martech/pkg/common/loader"
)

type FLResultApplyAction struct {
	rule *loader.FileDoubleBuffer
}

func NewFLResultApplyAction() *FLResultApplyAction {
	return &FLResultApplyAction{}
}

func (action *FLResultApplyAction) Init(ruleFile string) error {
	if ruleFile == "" {
		return nil
	}

	ruleLoader, wait := loader.CreateFileLoaderWithWaitDone(
		ruleFile,
		func() interface{} {
			return &ApplyRule{}
		},
		loadApplyRule,
		time.Second * 10,
		)

	action.rule = loader.NewFileDoubleBuffer(ruleLoader)

	if err := wait(); err != nil {
		return err
	}

	return nil
}

func (action *FLResultApplyAction) Run(_ interface{}) {

}

func (action *FLResultApplyAction) getApplyRule() *ApplyRule {
	return action.rule.Data().(*ApplyRule)
}

type ApplyRule struct {
	ModelName    string   `json:"model_name"`
	StrategyID   []string `json:"strategy_id,omitempty"`
	AdvertiserID []string `json:"advertiser_id,omitempty"`
	CampaignID   []string `json:"campaign_id,omitempty"`
}

func loadApplyRule(reader io.Reader, i interface{}) error {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, i)
}
