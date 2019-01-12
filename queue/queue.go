// Package queue implement a simple fix size FILO queue
package queue

import (
	"container/list"
	"sync"
)

const (
	Default_Queue_Size = 1024
)

// a fixed size FILO queue
type Queue struct {
	lock     sync.RWMutex
	capacity int
	items    *list.List
}

// NewQueue return a given size queue
func NewQueue(size int) *Queue {
	if size <= 0 {
		size = Default_Queue_Size
	}
	return &Queue{
		items:    list.New(),
		capacity: size,
	}
}

// Push a new item into queue
func (q *Queue) Push(item interface{}) (evicted bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items.PushBack(item)
	if q.items.Len() > q.capacity {
		q.items.Remove(q.items.Front())
		return true
	}
	return false
}

// Pop a item from queue
func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.items.Len() == 0 {
		return nil
	}
	return q.items.Remove(q.items.Front())
}

// return the queue length
func (q *Queue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.items.Len()
}

// return the queue capacity
func (q *Queue) Cap() int {
	return q.capacity
}

// Get the item at the head of queue but don't remove it from queue
func (q *Queue) GetHead() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.items.Len() == 0 {
		return nil
	}
	return q.items.Front().Value
}

// Get the item at the tail of queue
func (q *Queue) GetTail() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.items.Len() == 0 {
		return nil
	}
	return q.items.Back().Value
}

// Check if the queue is empty
func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.items.Len() == 0
}

// Check if the queue is full
func (q *Queue) IsFull() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.items.Len() == q.capacity
}
