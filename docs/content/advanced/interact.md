---
title: "Event Channel"
date: 2021-01-19T16:43:41+08:00
draft: false
---

## Event Channel

在 ZeroBot 中，提供了用于异步获取指定事件的方法，通过返回channel，搭配 go 语言的select和其他语言特性，
可以很轻松的实现交互式机器人。

其核心为 FutureEvent，其定义如下

```go
// FutureEvent 是 ZeroBot 交互式的核心，用于异步获取指定事件
type FutureEvent struct {
    // 需要获取的事件类型
    // 
    // 形如 message/group, 具体为 post_type / detail_type / sub_type
    Type     string
    // 优先级， 同 Matcher
    Priority int
    // 同 Matcher
    Rule     []Rule
    // 同 Matcher
    Block    bool
}
```

你可以使用 `zero.NewFutureEvent` 创建，或者使用`matcher.FutureEvent`来创建一个与当前`matcher`
优先级和阻断性相同的`FutureEvent`

`FutureEvent`提供了两个基本方法，用于获取符合条件的事件。

### Next

Next 返回一个 `channel` 用于接收下一个指定事件，并且该事件传输完成后，就会关闭该 `channel`,

该 chan 必须接收，如需手动取消监听，请使用 Repeat 方法

### Repeat

 Repeat 返回一个 `channel` 用于接收无穷个指定事件，和一个取消监听的函数，
 如果没有取消监听，将不断监听指定事件

### 灵活使用

可以使用这些`channel`配合go强大的channel支持，组合出一些其他用法，如下面这个例子

```go
// Take 基于 Repeat 封装，返回一个 chan 接收指定数量的事件
//
// 该 chan 对象必须接收，否则将有 goroutine 泄漏，如需手动取消请使用 Repeat
func (n *FutureEvent) Take(num int) <-chan Event {
    recv, cancel := n.Repeat()
    ch := make(chan Event, num)
    go func() {
        defer close(ch)
        for i := 0; i < num; i++ {
            ch <- <-recv
        }
        cancel()
    }()
    return ch
}
```

### 实战

例如 example 中的 复读例子

```go
zero.OnCommand("开启复读").SetBlock(true).SetPriority(10).
    Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
        stop := zero.NewFutureEvent("message/group", 8, true, // 关闭指令要比复读指令优先级高
            zero.CommandRule("关闭复读"),      // 关闭复读指令
            zero.CheckUser(event.UserID)). // 只有开启者可以关闭复读模式
            Next()                         // 关闭需要一次

        echo, cancel := matcher.FutureEvent("message/group", // 优先级 和 开启复读指令相同
            zero.CheckUser(event.UserID)). // 只复读开启复读模式的人的消息
            Repeat()                       // 不断监听复读

        zero.Send(event, "已开启复读模式!")
        for {
            select {
            case e := <-echo: // 接收到需要复读的消息
                zero.Send(event, e.RawMessage)
            case <-stop: // 收到关闭复读指令
                cancel() // 取消复读监听
                zero.Send(event, "已关闭复读模式!")
                return zero.FinishResponse // 返回
            }
        }
    })
```

这里我们通过 for select 与 channel相配合，轻松的写出了一个可开启，关闭的复读插件
