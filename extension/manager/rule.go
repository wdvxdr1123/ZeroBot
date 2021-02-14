// Package manager provides a simple group plugin manager.
package manager

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
)

var (
	bucket   = kv.New("manager")
	managers = map[string]Manager{}
	mu       = sync.RWMutex{}
)

// Manager is the interface of plugin group manager.
type Manager interface {
	zero.Hooker
	Enable(groupID int64)
	Disable(groupID int64)
}

// New returns Manager with settings.
func New(service string, o *Options) Manager {
	data, _ := bucket.Get([]byte(service))
	m := &manager{
		service: service,
		options: func() Options {
			if o == nil {
				return Options{}
			}
			return *o
		}(),
		states: unpack(data),
	}
	mu.Lock()
	defer mu.Unlock()
	managers[service] = m
	return m
}

type manager struct {
	sync.RWMutex
	service string
	options Options
	states  map[int64]bool
}

// Enable enables a group to pass the manager.
func (m *manager) Enable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = true
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

// Disable disables a group to pass the manager.
func (m *manager) Disable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = false
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

// Hook impls the zero.Hooker, so you can use zero.AddHook to the file.
func (m *manager) Hook() zero.Rule {
	return func(event *zero.Event, state zero.State) bool {
		m.RLock()
		state["manager"] = Manager(m)
		if st, ok := m.states[event.GroupID]; ok {
			m.RUnlock()
			return st
		}
		m.RUnlock()
		if m.options.DisableOnDefault {
			m.Disable(event.GroupID)
		} else {
			m.Enable(event.GroupID)
		}
		return !m.options.DisableOnDefault
	}
}

// Lookup returns a Manager by the service name, if
// not exist, it will returns nil.
func Lookup(service string) Manager {
	mu.RLock()
	defer mu.RUnlock()
	return managers[service]
}

// ForEach iterates through managers.
func ForEach(iterator func(key string, manager Manager) bool) {
	mu.RLock()
	m := copyMap(managers)
	mu.RUnlock()
	for k, v := range m {
		if !iterator(k, v) {
			return
		}
	}
}

func copyMap(m map[string]Manager) map[string]Manager {
	var ret = make(map[string]Manager, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func pack(m map[int64]bool) []byte {
	var (
		buf bytes.Buffer
		b   = make([]byte, 8)
	)
	for k, v := range m {
		binary.LittleEndian.PutUint64(b, uint64(k))
		if v {
			b[7] = b[7] | 0x80
		}
		buf.Write(b[:8])
	}
	return buf.Bytes()
}

func unpack(v []byte) map[int64]bool {
	var (
		m      = make(map[int64]bool)
		b      = make([]byte, 8)
		reader = bytes.NewReader(v)
		k      uint64
	)
	for {
		_, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		k = binary.LittleEndian.Uint64(b)
		m[int64(k&0x7fff_ffff_ffff_ffff)] = k&8000_0000_0000_0000 != 0
	}
	return m
}
