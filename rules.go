package zero

import (
	"regexp"
	"strconv"
	"strings"
)

// Type check the event's type
func Type(type_ string) Rule {
	t := strings.SplitN(type_, "/", 3)
	return func(event *Event, _ State) bool {
		if len(t) > 0 && t[0] != event.PostType {
			return false
		}
		if len(t) > 1 && t[1] != event.DetailType {
			return false
		}
		if len(t) > 2 && t[2] != event.SubType {
			return false
		}
		return true
	}
}

// PrefixRule check if the message has the prefix and trim the prefix
func PrefixRule(prefixes ...string) Rule {
	return func(event *Event, state State) bool {
		if event.Message == nil || event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		first := event.Message[0]
		firstMessage := first.Data["text"]
		for _, prefix := range prefixes {
			if strings.HasPrefix(firstMessage, prefix) {
				state["prefix"] = prefix
				arg := strings.TrimLeft(firstMessage[len(prefix):], " ")
				if len(event.Message) > 1 {
					arg += event.Message[1:].ExtractPlainText()
				}
				state["args"] = arg
				return true
			}
		}
		return false
	}
}

// SuffixRule check if the message has the suffix and trim the suffix
func SuffixRule(suffixes ...string) Rule {
	return func(event *Event, state State) bool {
		mLen := len(event.Message)
		if mLen <= 0 { // 确保无空指针
			return false
		}
		last := event.Message[mLen-1]
		if last.Type != "text" {
			return false
		}
		lastMessage := last.Data["text"]
		for _, suffix := range suffixes {
			if strings.HasSuffix(lastMessage, suffix) {
				state["suffix"] = suffix
				arg := strings.TrimRight(lastMessage[:len(lastMessage)-len(suffix)], " ")
				if mLen >= 2 {
					arg += event.Message[:mLen].ExtractPlainText()
				}
				state["args"] = arg
				return true
			}
		}
		return false
	}
}

// CommandRule check if the message is a command and trim the command name
func CommandRule(commands ...string) Rule {
	return func(event *Event, state State) bool {
		if event.Message == nil || event.Message[0].Type != "text" {
			return false
		}
		first := event.Message[0]
		firstMessage := first.Data["text"]
		if !strings.HasPrefix(firstMessage, BotConfig.CommandPrefix) {
			return false
		}
		cmdMessage := firstMessage[len(BotConfig.CommandPrefix):]
		for _, command := range commands {
			if strings.HasPrefix(cmdMessage, command) {
				state["command"] = command
				arg := strings.TrimLeft(cmdMessage[len(command):], " ")
				if len(event.Message) > 1 {
					arg += event.Message[1:].ExtractPlainText()
				}
				state["args"] = arg
				return true
			}
		}
		return false
	}
}

// RegexRule check if the message can be matched by the regex pattern
func RegexRule(regexPattern string) Rule {
	regex := regexp.MustCompile(regexPattern)
	return func(event *Event, state State) bool {
		msg := event.RawMessage
		if regex.MatchString(msg) {
			state["regex_matched"] = regex.FindStringSubmatch(msg)
			return true
		}
		return false
	}
}

// ReplyRule check if the message is replying some message
func ReplyRule(messageID int64) Rule {
	var mid = strconv.FormatInt(messageID, 10)
	return func(event *Event, state State) bool {
		if len(event.Message) <= 0 {
			return false
		}
		if event.Message[0].Type != "reply" {
			return false
		}
		return event.Message[0].Data["id"] == mid
	}
}

// KeywordRule check if the message has a keyword or keywords
func KeywordRule(src ...string) Rule {
	return func(event *Event, state State) bool {
		msg := event.Message.CQString()
		for _, str := range src {
			if strings.Contains(msg, str) {
				state["keyword"] = str
				return true
			}
		}
		return false
	}
}

// FullMatchRule check if src has the same copy of the message
func FullMatchRule(src ...string) Rule {
	return func(event *Event, state State) bool {
		msg := event.Message.CQString()
		for _, str := range src {
			if str == msg {
				state["matched"] = msg
				return true
			}
		}
		return false
	}
}

// OnlyToMe only triggered in conditions of @bot or begin with the nicknames
func OnlyToMe(event *Event, _ State) bool {
	return event.IsToMe
}

// CheckUser only triggered by specific person
func CheckUser(userId ...int64) Rule {
	return func(event *Event, state State) bool {
		for _, uid := range userId {
			if event.UserID == uid {
				return true
			}
		}
		return false
	}
}

// OnlyPrivate requires that the event is private message
func OnlyPrivate(event *Event, _ State) bool {
	return event.PostType == "message" && event.DetailType == "private"
}

// OnlyGroup requires that the event is public/group message
func OnlyGroup(event *Event, _ State) bool {
	return event.PostType == "message" && event.DetailType == "group"
}

// SuperUserPermission only triggered by the bot's owner
func SuperUserPermission(event *Event, _ State) bool {
	uid := strconv.FormatInt(event.UserID, 10)
	for _, su := range BotConfig.SuperUsers {
		if su == uid {
			return true
		}
	}
	return false
}

// AdminPermission only triggered by the group admins or higher permission
func AdminPermission(event *Event, state State) bool {
	return SuperUserPermission(event, state) || event.Sender.Role != "member"
}

// OwnerPermission only triggered by the group owner or higher permission
func OwnerPermission(event *Event, state State) bool {
	return SuperUserPermission(event, state) ||
		(event.Sender.Role != "member" && event.Sender.Role != "admin")
}
