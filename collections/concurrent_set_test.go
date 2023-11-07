package collections

import "testing"

func TestConcurrentSet_Contains(t *testing.T) {
	tests := []int{1, 2, 3}
	set := NewConcurrentSetCapacity[int](len(tests))
	for _, tt := range tests {
		set.Add(tt)
	}
	for _, tt := range tests {
		actual := set.Contains(tt)
		if !actual {
			t.Fatalf("set must contains the value %d", tt)
		}
	}

}

func TestConcurrentSet_Add(t *testing.T) {
	set := NewConcurrentSet[int]()
	for i := 1; i <= 3; i++ {
		if ok := set.Add(i); !ok {
			t.Fatalf("unexpected return value: %v, expected: %v", ok, true)
		}
	}
	if set.Size() != 3 {
		t.Fatalf("unexpected set size: %v, expected: %v", set.Size(), 3)
	}
}

func TestNewConcurrentSetCapacity(t *testing.T) {
	const capacity = 123
	set := NewConcurrentSetCapacity[string](capacity)
	if set.capacity != capacity {
		t.Fatalf("wrong capacity: %d, expected: %d", set.capacity, capacity)
	}
}

func TestNewConcurrentSet(t *testing.T) {
	set := NewConcurrentSet[string]()
	if set.capacity != 0 {
		t.Fatalf("wrong capacity: %d, expected: %d", set.capacity, 0)
	}
}
