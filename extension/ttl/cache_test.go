package ttl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	cache := NewCache[string, string](time.Second)

	data := cache.Get("hello")
	assert.Equal(t, "", data)

	cache.Set("hello", "world")
	data = cache.Get("hello")
	assert.Equal(t, "world", data)
}

func TestExpiration(t *testing.T) {
	cache := NewCache[string, string](time.Second)

	cache.Set("x", "1")
	cache.Set("y", "z")
	cache.Set("z", "3")

	<-time.After(500 * time.Millisecond)
	val := cache.Get("x")
	assert.Equal(t, "1", val)

	<-time.After(time.Second)
	val = cache.Get("x")
	assert.Equal(t, "", val)
	val = cache.Get("y")
	assert.Equal(t, "", val)
	val = cache.Get("z")
	assert.Equal(t, "", val)
	assert.Equal(t, 0, len(cache.items))
}
