package command

import zero "github.com/wdvxdr1123/ZeroBot"

type Ctx struct {
	zero.Ctx

	Cmd  *Command
	Args []string
}
