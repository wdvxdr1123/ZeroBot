package path

import (
	"strings"
)

type matcher struct {
	r      *Route
	step   int
	s      string
	result map[string]string
}

func (r *Route) Match(s string) (map[string]string, bool) {
	m := &matcher{r: r, s: s, result: make(map[string]string)}
	if ok := m.match(s); ok {
		return m.result, true
	}
	return nil, false
}

func (m *matcher) match(s string) bool {
	var i = 0
	var sb = strings.Builder{}

	for {
		if m.step >= len(m.r.fields) { // the check the left string is empty.
			return s == ""
		} else if i == len(s) {
			if ok := m.step == len(m.r.fields)-1; ok {
				m.result[m.r.fields[m.step].pattern] = sb.String() // save
				return true
			}
			return false
		}

		seg := m.r.fields[m.step]
		switch seg.kind {
		case constPart:
			if !strings.HasPrefix(s, seg.pattern) {
				return false
			}

			m.step++
			return m.match(s[len(seg.pattern):])

		case requiredParam:
			m.step++
			ok := m.match(s[i:])
			if ok && sb.Len() > 0 {
				m.result[seg.pattern] = sb.String()
				return true
			}
			m.step--
			sb.WriteByte(s[i])
			i++

		case optionalParam:
			m.step++
			ok := m.match(s[i:])
			if ok {
				m.result[seg.pattern] = sb.String()
				return true
			}
			m.step--
			sb.WriteByte(s[i])
			i++
		}
	}
}
