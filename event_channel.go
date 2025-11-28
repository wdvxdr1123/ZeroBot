package zero

// FutureEvent 是 ZeroBot 交互式的核心，用于异步获取指定事件
type FutureEvent struct {
	Type     string
	Priority int
	Rule     []Rule
	Block    bool
}

// NewFutureEvent 创建一个FutureEvent, 并返回其指针
func NewFutureEvent(typ string, priority int, block bool, rule ...Rule) *FutureEvent {
	return &FutureEvent{
		Type:     typ,
		Priority: priority,
		Rule:     rule,
		Block:    block,
	}
}

// FutureEvent 返回一个 FutureEvent 实例指针，用于获取满足 Rule 的 未来事件
//
// 此 FutureEvent 必然比 Matcher 之优先级少 1
func (m *Matcher) FutureEvent(typ string, rule ...Rule) *FutureEvent {
	return &FutureEvent{
		Type:     typ,
		Priority: m.Priority - 1,
		Block:    m.Block,
		Rule:     rule,
	}
}

// Next 返回一个 chan 用于接收下一个指定事件
//
// 该 chan 必须接收，如需手动取消监听，请使用 Repeat 方法
func (n *FutureEvent) Next() <-chan *Ctx {
	ch := make(chan *Ctx, 1)
	StoreTempMatcher(&Matcher{
		Type:     Type(n.Type),
		Block:    n.Block,
		Priority: n.Priority,
		Rules:    n.Rule,
		Engine:   defaultEngine,
		Handler: func(ctx *Ctx) {
			// 使用 go func 异步发送，确保不阻塞主线程
			go func() {
				defer func() { _ = recover() }()
				ch <- ctx
				close(ch)
			}()
		},
	})
	return ch
}

// Repeat 返回一个 chan 用于接收无穷个指定事件，和一个取消监听的函数
//
// 如果没有取消监听，将不断监听指定事件
func (n *FutureEvent) Repeat() (recv <-chan *Ctx, cancel func()) {
	// 保留扩容到 100，应对突发消息
	ch, done := make(chan *Ctx, 100), make(chan struct{})
	go func() {
		defer close(ch)
		in := make(chan *Ctx, 1)
		matcher := StoreMatcher(&Matcher{
			Type:     Type(n.Type),
			Block:    n.Block,
			Priority: n.Priority,
			Rules:    n.Rule,
			Engine:   defaultEngine,
			Handler: func(ctx *Ctx) {
				// 只要 Consumer 处理不是极度滞后，这种方式就能防止 Bot 核心被阻塞
				go func() {
					defer func() { _ = recover() }()
					in <- ctx
				}()
			},
		})
		for {
			select {
			case e := <-in:
				ch <- e
			case <-done:
				matcher.Delete()
				close(in)
				return
			}
		}
	}()
	return ch, func() {
		close(done)
	}
}

// Take 基于 Repeat 封装，返回一个 chan 接收指定数量的事件
//
// 该 chan 对象必须接收，否则将有 goroutine 泄漏，如需手动取消请使用 Repeat
func (n *FutureEvent) Take(num int) <-chan *Ctx {
	recv, cancel := n.Repeat()
	ch := make(chan *Ctx, num)
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			ch <- <-recv
		}
		cancel()
	}()
	return ch
}