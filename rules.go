package zero

import (
	"regexp"
	"strconv"
	"strings"
)

// 是否含有前缀
func PrefixRule(prefixes ...string) func(event *Event, state State) bool {
	return func(event *Event, state State) bool {
		if event.Message == nil && event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		first := event.Message[0]
		firstMessage := first.Data["text"]
		for _, prefix := range prefixes {
			first.Data["text"] = strings.TrimPrefix(firstMessage, prefix)
			if first.Data["text"] != firstMessage {
				state["prefix"] = prefix
				first.Data["text"] = strings.TrimLeft(first.Data["text"], " ")
				return true
			}
		}
		return false
	}
}

// 是否含有后缀
func SuffixRule(suffixes ...string) func(event *Event, state State) bool {
	return func(event *Event, state State) bool {
		if event.Message == nil { // 确保无空指针
			return false
		}
		last := event.Message[len(event.Message)-1]
		if last.Type != "text" {
			return false
		}
		lastMessage := last.Data["text"]
		for _, suffix := range suffixes {
			last.Data["text"] = strings.TrimSuffix(lastMessage, suffix)
			if last.Data["text"] != lastMessage {
				state["suffix"] = suffix
				last.Data["text"] = strings.TrimRight(last.Data["text"], " ")
				return true
			}
		}
		return false
	}
}

// command trigger
func CommandRule(commands ...string) func(event *Event, state State) bool {
	return func(event *Event, state State) bool {
		if event.Message == nil && event.Message[0].Type != "text" {
			return false
		}
		first := event.Message[0]
		firstMessage := first.Data["text"]
		cmdMessage := strings.TrimPrefix(firstMessage, zeroBot.commandPrefix)
		if cmdMessage == firstMessage {
			return false
		}
		for _, prefix := range commands {
			if strings.HasPrefix(cmdMessage,prefix) {
				state["command"] = prefix
				first.Data["text"] = cmdMessage[len(prefix):]
				first.Data["text"] = strings.TrimLeft(first.Data["text"], " ")
				return true
			}
		}
		return false
	}
}

// 正则匹配
func RegexRule(regexPattern string) func(event *Event, state State) bool {
	regex := regexp.MustCompile(regexPattern)
	return func(event *Event, state State) bool {
		msg := event.Message.CQString()
		state["regex_matched"] = regex.FindAllString(msg, -1)
		if state["regex_matched"] != nil {
			return false
		}
		return true
	}
}

// 关键词匹配
func KeywordRule(src ...string) func(event *Event, state State) bool {
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

// 完全匹配
func FullMatchRule(src ...string) func(event *Event, state State) bool {
	return func(event *Event, state State) bool {
		msg := event.Message.CQString()
		for _, str := range src {
			if str == msg {
				return true
			}
		}
		return false
	}
}

// only triggered in conditions of @bot or begin with the nicknames
func OnlyToMe(event *Event, _ State) bool {
	return event.IsToMe == true
}

// only triggered by specific person
func CheckUser(userId ...int64) func(event *Event, state State) bool {
	return func(event *Event, state State) bool {
		for _, uid := range userId {
			if event.UserID == uid {
				return true
			}
		}
		return false
	}
}

// only triggered in private message
func OnlyPrivate(event *Event, _ State) bool {
	return event.PostType == "message" && event.DetailType == "private"
}

// only triggered in public/group message
func OnlyGroup(event *Event, _ State) bool {
	return event.PostType == "message" && event.DetailType == "group"
}

func SuperUserPermission(event *Event, _ State) bool {
	uid := strconv.FormatInt(event.UserID, 10)
	for _, su := range zeroBot.SuperUsers {
		if su == uid {
			return true
		}
	}
	return false
}

// only triggered by the group admins or higher permission
func AdminPermission(event *Event, state State) bool {
	return SuperUserPermission(event, state) || event.Sender.Role != "member"
}

// only triggered by the group owner or higher permission
func OwnerPermission(event *Event, state State) bool {
	return SuperUserPermission(event, state) ||
		(event.Sender.Role != "member" && event.Sender.Role != "admin")
}
