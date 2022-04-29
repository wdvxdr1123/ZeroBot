package zero

import (
	"sort"
	"sync"
)

type (
	// Rule filter the event
	Rule func(ctx *Ctx) bool
	// Handler 事件处理函数
	Handler func(ctx *Ctx)
)

// Matcher 是 ZeroBot 匹配和处理事件的最小单元
type Matcher struct {
	// Temp 是否为临时Matcher，临时 Matcher 匹配一次后就会删除当前 Matcher
	Temp bool
	// Block 是否阻断后续 Matcher，为 true 时当前Matcher匹配成功后，后续Matcher不参与匹配
	Block bool
	// Priority 优先级，越小优先级越高
	Priority int
	// Event 当前匹配到的事件
	Event *Event
	// Type 匹配的事件类型
	Type Rule
	// Rules 匹配规则
	Rules []Rule
	// Handler 处理事件的函数
	Handler Handler
	// Engine 注册 Matcher 的 Engine，Engine可为一系列 Matcher 添加通用 Rule 和 其他钩子
	Engine *Engine
}

var (
	// 所有主匹配器列表
	matcherList = make([]*Matcher, 0)
	// Matcher 修改读写锁
	matcherLock = sync.RWMutex{}
	// 用于迭代的所有主匹配器列表
	matcherListForRanging []*Matcher
	// 是否 matcherList 已经改变
	// 如果改变，下次迭代需要更新
	// matcherListForRanging
	hasMatcherListChanged bool
)

// State store the context of a matcher.
type State map[string]interface{}

func sortMatcher() {
	sort.Slice(matcherList, func(i, j int) bool { // 按优先级排序
		return matcherList[i].Priority < matcherList[j].Priority
	})
	hasMatcherListChanged = true
}

// SetBlock 设置是否阻断后面的 Matcher 触发
func (m *Matcher) SetBlock(block bool) *Matcher {
	m.Block = block
	return m
}

// SetPriority 设置当前 Matcher 优先级
func (m *Matcher) SetPriority(priority int) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	m.Priority = priority
	sortMatcher()
	return m
}

// FirstPriority 设置当前 Matcher 优先级 - 0
func (m *Matcher) FirstPriority() *Matcher {
	return m.SetPriority(0)
}

// SecondPriority 设置当前 Matcher 优先级 - 1
func (m *Matcher) SecondPriority() *Matcher {
	return m.SetPriority(1)
}

// ThirdPriority 设置当前 Matcher 优先级 - 2
func (m *Matcher) ThirdPriority() *Matcher {
	return m.SetPriority(2)
}

// BindEngine bind the matcher to a engine
func (m *Matcher) BindEngine(e *Engine) *Matcher {
	m.Engine = e
	return m
}

// StoreMatcher store a matcher to matcher list.
func StoreMatcher(m *Matcher) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	// todo(wdvxdr): move to engine.
	if m.Engine != nil {
		m.Block = m.Block || m.Engine.block
	}
	matcherList = append(matcherList, m)
	sortMatcher()
	return m
}

// StoreTempMatcher store a matcher only triggered once.
func StoreTempMatcher(m *Matcher) *Matcher {
	m.Temp = true
	StoreMatcher(m)
	return m
}

// Delete remove the matcher from list
func (m *Matcher) Delete() {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	for i, matcher := range matcherList {
		if m == matcher {
			matcherList = append(matcherList[:i], matcherList[i+1:]...)
			hasMatcherListChanged = true
		}
	}
}

func (m *Matcher) copy() *Matcher {
	return &Matcher{
		Type:     m.Type,
		Rules:    m.Rules,
		Block:    m.Block,
		Priority: m.Priority,
		Handler:  m.Handler,
		Temp:     m.Temp,
		Engine:   m.Engine,
	}
}

// Handle 直接处理事件
func (m *Matcher) Handle(handler Handler) *Matcher {
	m.Handler = handler
	return m
}
