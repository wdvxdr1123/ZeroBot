package ttl

import "time"

// Item represents a record in the cache map
type Item[V any] struct {
	exp   time.Time // expired time
	value V         // value of the item
}

func (it *Item[V]) expired() bool {
	return it.exp.Before(time.Now())
}
