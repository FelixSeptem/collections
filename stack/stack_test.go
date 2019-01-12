package stack

import "testing"

func TestQueue_Cap(t *testing.T) {
	q := NewStack(32)
	if v := q.Cap(); v != 32 {
		t.Errorf("expect 32,got %d", v)
	}
}

func TestQueue_GetHead(t *testing.T) {
	q := NewStack(32)
	if i := q.GetTop(); i != nil {
		t.Errorf("expect nil,got %v", i)
	}
	q.Push(1)
	q.Push(2)
	if i := q.GetTop(); i != 2 {
		t.Errorf("expect 2,got %d", i)
	}
}

func TestQueue_GetTail(t *testing.T) {
	q := NewStack(32)
	if i := q.GetBottom(); i != nil {
		t.Errorf("expect nil,got %v", i)
	}
	q.Push(1)
	q.Push(2)
	if i := q.GetBottom(); i != 1 {
		t.Errorf("expect 1,got %d", i)
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	q := NewStack(32)
	if v := q.IsEmpty(); !v {
		t.Errorf("expect true, got %v", v)
	}
	q.Push(1)
	if v := q.IsEmpty(); v {
		t.Errorf("expect false, got %v", v)
	}
}

func TestQueue_IsFull(t *testing.T) {
	q := NewStack(32)
	if v := q.IsFull(); v {
		t.Errorf("expect false, got %v", v)
	}
	for i := 0; i <= 128; i++ {
		q.Push(i)
	}
	if v := q.IsFull(); !v {
		t.Errorf("expect true, got %v", v)
	}
}

func TestQueue_Len(t *testing.T) {
	q := NewStack(32)
	if l1 := q.Len(); l1 != 0 {
		t.Errorf("expect 0,got %d", l1)
	}
	q.Push(1)
	if l1 := q.Len(); l1 != 1 {
		t.Errorf("expect 1,got %d", l1)
	}
}

func TestQueue_Push(t *testing.T) {
	q := NewStack(32)
	if e := q.Push(1); !e {
		t.Errorf("expect true,got %v", e)
	}
}

func TestQueue_Pop(t *testing.T) {
	q := NewStack(32)
	q.Push(1)
	if v := q.Pop(); v != 1 {
		t.Errorf("expect 1,got %d", v)
	}
	q.Push(2)
	q.Push(3)
	if v := q.Pop(); v != 3 {
		t.Errorf("expect 3,got %d", v)
	}
}

func BenchmarkQueue_Push(b *testing.B) {
	b.StopTimer()
	q := NewStack(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkQueue_Pop(b *testing.B) {
	b.StopTimer()
	q := NewStack(8096)
	for i := 0; i < 8096; i++ {
		q.Push(i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
