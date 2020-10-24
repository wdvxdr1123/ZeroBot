package plugin

import "github.com/wdvxdr1123/ZeroBot"

func init() {
	a := testPlugin{}
	ZeroBot.RegisterPlugin(a)
}

type testPlugin struct{}

func (testPlugin) GetPluginInfo() ZeroBot.PluginInfo {
	return ZeroBot.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "test",
		Version:    "0.1.0",
		Details:    "这是一个测试插件",
	}
}

func (testPlugin) Start() {
	ZeroBot.OnPrefix("复读", "echo", "fudu").
		Got("echo", "请输入复读内容",
			func(event ZeroBot.Event, matcher *ZeroBot.Matcher) ZeroBot.Response {
				ZeroBot.Send(event, matcher.State["echo"])
				return ZeroBot.SuccessResponse
			},
		)
}
