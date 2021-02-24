package zero

import (
	"sort"
	"sync"
)

type (
	Response uint8
	// Rule filter the event
	Rule    func(event *Event, state State) bool
	Handler func(matcher *Matcher, event Event, state State) Response
)

const (
	SuccessResponse Response = iota
	RejectResponse
	FinishResponse
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

// StoreMatcher store a matcher to matcher list.
func StoreMatcher(m *Matcher) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
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

func (m *Matcher) run(event Event) {
	m.Event = &event
	if m.Handler == nil {
		return
	}
	switch m.Handler(m, event, nil) {
	}
}

// Get ..
func (m *Matcher) Get(prompt string) string {
	event := m.Event
	if prompt != "" {
		Send(*event, prompt)
	}
	return (<-m.FutureEvent("message", CheckUser(event.UserID)).Next()).RawMessage
}

func (m *Matcher) copy() *Matcher {
	return &Matcher{
		Type:     m.Type,
		Rules:    m.Rules,
		Block:    m.Block,
		Priority: m.Priority,
		Handler:  m.Handler,
		Temp:     m.Temp,
	}
}

// 拷贝字典
func copyState(src State) State {
	dst := make(State)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// Handle 直接处理事件
func (m *Matcher) Handle(handler Handler) *Matcher {
	m.Handler = handler
	return m
}
