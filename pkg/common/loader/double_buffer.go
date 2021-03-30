package loader

import (
	"sync"
	"sync/atomic"
	"time"
)

var m FileLoadManager

// FileLoader 文件加载器
type FileLoader interface {
	// DetectNewFile 判断是否有新文件
	DetectNewFile() (string, bool)
	// Load 加载文件到结构体
	Load(filePath string, i interface{}) error
	// Alloc 申请结构体内存
	Alloc() interface{}
}

// allocFunc 类结构构造函数
type allocFunc func() interface{}

// FileDoubleBuffer 双buffer，热加载
type FileDoubleBuffer struct {
	Loader     FileLoader
	bufferData []interface{}
	curIndex   int32
	notify     func(error)
}

// NewFileDoubleBuffer 创建双buffer
func NewFileDoubleBuffer(loader FileLoader) *FileDoubleBuffer {
	b := &FileDoubleBuffer{
		Loader:   loader,
		curIndex: 0,
	}
	b.bufferData = append(b.bufferData, loader.Alloc(), loader.Alloc())
	register(b)
	return b
}

// SetNotify 设置提示函数，在发现新文件时获取触发执行
func (b *FileDoubleBuffer) SetNotify(fn func(error)) {
	b.notify = fn
}

func (b *FileDoubleBuffer) load() {
	file, newFound := b.Loader.DetectNewFile()
	if newFound {
		ci := 1 - atomic.LoadInt32(&b.curIndex)
		err := b.Loader.Load(file, b.bufferData[ci])
		if err == nil {
			atomic.StoreInt32(&b.curIndex, ci)
		}

		if b.notify != nil {
			b.notify(err)
		}
	}
}

// Data 获取数据
func (b *FileDoubleBuffer) Data() interface{} {
	ci := atomic.LoadInt32(&b.curIndex)
	return b.bufferData[ci]
}

// FileLoadManager 管理所有的双buffer
type FileLoadManager struct {
	doubleBuffers []*FileDoubleBuffer
	mutex         sync.Mutex
}

func (m *FileLoadManager) load() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, b := range m.doubleBuffers {
		b.load()
	}
}

func (m *FileLoadManager) reload(reloadInterval int) {
	for {
		time.Sleep(time.Duration(reloadInterval) * time.Second)
		m.load()
	}
}

func register(doubleBuffer *FileDoubleBuffer) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.doubleBuffers = append(m.doubleBuffers, doubleBuffer)
}

// StartDoubleBufferLoad 开启后台热更新
func StartDoubleBufferLoad(reloadInterval int) {
	m.load()
	go m.reload(reloadInterval)
}
