package music

import (
	"fmt"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var _ = zero.OnCommandGroup([]string{"music", "点歌"}).
	SetBlock(true).
	SetPriority(8).
	Handle(func(matcher *Matcher, event Event, state State) Response {
		var cmd extension.CommandModel
		err := state.Parse(&cmd)
		if err != nil {
			zero.Send(event, fmt.Sprintf("处理 %v 命令发生错误: %v", cmd.Command, err))
		}

		if cmd.Args == "" { // 未填写歌曲名,索取歌曲名
			zero.Send(event, message.Message{message.Text("请输入要点的歌曲!")})
			next := matcher.FutureEvent("message", zero.CheckUser(event.UserID))
			recv, cancel := next.Repeat()
			for e := range recv {
				msg := e.Message.ExtractPlainText()
				if msg != "" {
					cmd.Args = msg
					cancel()
					continue
				}
				zero.Send(event, "歌曲名不合法oxo")
			}
		}

		zero.Send(event, message.Music("163", QueryNeteaseMusic(cmd.Args)))
		return zero.FinishResponse
	})
