package workflow

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type add struct {
	i int32
}

// Run
func (r *add) Run(i interface{}) {
	atomic.AddInt32(i.(*int32), r.i)
}

type multiply struct {
	i int32
}

// Run
func (r *multiply) Run(i interface{}) {
	//time.Sleep(time.Second)
	*(i.(*int32)) = *(i.(*int32)) * r.i
}

type printer struct {
}

// Run
func (p *printer) Run(_ interface{}) {
	//fmt.Println(*i.(*int32))
}

func TestWorkFlow(t *testing.T) {
	wf := buildWorkflow()

	isDAG := wf.CheckDAG()
	assert.True(t, isDAG)

	if !wf.CheckDAG() {
		t.Errorf("check DAG failed")
	}

	var i int32 = 10
	wf.Start(&i)
	wf.WaitDone()
	assert.EqualValues(t, 42, i)
}

// BenchmarkTaskNode_ExecuteTask
func BenchmarkTaskNode_ExecuteTask(b *testing.B) {
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			wf := buildWorkflow()
			var p int32 = 10
			wf.Start(&p)
			wf.WaitDone()
			wg.Done()
		}()
	}
	wg.Wait()
}

// BenchmarkTaskNode_SubmitTask
func BenchmarkTaskNode_SubmitTask(b *testing.B) {
	//b.ResetTimer()
	queue := NewDefaultJobQueue(&QueueOption{
		WorkerCount: 50,
		QueueSize:   1000,
		PushTimeout: time.Second * 3,
	})
	queue.Start()
	defer queue.Stop()
	for i := 0; i < b.N; i++ {
		wf := buildWorkflow()

		var p int32 = 10
		wf.StartWithJobQueue(queue, context.Background(), &p)
		wf.WaitDone()
	}
}

func TestWorkFlowWithContext(t *testing.T) {
	var i int32 = 10
	var wf *WorkFlow
	var ctx context.Context
	var cancel context.CancelFunc

	wf = buildWorkflow()
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	time.Sleep(time.Millisecond * 60)
	wf.StartWithContext(ctx, &i)
	wf.WaitDone()
	assert.True(t, ctx.Err() == context.DeadlineExceeded)
	assert.Error(t, ctx.Err())
	cancel()

	wf = buildWorkflow()
	ctx, cancel = context.WithCancel(context.Background())
	wf.StartWithContext(ctx, &i)
	cancel()
	wf.WaitDone()
	assert.True(t, ctx.Err() == context.Canceled)
	assert.Error(t, ctx.Err())
	cancel()
}

func buildWorkflow() *WorkFlow {
	add1 := &add{i: 1}
	add2 := &add{i: 2}
	multiply := &multiply{i: 3}

	add1Task := &TaskNode{
		Task: add1,
	}

	add2Task := &TaskNode{
		Task: add2,
	}

	multiplyTask := &TaskNode{
		Task: multiply,
	}

	printTask := &TaskNode{
		Task: &printer{},
	}
	add1Task2 := NewTaskNode(add1)
	add2Task2 := NewTaskNode(add2)

	// 实现 (i + 1 + 2) * 3 + 1 + 2，然后打印
	wf := NewWorkFlow()
	wf.AddTaskNode(add1Task)
	wf.AddTaskNode(add2Task)
	wf.AddTaskNode(multiplyTask, add1Task, add2Task)
	wf.AddTaskNode(add1Task2, multiplyTask)
	wf.AddTaskNode(add2Task2, multiplyTask)
	wf.AddTaskNode(printTask, add1Task2, add2Task2)

	return wf
}
