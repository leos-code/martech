package retrieval

import (
	"context"
	"fmt"
	"net"

	pb "github.com/tencentad/martech/api/proto/retrieval"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	serviceSingleton *ServiceImpl
)

// ServiceImplOption
type ServiceImplOption struct {
	SearcherManager *SearcherManagerOption `json:"searcher_manager"`
	ListenAddr      string                 `json:"listen_addr"`     // 参数为空时，不启动grpc服务
	QueueSize       int                    `json:"queue_size"`      // 队列长度
	ProcessorCount  int                    `json:"processor_count"` // 处理线程数
}

// InitServiceImpl 初始化retrieval服务实现
func InitServiceImpl(option *ServiceImplOption) error {
	sm := getSearcherManager(option.SearcherManager)
	if sm == nil {
		return fmt.Errorf("searcher manager is nil")
	}
	serviceSingleton = newServiceImpl(option)
	return nil
}

// GetServiceImpl
func GetServiceImpl() *ServiceImpl {
	return serviceSingleton
}

// ServiceImpl retrieval实现
type ServiceImpl struct {
	option          *ServiceImplOption
	searcherManager *searcherManager
	processors      []*processor
	queue           chan *reqWrapper
}

func newServiceImpl(option *ServiceImplOption) *ServiceImpl {
	srv := &ServiceImpl{
		searcherManager: getSearcherManager(),
		queue:           make(chan *reqWrapper, option.QueueSize),
		processors:      make([]*processor, option.ProcessorCount),
		option:          option,
	}

	return srv
}

// Start 启动服务
func (srv *ServiceImpl) Start() error {
	if err := srv.searcherManager.waitLoadDone(); err != nil {
		return err
	}

	for i := 0; i < srv.option.ProcessorCount; i++ {
		srv.processors[i] = newProcessor(srv.queue)
		go srv.processors[i].run()
	}

	if err := srv.serve(srv.option); err != nil {
		return err
	}
	log.Info("retrieval service start done")

	return nil
}

func (srv *ServiceImpl) serve(option *ServiceImplOption) error {
	if option.ListenAddr == "" {
		return nil
	}

	server := grpc.NewServer()
	pb.RegisterRetrievalServiceServer(server, srv)

	listener, err := net.Listen("tcp", option.ListenAddr)
	if err != nil {
		return err
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to serve retrieval server, err: %v", err)
		}
	}()

	return nil
}

// Retrieve 实现检索
func (srv *ServiceImpl) Retrieve(_ context.Context, req *pb.RetrievalRequest) (*pb.RetrievalResponse, error) {
	reqWrapper := &reqWrapper{
		req:      req,
		respChan: make(chan *pb.RetrievalResponse, 1),
	}

	select {
	case srv.queue <- reqWrapper:
		resp := <-reqWrapper.respChan
		return resp, nil
	default:
		resp := &pb.RetrievalResponse{
			RetCode: 1,
		}
		log.Errorf("too busy")
		return resp, nil
	}
}

// QueryInfo 实现信息查询
func (srv *ServiceImpl) QueryInfo(_ context.Context, req *pb.QueryInfoRequest) (*pb.QueryInfoResponse, error) {
	searcher := srv.searcherManager.latestSearcher()
	resp := &pb.QueryInfoResponse{}
	switch req.Cmd {
	case pb.QueryInfoRequest_GetTargetingID:
		idIndexVec := searcher.se.IdIndexVec()
		for _, idIndex := range idIndexVec {
			resp.TargetingId = append(resp.TargetingId, idIndex.TargetingId)
		}

	case pb.QueryInfoRequest_GetTargeting:
		tid := req.TargetingId
		localId, ok := searcher.se.GetLocalId(tid)
		if ok {
			resp.RtaInfo = searcher.GetRTAInfo(localId)
		}
	}
	return resp, nil
}
