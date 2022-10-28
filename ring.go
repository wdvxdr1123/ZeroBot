package zero

import (
	"container/ring"
	"sync"
)

type eventRing struct {
	sync.Mutex
	r *ring.Ring
}

type eventRingItem struct {
	sync.Mutex
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
	it.response = response
	it.caller = caller
	it.Unlock() // 开始处理事件
	evr.r = r.Next()
}

// handle 循环处理事件
func (evr *eventRing) handle() {
	r := evr.r
	for {
		it := r.Value.(*eventRingItem)
		it.Lock()
		processEventAsync(it.response, it.caller)
		it.Unlock()
		r = r.Next()
	}
}
