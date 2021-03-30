package retrieval

import (
	pb "github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/pkg/matchengine"
)

type processor struct {
	sc              *matchengine.SearchContext
	queue           chan *reqWrapper
	searcherManager *searcherManager
}

type reqWrapper struct {
	req      *pb.RetrievalRequest
	respChan chan *pb.RetrievalResponse
}

func newProcessor(queue chan *reqWrapper) *processor {
	return &processor{
		queue:           queue,
		searcherManager: getSearcherManager(),
		sc:              matchengine.NewSearchContext(),
	}
}

func (p *processor) run() {
	for {
		reqWrapper := <-p.queue
		p.sc.Clear()
		resp := p.searcherManager.latestSearcher().search(p.sc, reqWrapper.req)
		reqWrapper.respChan <- resp
	}
}
