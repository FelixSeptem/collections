package lru

import "testing"

func TestLRU_Set(t *testing.T) {
	s := NewLRUCache(32)
	if evicted := s.Set("key", "value"); evicted {
		t.Errorf("expect got false, got %v", evicted)
	}
}

func TestLRU_Get(t *testing.T) {
	s := NewLRUCache(32)
	if v, ok := s.Get("key"); ok {
		t.Errorf("expect got false,got %v with %v", ok, v)
	}
	s.Set("key", "value")
	if v, ok := s.Get("key"); !ok || v != "value" {
		t.Errorf("expect got 'value' with true,got %v with %v", v, ok)
	}
}

func TestLRU_Cap(t *testing.T) {
	s := NewLRUCache(32)
	if v := s.Cap(); v != 32 {
		t.Errorf("expect got 32,got %d", v)
	}
}

func TestLRU_Contains(t *testing.T) {
	s := NewLRUCache(32)
	if ok := s.Contains("key"); ok {
		t.Errorf("expect got false, got %v", ok)
	}
	s.Set("key", "value")
	if ok := s.Contains("key"); !ok {
		t.Errorf("expect got true, got %v", ok)
	}
}

func TestLRU_Info(t *testing.T) {
	s := NewLRUCache(32)
	if hits, misses, maxSize, currentSize := s.Info(); hits != 0 || misses != 0 || maxSize != 32 || currentSize != 0 {
		t.Errorf("expect got 0,0,32,0;got %d %d %d %d", hits, misses, maxSize, currentSize)
	}
}

func TestLRU_GetOrSet(t *testing.T) {
	s := NewLRUCache(32)
	if v, ok := s.GetOrSet("key1", "value1"); ok || v != "value1" {
		t.Errorf("expect 'value1' with false, got %v with %v", v, ok)
	}
	if v, ok := s.GetOrSet("key1", "xxx"); !ok || v != "value1" {
		t.Errorf("expect 'value1' with false, got %v with %v", v, ok)
	}
}

func TestLRU_Keys(t *testing.T) {
	s := NewLRUCache(32)
	s.Set("k", "v")
	s.Set(1, 2)
	keys := s.Keys()
	if keys[0] != "k" || keys[1] != 1 {
		t.Errorf("expect 'k',1;got %v,%v", keys[0], keys[1])
	}
}

func TestLRU_Len(t *testing.T) {
	s := NewLRUCache(32)
	if l := s.Len(); l != 0 {
		t.Errorf("expect 0 got %d", l)
	}
	s.Set("k", "v")
	if l := s.Len(); l != 1 {
		t.Errorf("expect 1 got %d", l)
	}
}

func TestLRU_PopOldest(t *testing.T) {
	s := NewLRUCache(32)
	s.Set(1, 2)
	if k, v := s.PopOldest(); k != 1 || v != 2 {
		t.Errorf("expect 1,2;got %v,%v", k, v)
	}
	s.Set(2, 4)
	s.Set(3, 6)
	if k, v := s.PopOldest(); k != 2 || v != 4 {
		t.Errorf("expect 2,4;got %v,%v", k, v)
	}
}

func TestLRU_Purge(t *testing.T) {
	s := NewLRUCache(32)
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	if v := s.Len(); v != 10 {
		t.Errorf("expect 10,got %d", v)
	}
	s.Purge()
	if v := s.Len(); v != 0 {
		t.Errorf("expect 0,got %d", v)
	}
}

func TestLRU_Remove(t *testing.T) {
	s := NewLRUCache(32)
	if _, ok := s.Get("key1"); ok {
		t.Errorf("expect false got %v", ok)
	}
	s.Set("key1", "value1")
	if _, ok := s.Get("key1"); !ok {
		t.Errorf("expect true got %v", ok)
	}
	s.Remove("key1")
	if _, ok := s.Get("key1"); ok {
		t.Errorf("expect false got %v", ok)
	}
}

func BenchmarkLRU_Set(b *testing.B) {
	b.StopTimer()
	s := NewLRUCache(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Set(i, i)
	}
}

func BenchmarkLRU_GetExist(b *testing.B) {
	b.StopTimer()
	s := NewLRUCache(8096)
	s.Set("key", "value")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get("key")
	}
}

func BenchmarkLRU_GetNotExist(b *testing.B) {
	b.StopTimer()
	s := NewLRUCache(8096)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get("key")
	}
}
