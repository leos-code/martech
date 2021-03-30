package retrieval

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/tencentad/martech/api/proto/rta"
	"github.com/tencentad/martech/pkg/common/loader"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	searcherManagerOnce      sync.Once
	searcherManagerSingleton *searcherManager
)

// searcherManager 检索管理器
type searcherManager struct {
	rtaSearcher *loader.FileDoubleBuffer
	incUpdater  *incUpdater
	mutex       sync.Mutex
	option      *SearcherManagerOption

	firstRound   bool
	loadDone     bool
	loadDoneChan chan struct{}
}

// SearcherManagerOption 索引管理器选项
type SearcherManagerOption struct {
	IndexPath      string `json:"index_path"`
	WaitLoadSecond int    `json:"wait_load_second"`
}

// getSearcherManager 获取检索管理器
func getSearcherManager(options ...*SearcherManagerOption) *searcherManager {
	searcherManagerOnce.Do(
		func() {
			var option *SearcherManagerOption
			if len(options) != 0 {
				option = options[0]
			}

			if option == nil {
				log.Errorf("call getSearcherManager without option for the first time")
				return
			}

			searcherManagerSingleton = newSearcherManager(option)
		})

	return searcherManagerSingleton
}

func newSearcherManager(option *SearcherManagerOption) *searcherManager {
	m := &searcherManager{
		incUpdater:   newIncUpdater(),
		option:       option,
		loadDoneChan: make(chan struct{}, 1),
		firstRound:   true,
		loadDone:     false,
	}

	m.rtaSearcher = loader.NewFileDoubleBuffer(loader.CreateFileLoader(option.IndexPath, func() interface{} {
		return &rtaSearcher{}
	}, m.load))

	m.rtaSearcher.SetNotify(func(err error) {
		if err != nil {
			log.Errorf("failed to load rta index, err: %v", err)
		}
	})

	go func() {
		for {
			time.Sleep(time.Second * 2)
			if !m.loadDone {
				continue
			}

			m.mutex.Lock()
			if err := m.incUpdater.updateFromDB(m.latestSearcher()); err != nil {
				log.Errorf("failed to update from db, err: %v", err)
			}
			m.mutex.Unlock()
		}

	}()

	return m
}

func (m *searcherManager) latestSearcher() *rtaSearcher {
	return m.rtaSearcher.Data().(*rtaSearcher)
}

func (m *searcherManager) load(reader io.Reader, i interface{}) error {
	current := i.(*rtaSearcher)

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	var rtaIndex rta.RTAIndex
	if err = proto.Unmarshal(content, &rtaIndex); err != nil {
		return err
	}
	current.init(&rtaIndex)

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.incUpdater.updateFromCache(current)
	if err = m.incUpdater.updateFromDB(current); err != nil {
		log.Errorf("failed to update form db, err: %v", err)
		return err
	}

	if m.firstRound {
		m.loadDoneChan <- struct{}{}
		m.firstRound = false
		m.loadDone = true
	}

	return nil
}

func (m *searcherManager) waitLoadDone() error {
	tm := time.NewTimer(time.Duration(m.option.WaitLoadSecond) * time.Second)

	select {
	case <-m.loadDoneChan:
		return nil
	case <-tm.C:
		return fmt.Errorf("load data timeout")
	}
}
