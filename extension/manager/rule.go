package manager

import (
	"bytes"
	"encoding/binary"
	"io"
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/kv"
)

var bucket = kv.New("manager")

type Manager interface {
	Hook() zero.Rule
	Enable(groupID int64)
	Disable(groupID int64)
}

func New(service string, o Options) Manager {
	data, err := bucket.Get([]byte(service))
	if err != nil {
		panic(err)
	}
	return &manager{
		service: service,
		options: o,
		states:  unpack(data),
	}
}

type manager struct {
	sync.RWMutex
	service string
	options Options
	states  map[int64]bool
}

func (m *manager) Enable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = true
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

func (m *manager) Disable(groupID int64) {
	m.Lock()
	defer m.Unlock()
	m.states[groupID] = true
	_ = bucket.Put([]byte(m.service), pack(m.states))
}

func (m *manager) Hook() zero.Rule {
	return func(event *zero.Event, state zero.State) bool {
		m.RLocker()
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

func pack(m map[int64]bool) []byte {
	var (
		buf bytes.Buffer
		b   = make([]byte, 8)
	)
	for k, v := range m {
		binary.LittleEndian.PutUint64(b, uint64(k))
		if v {
			b[7] = b[7] | 0x8
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
