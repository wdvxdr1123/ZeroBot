package single

import (
	"github.com/RomiChan/syncx"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// Option 配置项
type Option[K comparable] func(*Single[K])

// Single 反并发
type Single[K comparable] struct {
	group syncx.Map[K, struct{}]
	key   func(ctx *zero.Ctx) K
	post  func(ctx *zero.Ctx)
}

// WithKeyFn 指定反并发的 Key
func WithKeyFn[K comparable](fn func(ctx *zero.Ctx) K) Option[K] {
	return func(s *Single[K]) {
		s.key = fn
	}
}

// WithPostFn 指定反并发拦截后的操作
func WithPostFn[K comparable](fn func(ctx *zero.Ctx)) Option[K] {
	return func(s *Single[K]) {
		s.post = fn
	}
}

// New 创建反并发中间件
func New[K comparable](op ...Option[K]) *Single[K] {
	s := Single[K]{}
	for _, option := range op {
		option(&s)
	}
	return &s
}

// Apply 为指定 Engine 添加反并发功能
func (s *Single[K]) Apply(engine *zero.Engine) {
	engine.UseMidHandler(func(ctx *zero.Ctx) bool {
		if s.key == nil {
			return true
		}
		key := s.key(ctx)
		if _, ok := s.group.Load(key); ok {
			if s.post != nil {
				defer s.post(ctx)
			}
			return false
		}
		s.group.Store(key, struct{}{})
		ctx.State["__single-key__"] = key
		return true
	})

	engine.UsePostHandler(func(ctx *zero.Ctx) {
		s.group.Delete(ctx.State["__single-key__"].(K))
	})
}
