package zero

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type eventRing struct {
	sync.Mutex
	c uintptr
	r []*eventRingItem
	i uintptr
	p []eventRingItem
}

type eventRingItem struct {
	response []byte
	caller   APICaller
}

func newring(ringLen uint) eventRing {
	return eventRing{
		r: make([]*eventRingItem, ringLen),
		p: make([]eventRingItem, ringLen+1),
	} // 同一节点, 每 ringLen*(ringLen+1) 轮将共用同一 buffer
}

// processEvent 同步向池中放入事件
//
//go:nosplit
func (evr *eventRing) processEvent(response []byte, caller APICaller) {
	evr.Lock()
	defer evr.Unlock()
	r := evr.c % uintptr(len(evr.r))
	p := evr.i % uintptr(len(evr.p))
	evr.p[p] = eventRingItem{
		response: response,
		caller:   caller,
	}
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&evr.r[r])), unsafe.Pointer(&evr.p[p]))
	evr.c++
	evr.i++
}

// loop 循环处理事件
//
//	latency 延迟 latency 再处理事件
func (evr *eventRing) loop(latency, maxwait time.Duration, process func([]byte, APICaller, time.Duration)) {
	go func(r []*eventRingItem) {
		c := uintptr(0)
		if latency < time.Millisecond {
			latency = time.Millisecond
		}
		for range time.NewTicker(latency).C {
			i := c % uintptr(len(r))
			it := (*eventRingItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r[i]))))
			if it == nil { // 还未有消息
				continue
			}
			process(it.response, it.caller, maxwait)
			it.response = nil
			it.caller = nil
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r[i])), unsafe.Pointer(nil))
			c++
			runtime.GC()
		}
	}(evr.r)
}
