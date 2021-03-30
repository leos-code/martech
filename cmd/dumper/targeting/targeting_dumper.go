package targeting

import (
	"fmt"
	"time"

	"github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/common/dumper"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// Dumper 导出所有的策略定向信息
type Dumper struct {
	option       *DumperOption
	lastDumpTime time.Time
	FullRTAInfo  *rta.FullRTAInfo
}

// NewDumper 创建定向dumper
func NewDumper(opt *DumperOption) *Dumper {
	opt.interval = time.Duration(opt.IntervalSecond) * time.Second
	return &Dumper{
		option: opt,
	}
}

// DumperOption 定向dumper参数选项
type DumperOption struct {
	OutputPath     string `json:"output_path"`
	IntervalSecond int    `json:"interval_second"`
	interval       time.Duration
}

// Dump 从DB中导出定向数据
func (d *Dumper) Dump() {
	for {
		if time.Now().After(d.lastDumpTime.Add(d.option.interval)) {
			if err := d.dump(); err != nil {
				glog.Errorf("failed to dump, err: %v", err)
				time.Sleep(time.Second * 5)
			}
		}
		time.Sleep(time.Second)
	}
}

func (d *Dumper) dump() error {
	d.clear()
	now := time.Now()
	if err := d.getFromDB(); err != nil {
		return err
	}

	d.FullRTAInfo.DumpTime = now.Unix()
	d.FullRTAInfo.Version = now.Unix()
	if err := dumper.DumpMessageToFile(d.FullRTAInfo, d.option.OutputPath); err != nil {
		return err
	}

	if log.GetLevel() >= log.DebugLevel {
		log.Debug(proto.MarshalTextString(d.FullRTAInfo))
	}
	d.lastDumpTime = now
	return nil
}

func (d *Dumper) getFromDB() error {
	db := orm.GetDB()

	allTargeting, err := orm.GetAllTargeting(db, nil)
	if err != nil {
		return err
	}
	for _, dbT := range allTargeting {
		if err = orm.LoadTargetingBindStrategy(db, dbT); err != nil {
			return err
		}
		pbT, err := convertDBTargeting(dbT)
		if err != nil {
			return err
		}
		d.FullRTAInfo.RtaInfo = append(d.FullRTAInfo.RtaInfo, pbT)
	}

	return nil
}

func (d *Dumper) clear() {
	d.FullRTAInfo = &rta.FullRTAInfo{}
}

func convertDBTargeting(dbT *types.Targeting) (*rta.RTAInfo, error) {
	pbT := &rta.RTAInfo{
		Id:         dbT.ID,
		Name:       dbT.Name,
		UpdateTime: dbT.UpdatedAt.Unix(),
	}

	betree, err := convertDBTargetingInfo(dbT.TargetingInfo)
	if err != nil {
		return nil, err
	}
	pbT.Betree = betree

	for _, bs := range dbT.BindStrategy {
		pbBS, err := convertDBBindStrategy(bs)
		if err != nil {
			return nil, err
		}

		pbT.BindStrategy = append(pbT.BindStrategy, pbBS)
	}

	return pbT, nil

}

func convertDBTargetingInfo(dbInfos types.TargetingInfos) (*targeting.BETree, error) {
	betree := &targeting.BETree{
		Op: targeting.LogicalOp_And,
	}

	for _, info := range dbInfos {
		pbVs, err := convertDBTargetingValue(info.Values)
		if err != nil {
			return nil, err
		}

		sub := &targeting.BETree{
			Predicate: &targeting.Predicate{
				Not:   info.Not,
				Field: info.Name,
				Value: pbVs,
			},
		}
		betree.Betree = append(betree.Betree, sub)
	}
	return betree, nil
}

func convertDBTargetingValue(dbV *types.TargetingValue) ([]*targeting.Predicate_Value, error) {
	 pbVs := make([]*targeting.Predicate_Value, 0)
	switch dbV.Type {
	case types.TargetingValueTypeString:
		for _, str := range dbV.String {
			pbVs = append(pbVs, &targeting.Predicate_Value{
				Type: targeting.Predicate_Value_String,
				Str:  str,
			})
		}
	case types.TargetingValueTypeRange:
		for _, r := range dbV.Range {
			pbVs = append(pbVs, &targeting.Predicate_Value{
				Type: targeting.Predicate_Value_RANGE,
				Range: &targeting.Predicate_Value_Range{
					Begin: r.Begin,
					End:   r.End,
				},
			})
		}
	default:
		return nil, fmt.Errorf("not valid value type[%s]", dbV.Type)
	}

	return pbVs, nil
}

func convertDBBindStrategy(dbBindStrategy *types.BindStrategy) (*rta.BindStrategy, error) {
	pbStrategy := &rta.BindStrategy{
		Id:       dbBindStrategy.ID,
		Name:     dbBindStrategy.Name,
		Platform: string(dbBindStrategy.Platform),
	}

	v, err := dbBindStrategy.Strategy.Value()
	if err != nil {
		return nil, err
	}
	vStr, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("'Strategy'[%v] not string type", v)
	}
	pbStrategy.Strategy = vStr
	return pbStrategy, nil
}
