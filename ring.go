package zero

import (
	"container/ring"
	"math/rand"
	"runtime"
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
	recvtime int64 // atomic field
	response []byte
	caller   APICaller
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
	if !atomic.CompareAndSwapInt64(&it.recvtime, 0, time.Now().UnixNano()) { // 池满, 丢弃事件
		return
	}
	it.response = response
	it.caller = caller
	it.Unlock() // 开始处理事件
	evr.r = r.Next()
}

// handle 循环处理事件
//
//	latency 延迟 latency + (0~100ms) 再处理事件
func (evr *eventRing) handle(latency, maxwait time.Duration) {
	r := evr.r
	for {
		it := r.Value.(*eventRingItem)
		it.Lock()
		if time.Now().UnixNano()-it.recvtime < int64(maxwait/time.Nanosecond) {
			time.Sleep(latency + time.Duration(rand.Intn(100))*time.Millisecond)
			processEventAsync(it.response, it.caller, maxwait)
		} // 等待时间太长，不做处理直接跳过
		it.response = nil
		it.caller = nil
		atomic.StoreInt64(&it.recvtime, 0)
		r = r.Next()
		runtime.GC()
	}
}
