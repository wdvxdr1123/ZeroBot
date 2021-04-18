package command

import (
	"github.com/wdvxdr1123/ZeroBot/command"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	root := &command.Command{
		Name: "test",
		Run: func(ctx *command.Ctx) {
			ctx.SendChain(
				message.Text("test"),
			)
		},
	}

	child1 := &command.Command{
		Name: "test1",
		Run: func(ctx *command.Ctx) {
			ctx.SendChain(
				message.Text("test1"),
			)
		},
	}

	root.AddCommand(child1)
	command.AddRootCommand(root).FirstPriority()
}
