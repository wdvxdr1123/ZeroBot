package ZeroBot

import "strings"

// 是否为消息事件
func IsMessage() func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.PostType == "message"
	}
}

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

// 是否为通知事件
func IsNotice() func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.PostType == "notice"
	}
}

// 是否为请求事件
func IsRequest() func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.PostType == "request"
	}
}

// 是否为元事件
func IsMetaEvent() func(event Event, state State) bool {
	return func(event Event, state State) bool {
		return event.PostType == "meta_event"
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
		panic("not implemented")
	}
}
