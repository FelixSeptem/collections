// Package deque implement a fixed size thread safe queue as a generalization of both queue and stack
package deque

import (
	"container/list"
	"sync"
)

const (
	Default_Deque_Size = 1024
)

// Deque implement a queue as a generalization of both queue and stack
type Deque struct {
	capacity int
	lock     sync.RWMutex
	data     *list.List
	counter  map[interface{}]int
}

// a fixed size deque
func NewDeque(size int) *Deque {
	if size <= 0 {
		size = Default_Deque_Size
	}
	return &Deque{
		capacity: size,
		data:     list.New(),
		counter:  make(map[interface{}]int),
	}
}

// push a new item into deque from left
func (q *Deque) PushLeft(item interface{}) (evicted bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data.PushFront(item)
	q.counter[item] += 1
	if q.data.Len() > q.capacity {
		n := q.data.Back()
		q.counter[n.Value] -= 1
		q.data.Remove(n)
		return true
	}
	return false
}

// push a new item into deque from right
func (q *Deque) PushRight(item interface{}) (evicted bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.data.PushBack(item)
	q.counter[item] += 1
	if q.data.Len() > q.capacity {
		n := q.data.Front()
		q.counter[n.Value] -= 1
		q.data.Remove(n)
		return true
	}
	return false
}

// pop a item from left
func (q *Deque) PopLeft() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.data.Len() == 0 {
		return nil
	}
	n := q.data.Front()
	q.counter[n.Value] -= 1
	return q.data.Remove(q.data.Front())
}

// get a item from left
func (q *Deque) GetLeft() interface{} {
	q.lock.RLock()
	defer q.lock.RUnlock()
	if q.data.Len() == 0 {
		return nil
	}
	return q.data.Front().Value
}

// pop a item from right
func (q *Deque) PopRight() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.data.Len() == 0 {
		return nil
	}
	n := q.data.Back()
	q.counter[n.Value] -= 1
	return q.data.Remove(n)
}

// get a item from right
func (q *Deque) GetRight() interface{} {
	q.lock.RLock()
	defer q.lock.RLock()
	if q.data.Len() == 0 {
		return nil
	}
	return q.data.Back().Value
}

// remove the first occurrence of value
func (q *Deque) Remove(value interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if v, ok := q.counter[value]; !ok || v <= 0 {
		return
	}
	for i := q.data.Front(); i != nil; i = i.Next() {
		if i.Value == value {
			q.data.Remove(i)
			q.counter[value] -= 1
		}
	}
}

// reverse the deque
func (q *Deque) Reverse() {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.data.Len() == 0 {
		return
	}
	newList := list.New()
	for i := q.data.Back(); i != nil; i = i.Prev() {
		newList.PushBack(i.Value)
	}
	q.data = newList
}

// return the data length of deque
func (q *Deque) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.data.Len()
}

// return the max capacity of deque
func (q *Deque) Cap() int {
	return q.capacity
}

// clear the deque
func (q *Deque) Purge() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.counter = make(map[interface{}]int)
	q.data = list.New()
}

// rotate the deque n steps to the right
func (q *Deque) Rotate(step int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if step == 0 {
		return
	}
	if step > 0 {
		n := step % q.data.Len()
		for i := 0; i < n; i++ {
			q.data.PushBack(q.data.Remove(q.data.Front()))
		}
	}
	if step < 0 {
		n := (0 - step) % q.data.Len()
		for i := 0; i < n; i++ {
			q.data.PushFront(q.data.Remove(q.data.Back()))
		}
	}
}

// count the given value in the deque
func (q *Deque) Count(value interface{}) int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	v := q.counter[value]
	return v
}

// check if the deque contain the given value
func (q *Deque) Contains(value interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	v, ok := q.counter[value]
	return ok && v > 0
}

// get all values from deque
func (q *Deque) GetAll() []interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	keys := make([]interface{}, q.data.Len())
	head := q.data.Front()
	for i := 0; i < q.data.Len(); i++ {
		keys[i] = head.Value
		head = head.Next()
	}
	return keys
}
