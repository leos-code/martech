// Package workflow 把业务流程抽象成dag
package workflow

import (
	"context"
	"flag"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tencentad/martech/pkg/common/metricutil"

	"github.com/golang/glog"
)

var (
	enableTaskTimeCostReport = flag.Bool("enable_task_time_cost_report", false, "")
)

// Runnable 可执行接口
type Runnable interface {
	Run(i interface{})
}

// TaskEdge 边
type TaskEdge struct {
	Prev *TaskNode
	Next *TaskNode
}

// TaskNode 点
type TaskNode struct {
	Dependency   []*TaskEdge
	DepCompleted int32
	Task         Runnable
	Children     []*TaskEdge
}

// NewTaskNode
func NewTaskNode(task Runnable) *TaskNode {
	return &TaskNode{
		Task: task,
	}
}

// NewNilTaskNode 用于等待前面所有的任务执行完成，然后同时执行后面的任务
func NewNilTaskNode() *TaskNode {
	return &TaskNode{
		Task: nil,
	}
}

// AddEdge 增加边
func AddEdge(prev *TaskNode, next *TaskNode) *TaskEdge {
	edge := &TaskEdge{
		Prev: prev,
		Next: next,
	}
	prev.Children = append(prev.Children, edge)
	next.Dependency = append(next.Dependency, edge)
	return edge
}

// ExecuteTask 开始执行，无法中止
func (n *TaskNode) ExecuteTask(i interface{}) {
	if n.dependencySatisfied() {
		if n.Task != nil {
			n.Task.Run(i)
		}
		if len(n.Children) >= 1 {
			for idx := 1; idx < len(n.Children); idx++ {
				go func(child *TaskEdge) {
					child.Next.ExecuteTask(i)
				}(n.Children[idx])
			}
			n.Children[0].Next.ExecuteTask(i)
		}
	}
}

// ExecuteTaskWithContext 开始执行，可以通过cancel context中止任务
func (n *TaskNode) ExecuteTaskWithContext(ctx context.Context, wf *WorkFlow, i interface{}) {
	if n.dependencySatisfied() {

		// 在程序已经出错，不需要继续执行的时候，直接退出
		if ctx.Err() != nil {
			wf.interruptDone()
			return
		}

		if n.Task != nil {
			var startTime time.Time
			if *enableTaskTimeCostReport {
				startTime = time.Now()
			}
			n.Task.Run(i)
			if *enableTaskTimeCostReport {
				TaskTimeCost.WithLabelValues(getTaskName(n.Task)).Observe(metricutil.CalcTimeUsedMicro(startTime))
			}
		}

		if len(n.Children) >= 1 {
			for idx := 1; idx < len(n.Children); idx++ {
				go func(child *TaskEdge) {
					child.Next.ExecuteTaskWithContext(ctx, wf, i)
				}(n.Children[idx])
			}
			n.Children[0].Next.ExecuteTaskWithContext(ctx, wf, i)
		}
	}
}

// SubmitTask 在任务队列中执行workflow
func (n *TaskNode) SubmitTask(ctx context.Context, wf *WorkFlow, i interface{}) {
	if !n.dependencySatisfied() {
		return
	}

	// 在程序已经出错，不需要继续执行的时候，直接退出
	if ctx.Err() != nil {
		wf.interruptDone()
		return
	}

	if wf.alreadyDone {
		return
	}

	job := func() {
		if n.Task != nil {
			var startTime time.Time
			if *enableTaskTimeCostReport {
				startTime = time.Now()
			}
			n.Task.Run(i)
			if *enableTaskTimeCostReport {
				TaskTimeCost.WithLabelValues(getTaskName(n.Task)).Observe(metricutil.CalcTimeUsedMicro(startTime))
			}
		}

		for _, child := range n.Children {
			child.Next.SubmitTask(ctx, wf, i)
		}
	}
	if err := wf.jobQueue.PushJob(job); err != nil {
		wf.interruptDone()
	}
}

func (n *TaskNode) dependencySatisfied() bool {
	return n.Dependency == nil || len(n.Dependency) == 1 ||
		atomic.AddInt32(&n.DepCompleted, 1) == int32(len(n.Dependency))
}

func getTaskName(task Runnable) string {
	return reflect.TypeOf(task).Elem().Name()
}

// WorkFlow
type WorkFlow struct {
	done        chan struct{}
	doneOnce    sync.Once
	alreadyDone bool

	root  *TaskNode
	End   *TaskNode
	edges []*TaskEdge

	jobQueue JobQueue
}

// NewWorkFlow
func NewWorkFlow() *WorkFlow {
	wf := &WorkFlow{
		root: NewNilTaskNode(),
		done: make(chan struct{}, 1), // 使用buffered channel防止完全同步执行时死锁
	}
	wf.End = NewTaskNode(&endWorkFlow{
		done: wf.done,
	})

	return wf
}

// Start 直接开始任务，所有的task必须都执行
func (wf *WorkFlow) Start(i interface{}) {
	wf.autoConnectToEnd()
	wf.root.ExecuteTask(i)
}

// StartWithContext 开始执行，如果ctx出现错误，中断workflow
func (wf *WorkFlow) StartWithContext(ctx context.Context, i interface{}) {
	wf.autoConnectToEnd()
	wf.root.ExecuteTaskWithContext(ctx, wf, i)
}

// StartWithJobQueue 将任务放到队列中执行，使用固定的goroutine执行task
func (wf *WorkFlow) StartWithJobQueue(jobQueue JobQueue, ctx context.Context, i interface{}) {
	wf.autoConnectToEnd()
	wf.jobQueue = jobQueue
	wf.root.SubmitTask(ctx, wf, i)
}

// WaitDone 等待任务结束
func (wf *WorkFlow) WaitDone() {
	<-wf.done
	close(wf.done)
}

// interruptDone 因为超时或者异常出错，中断workflow，节省资源
func (wf *WorkFlow) interruptDone() {
	wf.doneOnce.Do(
		func() {
			wf.done <- struct{}{}
			wf.alreadyDone = true
		},
	)
}

type endWorkFlow struct {
	done chan struct{}
}

// Run implement Runnable
func (end *endWorkFlow) Run(_ interface{}) {
	end.done <- struct{}{}
}

// AddTaskNode 添加节点，同时设置依赖
func (wf *WorkFlow) AddTaskNode(node *TaskNode, deps ...*TaskNode) {
	if len(deps) == 0 {
		wf.edges = append(wf.edges, AddEdge(wf.root, node))
	}

	for _, dep := range deps {
		wf.edges = append(wf.edges, AddEdge(dep, node))
	}
}

func (wf *WorkFlow) autoConnectToEnd() {
	endNode := make(map[*TaskNode]struct{})

	for _, edge := range wf.edges {
		if len(edge.Next.Children) == 0 {
			endNode[edge.Next] = struct{}{}
		}
	}

	for n := range endNode {
		wf.edges = append(wf.edges, AddEdge(n, wf.End))
	}
}

// CheckDAG 检查是否是正常的DAG, 用于检查逻辑,实际运行可以关闭
func (wf *WorkFlow) CheckDAG() bool {
	visitEdge := make(map[*TaskEdge]bool)
	var sortedNode []*TaskNode
	var rootNodes []*TaskNode
	rootNodes = append(rootNodes, wf.root)
	for len(rootNodes) != 0 {
		r := rootNodes[0]
		rootNodes = rootNodes[1:]
		sortedNode = append(sortedNode, r)
	OUTER:
		for _, edge := range r.Children {
			visitEdge[edge] = true
			m := edge.Next
			for _, edge := range m.Dependency {
				if _, ok := visitEdge[edge]; !ok {
					continue OUTER
				}
			}
			rootNodes = append(rootNodes, m)
		}
	}

	if glog.V(100) {
		for _, node := range sortedNode {
			glog.V(100).Infof("%T%+v\n", node.Task, node.Task)
		}
	}

	if len(visitEdge) != len(wf.edges) {
		return false
	} else {
		return true
	}
}
