package updater

import (
	"github.com/sirupsen/logrus"
	"time"
)

// TableUpdateHelperManager 状态更新任务统一管理
type TableUpdateHelperManager struct {
	HelperList []*TableUpdateHelper
}

// NewTableUpdateHelperManager 构造函数
func NewTableUpdateHelperManager() *TableUpdateHelperManager {
	return &TableUpdateHelperManager{
		HelperList: make([]*TableUpdateHelper, 0),
	}
}

// Run 初始化状态并启动任务
func (m *TableUpdateHelperManager) Run() error {
	for _, helper := range m.HelperList {
		if err := helper.Init(); err != nil {
			return err
		}
	}
	go m.updateRoutine()
	return nil
}

// updateRoutine 启动所有helper的轮询任务
func (m *TableUpdateHelperManager) updateRoutine() {
	for {
		for _, helper := range m.HelperList {
			if err := helper.Check(); err != nil {
				logrus.Errorf("failed to inc sync table, err: %v", err)
			}
			time.Sleep(time.Second)
		}
	}
}

// AppendTableUpdateHelper 添加一个状态更新任务
func (m *TableUpdateHelperManager) AppendTableUpdateHelper(helper ...*TableUpdateHelper) {
	m.HelperList = append(m.HelperList, helper...)
}
