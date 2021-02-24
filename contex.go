package zero

// Ctx represents the Context which hold the event.
// 代表上下文
type Ctx struct {
	base  *Matcher
	Event Event
	State H
}
