package utils

import (
	"bytes"
	"math/rand"
	"time"
)

func RandomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(RandInt(65, 90)) != temp {
			temp = string(RandInt(65, 90))
			result.WriteString(temp)
			i++
		}
	}
	return result.String()
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// 拷贝字典
func CopyMap(src map[interface{}]interface{}) map[interface{}]interface{} {
	dst := make(map[interface{}]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
