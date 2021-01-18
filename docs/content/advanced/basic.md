---
title: "Rule"
date: 2020-12-29T13:55:32+08:00
draft: false
---

## Rule

Rule 是ZeroBot过滤事件的核心,其定义如下所示

```golang
// State store the context of a matcher.
type State map[string]interface{}
// Rule filter the event
type Rule func(event *Event, state State) bool
```

其中 State 保存了在进行过滤事件中运算的结果，比如ZeroBot自带的
`KeywordRule`中，就将匹配到的关键词保存到了`state["keyword"]`中

```golang
// KeywordRule check if the message has a keyword or keywords
func KeywordRule(src ...string) Rule {
    return func(event *Event, state State) bool {
        msg := event.Message.CQString()
        for _, str := range src {
            if strings.Contains(msg, str) {
                state["keyword"] = str
                return true
            }
        }
        return false
    }
}
```

ZeroBot中自带了一些Rule，你可以在`rule.go`中找到它们

## Matcher

Matcher是 ZeroBot 匹配事件的最小单元，其定义如下

```golang
type Matcher struct {
    // Temp 是否为临时Matcher，临时 Matcher 匹配一次后就会删除当前 Matcher
    Temp     bool
    // Block 是否阻断后续 Matcher，为 true 时当前Matcher匹配成功后，后续Matcher不参与匹配
    Block    bool
    // Priority 优先级，越小优先级越高
    Priority int
    // State 上下文
    State    State
    // Event 当前匹配到的事件
    Event    *Event
    // Type 匹配的事件类型
    Type     Rule
    // Rules 匹配规则
    Rules    []Rule
    // Handler 处理事件的函数
    Handler  Handler
}
```

在 ZeroBot 中，维护了一个有序的 `MatcherList` 每当接收一个新事件时，都会对`MatcherList`
中的每一`Matcher`逐一匹配，你可以通过`zero.On`函数向`MatcherList`添加一个`Matcher`

同时 ZeroBot 也提供了一些其他函数，添加一些自带指定Rule的`Matcher`,例如

```go
zero.OnSuffix("复读") // 该Matcher自带 zero.SuffixRule("复读“) Rule
```

你可以在{{< button href="https://pkg.go.dev/github.com/wdvxdr1123/ZeroBot#On" >}}这里{{< /button >}}
查看其他自带的添加`Matcher`的函数

## Handler

togu
