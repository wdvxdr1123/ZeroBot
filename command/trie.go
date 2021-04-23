package command

// Trie is a path trie of the command system.
type Trie struct {
	command *Command

	children map[string]*Trie
}

// check checks the if child trie has specific item.
func (t *Trie) check(name string) bool {
	if t.children == nil {
		return false
	}
	_, ok := t.children[name]
	return ok
}

// add adds a command to the trie, if name conflicted
// it panics, the child trie will not alloc at first,
// to save the space, the child trie will alloc when necessary.
func (t *Trie) add(cmd *Command) {
	if t.check(cmd.Name) {
		panic("command name conflicted")
	}
	if t.children == nil {
		t.children = make(map[string]*Trie)
	}

	cmd.trie = &Trie{
		command: cmd,
		// lazy init the map to save some space.
		// children: make(map[string]*Trie),
	}
	t.children[cmd.Name] = cmd.trie
}

func (t *Trie) visit(name string) *Trie {
	if t.children == nil {
		panic("no children command")
	}
	return t.children[name]
}

// end returns if the trie has no child.
func (t *Trie) end() bool {
	return t.children == nil
}
