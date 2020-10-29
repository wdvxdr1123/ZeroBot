package ZeroBot

import (
	"sync"
)

type (
	Response uint8
	Rule     func(event Event, state State) bool
	Handler  func(event Event, matcher *Matcher) Response
)

const (
	SuccessResponse Response = iota
	RejectResponse
	FinishResponse
)

type Matcher struct {
	State        State
	Rules        []Rule
	defaultState State
	handlers     []Handler
}

var (
	// 所有主匹配器列表
	matcherList = make([]*Matcher, 0)
	// 临时匹配器
	tempMatcherList = sync.Map{}
)

type State map[string]interface{}

func addTempMatcher(matcher *Matcher) {
	tempMatcherList.Store(getSeq(), matcher)
}

// 添加新的主匹配器
func On(rules ...Rule) *Matcher {
	var matcher = &Matcher{
		State:    map[string]interface{}{},
		Rules:    rules,
		handlers: []Handler{},
	}
	matcherList = append(matcherList, matcher)
	return matcher
}

func (m *Matcher) run(event Event) {
	for _, handler := range m.handlers {
		switch handler(event, m) {
		case SuccessResponse:
			continue
		case FinishResponse:
			return
		}
	}
}

func runMatcher(matcher *Matcher, event Event) {
	for _, rule := range matcher.Rules {
		var state = copyState(matcher.defaultState)
		if rule(event, state) == false {
			return
		}
	}
	m := matcher.copy()
	m.run(event)
}

func (m *Matcher) Get(event Event, prompt string) string {
	ch := make(chan string)
	Send(event, prompt)
	tempMatcherList.Store(getSeq(), &Matcher{
		State: map[string]interface{}{},
		Rules: []Rule{
			IsMessage(),
			CheckUser(event.UserID),
		},
		handlers: []Handler{
			func(ev Event, m *Matcher) Response {
				ch <- ev.RawMessage
				return SuccessResponse
			},
		},
	})
	return <-ch
}

func (m *Matcher) copy() *Matcher {
	newHandlers := make([]Handler, len(m.handlers))
	copy(newHandlers, m.handlers) // 复制
	return &Matcher{
		State:    copyState(m.State),
		Rules:    m.Rules,
		handlers: newHandlers,
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

// 直接处理事件
func (m *Matcher) Handle(handler Handler) *Matcher {
	m.handlers = append(m.handlers, handler)
	return m
}

// 判断State是否含有"name"键，若无则向用户索取
func (m *Matcher) Got(name, prompt string, handler Handler) *Matcher {
	m.handlers = append(m.handlers, func(event Event, matcher *Matcher) Response {
		if _, ok := matcher.State[name]; ok == false {
			matcher.State[name] = m.Get(event, prompt)
		}
		return handler(event, matcher)
	})
	return m
}

func OnMessage(rules ...Rule) *Matcher {
	return On(append(rules, IsMessage())...)
}

func OnNotice(rules ...Rule) *Matcher {
	return On(append(rules, IsNotice())...)
}

func OnRequest(rules ...Rule) *Matcher {
	return On(append(rules, IsRequest())...)
}

func OnMetaEvent(rules ...Rule) *Matcher {
	return On(append(rules, IsMetaEvent())...)
}

// 前缀触发器
func OnPrefix(prefix []string, rules ...Rule) *Matcher {
	return OnMessage(append(rules, IsPrefix(prefix...))...)
}

// 后缀触发器
func OnSuffix(suffix []string, rules ...Rule) *Matcher {
	return OnMessage(append(rules, IsSuffix(suffix...))...)
}

// 命令触发器
func OnCommand(commands []string, rules ...Rule) *Matcher {
	return OnMessage(append(rules, IsCommand(commands...))...)
}
