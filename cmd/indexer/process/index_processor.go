package process

import (
	"time"

	"github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/pkg/common/dumper"
	"github.com/tencentad/martech/pkg/common/loader"
	"github.com/tencentad/martech/pkg/matchengine"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// IndexProcessor 索引处理器
type IndexProcessor struct {
	option      *Option           // 选项
	lastDone    time.Time         // 上次更新完成的时间
	loader      loader.FileLoader // 加载定向数据文件的loader
	FullRTAInfo *rta.FullRTAInfo  // 输入, RTA定向数据
	RTAIndex    *rta.RTAIndex     // 输出，检索结构
	doneOnce    bool              // 是否已经成功加载完一次
}

// NewIndexProcessor
func NewIndexProcessor(option *Option) *IndexProcessor {
	option.MaxWaitDuration = time.Duration(option.MaxWaitSecond) * time.Second
	p := &IndexProcessor{
		option: option,
		loader: loader.CreatePBFileLoader(option.RTATargetingPath, func() interface{} {
			return nil
		}),
	}

	return p
}

// Option 检索处理选项
type Option struct {
	RTATargetingPath string `json:"rta_targeting_path"`
	OutputPath       string `json:"output_path"`
	MaxWaitSecond    int    `json:"max_wait_second"`
	MaxWaitDuration  time.Duration
}

// Process 处理函数
func (p *IndexProcessor) Process() {
	for {
		if err := p.process(); err != nil {
			log.Errorf("failed to process index: %v", err)
		}

		if !p.doneOnce {
			log.Error("failed to load for the first time")
		}

		if p.option.MaxWaitDuration > 0 && time.Now().Sub(p.lastDone) > p.option.MaxWaitDuration {
			log.Warnf("not generate index longer than %s", p.option.MaxWaitDuration.String())
		}
		time.Sleep(time.Second)
	}
}

func (p *IndexProcessor) process() error {
	filepath, detectNew := p.detectNewFile()
	if !detectNew {
		return nil
	}
	log.Info("start process index")

	p.clear()

	var err error
	if err = p.loader.Load(filepath, p.FullRTAInfo); err != nil {
		return err
	}

	if err = p.buildIndex(); err != nil {
		return err
	}

	if log.GetLevel() >= log.DebugLevel {
		log.Debug(proto.MarshalTextString(p.RTAIndex))
	}

	if err = p.sinkToFile(); err != nil {
		return err
	}

	p.done()

	log.Info("process index done")
	return nil
}

func (p *IndexProcessor) detectNewFile() (string, bool) {
	return p.loader.DetectNewFile()
}

func (p *IndexProcessor) clear() {
	p.FullRTAInfo = &rta.FullRTAInfo{}
	p.RTAIndex = &rta.RTAIndex{}
}

func (p *IndexProcessor) buildIndex() error {
	// 倒排信息
	targetingVec := make([]*targeting.Targeting, 0, len(p.FullRTAInfo.RtaInfo))
	for _, info := range p.FullRTAInfo.RtaInfo {
		targetingVec = append(targetingVec, &targeting.Targeting{
			Id:         info.Id,
			Betree:     info.Betree,
			Offline:    false,
			UpdateTime: info.UpdateTime,
		})
	}
	var index targeting.TargetingIndex

	builder := matchengine.NewTargetingIndexBuilder(&index)
	builder.Build(targetingVec)
	p.RTAIndex.Index = &index

	// 正排信息
	p.RTAIndex.RtaTargeting = p.FullRTAInfo.RtaInfo
	p.RTAIndex.DumpTime = p.FullRTAInfo.DumpTime
	p.RTAIndex.Version = p.FullRTAInfo.Version

	return nil
}

func (p *IndexProcessor) sinkToFile() error {
	return dumper.DumpMessageToFile(p.RTAIndex, p.option.OutputPath)
}

func (p *IndexProcessor) done() {
	p.lastDone = time.Now()
	p.doneOnce = true
}
