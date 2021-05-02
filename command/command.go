package command

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

// CommandPrefix is the prefix of the command system.
// For Example:
//     CommandPrefix: `>` (default value)
//
//     `>github issue #633` will trigger the command system.
//     `?github pr #114` won't trigger the command system.
//
// You should set this value to your favorite prefix.
var CommandPrefix = ">"

// trigger is the Command system entrypoint.
var trigger *zero.Matcher

// the root of the Command system Trie.
var root *Trie

func init() {
	root = &Trie{ // init the root.
		command:  nil,
		children: map[string]*Trie{},
	}

	trigger = &zero.Matcher{
		Temp:     false,
		Block:    true,
		Priority: 0,
		Type:     zero.Type("message"),
		Rules: []zero.Rule{func(ctx *zero.Ctx) bool {
			return zero.PrefixRule(CommandPrefix)(ctx)
		}},
		Handler: handle,
		Engine:  nil,
	}
	zero.StoreMatcher(trigger)
}

type Command struct {
	// Name is the name of the command
	Name string

	Short string
	Long  string

	Run func(ctx *Ctx)

	trie *Trie
}

// AddRootCommand 添加根指令
func AddRootCommand(cmd *Command) {
	root.add(cmd)
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmd *Command) {
	c.trie.add(cmd)
}

func isFlag(s string) bool { return len(s) > 0 && s[0] == '-' } // nolint

// Help get the help info of the command
func Help(c *Command) string {
	// todo(wdvxdr): impl this
	return ""
}

func handle(ctx *zero.Ctx) {
	args := Parse(ctx.State["args"].(string))
	var (
		i    int
		trie = root
	)
	for i = 0; i < len(args); i++ {
		if trie.check(args[i]) {
			trie = trie.visit(args[i])
		} else {
			break
		}
	}

	// If the trie goes to the end, it means
	// that the command has no child command,
	// so simply call the run function. Else
	// print the help information or prompt
	// the possible command.
	// todo(wdvxdr): impl the not end situation.
	if trie.end() {
		if i == len(args) {
			args = []string{}
		} else {
			args = args[i+1:]
		}
		trie.command.Run(&Ctx{
			Ctx:  *ctx,
			Cmd:  trie.command,
			Args: args,
		})
	} else {
		panic("todo!")
	}
}
