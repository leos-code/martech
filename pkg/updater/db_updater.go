package updater

import (
	"time"

	"gorm.io/gorm"
)

// TableSnapshotSync 同步DB表接口
type TableSnapshotSync interface {
	// GetAll 获取DB全量的数据
	GetAll(db *gorm.DB) ([]Record, error)

	// GetInc 获取增量的数据
	GetInc(db *gorm.DB, updateFrom time.Time) ([]Record, error)
}

// Record DB表记录接口
type Record interface {
	// UpdateTime 记录更新时间
	GetUpdateTime() time.Time
	// ID 记录ID
	GetID() interface{}
	// Update 更新记录，需要在Update内部进行检查判断是否需要更新
	Update(db *gorm.DB) error
}

// TableUpdateHelper 更新表数据Helper
type TableUpdateHelper struct {
	TableSync TableSnapshotSync

	db           *gorm.DB
	m            map[interface{}]Record
	latestModify time.Time
}

// NewTableUpdateHelper
func NewTableUpdateHelper(sync TableSnapshotSync, db *gorm.DB) *TableUpdateHelper {
	return &TableUpdateHelper{
		TableSync: sync,
		m:         make(map[interface{}]Record),
		db:        db,
	}
}

// Init 初始化
func (h *TableUpdateHelper) Init() error {
	all, err := h.TableSync.GetAll(h.db)
	if err != nil {
		return err
	}

	for _, r := range all {
		h.m[r.GetID()] = r
		h.updateLatestModify(r)
	}

	return nil
}

// Check 检查所有的表记录，并进行更新
func (h *TableUpdateHelper) Check() error {
	inc, err := h.TableSync.GetInc(h.db, h.latestModify)
	if err != nil {
		return err
	}

	incLastModify := h.latestModify
	for _, r := range inc {
		h.m[r.GetID()] = r
		if r.GetUpdateTime().After(incLastModify) {
			incLastModify = r.GetUpdateTime()
		}
	}

	if incLastModify == h.latestModify {
		h.latestModify.Add(time.Second)
	} else {
		h.latestModify = incLastModify
	}

	for _, r := range h.m {
		if err = r.Update(h.db); err != nil {
			return err
		}
	}

	return nil
}

func (h *TableUpdateHelper) updateLatestModify(r Record) {
	if r.GetUpdateTime().After(h.latestModify) {
		h.latestModify = r.GetUpdateTime()
	}
}
