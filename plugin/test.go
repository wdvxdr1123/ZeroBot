package plugin

import (
	"github.com/wdvxdr1123/ZeroBot"
)

func init() {
	a := testPlugin{}
	ZeroBot.RegisterPlugin(a) // 注册插件
}

type testPlugin struct{}

func (testPlugin) GetPluginInfo() ZeroBot.PluginInfo { // 返回插件信息
	return ZeroBot.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "test",
		Version:    "0.1.0",
		Details:    "这是一个测试插件",
	}
}

func (testPlugin) Start() { // 插件主体
	ZeroBot.OnPrefix([]string{"复读", "echo", "fudu"}, ZeroBot.OnlyToMe()).
		Got("echo", "请输入复读内容",
			func(event ZeroBot.Event, matcher *ZeroBot.Matcher) ZeroBot.Response {
				event.Message.Reply(matcher.State["echo"])
				return ZeroBot.SuccessResponse
			},
		)
}
