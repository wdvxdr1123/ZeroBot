package zero

import "time"

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
	// 稍微增加一点缓冲，防止极端情况死锁
	ch := make(chan *Ctx, 5)
	StoreTempMatcher(&Matcher{
		Type:     Type(n.Type),
		Block:    n.Block,
		Priority: n.Priority,
		Rules:    n.Rule,
		Engine:   defaultEngine,
		Handler: func(ctx *Ctx) {
			// 使用 select 防止阻塞，虽然 Next 只取一次，但非阻塞是好习惯
			select {
			case ch <- ctx:
			default:
			}
			// Next 是一次性的，发送完最好关闭，但由 StoreTempMatcher 机制决定
			// 这里不手动关闭 ch，避免多次触发 panic，让调用方接收
			close(ch)
		},
	})
	return ch
}

// Repeat 返回一个 chan 用于接收无穷个指定事件，和一个取消监听的函数
//
// 如果没有取消监听，将不断监听指定事件
func (n *FutureEvent) Repeat() (recv <-chan *Ctx, cancel func()) {
	// 【核心修改1】大幅增加缓冲区
	// 100 的缓冲区足够应对你在处理图片(1~3秒)期间群里的消息爆发
	ch := make(chan *Ctx, 100)
	done := make(chan struct{})

	go func() {
		defer close(ch)
		// 内部通道也给足缓冲
		in := make(chan *Ctx, 100)

		matcher := StoreMatcher(&Matcher{
			Type:     Type(n.Type),
			Block:    n.Block,
			Priority: n.Priority,
			Rules:    n.Rule,
			Engine:   defaultEngine,
			Handler: func(ctx *Ctx) {
				// 【核心修改2】非阻塞写入
				// 如果 consumer 处理太慢导致 in 满了，这里会走 default 分支
				// 从而“丢弃”消息，而不是“卡死”整个 Bot 引擎
				select {
				case in <- ctx:
				default:
					// 缓冲区已满，丢弃该消息以保护主进程
				}
			},
		})

		// 确保退出时清理 Matcher
		defer matcher.Delete()

		for {
			select {
			case e := <-in:
				// 将消息转发给用户通道
				select {
				case ch <- e:
				case <-done:
					// 收到外部取消信号，退出
					return
				}
			case <-done:
				// 收到外部取消信号，退出
				return
			}
		}
	}()

	return ch, func() {
		// 防止多次调用 cancel 导致 panic
		select {
		case <-done:
		default:
			close(done)
		}
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
		defer cancel() // 确保任务完成后取消监听
		for i := 0; i < num; i++ {
			select {
			case e := <-recv:
				ch <- e
			case <-time.After(time.Minute * 10):
				// 加上一个超长超时防止彻底泄露（可选优化）
				return
			}
		}
	}()
	return ch
}
