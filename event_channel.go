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
	// 【修改点1】保留扩容到 100，应对突发消息
	ch := make(chan *Ctx, 100)
	matcher := StoreMatcher(&Matcher{
		Type:     Type(n.Type),
		Block:    n.Block,
		Priority: n.Priority,
		Rules:    n.Rule,
		Engine:   defaultEngine,
		Handler: func(ctx *Ctx) {
			// 【修改点2】使用 go func 异步发送
			// 只要业务层（Consumer）处理速度基本正常，这就不会阻塞 Bot 核心
			// 即使业务层处理慢，消息也会先堆积在 goroutine 中，而不会卡死主线程
			go func() {
				// 防止 ch 被 close 后写入导致 panic
				defer func() { _ = recover() }()
				ch <- ctx
			}()
		},
	})

	return ch, func() {
		matcher.Delete()
		// 防止多次调用 cancel 导致关闭已关闭通道的 panic
		defer func() { _ = recover() }()
		close(ch)
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