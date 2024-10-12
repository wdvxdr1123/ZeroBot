package filter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type (
	Func  func(gjson.Result) bool
	Field struct {
		key string
	}
)

// New return a rule filter the message.
func New[Ctx any](getevent func(Ctx) gjson.Result, filters ...Func) func(ctx Ctx) bool {
	return func(ctx Ctx) bool {
		return And(filters...)(getevent(ctx))
	}
}

// Or ...
func Or(filters ...Func) Func {
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
func And(filters ...Func) Func {
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
func Not(filter Func) Func {
	return func(result gjson.Result) bool {
		return !filter(result)
	}
}

// NewField ...
func NewField(str string) *Field {
	return &Field{key: str}
}

// Any ...
func (f *Field) Any(filter ...Func) Func {
	return func(result gjson.Result) bool {
		return Or(filter...)(result.Get(f.key))
	}
}

// All ...
func (f *Field) All(filter ...Func) Func {
	return func(result gjson.Result) bool {
		return And(filter...)(result.Get(f.key))
	}
}

// Equal ...
func Equal(str string) Func {
	return func(result gjson.Result) bool {
		return str == result.String()
	}
}

// NotEqual ...
func NotEqual(str string) Func {
	return func(result gjson.Result) bool {
		return str != result.String()
	}
}

// In ...
func In(i ...interface{}) Func {
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
func Contain(str string) Func {
	return func(result gjson.Result) bool {
		return strings.Contains(result.String(), str)
	}
}

// Regex ...
func Regex(str string) Func {
	pat := regexp.MustCompile(str)
	return func(result gjson.Result) bool {
		return pat.MatchString(result.String())
	}
}
