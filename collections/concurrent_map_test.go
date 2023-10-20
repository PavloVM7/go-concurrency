package collections

import (
	"testing"
)

func TestConcurrentMap_PutIfNotExistsDoubleCheck(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key, val := "string key", 357
	if ok, _ := cm.PutIfNotExistsDoubleCheck(key, val); !ok {
		t.Fatalf("PutIfNotExistsDoubleCheck(), the value (%v) was not added for the key (%v)", val, key)
	}
	newVal := 123
	ok1, old := cm.PutIfNotExistsDoubleCheck(key, newVal)
	if ok1 {
		t.Fatalf("PutIfNotExistsDoubleCheck(), the value %v for the key %v was unexpectedly added", newVal, key)
	}
	if old != val {
		t.Fatalf("PutIfNotExistsDoubleCheck(), expected: %v, actual: %v", val, old)
	}
}

func TestConcurrentMap_PutIfNotExists(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key, val := "string key", 357
	if ok, _ := cm.PutIfNotExists(key, val); !ok {
		t.Fatalf("PutIfNotExists(), the value (%v) was not added for the key (%v)", val, key)
	}
	newVal := 123
	ok1, old := cm.PutIfNotExists(key, newVal)
	if ok1 {
		t.Fatalf("PutIfNotExists(), the value %v for the key %v was unexpectedly added", newVal, key)
	}
	if old != val {
		t.Fatalf("PutIfNotExists(), expected: %v, actual: %v", val, old)
	}
}

func TestConcurrentMap_RemoveIfExistsDoubleCheck(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key, val := "key string", 123
	cm.Put(key, val)
	if _, ok := cm.Get(key); !ok {
		t.Fatalf("value %v was not added for key %v", val, key)
	}
	ok, actual := cm.RemoveIfExistsDoubleCheck(key)
	if !ok {
		t.Fatalf("value not exists for key %v", key)
	}
	if actual != val {
		t.Fatalf("wrong value, expected: %v, actual: %v", val, actual)
	}
	ok1, actual1 := cm.RemoveIfExistsDoubleCheck(key)
	if ok1 {
		t.Fatalf("the value (%v) for the key (%v) suddenly exists", actual1, key)
	}
}

func TestConcurrentMap_RemoveIfExists(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key, val := "key string", 123
	cm.Put(key, val)
	if _, ok := cm.Get(key); !ok {
		t.Fatalf("value %v was not added for key %v", val, key)
	}
	ok, actual := cm.RemoveIfExists(key)
	if !ok {
		t.Fatalf("value not exists for key %v", key)
	}
	if actual != val {
		t.Fatalf("wrong value, expected: %v, actual: %v", val, actual)
	}
	ok1, actual1 := cm.RemoveIfExists(key)
	if ok1 {
		t.Fatalf("the value (%v) for the key (%v) suddenly exists", actual1, key)
	}
}
func TestConcurrentMap_Remove(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key := "string key"
	cm.Put(key, 123)
	if _, ok := cm.Get(key); !ok {
		t.Fatalf("value was not added for key %v", key)
	}
	cm.Remove(key)
	if _, ok := cm.Get(key); ok {
		t.Fatalf("value was not removed for key %v", key)
	}
}

func TestConcurrentMap_Put(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key := "key string"
	cm.Put(key, 1)
	actual, ok := cm.Get(key)
	if !ok {
		t.Fatal("Put(), value not exists")
	}
	if actual != 1 {
		t.Fatalf("Put(), wrong value, expected: %v, actual: %v", 1, actual)
	}
	cm.Put(key, 3)
	actual1, ok1 := cm.Get(key)
	if !ok1 {
		t.Fatal("Put(), the value has not been replaced")
	}
	if actual1 != 3 {
		t.Fatalf("Put(), the old value %v has not be replaced by a new value %v", actual1, 3)
	}
}

func TestConcurrentMap_Get(t *testing.T) {
	tests := []struct {
		key string
		val int
	}{
		{"string1", 1},
		{"string2", 2},
		{"string3", 3},
	}
	cm := NewConcurrentMapCapacity[string, int](3)

	for _, tt := range tests {
		cm.Put(tt.key, tt.val)
	}
	if cm.Size() != len(tests) {
		t.Fatalf("wrong size, want: %d, got: %d", len(tests), cm.Size())
	}

	for _, tt := range tests {
		got, ok := cm.Get(tt.key)
		if !ok {
			t.Fatalf("the value %v for the key %v not exists", tt.val, tt.key)
		}
		if got != tt.val {
			t.Fatalf("wrong value, expected: %v, actual: %v", tt.val, got)
		}
	}
}

func TestConcurrentMap_Copy(t *testing.T) {
	tests := []struct {
		key string
		val int
	}{
		{"string1", 1},
		{"string2", 2},
		{"string3", 3},
	}
	cm := NewConcurrentMapCapacity[string, int](3)

	for _, tt := range tests {
		cm.Put(tt.key, tt.val)
	}
	cpy := cm.Copy()
	if len(cpy) != 3 {
		t.Fatalf("wrong len, expected: %v, actual: %v", 3, len(cpy))
	}
	for _, tt := range tests {
		actual := cpy[tt.key]
		if actual != tt.val {
			t.Fatalf("wrong value, expected: %v, actual: %v", tt.val, actual)
		}
	}
}

func TestConcurrentMap_Clear(t *testing.T) {
	cm := NewConcurrentMap[int, int]()
	if cm.capacity != 0 {
		t.Fatal("wrong capacity")
	}
	cm.Put(1, 1)
	cm.Put(2, 2)
	cm.Put(3, 3)
	if cm.Size() != 3 {
		t.Fatal("wrong map size")
	}
	cm.Clear()
	if cm.Size() != 0 {
		t.Fatal("the map is not cleared")
	}
}
func TestConcurrentMap_Clear_capacity(t *testing.T) {
	cm := NewConcurrentMapCapacity[int, string](123)
	if cm.capacity != 123 {
		t.Fatal("wrong capacity")
	}
	cm.Put(1, "str")
	cm.Put(2, "str")
	cm.Put(3, "str")
	if cm.Size() != 3 {
		t.Fatal("wrong map size")
	}
	cm.Clear()
	if cm.Size() != 0 {
		t.Fatal("the map is not cleared")
	}
}

func TestConcurrentMap_Size(t *testing.T) {
	cm := NewConcurrentMapCapacity[int, string](123)
	if cm.Size() != 0 {
		t.Fatalf("wrong size: expected %d, actual: %d", 0, cm.Size())
	}
}
