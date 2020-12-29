package zero

import (
	"sort"
	"sync"
)

type (
	Response uint8
	Rule     func(event *Event, state State) bool
	Handler  func(matcher *Matcher, event Event, state State) Response
)

const (
	SuccessResponse Response = iota
	RejectResponse
	FinishResponse
)

type Matcher struct {
	Temp     bool
	Block    bool
	Priority int
	State    State
	Event    *Event
	Type     Rule
	Rules    []Rule
	Handlers []Handler
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

// On 添加新的主匹配器
func On(type_ string, rules ...Rule) *Matcher {
	var matcher = &Matcher{
		State:    map[string]interface{}{},
		Type:     Type(type_),
		Rules:    rules,
		Handlers: []Handler{},
	}
	StoreMatcher(matcher)
	return matcher
}

// StoreMatcher store a matcher to matcher list.
func StoreMatcher(m *Matcher) {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	matcherList = append(matcherList, m)
	sortMatcher()
}

// StoreTempMatcher store a matcher only triggered once.
func StoreTempMatcher(m *Matcher) {
	m.Temp = true
	StoreMatcher(m)
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
	for _, handler := range m.Handlers {
		m.Handlers = m.Handlers[1:] // delete the handling handler
		switch handler(m, event, m.State) {
		case SuccessResponse:
			continue
		case FinishResponse:
			return
		case RejectResponse:
			StoreTempMatcher(&Matcher{
				Type:     Type("message"),
				Block:    m.Block,
				Priority: m.Priority,
				State:    m.State,
				Rules: []Rule{
					CheckUser(event.UserID),
				},
				Handlers: append([]Handler{handler}, m.Handlers...),
			})
			return
		}
	}
}

// Get ..
func (m *Matcher) Get(prompt string) string {
	ch := make(chan string)
	event := m.Event
	Send(*event, prompt)
	StoreTempMatcher(&Matcher{
		Priority: m.Priority,
		Block:    m.Block,
		Type:     Type("message"),
		State:    map[string]interface{}{},
		Rules: []Rule{
			CheckUser(event.UserID),
		},
		Handlers: []Handler{
			func(_ *Matcher, ev Event, _ State) Response {
				ch <- ev.RawMessage
				return SuccessResponse
			},
		},
	})
	return <-ch
}

func (m *Matcher) copy() *Matcher {
	newHandlers := make([]Handler, len(m.Handlers))
	copy(newHandlers, m.Handlers) // 复制
	return &Matcher{
		State:    copyState(m.State),
		Type:     m.Type,
		Rules:    m.Rules,
		Handlers: newHandlers,
		Block:    m.Block,
		Priority: m.Priority,
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
	m.Handlers = append(m.Handlers, handler)
	return m
}

// Receive 接收一条消息后处理事件
func (m *Matcher) Receive(handler Handler) *Matcher {
	m.Handlers = append(m.Handlers, func(matcher *Matcher, event Event, state State) Response {
		StoreTempMatcher(&Matcher{
			Type:     Type("message"),
			Priority: matcher.Priority,
			Block:    matcher.Block,
			State:    matcher.State,
			Rules: []Rule{
				CheckUser(event.UserID),
			},
			Handlers: append([]Handler{handler}, m.Handlers...),
		})
		return FinishResponse
	})
	return m
}

// Got 判断State是否含有"name"键，若无则向用户索取
func (m *Matcher) Got(key, prompt string, handler Handler) *Matcher {
	m.Handlers = append(
		m.Handlers,
		// Got Handler
		func(matcher *Matcher, event Event, state State) Response {
			if _, ok := state[key]; !ok {
				state[key] = matcher.Get(prompt)
			}
			return handler(matcher, event, matcher.State)
		},
	)
	return m
}

// OnMessage 消息触发器
func OnMessage(rules ...Rule) *Matcher {
	return On("message", rules...)
}

// OnNotice 系统提示触发器
func OnNotice(rules ...Rule) *Matcher {
	return On("notice", rules...)
}

// OnRequest 请求消息触发器
func OnRequest(rules ...Rule) *Matcher {
	return On("request", rules...)
}

// OnMetaEvent 元事件触发器
func OnMetaEvent(rules ...Rule) *Matcher {
	return On("meta_event", rules...)
}

// OnPrefix 前缀触发器
func OnPrefix(prefix string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{PrefixRule(prefix)}, rules...)...)
}

// OnSuffix 后缀触发器
func OnSuffix(suffix string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{SuffixRule(suffix)}, rules...)...)
}

// OnCommand 命令触发器
func OnCommand(commands string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{CommandRule(commands)}, rules...)...)
}

// OnRegex 正则触发器
func OnRegex(regexPattern string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{RegexRule(regexPattern)}, rules...)...)
}

// OnKeyword 关键词触发器
func OnKeyword(keyword string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{KeywordRule(keyword)}, rules...)...)
}

// OnFullMatch 完全匹配触发器
func OnFullMatch(src string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{FullMatchRule(src)}, rules...)...)
}

// OnFullMatchGroup 完全匹配触发器组
func OnFullMatchGroup(src []string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{FullMatchRule(src...)}, rules...)...)
}

// OnKeywordGroup 关键词触发器组
func OnKeywordGroup(keywords []string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{KeywordRule(keywords...)}, rules...)...)
}

// OnCommandGroup 命令触发器组
func OnCommandGroup(commands []string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{CommandRule(commands...)}, rules...)...)
}

// OnPrefixGroup 前缀触发器组
func OnPrefixGroup(prefix []string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{PrefixRule(prefix...)}, rules...)...)
}

// OnSuffixGroup 后缀触发器组
func OnSuffixGroup(suffix []string, rules ...Rule) *Matcher {
	return OnMessage(append([]Rule{SuffixRule(suffix...)}, rules...)...)
}
