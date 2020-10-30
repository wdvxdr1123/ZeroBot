package ZeroBot

import "strings"

// 是否含有前缀,使用时请确保该事件为消息事件
func IsPrefix(prefixes ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil { // 确保无空指针
			return false
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(event.Message.StringMessage, prefix) { // 只要有一个前缀就行了
				return true
			}
		}
		return false
	}
}

// 是否含有前缀,使用时请确保该事件为消息事件
func IsSuffix(prefixs ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil { // 确保无空指针
			return false
		}
		for _, prefix := range prefixs {
			if strings.HasSuffix(event.Message.StringMessage, prefix) { // 只要有一个前缀就行了
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
		if event.Message == nil {
			return false
		}
		return event.Message.IsToMe == true
	}
}

func IsCommand(commands ...string) func(event Event, state State) bool {
	return func(event Event, state State) bool {
		if event.Message == nil { // 确保无空指针
			return false
		}
		// if event.
		for _, prefix := range commands {
			if strings.HasPrefix(event.Message.StringMessage, prefix) { // 只要有一个前缀就行了
				state["command"] = prefix
				return true
			}
		}
		return false
	}
}
