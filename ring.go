package zero

import (
	"container/ring"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type eventRing struct {
	sync.Mutex
	r *ring.Ring
}

type eventRingItem struct {
	sync.Mutex
	response     []byte
	caller       APICaller
	isprocessing uintptr
}

func newring(ringLen uint) eventRing {
	n := int(ringLen)
	r := ring.New(n)
	// Initialize the ring with locked eventRing
	for i := 0; i < n; i++ {
		evr := &eventRingItem{}
		evr.Lock()
		r.Value = evr
		r = r.Next()
	}
	return eventRing{r: r}
}

// processEvent 同步向池中放入事件
func (evr *eventRing) processEvent(response []byte, caller APICaller) {
	evr.Lock()
	defer evr.Unlock()
	r := evr.r
	it := r.Value.(*eventRingItem)
	if atomic.LoadUintptr(&it.isprocessing) > 0 { // 池满, 丢弃事件
		return
	}
	it.response = response
	it.caller = caller
	atomic.StoreUintptr(&it.isprocessing, 1)
	it.Unlock() // 开始处理事件
	evr.r = r.Next()
}

// handle 循环处理事件
//
//	latency 延迟 latency + (0~1000ms) 再处理事件
func (evr *eventRing) handle(latency time.Duration) {
	r := evr.r
	for {
		it := r.Value.(*eventRingItem)
		it.Lock()
		if latency > 0 {
			time.Sleep(latency + time.Duration(rand.Intn(1000))*time.Millisecond)
		}
		processEventAsync(it.response, it.caller)
		atomic.StoreUintptr(&it.isprocessing, 1)
		r = r.Next()
	}
}
