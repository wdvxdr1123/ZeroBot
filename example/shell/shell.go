package shell

import (
	"github.com/sirupsen/logrus"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// ShellRule Example
// 本插件仅作为演示
// Note: 只有带 flag 的Tag的字段才会注册,
// 支支持 bool, int, string, float64 四种类型

type Ping struct {
	T       bool   `flag:"t"`
	Timeout int    `flag:"w"`
	Host    string `flag:"host"`
}

func init() {
	zero.OnShell("ping", Ping{}).Handle(func(ctx *zero.Ctx) {
		ping := ctx.State["flag"].(*Ping) // Note: 指针类型
		logrus.Infoln("ping host:", ping.Host)
		logrus.Infoln("ping timeout:", ping.Timeout)
		logrus.Infoln("ping t:", ping.T)
		for i, v := range ctx.State["args"].([]string) {
			logrus.Infoln("args", i, ":", v)
		}
	})
}
