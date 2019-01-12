// Package stack implement a simple fix size FIFO stack
package stack

import (
	"container/list"
	"sync"
)

const (
	Default_Stack_Size = 1024
)

// a fixed size FIFO stack
type Stack struct {
	lock     sync.RWMutex
	capacity int
	items    *list.List
}

// NewStack return a given size stack
func NewStack(size int) *Stack {
	if size <= 0 {
		size = Default_Stack_Size
	}
	return &Stack{
		capacity: size,
		items:    list.New(),
	}
}

// Push a new item into stack
func (s *Stack) Push(item interface{}) (isSuccess bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items.Len() == s.capacity {
		return false
	}
	s.items.PushFront(item)
	return true
}

// Pop a item from stack
func (s *Stack) Pop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items.Len() == 0 {
		return nil
	}
	return s.items.Remove(s.items.Front())
}

// return the stack length
func (s *Stack) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items.Len()
}

// return the stack capacity
func (s *Stack) Cap() int {
	return s.capacity
}

// Get the item at the top of stack but don't remove it from stack
func (s *Stack) GetTop() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items.Len() == 0 {
		return nil
	}
	return s.items.Front().Value
}

// Get the item at the bottom of stack
func (s *Stack) GetBottom() interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items.Len() == 0 {
		return nil
	}
	return s.items.Back().Value
}

// Check if the stack is empty
func (s *Stack) IsEmpty() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.items.Len() == 0
}

// Check if the stack is full
func (s *Stack) IsFull() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items.Len() == s.capacity
}
