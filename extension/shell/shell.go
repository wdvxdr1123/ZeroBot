// Package shell provides a simple shell parser for zerobot.
package shell

import (
	"errors"
	"strings"
)

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}

type argType int

const (
	argNo argType = iota
	argSingle
	argQuoted
)

// Parse 将指令转换为指令参数.
// modified from https://github.com/mattn/go-shellwords
func Parse(line string) ([]string, error) {
	var args []string
	buf := ""
	var escaped, doubleQuoted, singleQuoted, backQuote, dollarQuote bool
	backtick := ""

	got := argNo

	for _, r := range line {
		if escaped {
			buf += string(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf += string(r)
			} else {
				escaped = true
			}
			continue
		}

		if isSpace(r) {
			if singleQuoted || doubleQuoted || backQuote || dollarQuote {
				buf += string(r)
				backtick += string(r)
			} else if got != argNo {
				args = append(args, buf)
				buf = ""
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted && !dollarQuote {
				backtick = ""
				backQuote = !backQuote
			}
		case ')':
			if !singleQuoted && !doubleQuoted && !backQuote {
				backtick = ""
				dollarQuote = !dollarQuote
			}
		case '(':
			if !singleQuoted && !doubleQuoted && !backQuote {
				if !dollarQuote && strings.HasSuffix(buf, "$") {
					dollarQuote = true
					buf += "("
					continue
				} else {
					return nil, errors.New("invalid command line string")
				}
			}
		case '"':
			if !singleQuoted && !dollarQuote {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
			}
		case '\'':
			if !doubleQuoted && !dollarQuote {
				if singleQuoted {
					got = argSingle
				}
				singleQuoted = !singleQuoted
			}
		default:
			got = argSingle
			buf += string(r)
			if backQuote || dollarQuote {
				backtick += string(r)
			}
		}
	}

	if got != argNo {
		args = append(args, buf)
	}

	if escaped || singleQuoted || doubleQuoted || backQuote || dollarQuote {
		return nil, errors.New("invalid command line string")
	}

	return args, nil
}
