package repeat

import (
	"github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/manager"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	var m = manager.New("echo", nil)
	zero.AddHook(m)
	zero.RegisterPlugin(testPlugin{}) // 注册插件
}

type testPlugin struct{}

func (_ testPlugin) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "test",
		Version:    "0.1.0",
		Details:    "这是一个测试复读插件",
	}
}

func (_ testPlugin) Start() { // 插件主体
	zero.OnCommand("开启复读").SetBlock(true).SetPriority(10).
		Handle(func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			stop := zero.NewFutureEvent("message/group", 8, true,
				zero.CommandRule("关闭复读"),      // 关闭复读指令
				zero.CheckUser(event.UserID)). // 只有开启者可以关闭复读模式
				Next()                         // 关闭需要一次

			echo, cancel := matcher.FutureEvent("message/group",
				zero.CheckUser(event.UserID)). // 只复读开启复读模式的人的消息
				Repeat()                       // 不断监听复读
			zero.Send(event, "已开启复读模式!")
			for {
				select {
				case e := <-echo: // 接收到需要复读的消息
					zero.Send(event, e.RawMessage)
				case <-stop: // 收到关闭复读指令
					cancel() // 取消复读监听
					zero.SendGroupForwardMessage(event.GroupID, []message.MessageSegment{
						message.CustomNode("bot", zero.BotConfig.SelfID, "取消复读"),
					})
					return zero.FinishResponse // 返回
				}
			}
		})
}
