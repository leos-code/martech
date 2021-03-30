package logic

import (
	"net/http"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/rta/server/data"
	"github.com/tencentad/martech/cmd/rta/server/media"
	"github.com/tencentad/martech/cmd/rta/server/metrics"
	"github.com/tencentad/martech/cmd/rta/server/userprofile"
	log "github.com/sirupsen/logrus"

	"github.com/tencentad/martech/cmd/rta/server/action"
	"github.com/tencentad/martech/pkg/common/workflow"
)

// RTAService RTA服务
type RTAService struct {
	option *RTAServiceOption

	UserProfileAction      *action.UserProfileAction      // 获取用户画像
	UserNoveltyFetchAction *action.UserNoveltyFetchAction // 获取用户新鲜度
	RetrievalAction        *action.RetrievalAction        // 检索满足的策略
	FLResultFetcherAction  *action.FLResultFetchAction    // 获取联邦模型结果
	FLResultApplyAction    *action.FLResultApplyAction    // 应用联邦模型结果

	WriteClickAction *action.WriteClickAction // 写入点击
}

// RTAServiceOption RTA服务配置选项
type RTAServiceOption struct {
	SchemaPath      string `json:"schema_path"`        // 画像schema路径
	UserProfileUri  string `json:"user_profile_uri"`   // 用户画像uri
	FLApplyRuleFile string `json:"fl_apply_rule_file"` // 联邦学习结果应用规则文件
}

// NewRTAService
func NewRTAService(option *RTAServiceOption) *RTAService {
	return &RTAService{
		option: option,

		UserProfileAction:      action.NewUserProfileAction(),
		UserNoveltyFetchAction: action.NewUserNoveltyFetchAction(),
		RetrievalAction:        action.NewRetrievalAction(),
		FLResultFetcherAction:  action.NewFLResultFetchAction(),
		FLResultApplyAction:    action.NewFLResultApplyAction(),
		WriteClickAction:       action.NewWriteClickAction(),
	}
}

func (s *RTAService) setUserProfileGetter(getter action.UserProfileGetter) {
	s.UserProfileAction.Getter = getter
}

// Init RTA服务初始化
func (s *RTAService) Init() error {
	if err := s.RetrievalAction.Init(s.option.SchemaPath); err != nil {
		return err
	}
	if s.option.UserProfileUri != "" {
		s.setUserProfileGetter(userprofile.NewHttpGetter(s.option.UserProfileUri))
	}
	return nil
}

// ServeHTTP 处理http请求
func (s *RTAService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.process(w, r); err != nil {
		log.Errorf("failed to process rta request, err: %v", err)
		metrics.RTAProcessErrorCount.Add(1)
	}
}

func (s *RTAService) process(w http.ResponseWriter, r *http.Request) error {
	platform := r.URL.Query().Get("platform")

	handler := media.GetMediaHandler(types.ADPlatformType(platform))
	c, err := handler.Parse(r)
	if err != nil {
		handler.SendRsp(w, c, err)
		return err
	}

	wf := s.buildRTAWorkflow()
	wf.StartWithContext(c.Ctx, c)
	wf.WaitDone()

	handler.SendRsp(w, c, c.Error)
	return c.Error
}

func (s *RTAService) buildWorkflow(t data.RequestType) *workflow.WorkFlow {
	if t == data.RequestTypeClick {
		return s.buildClickWorkflow()
	} else {
		return s.buildRTAWorkflow()
	}
}

func (s *RTAService) buildRTAWorkflow() *workflow.WorkFlow {
	wf := workflow.NewWorkFlow()

	userProfileTask := workflow.NewTaskNode(s.UserProfileAction)
	userNoveltyFetchTask := workflow.NewTaskNode(s.UserNoveltyFetchAction)
	retrievalTask := workflow.NewTaskNode(s.RetrievalAction)
	flResultFetcherTask := workflow.NewTaskNode(s.FLResultFetcherAction)
	flResultApplyTask := workflow.NewTaskNode(s.FLResultApplyAction)

	wf.AddTaskNode(userProfileTask)
	wf.AddTaskNode(flResultFetcherTask)
	wf.AddTaskNode(userNoveltyFetchTask)
	wf.AddTaskNode(retrievalTask, userProfileTask)
	wf.AddTaskNode(flResultApplyTask, flResultFetcherTask, retrievalTask)

	return wf
}

func (s *RTAService) buildClickWorkflow() *workflow.WorkFlow {
	wf := workflow.NewWorkFlow()

	writeClickTask := workflow.NewTaskNode(s.WriteClickAction)
	wf.AddTaskNode(writeClickTask)

	return wf
}
