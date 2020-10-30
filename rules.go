package zero

import "strings"

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

func CheckUser(userId int64) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.UserID == userId
	}
}

func OnlyToMe() func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.IsToMe == true
	}
}

func IsCommand(commands ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil && event.Message[0].Type != "text" { // 确保无空指针
			return false
		}
		firstMessage := event.Message[0].Data["text"]
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
