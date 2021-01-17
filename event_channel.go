package zero

type nextEvent struct {
	eventType string
	matcher   *Matcher
	rule      []Rule
}

//NextEvent
func (m *Matcher) NextEvent(eventType string, rule ...Rule) (next *nextEvent) {
	next = &nextEvent{
		eventType: eventType,
		matcher:   m,
		rule:      rule,
	}
	return
}

// Recv returns a channel to receive next
func (n *nextEvent) Recv() <-chan Event {
	ch := make(chan Event)
	StoreTempMatcher(&Matcher{
		Type:     Type(n.eventType),
		Block:    n.matcher.Block,
		Priority: n.matcher.Priority,
		Rules:    n.rule,
		Handler: func(_ *Matcher, e Event, _ State) Response {
			ch <- e
			close(ch)
			return FinishResponse
		},
	})
	return ch
}

func (n *nextEvent) Repeat() (recv <-chan Event, cancel func()) {
	ch, done := make(chan Event), make(chan struct{})
	go func() {
		defer close(ch)
		in := make(chan Event)
		matcher := StoreMatcher(&Matcher{
			Type:     Type(n.eventType),
			Block:    n.matcher.Block,
			Priority: n.matcher.Priority,
			Rules:    n.rule,
			Handler: func(_ *Matcher, e Event, _ State) Response {
				in <- e
				return FinishResponse
			},
		})
		for {
			select {
			case e := <-in:
				ch <- e
			case <-done:
				matcher.Delete()
				close(in)
				return
			}
		}
	}()
	return ch, func() {
		close(done)
	}
}

func (n *nextEvent) Take(num int) <-chan Event {
	recv, cancel := n.Repeat()
	ch := make(chan Event, num)
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			ch <- <-recv
		}
		cancel()
	}()
	return ch
}
