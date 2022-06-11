package filter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type (
	FilterFunc func(gjson.Result) bool
	field      struct {
		key string
	}
)

// Filter return a rule filter the message.
func Filter[Ctx any](getevent func(Ctx) gjson.Result, filters ...FilterFunc) func(ctx Ctx) bool {
	return func(ctx Ctx) bool {
		return And(filters...)(getevent(ctx))
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
			if !filter(result) {
				return false
			}
		}
		return true
	}
}

// Not ...
func Not(filter FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return !filter(result)
	}
}

// Field ...
func Field(str string) *field {
	return &field{key: str}
}

// Any ...
func (f *field) Any(filter ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return Or(filter...)(result.Get(f.key))
	}
}

// All ...
func (f *field) All(filter ...FilterFunc) FilterFunc {
	return func(result gjson.Result) bool {
		return And(filter...)(result.Get(f.key))
	}
}

// Equal ...
func Equal(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return str == result.String()
	}
}

// NotEqual ...
func NotEqual(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return str != result.String()
	}
}

// In ...
func In(i ...interface{}) FilterFunc {
	ss := make([]string, 0)
	for _, v := range i {
		ss = append(ss, fmt.Sprint(v))
	}
	return func(result gjson.Result) bool {
		for _, s := range ss {
			if s == result.String() {
				return true
			}
		}
		return false
	}
}

// Contain ...
func Contain(str string) FilterFunc {
	return func(result gjson.Result) bool {
		return strings.Contains(result.String(), str)
	}
}

// Regex ...
func Regex(str string) FilterFunc {
	pat := regexp.MustCompile(str)
	return func(result gjson.Result) bool {
		return pat.MatchString(result.String())
	}
}
