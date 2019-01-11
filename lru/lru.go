// Package lru implement a thread safe lru cache inspired by https://github.com/golang/groupcache/tree/master/lru
package lru

import (
	"container/list"
	"sync"
)

const (
	// default LRU size
	Default_LRU_Size = 1024
)

// LRU implements a thread safe fixed size LRU cache
type LRU struct {
	lock      sync.RWMutex
	capacity  int
	evictList *list.List
	items     map[interface{}]*list.Element
	misses    int
	hits      int
}

// payload contains the value evictList hold
type payload struct {
	key   interface{}
	value interface{}
}

// NewLRUCache return a given size LRU
func NewLRUCache(size int) *LRU {
	if size <= 0 {
		size = Default_LRU_Size
	}
	return &LRU{
		capacity:  size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
	}
}

// return the LRU running information
func (l *LRU) Info() (hits int, misses int, maxSize int, currentSize int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.hits, l.misses, l.capacity, l.Len()
}

// return the LRU max capacity
func (l *LRU) Cap() int {
	return l.capacity
}

// Add a new item into LRU
func (l *LRU) Set(key, value interface{}) (evicted bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// key has exists, update it to new value
	if v, ok := l.items[key]; ok {
		l.evictList.MoveToFront(v)
		v.Value.(*payload).value = value
		return false
	}

	v := &payload{
		key:   key,
		value: value,
	}
	item := l.evictList.PushFront(v)
	l.items[key] = item
	if l.evictList.Len() > l.capacity {
		l.removeItem(l.evictList.Back())
		return true
	}
	return false
}

// Get value from LRU by key
func (l *LRU) Get(key interface{}) (value interface{}, ok bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	v, ok := l.items[key]
	if ok {
		l.evictList.MoveToFront(v)
		l.hits += 1
	} else {
		l.misses += 1
		return nil, ok
	}
	return v.Value.(*payload).value, ok
}

// Cotains check if the LRU contains the given key
func (l *LRU) Contains(key interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, ok := l.items[key]
	return ok
}

// Remove the given key item return if the key has existed before
func (l *LRU) Remove(key interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	v, ok := l.items[key]
	if ok {
		l.removeItem(v)
	}
	return ok
}

// Remove and return the oldest item from LRU
func (l *LRU) PopOldest() (key, value interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.evictList.Len() == 0 {
		return nil, nil
	}
	v := l.evictList.Back()
	delete(l.items, v.Value.(*payload).key)
	l.removeItem(l.evictList.Back())
	return v.Value.(*payload).key, v.Value.(*payload).value
}

// return the value if the key exist, otherwise update the key by given value similar with redis SETNX
func (l *LRU) GetOrSet(key, value interface{}) (newValue interface{}, isGet bool) {
	if v, ok := l.Get(key); ok {
		return v, ok
	}
	l.Set(key, value)
	l.misses += 1
	return value, false
}

// return all keys the LRU hold from oldest to newest
func (l *LRU) Keys() []interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()
	keys := make([]interface{}, 0)
	for v := l.evictList.Back(); v != nil; v = v.Prev() {
		keys = append(keys, v.Value.(*payload).key)
	}
	return keys
}

// return the LRU length
func (l *LRU) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.evictList.Len()
}

// Purge use to clear all items in LRU
func (l *LRU) Purge() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for k := range l.items {
		delete(l.items, k)
	}
	l.evictList.Init()
	l.misses = 0
	l.hits = 0
}

// remove item from lru
func (l *LRU) removeItem(e *list.Element) {
	l.evictList.Remove(e)
	delete(l.items, e.Value.(*payload).key)
}
