package extension

import (
	"sync"
	"sync/atomic"

	"github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type (
	nextMessage struct {
		rule []zero.Rule
		fn   func(m message.Message)
	}

	forMessage struct {
		rule []zero.Rule
		fn   func(m message.Message) zero.Response
	}

	selectMessage struct {
		cases []messageCase
	}

	messageCase struct {
	}
)

// NextMessage is a basic interact method.
func NextMessage() *nextMessage {
	return &nextMessage{
		rule: []zero.Rule{},
	}
}

// Rule is the next message trigger condition.
func (n *nextMessage) Rule(rule ...zero.Rule) *nextMessage {
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
	zero.StoreTempMatcher(&zero.Matcher{
		Type:     zero.Type("message"),
		Block:    true,
		Priority: 0,
		Rules:    n.rule,
		Handlers: []zero.Handler{
			func(_ *zero.Matcher, e zero.Event, _ zero.State) zero.Response {
				ch <- e.Message
				return zero.FinishResponse
			},
		},
	})
	n.fn(<-ch)
}

// ForMessage is a loop of NextMessage
func ForMessage() *forMessage {
	return &forMessage{}
}

// Rule is the next message trigger condition.
func (n *forMessage) Rule(rule ...zero.Rule) *forMessage {
	n.rule = append(n.rule, rule...)
	return n
}

// Handle is the logic of handle next message.
func (n *forMessage) Handle(fn func(m message.Message) zero.Response) *forMessage {
	n.fn = fn
	return n
}

// Do start wait next message.
func (n *forMessage) Do() {
	cond := sync.NewCond(&sync.Mutex{})
	var state uint32 = 0
	waitNextMessage := NextMessage().Rule(n.rule...).Handle(func(m message.Message) {
		if n.fn(m) == zero.FinishResponse {
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
func Select() *selectMessage {
	return &selectMessage{}
}

func (s *selectMessage) AddCase(cases ...messageCase) *selectMessage {
	s.cases = append(s.cases, cases...)
	return s
}

func (s *selectMessage) Do() {
	panic("impl me")
}
