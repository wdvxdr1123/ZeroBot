package message

import (
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
)

// Modified from https://github.com/catsworld/qq-bot-api

var (
	matchReg = regexp.MustCompile(`\[CQ:\w+?.*?]`)
	typeReg  = regexp.MustCompile(`\[CQ:(\w+)`)
	paramReg = regexp.MustCompile(`,([\w\-.]+?)=([^,\]]+)`)
)

// StrictCommand indicates that whether a command must start with a specified command prefix, default to "/".
// See function #Command
var StrictCommand = false

// CommandPrefix is the prefix to identify a message as a command.
// See function #Command
var CommandPrefix = "/"

// ParseMessage parses msg, which might have 2 types, string or array,
// depending on the configuration of cqhttp, to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessage(msg []byte) Message {
	x := gjson.ParseBytes(msg)
	if x.IsArray() {
		return ParseMessageFromArray(x)
	} else {
		return ParseMessageFromString(x.String())
	}
}

// ParseMessageFromArray parses msg as type array to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessageFromArray(msgs gjson.Result) Message {
	var message = Message{}
	parse2map := func(val gjson.Result) map[string]string {
		var m = map[string]string{}
		val.ForEach(func(key, value gjson.Result) bool {
			m[key.String()] = value.String()
			return true
		})
		return m
	}
	msgs.ForEach(func(_, item gjson.Result) bool {
		message = append(message, MessageSegment{
			Type: item.Get("type").String(),
			Data: parse2map(item.Get("data")),
		})
		return true
	})
	return message
}

// ParseMessageSegmentsFromString parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessageFromString(str string) Message {
	var m = Message{}
	i := matchReg.FindAllStringSubmatchIndex(str, -1)
	si := 0
	for _, idx := range i {
		if idx[0] > si {
			text := str[si:idx[0]]
			m = append(m, Text(UnescapeCQText(text)))
		}
		code := str[idx[0]:idx[1]]
		si = idx[1]
		t := typeReg.FindAllStringSubmatch(code, -1)[0][1]
		ps := paramReg.FindAllStringSubmatch(code, -1)
		d := make(map[string]string)
		for _, p := range ps {
			d[p[1]] = UnescapeCQCodeText(p[2])
		}
		m = append(m, MessageSegment{
			Type: t,
			Data: d,
		})
	}
	if si != len(str) {
		m = append(m, Text(str[si:]))
	}
	return m
}

//todo: modify this command system

// IsCommand indicates whether a Message is a command.
// If #StrictCommand is true, only messages start with #CommandPrefix will be regard as command.
func (m *Message) IsCommand() bool {
	str := m.CQString()
	return IsCommand(str)
}

// Command parses a command message and returns the command with command arguments.
// In a StrictCommand mode, the initial #CommandPrefix in a command will be stripped off.
func (m *Message) Command() (cmd string, args []string) {
	str := m.CQString()
	return Command(str)
}

// IsCommand indicates whether a string is a command.
// If #StrictCommand is true, only strings start with #CommandPrefix will be regard as command.
func IsCommand(str string) bool {
	if len(str) == 0 {
		return false
	}
	if StrictCommand && (len(str) < len(CommandPrefix) || str[:len(CommandPrefix)] != CommandPrefix) {
		return false
	}
	return true
}

// Command parses a command string and returns the command with command arguments.
// In a StrictCommand mode, the initial #CommandPrefix in a command will be stripped off.
func Command(str string) (cmd string, args []string) {
	lcp := len(CommandPrefix)
	str = strings.Replace(str, `\\`, `\0x5c`, -1)
	str = strings.Replace(str, `\"`, `\0x22`, -1)
	str = strings.Replace(str, `\'`, `\0x27`, -1)
	strs := regexp.MustCompile(`'[\s\S]*?'|"[\s\S]*?"|\S*\[CQ:[\s\S]*?\]\S*|\S+`).FindAllString(str, -1)
	if len(strs) == 0 || len(strs[0]) == 0 {
		return
	}
	if StrictCommand {
		if len(strs[0]) < lcp || strs[0][:lcp] != CommandPrefix {
			return
		}
		cmd = strs[0][lcp:]
	} else {
		cmd = strs[0]
	}
	for _, arg := range strs[1:] {
		arg = strings.Trim(arg, `'"`)
		arg = strings.Replace(arg, `\0x27`, `'`, -1)
		arg = strings.Replace(arg, `\0x22`, `"`, -1)
		arg = strings.Replace(arg, `\0x5c`, `\`, -1)
		args = append(args, arg)
	}
	return
}

// CQString returns the CQEncoded string. All media in the message will be converted
// to its CQCode.
func (m Message) CQString() string {
	var str string
	for _, media := range m {
		if media.Type != "text" {
			str += media.CQCode()
		} else {
			str += media.Data["text"]
		}
	}
	return str
}
