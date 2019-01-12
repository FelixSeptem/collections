// Package lru implement a thread safe lfu cache
package lfu

import (
	"container/list"
	"sync"
)

const (
	// default LFU size
	Default_LFU_Size = 1024
)

// LFU implements a thread safe fixed size LFU cache
type LFU struct {
	lock      sync.RWMutex
	capacity  int
	evictList *list.List
	items     map[interface{}]*list.Element
	misses    int
	hits      int
}

// payload contains the value evictList hold
type payload struct {
	key       interface{}
	value     interface{}
	frequency uint
}

// NewLFUCache return a given size LFU
func NewLFUCache(size int) *LFU {
	if size <= 0 {
		size = Default_LFU_Size
	}
	return &LFU{
		capacity:  size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
	}
}

// return the LFU running information
func (l *LFU) Info() (hits int, misses int, maxSize int, currentSize int) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.hits, l.misses, l.capacity, l.Len()
}

// return the LFU max capacity
func (l *LFU) Cap() int {
	return l.capacity
}

// Add a new item into LFU
func (l *LFU) Set(key, value interface{}) (evicted bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// key has exists, update it to new value
	if v, ok := l.items[key]; ok {
		v.Value.(*payload).frequency += 1
		v.Value.(*payload).value = value
		l.adjust(v)
		return false
	}
	v := &payload{
		key:   key,
		value: value,
	}
	item := l.evictList.PushBack(v)
	i := l.adjust(item)
	l.items[key] = i
	if l.evictList.Len() > l.capacity {
		l.removeItem(l.evictList.Back())
		return true
	}
	return false
}

// Get value from LFU by key
func (l *LFU) Get(key interface{}) (value interface{}, ok bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	v, ok := l.items[key]
	if ok {
		v.Value.(*payload).frequency += 1
		l.adjust(v)
		l.hits += 1
	} else {
		l.misses += 1
		return nil, ok
	}
	return v.Value.(*payload).value, ok
}

// Cotains check if the LRU contains the given key
func (l *LFU) Contains(key interface{}) bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	_, ok := l.items[key]
	return ok
}

// Remove the given key item return if the key has existed before
func (l *LFU) Remove(key interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	v, ok := l.items[key]
	if ok {
		l.removeItem(v)
	}
	return ok
}

// Remove and return the oldest item from LFU
func (l *LFU) PopOldest() (key, value interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.evictList.Len() == 0 {
		return nil, nil
	}
	v := l.evictList.Back()
	delete(l.items, v.Value.(*payload).key)
	l.removeItem(v)
	return v.Value.(*payload).key, v.Value.(*payload).value
}

// return the value if the key exist, otherwise update the key by given value similar with redis SETNX
func (l *LFU) GetOrSet(key, value interface{}) (newValue interface{}, isGet bool) {
	if v, ok := l.Get(key); ok {
		return v, ok
	}
	l.Set(key, value)
	l.misses += 1
	return value, false
}

// return all keys the LRU hold from oldest to newest
func (l *LFU) Keys() []interface{} {
	l.lock.RLock()
	defer l.lock.RUnlock()
	keys := make([]interface{}, 0)
	for v := l.evictList.Back(); v != nil; v = v.Prev() {
		keys = append(keys, v.Value.(*payload).key)
	}
	return keys
}

// return the LFU length
func (l *LFU) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.evictList.Len()
}

// Purge use to clear all items in LFU
func (l *LFU) Purge() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for k := range l.items {
		delete(l.items, k)
	}
	l.evictList.Init()
}

// adjust the list element to correct location
func (l *LFU) adjust(i *list.Element) *list.Element {
	if i.Prev() == nil {
		return i
	}
	if i.Value.(*payload).frequency >= l.evictList.Front().Value.(*payload).frequency {
		l.evictList.MoveBefore(i, l.evictList.Front())
		return i
	}
	for n := i; n != nil; n = n.Prev() {
		if i.Value.(*payload).frequency < n.Value.(*payload).frequency {
			l.evictList.MoveAfter(i, n)
			return i
		}
	}
	l.evictList.MoveBefore(i, l.evictList.Front())
	return nil
}

// remove item from lru
func (l *LFU) removeItem(e *list.Element) {
	l.evictList.Remove(e)
	delete(l.items, e.Value.(*payload).key)
}
