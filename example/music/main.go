package music

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	zero.RegisterPlugin(music{}) // 注册插件
}

type music struct{}

func (music) GetPluginInfo() zero.PluginInfo { // 返回插件信息
	return zero.PluginInfo{
		Author:     "wdvxdr1123",
		PluginName: "music",
		Version:    "0.1.0",
		Details:    "点歌",
	}
}

func (m music) Start() {
	zero.OnCommand("点歌", zero.OnlyGroup).SetBlock(true).SetPriority(8).
		Handle(
			func(matcher *Matcher, event Event, state State) Response {
				if songName, ok := state["args"].(string); ok {
					if songName != "" {
						state["song_name"] = songName
					}
					return zero.SuccessResponse
				} else {
					return zero.FinishResponse
				}
			},
		).
		Got("song_name", "请输入要点的歌曲!",
			func(matcher *Matcher, event Event, state State) Response {
				MusicID := QueryNeteaseMusic(state["song_name"].(string))
				if MusicID != "" {
					zero.Send(event, message.Music("163", MusicID))
				} else {
					zero.Send(event, "没有找到那首歌哦OxO")
				}
				return zero.FinishResponse
			},
		)

	zero.OnCommand("music").SetBlock(true).SetPriority(8).
		Handle(
			func(matcher *Matcher, event Event, state State) Response {
				if songName, ok := state["args"].(string); ok {
					if songName == "" {
						zero.Send(event, "请输入要点的歌曲!")
						zero.ForMessage().Rule(zero.CheckUser(event.UserID)).Handle(
							func(m message.Message) Response {
								msg := ""
								for _, val := range m {
									if val.Type == "text" {
										msg += val.Data["text"]
									}
								}
								if msg != "" {
									songName = msg
									return zero.FinishResponse
								}
								zero.Send(event, "歌曲名不合法oxo")
								return zero.RejectResponse
							},
						).Do()
					}
					zero.Send(event, message.Music("163", QueryNeteaseMusic(songName)))
					return zero.FinishResponse
				}
				return zero.FinishResponse
			},
		)
}
