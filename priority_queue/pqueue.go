// Package priority_queue implement a thread safe priority powered by "container/heap"
package priority_queue

import (
	"container/heap"
	"sync"
)

const (
	Default_PQueue_Size = 1024
)

// PQueue implement a priority queue
type PQueue struct {
	lock     sync.RWMutex
	capacity int
	data     []*Payload
}

// payload holds the data in priority queue
type Payload struct {
	Value    interface{}
	Priority int
	index    int
}

// return a fix size priority queue
func NewPQueue(size int) *PQueue {
	if size <= 0 {
		size = Default_PQueue_Size
	}
	pq := &PQueue{
		capacity: size,
	}
	heap.Init(pq)
	return pq
}

// Push a item into priority queue
func (pq *PQueue) PushItem(v *Payload) (evicted bool) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	heap.Push(pq, v)
	if len(pq.data) > pq.capacity {
		heap.Pop(pq)
		return true
	}
	return false
}

// Pop a item from priority queue
func (pq *PQueue) PopItem() (interface{}, bool) {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if len(pq.data) == 0 {
		return nil, false
	}
	return heap.Pop(pq), true
}

// return queue size
func (pq *PQueue) Cap() int {
	return pq.capacity
}

// return queue size
func (pq *PQueue) Length() int {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.data)
}

// check if queue is empty
func (pq *PQueue) IsEmpty() bool {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.data) == 0
}

// check if queue if full(reach max capacity)
func (pq *PQueue) IsFull() bool {
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return len(pq.data) == pq.capacity
}

// below is used to implement internal interface ref:https://godoc.org/container/heap shall not use it directly
func (pq PQueue) Len() int {
	return len(pq.data)
}

func (pq PQueue) Less(i, j int) bool {
	return pq.data[i].Priority > pq.data[j].Priority
}

func (pq PQueue) Swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
	pq.data[i].index, pq.data[j].index = i, j
}

func (pq *PQueue) Push(v interface{}) {
	item := v.(*Payload)
	item.index = len(pq.data)
	pq.data = append(pq.data, item)
	heap.Fix(pq, item.index)
}

func (pq *PQueue) Pop() interface{} {
	old := *pq
	n := len(old.data)
	item := old.data[n-1]
	item.index = -1
	pq.data = old.data[0 : n-1]
	return item
}
