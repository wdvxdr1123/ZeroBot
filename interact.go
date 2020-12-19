package zero

import "github.com/wdvxdr1123/ZeroBot/message"

type (
	nextMessage struct {
		rule []Rule
		fn   func(m message.Message)
	}

	forMessage struct {
		rule []Rule
		fn func(m message.Message) Response
	}

	selectNextMessage struct {
	}
)

// NextMessage is a basic interact method.
func NextMessage() *nextMessage {
	return &nextMessage{
		rule: []Rule{},
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
	tempMatcherList.Store(getSeq(), &Matcher{
		Block:    true,
		Type:     "message",
		Priority: 0,
		Rules:    n.rule,
		handlers: []Handler{
			func(_ *Matcher, e Event, _ State) Response {
				ch <- e.Message
				return FinishResponse
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
func (n *forMessage) Rule(rule ...Rule) *forMessage {
	n.rule = append(n.rule, rule...)
	return n
}

// Handle is the logic of handle next message.
func (n *forMessage) Handle(fn func(m message.Message) Response) *forMessage {
	n.fn = fn
	return n
}

// Do start wait next message.
func (n *forMessage) Do()  {
	panic("impl me")
}

// Select
func Select() {
	panic("impl me")
}