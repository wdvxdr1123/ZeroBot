package ZeroBot

import (
	"sort"
	"sync"
)

type Rule func(Event) bool

type Matcher struct {
	sync.RWMutex // todo: 并发安全
	Priority     int64
	Block        bool
	State        State
	Rules        []Rule
	isTemp       bool
}

// 所有匹配器列表
var MatcherList []*Matcher //todo: 替换为并发安全的链表

type State map[string]interface{}

func addMatcher(matcher *Matcher) {
	MatcherList = append(MatcherList, matcher)
	sort.Slice(MatcherList, func(i, j int) bool { // 按照优先级排序
		return MatcherList[i].Priority < MatcherList[j].Priority
	})
}

func On(priority int64, block bool, defaultState State, rules ...Rule) *Matcher {
	var matcher = &Matcher{
		Priority: priority,
		Block:    block,
		State:    defaultState,
		Rules:    rules,
		isTemp:   false,
	}
	if MatcherList != nil {
		MatcherList = []*Matcher{}
	}
	return matcher
}

func (m *Matcher) run(event Event) error {
	for _, rule := range m.Rules {
		if rule(event) == false {
			// return
		}
	}
	// 满足所有条件，创建一个新会话
	panic("impl me")
}

func (m *Matcher) Get() string {
	ch := make(chan string)
	seqMap.Store(getSeq(),ch)
	// todo:处理
	return<-ch
}

func (m *Matcher) copy() *Matcher {
	return &Matcher{
		Priority: m.Priority,
		Block:    m.Block,
		State:    m.State, // Fixme:copy
		Rules:    m.Rules,
		isTemp:   m.isTemp,
	}
}
