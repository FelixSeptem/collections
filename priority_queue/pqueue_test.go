package priority_queue

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPQueue_Cap(t *testing.T) {
	pq := NewPQueue(32)
	if v := pq.Cap(); v != 32 {
		t.Errorf("expect 32,got %d", v)
	}
}

func TestPQueue_IsEmpty(t *testing.T) {
	pq := NewPQueue(32)
	if ok := pq.IsEmpty(); !ok {
		t.Errorf("expect true, got %v", ok)
	}
	pq.PushItem(&Payload{
		Value:    1,
		Priority: 2,
	})
	if ok := pq.IsEmpty(); ok {
		t.Errorf("expect false, got %v", ok)
	}
}

func TestPQueue_IsFull(t *testing.T) {
	pq := NewPQueue(32)
	if ok := pq.IsFull(); ok {
		t.Errorf("expect false, got %v", ok)
	}
	for i := 0; i < pq.Cap(); i++ {
		pq.PushItem(&Payload{
			Value:    i,
			Priority: i,
		})
	}
	if ok := pq.IsFull(); !ok {
		t.Errorf("expect true, got %v", ok)
	}
}

func TestPQueue_Length(t *testing.T) {
	pq := NewPQueue(32)
	if l := pq.Length(); l != 0 {
		t.Errorf("expect 0,got %v", l)
	}
	pq.PushItem(&Payload{
		Value:    1,
		Priority: 2,
	})
	if l := pq.Length(); l != 1 {
		t.Errorf("expect 1,got %v", l)
	}
}

func TestPQueue_PushItem(t *testing.T) {
	pq := NewPQueue(32)
	pq.PushItem(&Payload{
		Value:    1,
		Priority: 1,
	})
}

func TestPQueue_PopItem(t *testing.T) {
	pq := NewPQueue(32)
	p := &Payload{
		Value:    1,
		Priority: 1,
	}
	pq.PushItem(p)
	v, ok := pq.PopItem()
	if !ok && !(cmp.Equal(v, p)) {
		t.Errorf("expect %+v, got %+v", *p, v)
	}
}

func BenchmarkPQueue_PushItem(b *testing.B) {
	b.StopTimer()
	pq := NewPQueue(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p := &Payload{
			Value:    i,
			Priority: i * 2,
		}
		pq.PushItem(p)
	}
}

func BenchmarkPQueue_PopItem(b *testing.B) {
	b.StopTimer()
	pq := NewPQueue(8096)
	for i := 0; i < 1000000; i++ {
		p := &Payload{
			Value:    1,
			Priority: i * 2,
		}
		pq.PushItem(p)
		pq.PushItem(p)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pq.PopItem()
	}
}
