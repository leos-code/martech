package workflow

import (
	"fmt"
	"sync"
	"time"
)

type job func()

var (
	ErrPushQueueTimeout = fmt.Errorf("push queue timeout")
)

// JobQueue 任务队列接口
type JobQueue interface {
	PushJob(j job) error
}

// DefaultJobQueue 默认任务队列
type DefaultJobQueue struct {
	queue  chan job
	once   sync.Once
	option *QueueOption
}

// QueueOption 队列参数
type QueueOption struct {
	WorkerCount int
	QueueSize   int
	PushTimeout time.Duration
}

// NewDefaultJobQueue
func NewDefaultJobQueue(option *QueueOption) *DefaultJobQueue {
	queue := &DefaultJobQueue{
		queue:  make(chan job, option.QueueSize),
		option: option,
	}

	return queue
}

// Start 启动处理队列数据的worker
func (queue *DefaultJobQueue) Start() {
	for i := 0; i < queue.option.WorkerCount; i++ {
		go func() {
			for job := range queue.queue {
				job()
			}
		}()
	}
}

// PushJob 添加任务
func (queue *DefaultJobQueue) PushJob(j job) error {
	timer := time.NewTimer(queue.option.PushTimeout)

	select {
	case queue.queue <- j:
		return nil
	case <-timer.C:
		return ErrPushQueueTimeout
	}
}

// Stop 终止队列
func (queue *DefaultJobQueue) Stop() {
	queue.once.Do(func() {
		close(queue.queue)
	})
}
