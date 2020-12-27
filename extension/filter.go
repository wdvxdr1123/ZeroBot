package extension

import zero "github.com/wdvxdr1123/ZeroBot"

type (
	FilterFunc func(e *zero.Event) bool
)

// Filter return a rule filter the message.
func Filter(filters ...FilterFunc) zero.Rule {
	return func(event *zero.Event, _ zero.State) bool {
		return And(filters...)(event)
	}
}

// Or ...
func Or(filters ...FilterFunc) FilterFunc {
	return func(event *zero.Event) bool {
		for _,filter := range filters {
			if filter(event) {
				return true
			}
		}
		return false
	}
}

// And ...
func And(filters ...FilterFunc) FilterFunc {
	return func(event *zero.Event) bool {
		for _,filter := range filters {
			if filter(event) {
				return true
			}
		}
		return false
	}
}

// Not ...
func Not(filter FilterFunc) FilterFunc {
	return func(event *zero.Event) bool {
		return !filter(event)
	}
}
