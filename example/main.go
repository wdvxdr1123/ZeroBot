package main

import "github.com/wdvxdr1123/ZeroBot"

func main() {
	ZeroBot.On(func(event ZeroBot.Event) bool {
		if tp, ok := event["post_type"]; !ok || tp.String() != "message" {
			return false
		}
		return event["raw_message"].Str == "复读"
	}).Got("echo", "请输入复读内容", func(event ZeroBot.Event, matcher *ZeroBot.Matcher) ZeroBot.Response {
		ZeroBot.Send(event, matcher.State["echo"])
		return ZeroBot.SuccessResponse
	})
	ZeroBot.Run("ws://127.0.0.1:6700", "")
	select {}
}
