// Package manager provides a simple group plugin Manager.
package manager

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	bucket   = kv.New("Manager")
	managers = map[string]*Manager{}
	mu       = sync.RWMutex{}
)

// New returns Manager with settings.
func New(service string, o *Options) *Manager {
	data, _ := bucket.Get([]byte(service))
	m := &Manager{
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

// Manager is the plugin group manager.
type Manager struct {
	sync.RWMutex
	service string
	options Options
	states  map[int64]bool
}

// Enable enables a group to pass the Manager.
func (m *Manager) Enable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = true
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

// Disable disables a group to pass the Manager.
func (m *Manager) Disable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = false
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

// Handler 返回 预处理器
func (m *Manager) Handler() zero.Rule {
	return func(ctx *zero.Ctx) bool {
		m.RLock()
		ctx.State["manager"] = m
		if st, ok := m.states[ctx.Event.GroupID]; ok {
			m.RUnlock()
			return st
		}
		m.RUnlock()
		if m.options.DisableOnDefault {
			m.Disable(ctx.Event.GroupID)
		} else {
			m.Enable(ctx.Event.GroupID)
		}
		return !m.options.DisableOnDefault
	}
}

// Lookup returns a Manager by the service name, if
// not exist, it will returns nil.
func Lookup(service string) (*Manager, bool) {
	mu.RLock()
	defer mu.RUnlock()
	m, ok := managers[service]
	return m, ok
}

// ForEach iterates through managers.
func ForEach(iterator func(key string, manager *Manager) bool) {
	mu.RLock()
	m := copyMap(managers)
	mu.RUnlock()
	for k, v := range m {
		if !iterator(k, v) {
			return
		}
	}
}

func copyMap(m map[string]*Manager) map[string]*Manager {
	ret := make(map[string]*Manager, len(m))
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
			b[7] |= 0x80
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

func init() {
	engine := zero.New()

	engine.OnCommandGroup([]string{"启用", "enable"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.Send("没有找到指定服务!")
			}
			service.Enable(ctx.Event.GroupID)
			ctx.Send(message.Text("已启用服务: " + model.Args))
		})

	engine.OnCommandGroup([]string{"禁用", "disable"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.Send("没有找到指定服务!")
			}
			service.Disable(ctx.Event.GroupID)
			ctx.Send(message.Text("已关闭服务: " + model.Args))
		})

	engine.OnCommandGroup([]string{"服务列表", "service_list"}, zero.AdminPermission, zero.OnlyGroup).
		Handle(func(ctx *zero.Ctx) {
			msg := `---服务列表---`
			i := 0
			ForEach(func(key string, _ *Manager) bool {
				i++
				msg += "\n" + strconv.Itoa(i) + `: ` + key
				return true
			})
			ctx.Send(message.Text(msg))
		})
}
