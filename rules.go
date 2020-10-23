package ZeroBot

import "strings"

// 是否为消息事件
func IsMessage() func(event Event) bool {
	return func(event Event) bool {
		return event["post_type"].Str == "message"
	}
}

// 是否含有前缀,使用时请确保该事件为消息事件
func IsPrefix(prefixs ...string) func(event Event) bool {
	return func(event Event) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(event["raw_message"].Str, prefix) { // 只要有一个前缀就行了
				return true
			}
		}
		return false
	}
}

// 是否为通知事件
func IsNotice() func(event Event) bool {
	return func(event Event) bool {
		return event["post_type"].Str == "notice"
	}
}

// 是否为请求事件
func IsRequest() func(event Event) bool {
	return func(event Event) bool {
		return event["post_type"].Str == "request"
	}
}

// 是否为元事件
func IsMetaEvent() func(event Event) bool {
	return func(event Event) bool {
		return event["post_type"].Str == "meta_event"
	}
}

func CheckUser(userId int64) func(event Event) bool {
	return func(event Event) bool {
		return event["user_id"].Int() == userId
	}
}

