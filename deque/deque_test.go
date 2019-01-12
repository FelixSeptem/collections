package deque

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDeque_Cap(t *testing.T) {
	q := NewDeque(32)
	if v := q.Cap(); v != 32 {
		t.Errorf("expect 32,got %d", v)
	}
}

func TestDeque_Contains(t *testing.T) {
	q := NewDeque(32)
	if ok := q.Contains("key"); ok {
		t.Errorf("expect false,got %v", ok)
	}
	q.PushLeft("key")
	if ok := q.Contains("key"); !ok {
		t.Errorf("expect true,got %v", ok)
	}
}

func TestDeque_Count(t *testing.T) {
	q := NewDeque(32)
	if v := q.Count("key"); v != 0 {
		t.Errorf("expect 0,got %d", v)
	}
	q.PushLeft("key")
	q.PushRight("key")
	if v := q.Count("key"); v != 2 {
		t.Errorf("expect 2,got %d", v)
	}
	q.PopLeft()
	if v := q.Count("key"); v != 1 {
		t.Errorf("expect 1,got %d", v)
	}
	q.PopRight()
	if v := q.Count("key"); v != 0 {
		t.Errorf("expect 0,got %d", v)
	}
}

func TestDeque_GetAll(t *testing.T) {
	q := NewDeque(32)
	data := []interface{}{
		1,
		3.14,
		"UserName",
		true,
	}
	for _, v := range data {
		q.PushRight(v)
	}
	if !cmp.Equal(q.GetAll(), data) {
		t.Errorf("expect %+v got %+v", data, q.GetAll())
	}
}

func TestDeque_GetLeft(t *testing.T) {
	q := NewDeque(32)
	q.PushRight(1)
	if v := q.GetLeft(); v != 1 {
		t.Errorf("expect 1,got %d", v)
	}
	q.PushLeft(2)
	if v := q.GetLeft(); v != 2 {
		t.Errorf("expect 2,got %d", v)
	}
}

func TestDeque_GetRight(t *testing.T) {
	q := NewDeque(32)
	q.PushLeft(1)
	if v := q.GetLeft(); v != 1 {
		t.Errorf("expect 1,got %d", v)
	}
	q.PushRight(2)
	if v := q.GetLeft(); v != 1 {
		t.Errorf("expect 2,got %d", v)
	}
}

func TestDeque_Len(t *testing.T) {
	q := NewDeque(32)
	for i := 0; i < 10; i++ {
		q.PushLeft(i)
	}
	if l := q.Len(); l != 10 {
		t.Errorf("expect 10,got %d", l)
	}
}

func TestDeque_PopLeft(t *testing.T) {
	q := NewDeque(32)
	if v := q.PopLeft(); v != nil {
		t.Errorf("expect nil,got %v", v)
	}
	q.PushLeft(1)
	if v := q.PopLeft(); v != 1 {
		t.Errorf("expect 1,got %v", v)
	}
	q.PushLeft(1)
	q.PushRight(2)
	q.PushLeft(3)
	q.PushRight(4)
	if v := q.PopLeft(); v != 3 {
		t.Errorf("expect 3,got %v", v)
	}
}

func TestDeque_PopRight(t *testing.T) {
	q := NewDeque(32)
	if v := q.PopRight(); v != nil {
		t.Errorf("expect nil,got %v", v)
	}
	q.PushLeft(1)
	if v := q.PopRight(); v != 1 {
		t.Errorf("expect 1,got %v", v)
	}
	q.PushLeft(1)
	q.PushRight(2)
	q.PushLeft(3)
	q.PushRight(4)
	if v := q.PopRight(); v != 4 {
		t.Errorf("expect 4,got %v", v)
	}
}

func TestDeque_Purge(t *testing.T) {
	q := NewDeque(32)
	for i := 0; i < 10; i++ {
		q.PushLeft(i)
		q.PushRight(i * 2)
	}
	q.Purge()
	if q.Len() > 0 {
		t.Errorf("expect 0,got %d", q.Len())
	}
}

func TestDeque_PushLeft(t *testing.T) {
	q := NewDeque(32)
	for i := 0; i < 24; i++ {
		q.PushLeft(i)
	}
	if q.GetLeft() != 23 || q.GetRight() != 0 {
		t.Errorf("expect 0,23;got %d,%d", q.GetLeft(), q.GetRight())
	}
}

func TestDeque_PushRight(t *testing.T) {
	q := NewDeque(32)
	for i := 0; i < 24; i++ {
		q.PushRight(i)
	}
	if q.GetLeft() != 0 || q.GetRight() != 23 {
		t.Errorf("expect 0,23;got %d,%d", q.GetLeft(), q.GetRight())
	}
}

func TestDeque_Remove(t *testing.T) {
	q := NewDeque(32)
	q.Remove(1)
	q.PushLeft(1)
	q.PushRight(3)
	q.Remove(1)
	if v := q.Count(1); v != 0 {
		t.Errorf("expect 0, got %d", v)
	}
}

func TestDeque_Reverse(t *testing.T) {
	q1 := NewDeque(32)
	q2 := NewDeque(32)
	for i := 0; i < 24; i++ {
		q1.PushLeft(i)
		q2.PushRight(i)
	}
	q1.Reverse()
	if !cmp.Equal(q1.GetAll(), q2.GetAll()) {
		t.Errorf("expect %+v,got %+v", q2.GetAll(), q1.GetAll())
	}
}

func TestDeque_Rotate(t *testing.T) {
	q := NewDeque(32)
	var (
		res1 = []interface{}{1, 2, 3, 4, 5, 6, 7}
		res2 = []interface{}{3, 4, 5, 6, 7, 1, 2}
		res3 = []interface{}{6, 7, 1, 2, 3, 4, 5}
	)
	for i := 1; i < 8; i++ {
		q.PushRight(i)
	}
	q.Rotate(0)
	if !cmp.Equal(res1, q.GetAll()) {
		t.Errorf("expect %+v;got %+v", res1, q.GetAll())
	}
	q.Rotate(2)
	if !cmp.Equal(res2, q.GetAll()) {
		t.Errorf("expect %+v;got %+v", res2, q.GetAll())
	}
	q.Rotate(-4)
	if !cmp.Equal(res3, q.GetAll()) {
		t.Errorf("expect %+v;got %+v", res3, q.GetAll())
	}
}

func BenchmarkDeque_Push(b *testing.B) {
	b.StopTimer()
	q := NewDeque(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			q.PushLeft(i)
		} else {
			q.PushRight(i)
		}
	}
}

func BenchmarkDeque_Pop(b *testing.B) {
	b.StopTimer()
	q := NewDeque(8096)
	for i := 0; i < 10000000; i++ {
		q.PushLeft(i)
		q.PushRight(i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		q.PopLeft()
		q.PopRight()
	}
}
