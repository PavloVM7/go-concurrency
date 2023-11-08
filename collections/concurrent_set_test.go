package collections

import (
	"fmt"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestConcurrentSet_ForeEach(t *testing.T) {
	set := NewConcurrentSetWithValues[int](1, 2, 3)
	var sum int
	set.ForEach(func(value int) {
		sum += value
	})
	expectedSum := 6
	if sum != expectedSum {
		t.Fatalf("wrong a sum value: %d, wanted: %d", sum, expectedSum)
	}
}

func TestConcurrentSet_ToSlice(t *testing.T) {
	tests := []int{1, 2, 3}
	set := NewConcurrentSetCapacity[int](len(tests))
	for _, tt := range tests {
		set.Add(tt)
	}
	actual := set.ToSlice()
	slices.Sort(actual)
	if !reflect.DeepEqual(tests, actual) {
		t.Fatalf("wrong slice: '%v', expected: '%v'", actual, tests)
	}
}

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

func TestConcurrentSet_AddAll_set_is_not_changing(t *testing.T) {
	tests := []string{"string 1", "string 2", "string 3"}
	set := NewConcurrentSetWithValues[string](tests...)
	prevSize := set.Size()
	actual := set.AddAll(tests...)
	if actual {
		t.Fatalf("AddAll(), set has changed: %v, expected: %v", actual, false)
	}
	if set.Size() != prevSize {
		t.Fatalf("AddAll(), invalid set size: %v, expected: %v", set.Size(), prevSize)
	}
	for _, str := range tests {
		if !set.Contains(str) {
			t.Fatalf("AddAll(), vaues added to the set incorrectly, value: '%s' is missing", str)
		}
	}
}
func TestConcurrentSet_AddAll_set_contains_some_values(t *testing.T) {
	tests := []string{"string 1", "string 2", "string 3"}
	set := NewConcurrentSetWithValues[string](tests[0], tests[1])
	actual := set.AddAll(tests...)
	if !actual {
		t.Fatalf("AddAll(), invalid return value: %v, expected: %v", actual, true)
	}
	if set.Size() != len(tests) {
		t.Fatalf("AddAll(), invalid set size: %v, expected: %v", set.Size(), len(tests))
	}
	for _, str := range tests {
		if !set.Contains(str) {
			t.Fatalf("AddAll(), vaues added to the set incorrectly, value: '%s' is missing", str)
		}
	}
}
func TestConcurrentSet_AddAll(t *testing.T) {
	tests := []string{"string 1", "string 2", "string 3"}
	set := NewConcurrentSet[string]()
	actual := set.AddAll(tests...)
	if !actual {
		t.Fatalf("AddAll(), invalid return value: %v, expected: %v", actual, true)
	}
	if set.Size() != len(tests) {
		t.Fatalf("AddAll(), invalid set size: %v, expected: %v", set.Size(), len(tests))
	}
	for _, str := range tests {
		if !set.Contains(str) {
			t.Fatalf("AddAll(), vaues added to the set incorrectly, value: '%s' is missing", str)
		}
	}
}

func TestConcurrentSet_IsEmpty_false(t *testing.T) {
	set := NewConcurrentSetWithValues[int](1, 2, 3)
	if set.IsEmpty() {
		t.Fatal("expected not empty set")
	}
}
func TestConcurrentSet_IsEmpty(t *testing.T) {
	set := NewConcurrentSetCapacity[int](123)
	if !set.IsEmpty() {
		t.Fatal("expected empty set")
	}
}

func TestConcurrentSet_Clear(t *testing.T) {
	set := NewConcurrentSetWithValues[int](1, 2, 3)
	if set.Size() != 3 {
		t.Fatalf("wrong set size: %d, want: %d", set.Size(), 3)
	}
	set.Clear()
	if !set.IsEmpty() {
		t.Fatal("expected empty set")
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

func TestNewConcurrentSetWithValues(t *testing.T) {
	tests := []string{"string 1", "string 2", "string 3"}
	set := NewConcurrentSetWithValues[string](tests...)
	if set.capacity != len(tests) {
		t.Fatalf("TestNewConcurrentSetWithValues(), wrong set capacity: %d, want: %d", set.capacity, len(tests))
	}
	for _, str := range tests {
		if !set.Contains(str) {
			t.Fatalf("TestNewConcurrentSetWithValues(), set created incorrectly, value: '%s' is missing", str)
		}
	}
}

func TestConcurrentSet(t *testing.T) {
	const (
		count   = 100_000
		threads = 100
	)
	adds := make([]int, threads)
	set := NewConcurrentSetCapacity[int](count)
	chStart := make(chan struct{})
	chEnd := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			<-chStart
			for j := 1; j <= count; j++ {
				if !set.Contains(j) && set.Add(j) {
					adds[num]++
				}
			}
			<-chEnd
			wg.Done()
		}(i)
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if set.Contains(count) {
				close(chEnd)
				return
			}
		}
	}()
	close(chStart)
	wg.Wait()
	sum := 0
	for _, v := range adds {
		sum += v
	}
	fmt.Println("sum=", sum)
	t.Log("sum=", sum, adds)
	if sum != count {
		t.Fatalf("wrong sum: %d, want: %d", sum, count)
	}
}
