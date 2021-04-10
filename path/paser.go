package path

import (
	"errors"
)

var (
	// InvalidPattern 不合法的模式串
	InvalidPattern = errors.New("invalid pattern")
	// InvalidParamName 不合法的参数名
	InvalidParamName = errors.New("invalid name")
	// InvalidQuestion 无效的?操作
	InvalidQuestion = errors.New("qustion only used before colon")
)

// segmentKind is the kind of route segment, see the consts below.
type segmentKind int

const (
	constPart segmentKind = iota
	requiredParam
	optionalParam
)

type segment struct {
	kind    segmentKind
	pattern string
}

// Route is a simple command route
type Route struct {
	fields []segment
}

type scanner struct {
	state     int
	pos, prev int
	pattern   string
}

const (
	constPath = iota
	paramPath
)

// save saves the segment and returns a required param or a const path.
func (s *scanner) save() (*segment, error) {
	if s.state == constPath {
		if s.pos > s.prev {
			return &segment{kind: constPart, pattern: s.pattern[s.prev:s.pos]}, nil
		}
		return nil, nil // None segment
	}

	pos, prev := s.pos-1, s.prev-1
	if pos > prev {
		return &segment{kind: requiredParam, pattern: s.pattern[prev:pos]}, nil
	}
	return nil, InvalidParamName // param without name
}

// Parse parses pattern to compiled Route
func Parse(pattern string) (*Route, error) {
	// todo: better error trance
	route := new(Route)
	s := &scanner{pattern: pattern}
	for s.pos < len(pattern) {
		switch pattern[s.pos] {
		// param pattern start&end
		case ':':
			field, err := s.save() // save the pattern
			if err != nil {
				return nil, err
			} else if field != nil {
				route.fields = append(route.fields, *field)
			}

			s.state = s.state ^ paramPath // reverse the state

		// optional param pattern
		case '?':
			if s.state != paramPath || // check the state
				s.pos+1 >= len(pattern) || // check the next char
				pattern[s.pos+1] != ':' {
				return nil, InvalidQuestion
			}
			field, err := s.save()
			if err != nil || field == nil { // invalid param
				return nil, err
			}

			field.kind = optionalParam
			route.fields = append(route.fields, *field)
			s.pos++
			s.state = constPath

		case '\t', '\r', '\n', ' ':
			if s.state == paramPath {
				return nil, InvalidParamName
			}

		default:
		}
		s.pos++
	}

	if s.state != 0 {
		return nil, InvalidPattern
	}

	field, err := s.save() // save the pattern
	if err != nil {
		return nil, err
	} else if field != nil {
		route.fields = append(route.fields, *field)
	}

	return route, nil
}
