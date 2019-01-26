// Package arc implement Adaptive Replacement Cache inspired by http://code.activestate.com/recipes/576532/
package arc

import (
	"github.com/FelixSeptem/collections/lfu"
	"github.com/FelixSeptem/collections/lru"
	"sync"
)

const (
	// default ARC size
	Default_ARC_Size = 1024
)

// a fixed size arc(Adaptive Replacement Cache) cache
type ARC struct {
	lock     sync.RWMutex
	capacity int
	p        int

	t1 *lru.LRU
	b1 *lru.LRU

	t2 *lfu.LFU
	b2 *lfu.LFU
}

// NewARC return a given size arc
func NewARCCache(size int) *ARC {
	if size <= 0 {
		size = Default_ARC_Size
	}
	return &ARC{
		capacity: size,
		p:        0,
		t1:       lru.NewLRUCache(size),
		b1:       lru.NewLRUCache(size),
		t2:       lfu.NewLFUCache(size),
		b2:       lfu.NewLFUCache(size),
	}
}

// return the ARC max capacity
func (a *ARC) Cap() int {
	return a.capacity
}

// Add a new item into arc
func (a *ARC) Set(key, value interface{}) bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.t1.Contains(key) {
		a.t1.Remove(key)
		a.t2.Set(key, value)
		return false
	}

	if a.t2.Contains(key) {
		a.t2.Set(key, value)
		return false
	}

	var evicted bool
	if a.b1.Contains(key) {
		var (
			delta = 1
			l1    = a.b1.Len()
			l2    = a.b2.Len()
		)
		if l2 > l1 {
			delta = l2 / l1
		}
		if a.p+delta >= a.capacity {
			a.p = a.capacity
		} else {
			a.p += a.capacity
		}

		if a.t1.Len()+a.t2.Len() > a.capacity {
			a.adaptEvict(false)
			evicted = true
		}
		a.b1.Remove(key)
		a.t2.Set(key, value)
		return evicted
	}

	if a.b2.Contains(key) {
		var (
			delta = 1
			l1    = a.b1.Len()
			l2    = a.b2.Len()
		)
		if l1 > l2 {
			delta = l1 / l2
		}
		if delta > a.p {
			a.p = 0
		} else {
			a.p -= delta
		}

		if a.t1.Len()+a.t2.Len() > a.capacity {
			a.adaptEvict(true)
			evicted = true
		}
		a.b2.Remove(key)
		a.t2.Set(key, value)
		return evicted
	}

	if a.t1.Len()+a.t2.Len() > a.capacity {
		a.adaptEvict(false)
		evicted = true
	}
	if a.b1.Len() > a.capacity-a.p {
		a.b1.PopOldest()
		evicted = true
	}
	if a.b2.Len() > a.p {
		a.b2.PopOldest()
		evicted = true
	}

	a.t1.Set(key, value)
	return evicted
}

// Get return the given key's value
func (a *ARC) Get(key interface{}) (value interface{}, ok bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if value, ok := a.t1.Get(key); ok {
		a.t1.Remove(key)
		a.t1.Set(key, value)
		return value, ok
	}

	if value, ok := a.t2.Get(key); ok {
		return value, ok
	}
	return nil, false
}

// return the ARC length
func (a *ARC) Len() int {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.t1.Len() + a.t2.Len()
}

// return all keys in cache
func (a *ARC) Keys() []interface{} {
	a.lock.RLock()
	defer a.lock.RUnlock()

	k1 := a.t1.Keys()
	k2 := a.t2.Keys()
	return append(k1, k2...)
}

// Cotains check if the ARC contains the given key
func (a *ARC) Contains(key interface{}) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.t1.Contains(key) || a.t2.Contains(key)
}

// Remove the item from cache by key
func (a *ARC) Remove(key interface{}) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	switch {
	case a.t1.Remove(key):
		return true
	case a.t2.Remove(key):
		return true
	case a.b1.Remove(key):
		return false
	case a.b2.Remove(key):
		return false
	default:
		return false
	}
}

// Purge use to clear all items in ARC
func (a *ARC) Purge() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.t1.Purge()
	a.t2.Purge()
	a.b1.Purge()
	a.b2.Purge()
}

// Remove and return the oldest item from ARC
func (a *ARC) PopOldest() (key, value interface{}) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.t1.Len() > 0 {
		return a.t1.PopOldest()
	}
	if a.t2.Len() > 0 {
		return a.t2.PopOldest()
	}
	return nil, nil
}

// return the value if the key exist, otherwise update the key by given value similar with redis SETNX
func (a *ARC) GetOrSet(key, value interface{}) (newValue interface{}, isGet bool) {
	if v, ok := a.Get(key); ok {
		return v, ok
	}
	a.Set(key, value)
	return value, false
}

// adaptEvict used to evicted value
func (a *ARC) adaptEvict(inB2 bool) {
	l1 := a.t1.Len()
	if l1 > 0 && (l1 > a.p || (inB2 && l1 == a.p)) {
		if a.t1.Len() >= 1 {
			a.b1.Set(a.t1.PopOldest())
		} else {
			a.t1.PopOldest()
		}
	} else {
		if a.t2.Len() >= 1 {
			a.b2.Set(a.t2.PopOldest())
		} else {
			a.t2.PopOldest()
		}
	}
}
