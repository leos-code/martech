package data

import (
	"context"
)

// BaseContext 基础context，包含错误和上下文
type BaseContext struct {
	Ctx    context.Context
	cancel context.CancelFunc
	Error  error
}

// NewBaseContext 创建基础context
func NewBaseContext() *BaseContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseContext{
		Ctx:    ctx,
		cancel: cancel,
	}
}

// StopWithError 出现错误时终止context
func (c *BaseContext) StopWithError(err error) {
	c.Error = err
	c.cancel()
}
