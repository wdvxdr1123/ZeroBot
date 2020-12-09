package priority

import zero "github.com/wdvxdr1123/ZeroBot"

func init() {
	a := testPlugin{}
	zero.RegisterPlugin(a) // 注册插件
}

type testPlugin struct{}

func (testPlugin) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "priority_test",
		Version:    "0.1.0",
		Details:    "这是优先级测试",
	}
}

func (testPlugin) Start() {
	zero.OnPrefix("1234").
		Handle(
			func(matcher *Matcher, event Event, state State) zero.Response {
				zero.Send(event, "这是触发器A")
				return zero.FinishResponse
			},
		).
		SetBlock(false).
		SetPriority(10)

	zero.OnPrefix("12345").
		Handle(
			func(matcher *Matcher, event Event, state State) zero.Response {
				zero.Send(event, "这是触发器B")
				return zero.FinishResponse
			},
		).
		SetBlock(true).SetPriority(12)

	zero.OnPrefix("123456").
		Handle(
			func(matcher *Matcher, event Event, state State) zero.Response {
				zero.Send(event, "这是触发器C")
				return zero.FinishResponse
			},
		).
		SetPriority(15)
}
