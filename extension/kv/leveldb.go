// Package kv provides a simple wrap of goleveldb for multi bucket database
package kv

import (
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var db *leveldb.DB

func init() {
	var err error
	db, err = leveldb.OpenFile(".db", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Bucket is the interface of the database bucket
type Bucket interface {
	Get(k []byte) ([]byte, error)
	Put(k []byte, v []byte) error
	Delete(k []byte) error
	Iterator(func(k, v []byte) bool)
}

type bucket struct {
	name []byte
}

var defaultBucket = New("\x01")

// New returns a Bucket with specific name.
func New(name string) Bucket {
	return &bucket{name: []byte(name)}
}

func pack(name []byte, k []byte) []byte { return append(append(name, 0x02), k...) }

// Get returns a value for the given key from the default bucket.
func Get(k []byte) ([]byte, error) { return defaultBucket.Get(k) }

// Get returns a value for the given key from the bucket.
func (b *bucket) Get(k []byte) ([]byte, error) {
	return db.Get(pack(b.name, k), nil)
}

// Put push/update a key value pair to the default bucket.
func Put(k []byte, v []byte) error { return defaultBucket.Put(k, v) }

// Put push/update a key value pair to the bucket.
func (b *bucket) Put(k []byte, v []byte) error {
	return db.Put(pack(b.name, k), v, nil)
}

// Delete deletes a key from the default bucket.
func Delete(k []byte) error { return defaultBucket.Delete(k) }

// Delete deletes a key from the bucket.
func (b *bucket) Delete(k []byte) error {
	return db.Delete(pack(b.name, k), nil)
}

func (b *bucket) Iterator(iter func(k, v []byte) bool) {
	iterator := db.NewIterator(util.BytesPrefix(append(b.name, 0x02)), nil)
	defer iterator.Release()
	for iterator.Next() {
		if !iter(iterator.Key()[len(b.name)+1:], iterator.Value()) {
			break
		}
	}
}
