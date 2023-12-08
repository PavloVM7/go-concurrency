// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentMap_TrimToSize(t *testing.T) {
	const amount = 1_000_0000
	const rest = 20
	value := func(i int) string {
		return fmt.Sprintf("value-%d", i)
	}
	cm := NewConcurrentMapCapacity[int, string](amount)
	for i := 1; i <= amount; i++ {
		cm.Put(i, value(i))
	}
	assert.Equal(t, amount, cm.Size())

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := rest + 1; i <= amount; i++ {
		cm.Remove(i)
	}
	assert.Equal(t, rest, cm.Size())

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	runtime.GC()

	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)

	cm.TrimToSize()

	var m4 runtime.MemStats
	runtime.ReadMemStats(&m4)

	runtime.GC()

	var m5 runtime.MemStats
	runtime.ReadMemStats(&m5)

	memToString := func(ms *runtime.MemStats) string {
		return fmt.Sprintf("%d Kb", ms.Alloc/1024)
	}

	t.Logf("Memory after fill: %s; after remove: %s (GC: %s); after trim: %s (GC: %s)",
		memToString(&m1), memToString(&m2), memToString(&m3), memToString(&m4), memToString(&m5))

	assert.Equal(t, rest, cm.Size())
	for i := 1; i <= rest; i++ {
		expected := value(i)
		actual, ok := cm.Get(i)
		assert.True(t, ok)
		assert.Equal(t, expected, actual)
	}
}

func TestConcurrentMap_ForEachRead(t *testing.T) {
	cm := NewConcurrentMap[int, int]()
	cm.Put(1, 3)
	cm.Put(3, 5)
	cm.Put(5, 7)
	sumK, sumV, sumSize, sumVget := 0, 0, 0, 0
	cm.ForEachRead(func(key int, value int) {
		sumK += key
		sumV += value
		sumSize += cm.Size()
		v, _ := cm.Get(value)
		sumVget += v
	})
	const expectedSumK = 9
	if sumK != expectedSumK {
		t.Fatal("ForEachRead() incorrect sum of keys", "expected:", expectedSumK, "actual:", sumK)
	}
	const expectedSumValues = 15
	if sumV != expectedSumValues {
		t.Fatal("ForEachRead() incorrect sum of values", "expected:", expectedSumValues, "actual:", sumV)
	}
	const expectedSumSizes = 9
	if sumSize != expectedSumSizes {
		t.Fatal("ForEachRead() incorrect sum of size", "expected:", expectedSumSizes, "actual:", sumSize)
	}
	const wantSumValsObtInsdFunc = 12
	if sumVget != wantSumValsObtInsdFunc {
		t.Fatal("ForEachRead() incorrect sum of values", "expected:", wantSumValsObtInsdFunc, "actual:", sumVget)
	}
}

func TestConcurrentMap_ForEach(t *testing.T) {
	type tstType struct {
		name  string
		value int
	}
	cm := NewConcurrentMap[int, *tstType]()
	cm.Put(2, &tstType{"tst 2", 2})
	cm.Put(3, &tstType{"tst 3", 3})
	cm.Put(5, &tstType{"tst 5", 5})
	sum := 0
	cm.ForEach(func(key int, value *tstType) {
		sum += key
		value.value *= 2
	})
	const expectedSum = 10
	if sum != expectedSum {
		t.Fatal("incorrect sum", "expected:", expectedSum, "actual:", sum)
	}
	expected2 := &tstType{"tst 2", 4}
	actual2, _ := cm.Get(2)
	if !reflect.DeepEqual(actual2, expected2) {
		t.Log("expected:", expected2, "actual:", actual2)
	}
	expected3 := &tstType{"tst 3", 6}
	actual3, _ := cm.Get(3)
	if !reflect.DeepEqual(actual3, expected3) {
		t.Log("expected:", expected3, "actual:", actual3)
	}
	expected5 := &tstType{"tst 5", 10}
	actual5, _ := cm.Get(5)
	if !reflect.DeepEqual(actual5, expected5) {
		t.Log("expected:", expected5, "actual:", actual5)
	}
}

func TestConcurrentMap_PutIfNotExistsDoubleCheck(t *testing.T) {
	cm := NewConcurrentMap[string, int]()
	key, val := "string strong key", 357
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
		t.Fatalf("incorrect value, expected: %v, actual: %v", val, actual)
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
		t.Fatalf("incorrect value, expected: %v, actual: %v", val, actual)
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
		t.Fatalf("Put(), incorrect value, expected: %v, actual: %v", 1, actual)
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
		t.Fatalf("incorrect size, want: %d, got: %d", len(tests), cm.Size())
	}

	for _, tt := range tests {
		got, ok := cm.Get(tt.key)
		if !ok {
			t.Fatalf("the value %v for the key %v not exists", tt.val, tt.key)
		}
		if got != tt.val {
			t.Fatalf("incorrect value, expected: %v, actual: %v", tt.val, got)
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
		t.Fatalf("incorrect len, expected: %v, actual: %v", 3, len(cpy))
	}
	for _, tt := range tests {
		actual := cpy[tt.key]
		if actual != tt.val {
			t.Fatalf("incorrect value, expected: %v, actual: %v", tt.val, actual)
		}
	}
}

func TestConcurrentMap_Keys(t *testing.T) {
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
	if cm.IsEmpty() {
		t.Fatal("map is empty")
	}
	keys := cm.Keys()
	if len(keys) != cm.Size() {
		t.Fatalf("incorrect key slice length: %d, expected: %d", len(keys), cm.Size())
	}
	contains := func(key string) bool {
		for _, k := range keys {
			if k == key {
				return true
			}
		}
		return false
	}
	for _, tt := range tests {
		if !contains(tt.key) {
			t.Fatalf("slice not contains key '%s'", tt.key)
		}
	}
}

func TestConcurrentMap_Clear(t *testing.T) {
	cm := NewConcurrentMap[int, int]()
	if cm.capacity != 0 {
		t.Fatal("incorrect capacity")
	}
	cm.Put(1, 1)
	cm.Put(2, 2)
	cm.Put(3, 3)
	if cm.Size() != 3 {
		t.Fatal("incorrect map size")
	}
	cm.Clear()
	if cm.Size() != 0 {
		t.Fatal("the map is not cleared")
	}
}
func TestConcurrentMap_Clear_capacity(t *testing.T) {
	cm := NewConcurrentMapCapacity[int, string](123)
	if cm.capacity != 123 {
		t.Fatal("incorrect capacity")
	}
	cm.Put(1, "str")
	cm.Put(2, "str")
	cm.Put(3, "str")
	if cm.Size() != 3 {
		t.Fatal("incorrect map size")
	}
	cm.Clear()
	if cm.Size() != 0 {
		t.Fatal("the map is not cleared")
	}
}

func TestConcurrentMap_Size(t *testing.T) {
	const capacity = 123
	cm := NewConcurrentMapCapacity[int, string](capacity)
	if cm.Size() != 0 {
		t.Fatalf("incorrect size: expected %d, actual: %d", 0, cm.Size())
	}
	if cm.capacity != capacity {
		t.Fatalf("invalid capacity: %d, want: %d", cm.capacity, cm.capacity)
	}
}

func TestConcurrentMap_IsEmpty(t *testing.T) {
	const capacity = 123
	cm := NewConcurrentMapCapacity[int, string](capacity)
	if cm.Size() != 0 {
		t.Fatalf("incorrect size: expected %d, actual: %d", 0, cm.Size())
	}
	if !cm.IsEmpty() {
		t.Fatal("expected empty map")
	}
}

func TestNewConcurrentMap(t *testing.T) {
	const (
		threads = 100
		count   = 100_000
	)

	cm := NewConcurrentMap[int, int]()
	counters := make([]int32, threads)
	var state int32
	var wg sync.WaitGroup
	fnc := func(num int) {
		//revive:disable:empty-block
		for atomic.LoadInt32(&state) == 0 {
			// waiting for a start
		} //revive:enable:empty-block
		for i := 0; i < count; i++ {
			if ok, _ := cm.PutIfNotExistsDoubleCheck(i, num); ok {
				atomic.AddInt32(&counters[num], 1)
			}
		}
		wg.Done()
	}
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go fnc(i)
	}
	atomic.StoreInt32(&state, 1)
	wg.Wait()
	size := cm.Size()
	if size != count {
		t.Errorf("incorrect map size: %d, expected: %d", size, count)
	}
	amounts := make([]int, threads)
	cm.ForEachRead(func(key int, value int) {
		amounts[value]++
	})
	var sum int32
	amount := 0
	for i, c := range counters {
		if c > 0 {
			sum += c
			amount++
		}
		t.Log(i, "=", c, "=", amounts[i])
	}
	if sum != int32(count) {
		t.Fatalf("incorrect count: %d, expected: %d", sum, count)
	}
	t.Log("size:", size, "sum:", sum, "amount:", amount)
}
