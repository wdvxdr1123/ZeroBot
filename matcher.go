package zero

import "sync"

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
	Block    bool
	Type     string
	Priority int
	State    State
	Event    *Event
	Rules    []Rule
	handlers []Handler
}

var (
	// 所有主匹配器列表
	matcherList = make([]*Matcher, 0)
	// 临时匹配器
	tempMatcherList = matcherMap{}
	// Matcher 修改读写锁
	matcherLock = sync.RWMutex{}
)

type State map[string]interface{}

// SetBlock 设置是否阻断后面的 Matcher 触发
func (m *Matcher) SetBlock(block bool) *Matcher {
	m.Block = block
	return m
}

// SetBlock 设置当前 Matcher 优先级
func (m *Matcher) SetPriority(priority int) *Matcher {
	m.Priority = priority
	return m
}

// On 添加新的主匹配器
func On(type_ string, rules ...Rule) *Matcher {
	matcherLock.Lock()
	defer matcherLock.Unlock()
	var matcher = &Matcher{
		Type:     type_,
		State:    map[string]interface{}{},
		Rules:    rules,
		handlers: []Handler{},
	}
	matcherList = append(matcherList, matcher)
	return matcher
}

func (m *Matcher) run(event Event) {
	m.Event = &event
	for _, handler := range m.handlers {
		m.handlers = m.handlers[1:] // delete the handling handler
		switch handler(m, event, m.State) {
		case SuccessResponse:
			continue
		case FinishResponse:
			return
		case RejectResponse:
			tempMatcherList.Store(getSeq(), &Matcher{
				Type:  "message",
				State: m.State,
				Rules: []Rule{
					CheckUser(event.UserID),
				},
				handlers: append([]Handler{handler}, m.handlers...),
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
	tempMatcherList.Store(getSeq(), &Matcher{
		Type:  "message",
		State: map[string]interface{}{},
		Rules: []Rule{
			CheckUser(event.UserID),
		},
		handlers: []Handler{
			func(_ *Matcher, ev Event, _ State) Response {
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

// Handle 直接处理事件
func (m *Matcher) Handle(handler Handler) *Matcher {
	m.handlers = append(m.handlers, handler)
	return m
}

// Receive 接收一条消息后处理事件
func (m *Matcher) Receive(handler Handler) *Matcher {
	m.handlers = append(m.handlers, func(matcher *Matcher, event Event, state State) Response {
		tempMatcherList.Store(getSeq(), &Matcher{
			Type:     "message",
			State:    matcher.State,
			Rules:    []Rule{CheckUser(event.UserID)},
			handlers: append([]Handler{handler}, m.handlers...),
		})
		return FinishResponse
	})
	return m
}

// Got 判断State是否含有"name"键，若无则向用户索取
func (m *Matcher) Got(key, prompt string, handler Handler) *Matcher {
	m.handlers = append(
		m.handlers,
		// Got Handler
		func(matcher *Matcher, event Event, state State) Response {
			if _, ok := matcher.State[key]; ok == false {
				// send message to notify the user
				if prompt != "" {
					Send(event, prompt)
				}

				gotKeyHandler := func(matcher *Matcher, event Event, state State) Response {
					state[key] = event.RawMessage
					return SuccessResponse
				}
				// add temp matcher to got and process the left handlers
				tempMatcherList.Store(getSeq(), &Matcher{
					Type:     "message",
					State:    matcher.State,
					Rules:    []Rule{CheckUser(event.UserID)},
					handlers: append([]Handler{gotKeyHandler}, m.handlers...),
				})
				return FinishResponse
			}
			return handler(matcher, event, matcher.State)
		},
	)
	return m
}

// 消息触发器
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
