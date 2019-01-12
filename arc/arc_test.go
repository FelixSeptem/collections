package arc

import "testing"

func TestARC_Set(t *testing.T) {
	a := NewARCCache(32)
	a.Set("key", "value")
}

func TestARC_Get(t *testing.T) {
	a := NewARCCache(32)
	if _, ok := a.Get("key"); ok {
		t.Errorf("expect false,got %v", ok)
	}
	a.Set("key", "value")
	if _, ok := a.Get("key"); !ok {
		t.Errorf("expect true,got %v", ok)
	}
}

func TestARC_Cap(t *testing.T) {
	a := NewARCCache(32)
	if c := a.Cap(); c != 32 {
		t.Errorf("expect 32,got %d", c)
	}
}

func TestARC_Contains(t *testing.T) {
	a := NewARCCache(32)
	if ok := a.Contains("key"); ok {
		t.Errorf("expect false,got %v", ok)
	}
	a.Set("key", "value")
	if ok := a.Contains("key"); !ok {
		t.Errorf("expect true,got %v", ok)
	}
}

func TestARC_GetOrSet(t *testing.T) {
	a := NewARCCache(32)
	if v, ok := a.GetOrSet("key1", "value1"); ok || v != "value1" {
		t.Errorf("expect 'value1' with false, got %v with %v", v, ok)
	}
	if v, ok := a.GetOrSet("key1", "xxx"); !ok || v != "value1" {
		t.Errorf("expect 'value1' with false, got %v with %v", v, ok)
	}
}

func TestARC_Keys(t *testing.T) {
	a := NewARCCache(32)
	a.Set("k", "v")
	a.Set(1, 2)
	keys := a.Keys()
	if keys[0] != "k" || keys[1] != 1 {
		t.Errorf("expect 'k',1;got %v,%v", keys[0], keys[1])
	}
}

func TestARC_Len(t *testing.T) {
	a := NewARCCache(32)
	if l := a.Len(); l != 0 {
		t.Errorf("expect 0 got %d", l)
	}
	a.Set("k", "v")
	if l := a.Len(); l != 1 {
		t.Errorf("expect 1 got %d", l)
	}
}

func TestARC_PopOldest(t *testing.T) {
	a := NewARCCache(32)
	a.Set(1, 2)
	if k, v := a.PopOldest(); k != 1 || v != 2 {
		t.Errorf("expect 1,2;got %v,%v", k, v)
	}
	a.Set(2, 4)
	a.Set(3, 6)
	if k, v := a.PopOldest(); k != 2 || v != 4 {
		t.Errorf("expect 2,4;got %v,%v", k, v)
	}
}

func TestARC_Purge(t *testing.T) {
	a := NewARCCache(32)
	for i := 0; i < 10; i++ {
		a.Set(i, i)
	}
	if v := a.Len(); v != 10 {
		t.Errorf("expect 10,got %d", v)
	}
	a.Purge()
	if v := a.Len(); v != 0 {
		t.Errorf("expect 0,got %d", v)
	}
}

func TestARC_Remove(t *testing.T) {
	a := NewARCCache(32)
	if _, ok := a.Get("key1"); ok {
		t.Errorf("expect false got %v", ok)
	}
	a.Set("key1", "value1")
	if _, ok := a.Get("key1"); !ok {
		t.Errorf("expect true got %v", ok)
	}
	a.Remove("key1")
	if _, ok := a.Get("key1"); ok {
		t.Errorf("expect false got %v", ok)
	}
}

func BenchmarkARC_Set(b *testing.B) {
	b.StopTimer()
	a := NewARCCache(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.Set(i, i)
	}
}

func BenchmarkARC_GetExist(b *testing.B) {
	b.StopTimer()
	a := NewARCCache(8096)
	a.Set("key", "value")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.Get("key")
	}
}

func BenchmarkARC_GetNotExist(b *testing.B) {
	b.StopTimer()
	a := NewARCCache(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a.Get("key")
	}
}
