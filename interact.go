package zero

import (
	"sync"
	"sync/atomic"

	"github.com/wdvxdr1123/ZeroBot/message"
)

type (
	nextMessage struct {
		matcher *Matcher
		rule    []Rule
		fn      func(m message.Message)
	}

	forMessage struct {
		matcher *Matcher
		next    *nextMessage
		rule    []Rule
		fn      func(m message.Message) Response
	}

	selectMessage struct {
		cases []messageCase
	}

	messageCase struct {
	}
)

// NextMessage is a basic interact method.
func (m *Matcher) NextMessage() *nextMessage {
	return &nextMessage{
		rule:    []Rule{},
		matcher: m,
	}
}

// Rule is the next message trigger condition.
func (n *nextMessage) Rule(rule ...Rule) *nextMessage {
	n.rule = append(n.rule, rule...)
	return n
}

// Handle is the logic of handle next message.
func (n *nextMessage) Handle(fn func(m message.Message)) *nextMessage {
	n.fn = fn
	return n
}

// Do start wait next message.
func (n *nextMessage) Do() {
	ch := make(chan message.Message)
	StoreTempMatcher(&Matcher{
		Type:     Type("message"),
		Block:    n.matcher.Block,
		Priority: n.matcher.Priority,
		Rules:    n.rule,
		Handler: func(_ *Matcher, e Event, _ State) Response {
			ch <- e.Message
			return FinishResponse
		},
	})
	n.fn(<-ch)
}

// ForMessage is a loop of NextMessage
func (m *Matcher) ForMessage() *forMessage {
	return &forMessage{
		next: m.NextMessage(),
	}
}

// Rule is the next message trigger condition.
func (n *forMessage) Rule(rule ...Rule) *forMessage {
	n.next.Rule(rule...)
	return n
}

// Handle is the logic of handle next message.
func (n *forMessage) Handle(fn func(m message.Message) Response) *forMessage {
	n.fn = fn
	return n
}

// Do start wait next message.
func (n *forMessage) Do() {
	cond := sync.NewCond(&sync.Mutex{})
	var state uint32 = 0
	waitNextMessage := n.next.Handle(func(m message.Message) {
		if n.fn(m) == FinishResponse {
			atomic.StoreUint32(&state, 1)
		}
		cond.Signal()
	})
	cond.L.Lock()
	for state != 1 {
		go waitNextMessage.Do()
		cond.Wait()
	}
	cond.L.Unlock()
}

// Select
func (m *Matcher) Select() *selectMessage {
	return &selectMessage{}
}

func (s *selectMessage) AddCase(cases ...messageCase) *selectMessage {
	s.cases = append(s.cases, cases...)
	return s
}

func (s *selectMessage) Do() {
	panic("impl me")
}
