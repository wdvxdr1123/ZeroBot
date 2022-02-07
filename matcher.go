package zero

import (
	"sort"
	"sync"
)

// Handler 事件处理函数
type Handler func(ctx *Ctx)

// Matcher 是 ZeroBot 匹配和处理事件的最小单元
type Matcher struct {
	// Temp 是否为临时Matcher，临时 Matcher 匹配一次后就会删除当前 Matcher
	Temp bool
	// Priority 优先级，越小优先级越高
	Priority int
	// Event 当前匹配到的事件
	Event *Event
	// Handlers 处理事件的函数
	Handlers []Handler
	// MatcherGroup 所属的 MatcherGroup
	MatcherGroup *MatcherGroup
}

var (
	// 所有主匹配器列表
	matcherList = make([]*Matcher, 0)
	// Matcher 修改读写锁
	matcherLock = sync.RWMutex{}
)

// State store the context of a matcher.
type State map[string]interface{}

func sortMatcher() {
	sort.Slice(matcherList, func(i, j int) bool { // 按优先级排序
		return matcherList[i].Priority < matcherList[j].Priority
	})
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

// StoreMatcher store a matcher to matcher list.
func StoreMatcher(m *Matcher) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	// todo(wdvxdr): move to engine.
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
		}
	}
}

func (m *Matcher) copy() *Matcher {
	return &Matcher{
		Priority:     m.Priority,
		Handlers:     m.Handlers,
		Temp:         m.Temp,
		Event:        m.Event,
		MatcherGroup: m.MatcherGroup,
	}
}

// Handle 直接处理事件
func (m *Matcher) Handle(handler Handler) *Matcher {
	m.Handlers = append(m.Handlers, handler)
	return m
}
