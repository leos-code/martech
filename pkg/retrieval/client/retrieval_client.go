package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	pb "github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/api/proto/targeting"
	"github.com/tencentad/martech/pkg/common/tool/press"
	"github.com/tencentad/martech/pkg/matchengine/tool"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Cmd string

const (
	GetTargetingID Cmd = "get_targeting_id"
	GetTargeting   Cmd = "get_targeting"
	Search         Cmd = "search"
	Match          Cmd = "match"
	Press          Cmd = "press"
)

var supportedCmd = "get_targeting_id|get_targeting|search|match|press"

var (
	address       = flag.String("address", ":9528", "")
	cmd           = flag.String("cmd", "search", supportedCmd)
	tid           = flag.Uint64("targeting_id", 0, "targeting id")
	reqFile       = flag.String("req_file", "", "")
	timeoutSecond = flag.Int("timeout_second", 5, "")

	qps        = flag.Int("qps", 500, "")
	agentCount = flag.Int("agent_count", 10, "")
	number     = flag.Int("number", 1000, "")

	req *pb.RetrievalRequest
)

func initRequest() error {
	if *reqFile == "" {
		return nil
	}

	req = &pb.RetrievalRequest{}
	text, err := ioutil.ReadFile(*reqFile)
	if err != nil {
		log.Errorf("read file err:%s", err.Error())
		return err
	}
	err = json.Unmarshal(text, req)
	if err != nil {
		log.Errorf("parse pb err:%s", err.Error())
		return err
	}

	return nil
}

type RetrievalPressAgent struct {
	client pb.RetrievalServiceClient
}

func RetrievalTestAgent(client pb.RetrievalServiceClient) press.Agent {
	return &RetrievalPressAgent{client: client}
}

func (agent *RetrievalPressAgent) Execute() error {
	resp, err := agent.client.Retrieve(context.Background(), req)

	if err != nil {
		log.Errorf("failed to request service: %s", err.Error())
		return err
	}

	if len(resp.HitRtaInfo) != 1 || resp.HitRtaInfo[0].Id != 1 {
		js, _ := json.Marshal(resp)
		log.Info(string(js))
		return fmt.Errorf("response not correct")
	}

	return nil
}

func main() {
	flag.Parse()

	if err := initRequest(); err != nil {
		log.Fatal(err)
	}

	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Second * time.Duration(*timeoutSecond))
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, *address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("init client err:%s", err.Error())
	}

	client := pb.NewRetrievalServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * time.Duration(*timeoutSecond))
	defer cancel()

	switch Cmd(*cmd) {
	case Search:
		resp, err := client.Retrieve(context.Background(), req)
		if err != nil {
			log.Fatalf("search err: %v", err)
		}
		log.Infof("%s", proto.MarshalTextString(resp))
	case GetTargetingID:
		resp, err := client.QueryInfo(ctx, &pb.QueryInfoRequest{
			Cmd: pb.QueryInfoRequest_GetTargetingID,
		})
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		log.Infof("%s", proto.MarshalTextString(resp))
	case GetTargeting:
		resp, err := client.QueryInfo(ctx, &pb.QueryInfoRequest{
			TargetingId: *tid,
			Cmd:         pb.QueryInfoRequest_GetTargeting,
		})
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		if resp.RtaInfo == nil {
			log.Fatalf("rta id[%d] not in index", *tid)
		}
		log.Infof(tool.BETreeDebugString(resp.RtaInfo.Betree))
	case Match:
		resp, err := client.QueryInfo(ctx, &pb.QueryInfoRequest{
			TargetingId: *tid,
			Cmd:         pb.QueryInfoRequest_GetTargeting,
		})
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		if resp.RtaInfo == nil {
			log.Fatalf("targeting id[%d] not in index", *tid)
		}

		log.Infof(tool.NewMatchDebug().MatchDebug(&targeting.Targeting{
			Betree: resp.RtaInfo.Betree,
		}, req))
	case Press:
		controller := press.NewController(func() press.Agent {
			return RetrievalTestAgent(client)
		}, &press.Option{
			AgentCount: *agentCount,
			QPS:        *qps,
			Number:     *number,
		})
		controller.Start()
	default:
		log.Fatalf("unsupported cmd [%s], only support %s", *cmd, supportedCmd)
	}

}
