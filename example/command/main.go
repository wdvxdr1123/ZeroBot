package command

import (
	"flag"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	zero.OnCommand("github").Handle(func(ctx *zero.Ctx) {
		fset := flag.FlagSet{}
		var (
			owner string
			repo  string
		)
		fset.StringVar(&owner, "o", "wdvxdr1123", "")
		fset.StringVar(&repo, "r", "ZeroBot", "")
		arguments := shell.Parse(ctx.State["args"].(string))
		err := fset.Parse(arguments)
		if err != nil {
			return
		}
		ctx.Send(message.Text("github\n" +
			"owner: " + owner + "\n" +
			"repo: " + repo,
		))
	})
}
