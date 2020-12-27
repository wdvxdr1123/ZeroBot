package extension

import (
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// Type check the event's type
func Type(type_ string) FilterFunc {
	t := strings.SplitN(type_, "/", 3)
	return func(event *zero.Event) bool {
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

// PostType check post_type
func PostType(postType string) FilterFunc {
	return func(event *zero.Event) bool {
		return postType == event.PostType
	}
}

// DetailType check request_type,message_type,notice_type
func DetailType(detailType string) FilterFunc {
	return func(event *zero.Event) bool {
		return detailType == event.DetailType
	}
}

// SubType check sub_type
func SubType(subType string) FilterFunc {
	return func(event *zero.Event) bool {
		return subType == event.SubType
	}
}