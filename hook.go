package zero

import (
	"runtime"
)

type Hooker interface {
	Hook() Rule
}

type multiHooker struct {
	hookers []Hooker
}

func (m *multiHooker) Hook() Rule {
	return func(event *Event, state State) bool {
		for i := range m.hookers {
			if m.hookers[i].Hook()(event, state) == false {
				return false
			}
		}
		return true
	}
}

func MultiHooker(h ...Hooker) Hooker {
	return &multiHooker{hookers: h}
}

var hooks = map[string]Hooker{}

func AddHook(hookers ...Hooker) Hooker {
	_, file, _, _ := runtime.Caller(1) // who calls this method
	h := MultiHooker(hookers...)
	hooks[file] = h
	return h
}
