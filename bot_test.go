package ZeroBot

import (
	"testing"
)

func TestRun(t *testing.T) {
	On(func(event Event) bool {
		if tp, ok := event["post_type"]; !ok || tp.String() != "message" {
			return false
		}
		return event["raw_message"].Str == "复读"
	}).Got("echo","请输入复读内容",func(event Event, matcher *Matcher) Response {
		Send(event, matcher.State["echo"])
		return SuccessResponse
	})
	Run("ws://127.0.0.1:6700", "")
	select {}
}
