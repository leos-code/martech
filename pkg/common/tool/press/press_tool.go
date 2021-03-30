package press

import (
	"fmt"
	"sync"
	"time"
)

// Agent 压测Agent，需要在SendRequest中完成所有操作
type Agent interface {
	Execute() error
}

// NewAgent 压测Agent构造函数
type NewAgent func() Agent

// Controller
type Controller struct {
	agents      []Agent
	agentCnt    int
	qps         int
	number      int
	sendChan    chan struct{}
	latencyChan chan float64
	errorChan   chan struct{}
	doneChan    chan struct{}
	sendWG      sync.WaitGroup

	totalLatency  float64
	successCnt    int64
	errorCnt      int64
	statisticDone chan struct{}
}

type Option struct {
	AgentCount int
	QPS        int
	Number     int
}

var (
	defaultQPS    = 500
	defaultAgent  = 10
	defaultNumber = 1000
)

func NewController(newAgent NewAgent, option *Option) *Controller {
	controller := &Controller{
		agentCnt:      defaultAgent,
		qps:           defaultQPS,
		number:        defaultNumber,
		sendChan:      make(chan struct{}),
		latencyChan:   make(chan float64),
		errorChan:     make(chan struct{}),
		doneChan:      make(chan struct{}),
		statisticDone: make(chan struct{}),
	}

	if option != nil {
		setIfNotZero(option.AgentCount, &controller.agentCnt)
		setIfNotZero(option.QPS, &controller.qps)
		setIfNotZero(option.Number, &controller.number)
	}

	for i := 0; i < controller.agentCnt; i++ {
		controller.agents = append(controller.agents, newAgent())
	}
	return controller
}

func setIfNotZero(src int, dst *int) {
	if src != 0 {
		*dst = src
	}
}

// 统计线程
func (c *Controller) statistic() {
	for {
		select {
		case <-c.doneChan:
			var latency float64
			if c.successCnt == 0 {
				latency = 0
			} else {
				latency = c.totalLatency / float64(c.successCnt)
			}
			fmt.Printf("success cnt[%d], error cnt[%d], average latency[%f]ms", c.successCnt,
				c.errorCnt, latency)
			c.statisticDone <- struct{}{}
		case latency, ok := <-c.latencyChan:
			if ok {
				c.successCnt += 1
				c.totalLatency += latency
			} else {
				fmt.Printf("latency channel closes before done")
				return
			}
		case _, ok := <-c.errorChan:
			if ok {
				c.errorCnt += 1
			} else {
				fmt.Printf("error channel closed before done")
				return
			}
		}
	}
}

// 实际发送请求线程
func (c *Controller) pressSend() {
	c.sendWG.Add(len(c.agents))
	for _, agent := range c.agents {
		go func(agent Agent) {
			for range c.sendChan {
				startTime := time.Now()
				err := agent.Execute()
				if err != nil {
					c.errorChan <- struct{}{}
				} else {
					c.latencyChan <- float64(time.Since(startTime)) / float64(time.Millisecond)
				}
			}
			c.sendWG.Done()
		}(agent)
	}
}

// 外部调用线程
func (c *Controller) Start() {
	c.pressSend()
	go c.statistic()

	var curCnt = 0
	startTime := time.Now()
	for curCnt < c.number {
		sleepTime := time.Duration(curCnt+1).Nanoseconds()*time.Second.Nanoseconds()/time.Duration(c.qps).Nanoseconds() - time.Since(startTime).Nanoseconds()
		if sleepTime >= 0 {
			time.Sleep(time.Duration(sleepTime))
		}
		c.sendChan <- struct{}{}
		curCnt += 1
	}
	timeUsed := time.Since(startTime)
	actualQps := float64(curCnt) / timeUsed.Seconds()
	fmt.Printf("acutual qps[%f], request cnt[%d], time used[%d]\n", actualQps, curCnt, timeUsed.Nanoseconds()/int64(time.Millisecond))
	close(c.sendChan)
	c.sendWG.Wait()

	c.doneChan <- struct{}{}
	<-c.statisticDone
}
