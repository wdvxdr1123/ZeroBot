package command

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
	buf := ""
	var escaped, doubleQuoted, singleQuoted, backQuote bool
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
			if singleQuoted || doubleQuoted || backQuote {
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
			buf += string(r)
			if backQuote {
				backtick += string(r)
			}
		}
	}

	if got != argNo {
		args = append(args, buf)
	}

	return args
}
