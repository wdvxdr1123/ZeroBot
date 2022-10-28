package zero

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

var buf [256]byte

func TestRing(t *testing.T) {
	r := newring(128)
	go r.testHandle(0)
	for i := 0; i < 256; i++ {
		r.testProcessEvent([]byte{byte(i), byte(i)}, nil)
	}
	time.Sleep(time.Second)
	for i := 0; i < 128; i++ {
		if buf[i] != byte(i) {
			t.Fatal("ring missed", i)
		}
		buf[i] = 0
	}
	for i := 128; i < 256; i++ {
		if buf[i] != 0 {
			t.Fatal("unexpected ring value at", i)
		}
	}
	for i := 0; i < 256; i++ {
		r.testProcessEvent([]byte{byte(i), byte(i)}, nil)
		time.Sleep(time.Millisecond)
	}
	for i := 0; i < 256; i++ {
		if buf[i] != byte(i) {
			t.Fatal("ring missed", i)
		}
	}
}

// processEvent 同步向池中放入事件
func (evr *eventRing) testProcessEvent(response []byte, caller APICaller) {
	evr.Lock()
	defer evr.Unlock()
	r := evr.r
	it := r.Value.(*eventRingItem)
	if !atomic.CompareAndSwapUintptr(&it.isprocessing, 0, 1) { // 池满, 丢弃事件
		return
	}
	it.response = response
	it.caller = caller
	it.Unlock() // 开始处理事件
	evr.r = r.Next()
}

// handle 循环处理事件
//
//	latency 延迟 latency + (0~1000ms) 再处理事件
func (evr *eventRing) testHandle(latency time.Duration) {
	r := evr.r
	for {
		it := r.Value.(*eventRingItem)
		it.Lock()
		if latency > 0 {
			time.Sleep(latency + time.Duration(rand.Intn(1000))*time.Millisecond)
		}
		buf[it.response[0]] = it.response[1]
		println(it.response[0], "processed")
		atomic.StoreUintptr(&it.isprocessing, 0)
		r = r.Next()
	}
}
