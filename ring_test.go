package zero

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

var buf [256]byte

func TestRing(t *testing.T) {
	r := newring(128)
	go r.testHandle(time.Millisecond, time.Second)
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
	it := r.Value.(*eventRingItem)
	if !atomic.CompareAndSwapInt64(&it.recvtime, 0, time.Now().UnixNano()) { // 池满, 丢弃事件
		fmt.Println("drop", response[0])
		return
	}
	it.response = response
	it.caller = caller
	fmt.Println("fill", response[0])
	it.Unlock() // 开始处理事件
	evr.r = r.Next()
}

// handle 循环处理事件
//
//	latency 延迟 latency + (0~1000ms) 再处理事件
func (evr *eventRing) testHandle(latency time.Duration, maxwait time.Duration) {
	r := evr.r
	for {
		it := r.Value.(*eventRingItem)
		it.Lock()
		fmt.Println("no:", it.response[0], "diff:", time.Now().UnixNano()-it.recvtime, "max:", int64(maxwait/time.Nanosecond))
		if time.Now().UnixNano()-it.recvtime < int64(maxwait/time.Nanosecond) {
			time.Sleep(latency)
			buf[it.response[0]] = it.response[1]
			fmt.Println(it.response[0], "processed")
		} else {
			fmt.Println("skip", it.response[0])
		} // 等待时间太长，不做处理直接跳过
		it.response = nil
		it.caller = nil
		atomic.StoreInt64(&it.recvtime, 0)
		r = r.Next()
		runtime.GC()
	}
}
