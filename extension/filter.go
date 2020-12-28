package extension

import (
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"regexp"
	"strings"
)

type (
	FilterFunc func(gjson.Result) bool
	field      struct {
		key string
	}
)

// Filter return a rule filter the message.
func Filter(filters ...FilterFunc) zero.Rule {
	return func(event *zero.Event, _ zero.State) bool {
		return And(filters...)(event.RawEvent)
	}
}

// Or ...
func Or(filters ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		for _, filter := range filters {
			if filter(result) {
				return true
			}
		}
		return false
	}
}

// And ...
func And(filters ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		for _, filter := range filters {
			if filter(result) {
				return true
			}
		}
		return false
	}
}

// Not ...
func Not(filter FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return !filter(result)
	}
}

// Field...
func Field(str string) *field {
	return &field{key: str}
}

// Select...
func (f *field) Select(filter ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return Or(filter...)(result.Get(f.key))
	}
}

// Match...
func (f *field) Match(filter ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return And(filter...)(result.Get(f.key))
	}
}

// Equal
func Equal(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return str == result.String()
	}
}

// NotEqual
func NotEqual(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return str != result.String()
	}
}

// In
func In(str ...string) FilterFunc {
	return func(result gjson.Result) bool {
		for _, s := range str {
			if s == result.Str {
				return true
			}
		}
		return false
	}
}

// Contain...
func Contain(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return strings.Contains(result.String(), str)
	}
}

// Regex...
func Regex(str string) FilterFunc {
	pat := regexp.MustCompile(str)
	return func(result gjson.Result) bool {
		return pat.MatchString(result.String())
	}
}
