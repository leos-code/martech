package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/pkg/common/dumper"
	"github.com/tencentad/martech/pkg/schema"
	log "github.com/sirupsen/logrus"
)

var (
	schemaPath = flag.String("schema_path", "", "")
	number     = flag.Int("number", 1, "")
	outputPath = flag.String("output_path", "", "")

	s schema.Schema
)

func initSchemaHelper() error {
	content, err := ioutil.ReadFile(*schemaPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, &s); err != nil {
		return err
	}

	return err
}

func generateFullRTAInfo() *rta.FullRTAInfo {
	info := &rta.FullRTAInfo{}
	for i := 0; i < *number; i++ {
		info.RtaInfo = append(info.RtaInfo, &rta.RTAInfo{
			Id:     uint64(i + 1),
			Betree: generateBETree(),
		})
	}

	//log.Infof(proto.MarshalTextString(info))

	return info
}

func generateBETree() *targeting.BETree {
	n := len(s.Fields)

	t := &targeting.BETree{
		Op: targeting.LogicalOp_And,
	}
	for i := 0; i < n; i++ {
		if rand.Float64() >= 0.5 {
			t.Betree = append(t.Betree, &targeting.BETree{
				Predicate: generatePredicate(s.Fields[i]),
			})
		}
	}

	return t
}

func generatePredicate(f *schema.Field) *targeting.Predicate {
	return &targeting.Predicate{
		Not:   rand.Float64() > 0.9,
		Value: []*targeting.Predicate_Value{generateValue(f)},
		Field: f.Name,
	}
}

const rangeMax uint64 = 50
func generateValue(f *schema.Field) *targeting.Predicate_Value {
	switch f.Type {
	case schema.FieldTypeInteger:
		r := rand.Uint64() % 100000
		if r < rangeMax {
			r = r + rangeMax
		}
		return &targeting.Predicate_Value{
			Type: targeting.Predicate_Value_RANGE,
			Range: &targeting.Predicate_Value_Range{
				Begin: r - rangeMax,
				End:   r - rangeMax + 1 + uint64(rand.Int31n(int32(rangeMax))),
			},
		}
	case schema.FieldTypeString:
		r := rand.Uint64() % 100000
		return &targeting.Predicate_Value{
			Type: targeting.Predicate_Value_String,
			Str:  strconv.FormatUint(r, 10),
		}
	case schema.FieldTypeEnum:
		enums := f.GetAllEnum()
		return &targeting.Predicate_Value{
			Type: targeting.Predicate_Value_String,
			Str:  enums[rand.Int31n(int32(len(enums)))],
		}
	default:
		return nil
	}
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().Unix())

	if err := initSchemaHelper(); err != nil {
		log.Fatalf("failed to init schema helper, err: %v", err)
	}

	info := generateFullRTAInfo()
	if err := dumper.DumpMessageToFile(info, *outputPath); err != nil {
		log.Fatalf("failed to dump to file: %v", err)
	}
}
