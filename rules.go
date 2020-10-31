package zero

import (
	"strconv"
	"strings"
)

// 是否含有前缀
func IsPrefix(prefixes ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil && event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		firstMessage := event.Message[0].Data
		for _, prefix := range prefixes {
			if strings.HasPrefix(firstMessage["text"], prefix) { // 只要有一个前缀就行了
				return true
			}
		}
		return false
	}
}

// 是否含有后缀
func IsSuffix(prefixes ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool { // todo
		if event.Message == nil && event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		firstMessage := event.Message[0].Data
		for _, prefix := range prefixes {
			if strings.HasPrefix(firstMessage["text"], prefix) { // 只要有一个前缀就行了
				return true
			}
		}
		return false
	}
}

func OnlyToMe(event Event, _ State) bool {
	return event.IsToMe == true
}

func IsCommand(commands ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil && event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		firstMessage := event.Message[0].Data["text"]
		if strings.HasPrefix(firstMessage, zeroBot.commandPrefix) {
			event.Message[0].Data["text"] = firstMessage[len(zeroBot.commandPrefix):]
			firstMessage = event.Message[0].Data["text"]
		} else {
			return false
		}
		for _, prefix := range commands {
			if strings.HasPrefix(firstMessage, prefix) { // 只要有一个前缀就行了
				state["command"] = prefix
				event.Message[0].Data["text"] = firstMessage[len(prefix):] // 去除指令
				return true
			}
		}
		return false
	}
}

// only triggered by specific person
func CheckUser(userId int64) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.UserID == userId
	}
}

// only triggered in private message
func OnlyPrivate(event Event, _ State) bool {
	return event.PostType=="message" && event.DetailType == "private"
}

// only triggered in public/group message
func OnlyGroup(event Event, _ State) bool {
	return event.PostType == "message" && event.DetailType == "group"
}

func SuperUserPermission(event Event, _ State) bool {
	uid := strconv.FormatInt(event.UserID,10)
	for _, su := range zeroBot.SuperUsers {
		if su == uid {
			return true
		}
	}
	return false
}

// only triggered by the group admins or higher permission
func AdminPermission(event Event, state State) bool {
	return SuperUserPermission(event, state) || event.Sender.Role != "member"
}

// only triggered by the group owner or higher permission
func OwnerPermission(event Event, state State) bool {
	return SuperUserPermission(event, state) ||
		(event.Sender.Role != "member" && event.Sender.Role != "admin")
}