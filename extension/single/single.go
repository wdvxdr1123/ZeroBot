package single

import (
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
)

var group sync.Map

// Option 配置项
type Option func(*Single)

// Single 反并发
type Single struct {
	key  func(ctx *zero.Ctx) interface{}
	post func(ctx *zero.Ctx)
}

// WithKeyFn 指定反并发的 Key
func WithKeyFn(fn func(ctx *zero.Ctx) interface{}) Option {
	return func(s *Single) {
		s.key = fn
	}
}

// WithPostFn 指定反并发拦截后的操作
func WithPostFn(fn func(ctx *zero.Ctx)) Option {
	return func(s *Single) {
		s.post = fn
	}
}

// New 创建反并发中间件
func New(op ...Option) *Single {
	s := Single{}
	for _, option := range op {
		option(&s)
	}
	return &s
}

// Apply 为指定 Engine 添加反并发功能
func (s *Single) Apply(engine *zero.Engine) {
	engine.UsePreHandler(func(ctx *zero.Ctx) bool {
		if s.key == nil {
			return true
		}
		key := s.key(ctx)
		if _, ok := group.Load(key); ok {
			if s.post != nil {
				defer s.post(ctx)
			}
			return false
		}
		group.Store(key, struct{}{})
		ctx.State["__single-key__"] = key
		return true
	})

	engine.UsePostHandler(func(ctx *zero.Ctx) {
		group.Delete(ctx.State["__single-key__"])
	})
}
