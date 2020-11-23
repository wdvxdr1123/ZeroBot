package repeat

import (
	"github.com/wdvxdr1123/ZeroBot"
)

func init() {
	a := testPlugin{}
	zero.RegisterPlugin(a) // 注册插件
}

type testPlugin struct{}

func (testPlugin) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "test",
		Version:    "0.1.0",
		Details:    "这是一个测试复读插件",
	}
}

func (testPlugin) Start() { // 插件主体
	zero.OnPrefixGroup([]string{"复读", "fudu"}, zero.OnlyToMe).
		Got(
			"echo",
			"请输入复读内容",
			func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
				zero.Send(event, matcher.State["echo"])
				return zero.FinishResponse
			},
		)

	zero.OnCommand("echo").Handle(handleCommand)

	zero.OnCommand("cmd_echo").Handle(handleCommand)

	zero.OnSuffix("复读").Handle(func(_ *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		zero.Send(event, state["args"])
		return zero.FinishResponse
	})

	zero.OnFullMatch("你是谁", zero.OnlyToMe).Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		zero.Send(event, "我是一个复读机~~~")
		echo := matcher.Get("我想要复读你的话!")
		zero.Send(event, echo)
		return zero.FinishResponse
	})

	zero.OnCommand("next_echo").Receive(receiveNextEcho)
}

func receiveNextEcho(_ *zero.Matcher, event zero.Event, _ zero.State) zero.Response {
	zero.Send(event, event.Message)
	if event.Message.CQString() == "哈哈哈" {
		return zero.RejectResponse
	}
	return zero.FinishResponse
}

func handleCommand(_ *zero.Matcher, event zero.Event, state zero.State) zero.Response {
	zero.Send(event, state["args"])
	return zero.FinishResponse
}
