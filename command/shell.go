package command

import "strings"

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
func Parse(line string) []string {
	var args []string
	buf := strings.Builder{}
	var escaped, doubleQuoted, singleQuoted, backQuote bool
	backtick := ""

	got := argNo

	for _, r := range line {
		if escaped {
			buf.WriteRune(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf.WriteRune(r)
			} else {
				escaped = true
			}
			continue
		}

		if isSpace(r) {
			if singleQuoted || doubleQuoted || backQuote {
				buf.WriteRune(r)
				backtick += string(r)
			} else if got != argNo {
				args = append(args, buf.String())
				buf.Reset()
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted {
				backtick = ""
				backQuote = !backQuote
			}
		case '"':
			if !singleQuoted {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
			}
		case '\'':
			if !doubleQuoted {
				if singleQuoted {
					got = argSingle
				}
				singleQuoted = !singleQuoted
			}
		default:
			got = argSingle
			buf.WriteRune(r)
			if backQuote {
				backtick += string(r)
			}
		}
	}

	if got != argNo {
		args = append(args, buf.String())
	}

	return args
}
