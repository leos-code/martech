package orm

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	NoPageError = errors.New("not page")
)

// StatementOption 查询选项
type StatementOption struct {
	Pagination *PaginationOption
	Filter     []*FilterOption
}

// PaginationOption 分页选项
type PaginationOption struct {
	Page     int // 第几页，从第1页开始
	PageSize int
}

// pagination 分页功能
func pagination(opt *PaginationOption) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if opt == nil {
			return db
		}

		page := opt.Page
		pageSize := opt.PageSize

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func getPaginationOption(option *StatementOption) *PaginationOption {
	if option == nil {
		return nil
	}

	return option.Pagination
}

// FilterOption 过滤选项
type FilterOption struct {
	Field     string
	Operation FilterOperationType
	Value     interface{}
}

// FilterOperationType 过滤操作类型
type FilterOperationType string

const (
	FilterOperationEqual = "="
	FilterOperationIn    = "in"
	FilterOperationLike  = "like"
)

func (opt *FilterOption) cond() string {
	if opt.Operation == "" {
		opt.Operation = FilterOperationEqual
	}

	return fmt.Sprintf("%s %s ?", opt.Field, opt.Operation)
}

func (opt *FilterOption) value() interface{} {
	if strings.ToLower(string(opt.Operation)) == FilterOperationLike {
		return fmt.Sprintf("%%%v%%", opt.Value)
	}
	return opt.Value
}

func filter(opts ...*FilterOption) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if opts == nil {
			return db
		}

		for _, opt := range opts {
			db = db.Where(opt.cond(), opt.value())
		}

		return db
	}
}

func getFilterOption(option *StatementOption) []*FilterOption {
	if option == nil {
		return nil
	}
	return option.Filter
}

func getDBWithOption(db *gorm.DB, option *StatementOption) *gorm.DB {
	db = getDBWithFilterOption(db, option)
	db = getDBWithPaginationOption(db, option)
	return db
}

func getDBWithPaginationOption(db *gorm.DB, option *StatementOption) *gorm.DB {
	return db.Scopes(pagination(getPaginationOption(option)))
}

func getDBWithFilterOption(db *gorm.DB, option *StatementOption) *gorm.DB {
	return db.Scopes(filter(getFilterOption(option)...))
}
