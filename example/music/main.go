package music

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var _ = zero.OnCommandGroup([]string{"music", "点歌"}).
	SetBlock(true).
	SetPriority(8).
	Handle(func(matcher *Matcher, event Event, state State) Response {
		if songName, ok := state["args"].(string); ok {
			if songName == "" {
				zero.Send(event, "请输入要点的歌曲!")
				next := matcher.FutureEvent("message", zero.CheckUser(event.UserID))
				recv, cancel := next.Repeat()
				for e := range recv {
					msg := e.Message.ExtractPlainText()
					if msg != "" {
						songName = msg
						cancel()
						continue
					}
					zero.Send(event, "歌曲名不合法oxo")
				}
			}
			zero.Send(event, message.Music("163", QueryNeteaseMusic(songName)))
			return zero.FinishResponse
		}
		return zero.FinishResponse
	})
