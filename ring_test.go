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

func TestNewRing(t *testing.T) {
	r := newring(128)
	rr := r.r
	for i := 0; i < 128; i++ {
		it := rr.Value.(*eventRingItem)
		if it != nil {
			t.Fatal("unexpected non nil value")
		}
	}
	rr = r.r
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&rr.Value), unsafe.Sizeof(uintptr(0)))),
		unsafe.Pointer(&eventRingItem{
			response: make([]byte, 2),
		}),
	)
	runtime.GC()
	it := rr.Value.(*eventRingItem)
	if it == nil {
		t.Fatal("unexpected nil value")
	}
	it.response = nil
	it.caller = nil
	it = nil
	runtime.GC()
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&rr.Value), unsafe.Sizeof(uintptr(0)))), unsafe.Pointer(nil))
	runtime.GC()
	it = rr.Value.(*eventRingItem)
	if it != nil {
		t.Fatal("unexpected non nil value")
	}
}

func TestRing(t *testing.T) {
	r := newring(128)
	r.loop(8*time.Millisecond, 0, testProcess)
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 256; i++ {
		r.processEvent([]byte{byte(i), byte(i)}, nil)
		time.Sleep(time.Duration(rand.Intn(10)+1) * time.Millisecond)
	}
	time.Sleep(time.Second)
	for i := 0; i < 256; i++ {
		if buf[i] != byte(i) {
			t.Fatal("ring missed", i)
		}
		buf[i] = 0
	}
	for i := 0; i < 256; i++ {
		r.processEvent([]byte{byte(i), byte(i)}, nil)
	}
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
}

func testProcess(response []byte, _ APICaller, _ time.Duration) {
	time.Sleep(time.Duration(rand.Intn(100)+1) * time.Microsecond)
	buf[response[0]] = response[1]
	fmt.Println(response[0], "processed")
}
