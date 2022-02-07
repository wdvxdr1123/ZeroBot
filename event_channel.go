package zero

// FutureEvent 是 ZeroBot 交互式的核心，用于异步获取指定事件
type FutureEvent struct {
	Type     string
	Priority int
	Handlers []Handler
	Block    bool
}

// NewFutureEvent 创建一个FutureEvent, 并返回其指针
func NewFutureEvent(Type string, Priority int, Block bool, handlers ...Handler) *FutureEvent {
	return &FutureEvent{
		Type:     Type,
		Priority: Priority,
		Handlers: handlers,
		Block:    Block,
	}
}

// FutureEvent 返回一个 FutureEvent 实例指针，用于获取满足 Rule 的 未来事件
func (m *Matcher) FutureEvent(Type string, handlers ...Handler) *FutureEvent {
	return &FutureEvent{
		Type:     Type,
		Priority: m.Priority,
		Handlers: handlers,
	}
}

// Next 返回一个 chan 用于接收下一个指定事件
//
// 该 chan 必须接收，如需手动取消监听，请使用 Repeat 方法
func (n *FutureEvent) Next() <-chan *Event {
	ch := make(chan *Event, 1)
	StoreTempMatcher(&Matcher{
		Priority: n.Priority,
		Handlers: append(n.Handlers, func(c *Ctx) {
			ch <- c.Event
			close(ch)
			if n.Block {
				c.Block()
			}
		}),
		MatcherGroup: defualtMatcherGroup,
	})
	return ch
}

// Repeat 返回一个 chan 用于接收无穷个指定事件，和一个取消监听的函数
//
// 如果没有取消监听，将不断监听指定事件
func (n *FutureEvent) Repeat() (recv <-chan *Event, cancel func()) {
	ch, done := make(chan *Event), make(chan struct{})
	go func() {
		defer close(ch)
		in := make(chan *Event, 1)
		matcher := StoreMatcher(&Matcher{
			Priority: n.Priority,
			Handlers: append(n.Handlers, func(ctx *Ctx) {
				in <- ctx.Event
				if n.Block {
					ctx.Block()
				}
			}),
			MatcherGroup: defualtMatcherGroup,
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
func (n *FutureEvent) Take(num int) <-chan *Event {
	recv, cancel := n.Repeat()
	ch := make(chan *Event, num)
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			ch <- <-recv
		}
		cancel()
	}()
	return ch
}
