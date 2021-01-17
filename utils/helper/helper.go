package helper

import (
	"bytes"
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(rune(RandInt(65, 90))) != temp {
			temp = string(rune(RandInt(65, 90)))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// BytesToString 没有内存开销的转换
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes 没有内存开销的转换
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
