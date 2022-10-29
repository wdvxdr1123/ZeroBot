package zero

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

var buf [256]byte

func TestRing(t *testing.T) {
	r := newring(128)
	for i := 0; i < 256; i++ {
		r.testProcessEvent([]byte{byte(i), byte(i)}, nil)
	}
	go r.testHandle(8 * time.Millisecond)
	time.Sleep(8 * time.Millisecond * 300)
	for i := 128; i < 256; i++ {
		if buf[i] != byte(i) {
			t.Fatal("ring missed", i)
		}
		buf[i] = 0
	}
	for i := 0; i < 128; i++ {
		if buf[i] != 0 {
			t.Fatal("unexpected ring value at", i)
		}
	}
	for i := 0; i < 256; i++ {
		r.testProcessEvent([]byte{byte(i), byte(i)}, nil)
		time.Sleep(10 * time.Millisecond)
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
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&r.Value), unsafe.Sizeof(uintptr(0)))),
		unsafe.Pointer(&eventRingItem{
			response: response,
			caller:   caller,
		}),
	)
	fmt.Println("fill", response[0])
	evr.r = r.Next()
}

// handle 循环处理事件
//
//	latency 延迟 latency + (0~1000ms) 再处理事件
func (evr *eventRing) testHandle(latency time.Duration) {
	r := evr.r
	for range time.NewTicker(latency).C {
		it := r.Value.(*eventRingItem)
		if it == nil { // 还未有消息
			continue
		}
		time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		buf[it.response[0]] = it.response[1]
		fmt.Println(it.response[0], "processed")
		it.response = nil
		it.caller = nil
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&r.Value), unsafe.Sizeof(uintptr(0)))), unsafe.Pointer(nil))
		r = r.Next()
		runtime.GC()
	}
}
