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
	command.AddRootCommand(root)

	child1 := &command.Command{
		Name: "test1",
		Run: func(ctx *command.Ctx) {
			ctx.SendChain(
				message.Text("test1"),
			)
		},
	}
	root.AddCommand(child1)

	// child2 := &command.Command{
	// 	 Name: "test1",
	//	 Run:  nil,
	// }
	// root.AddCommand(child2) // should panic: conflicted name
}
