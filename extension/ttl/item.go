package ttl

import "time"

// Item represents a record in the cache map
type Item struct {
	exp   time.Time   // expired time
	value interface{} // value of the item
}

func (it *Item) expired() bool {
	return it.exp.Before(time.Now())
}
