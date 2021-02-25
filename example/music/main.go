package music

import (
	"fmt"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var limit = rate.NewManager(time.Minute*1, 1)

var _ = zero.OnCommandGroup([]string{"music", "点歌"}).
	SetBlock(true).
	SetPriority(8).
	Handle(func(ctx *zero.Ctx) {
		if !limit.Load(ctx.Event.UserID).Acquire() {
			ctx.Send("您的请求太快，请稍后重试0x0...")
			return
		}
		var cmd extension.CommandModel
		err := ctx.Parse(&cmd)
		if err != nil {
			ctx.Send(fmt.Sprintf("处理 %v 命令发生错误: %v", cmd.Command, err))
		}

		if cmd.Args == "" { // 未填写歌曲名,索取歌曲名
			ctx.Send(message.Message{message.Text("请输入要点的歌曲!")})
			next := ctx.FutureEvent("message", ctx.CheckSession())
			recv, cancel := next.Repeat()
			for e := range recv {
				msg := e.Message.ExtractPlainText()
				if msg != "" {
					cmd.Args = msg
					cancel()
					continue
				}
				ctx.Send("歌曲名不合法oxo")
			}
		}
		zero.RangeBot(func(id int64, ctx2 *zero.Ctx) bool { // test the range bot function
			ctx2.SendGroupMessage(ctx.Event.GroupID, message.Music("163", queryNeteaseMusic(cmd.Args)))
			return true
		})
		//ctx.Send(message.Music("163", queryNeteaseMusic(cmd.Args)))
	})
