// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

func TestLinkedList_example(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	using := func(funcs string) {
		fmt.Printf("=== using %s\n", funcs)
	}
	var wg sync.WaitGroup
	using("AddLast() and AddFirst()")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 11; i <= 20; i++ {
			list.AddLast(i) // adds items to the end of the list
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10; i > 0; i-- {
			list.AddFirst(i) // adds items to the head of the list
		}
	}()
	wg.Wait()
	showList := func() {
		fmt.Printf(">>> list size: %d, items: %v\n", list.Size(), list.ToArray())
	}
	showList()

	using("Get() and Remove()")
	item10, err := list.Get(10)
	fmt.Printf("before remove 10th item = %d, err = %v\n", item10, err)
	item10, err = list.Remove(10) // removes 10th item
	fmt.Printf("removed item10 = %d, err = %v\n", item10, err)
	item10, err = list.Get(10)
	fmt.Printf("after remove 10th item = %d, err = %v\n", item10, err)
	showList()

	using("GetFirst() and RemoveFirst()")
	first, firstOk := list.GetFirst()
	fmt.Printf("before remove first element: %d, exists: %t\n", first, firstOk)
	first, firstOk = list.RemoveFirst()
	fmt.Printf("first element: %d, removed: %t\n", first, firstOk)
	first, firstOk = list.GetFirst()
	fmt.Printf("current first element: %d, exists: %t\n", first, firstOk)
	showList()

	using("GetLast() and RemoveLast()")
	last, lastOk := list.GetLast()
	fmt.Printf("before remove last element: %d, exists: %t\n", last, lastOk)
	last, lastOk = list.RemoveLast()
	fmt.Printf("last element: %d, removed: %t\n", last, lastOk)
	last, lastOk = list.GetLast()
	fmt.Printf("current last element: %d, exists: %t\n", last, lastOk)
	showList()

	using("RemoveFirstOccurrence()")
	rFirst, fIndex := list.RemoveFirstOccurrence(func(value int) bool {
		return value%2 != 0
	})
	fmt.Printf("removed first odd value: %d, index: %d\n", rFirst, fIndex)
	showList()

	using("RemoveLastOccurrence()")
	rLast, lIndex := list.RemoveLastOccurrence(func(value int) bool {
		return value%2 == 0
	})
	fmt.Printf("removed last even value: %d, index: %d\n", rLast, lIndex)
	showList()

	using("RemoveAll()")
	count := list.RemoveAll(func(value int) bool {
		return value%3 == 0
	})
	fmt.Printf("%d elements that are dividable by 3 have been removed\n", count)
	showList()

	using("Clear()")
	list.Clear()
	showList()
}
func TestLinkedList_RemoveAll_duplicates(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	fill := func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				list.AddFirst(i)
			} else {
				list.AddLast(i)
			}
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fill()
		}()
	}
	wg.Wait()
	assert.Equal(t, 100, list.Size())
	present := make(map[int]struct{})
	list.RemoveAll(func(value int) bool {
		if _, ok := present[value]; !ok {
			present[value] = struct{}{}
			return false
		}
		return true
	})
	t.Log("list=", list.ToArray())
	assert.Equal(t, 10, list.Size())
}
func TestLinkedList_RemoveAll(t *testing.T) {
	type testCase[T any] struct {
		name       string
		list       *ConcurrentLinkedList[T]
		needRemove func(value T) bool
		want       T
		wantArray  []T
	}
	tests := []testCase[int]{
		{
			name:       "empty",
			list:       NewConcurrentLinkedList[int](),
			needRemove: func(value int) bool { return value > 0 },
			want:       0,
			wantArray:  []int{},
		},
		{
			name:       "not found",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3),
			needRemove: func(value int) bool { return value == 5 },
			want:       0,
			wantArray:  []int{1, 2, 3},
		},
		{
			name:       "single value",
			list:       NewConcurrentLinkedListItems[int](1),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantArray:  []int{},
		},
		{
			name:       "first value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantArray:  []int{2},
		},
		{
			name:       "last value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       1,
			wantArray:  []int{1},
		},
		{
			name:       "middle value",
			list:       NewConcurrentLinkedListItems[int](2, 1, 2, 3, 2, 5, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       4,
			wantArray:  []int{1, 3, 5},
		},
		{
			name:       "middle double values",
			list:       NewConcurrentLinkedListItems[int](2, 2, 1, 2, 2, 3, 2, 2, 5, 2, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       8,
			wantArray:  []int{1, 3, 5},
		},
		{
			name:       "even values",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
			needRemove: func(value int) bool { return value%2 == 0 },
			want:       5,
			wantArray:  []int{1, 3, 5, 7, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.RemoveAll(tt.needRemove)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveAll() got = %v, want %v", got, tt.want)
			}
			gotArray := tt.list.ToArray()
			if !reflect.DeepEqual(gotArray, tt.wantArray) {
				t.Errorf("RemoveAll() gotArray = %v, wantArray %v", gotArray, tt.wantArray)
			}
			if tt.list.Size() != len(gotArray) {
				t.Errorf("RemoveAll() gotSize = %v, wantSize %v", len(gotArray), tt.list.Size())
			}
		})
	}
}
func TestLinkedList_RemoveLastOccurrence(t *testing.T) {
	type testCase[T any] struct {
		name       string
		list       *ConcurrentLinkedList[T]
		needRemove func(value T) bool
		want       T
		wantIndex  int
		wantArray  []T
	}
	tests := []testCase[int]{
		{
			name:       "empty",
			list:       NewConcurrentLinkedList[int](),
			needRemove: func(value int) bool { return value > 0 },
			want:       0,
			wantIndex:  -1,
			wantArray:  []int{},
		},
		{
			name:       "not found",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3),
			needRemove: func(value int) bool { return value == 5 },
			want:       0,
			wantIndex:  -1,
			wantArray:  []int{1, 2, 3},
		},
		{
			name:       "single value",
			list:       NewConcurrentLinkedListItems[int](1),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantIndex:  0,
			wantArray:  []int{},
		},
		{
			name:       "first value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantIndex:  0,
			wantArray:  []int{2},
		},
		{
			name:       "last value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       2,
			wantIndex:  1,
			wantArray:  []int{1},
		},
		{
			name:       "middle value",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3, 2, 5),
			needRemove: func(value int) bool { return value == 2 },
			want:       2,
			wantIndex:  3,
			wantArray:  []int{1, 2, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotIndex := tt.list.RemoveLastOccurrence(tt.needRemove)
			if !reflect.DeepEqual(gotValue, tt.want) {
				t.Errorf("RemoveLastOccurrence() got = %v, want %v", gotValue, tt.want)
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("RemoveLastOccurrence() gotIndex = %v, wantIndex %v", gotIndex, tt.wantIndex)
			}
			gotArray := tt.list.ToArray()
			if !reflect.DeepEqual(gotArray, tt.wantArray) {
				t.Errorf("RemoveLastOccurrence() gotArray = %v, wantArray %v", gotArray, tt.wantArray)
			}
			if tt.list.Size() != len(gotArray) {
				t.Errorf("RemoveLastOccurrence() gotSize = %v, wantSize %v", len(gotArray), tt.list.Size())
			}
		})
	}
}
func TestLinkedList_RemoveFirstOccurrence(t *testing.T) {
	type testCase[T any] struct {
		name       string
		list       *ConcurrentLinkedList[T]
		needRemove func(value T) bool
		want       T
		wantIndex  int
		wantArray  []T
	}
	tests := []testCase[int]{
		{
			name:       "empty",
			list:       NewConcurrentLinkedList[int](),
			needRemove: func(value int) bool { return value > 0 },
			want:       0,
			wantIndex:  -1,
			wantArray:  []int{},
		},
		{
			name:       "not found",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3),
			needRemove: func(value int) bool { return value == 5 },
			want:       0,
			wantIndex:  -1,
			wantArray:  []int{1, 2, 3},
		},
		{
			name:       "single value",
			list:       NewConcurrentLinkedListItems[int](1),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantIndex:  0,
			wantArray:  []int{},
		},
		{
			name:       "first value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 1 },
			want:       1,
			wantIndex:  0,
			wantArray:  []int{2},
		},
		{
			name:       "last value",
			list:       NewConcurrentLinkedListItems[int](1, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       2,
			wantIndex:  1,
			wantArray:  []int{1},
		},
		{
			name:       "middle value",
			list:       NewConcurrentLinkedListItems[int](1, 2, 3, 2),
			needRemove: func(value int) bool { return value == 2 },
			want:       2,
			wantIndex:  1,
			wantArray:  []int{1, 3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotIndex := tt.list.RemoveFirstOccurrence(tt.needRemove)
			if !reflect.DeepEqual(gotValue, tt.want) {
				t.Errorf("RemoveFirstOccurrence() got = %v, want %v", gotValue, tt.want)
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("RemoveFirstOccurrence() gotIndex = %v, wantIndex %v", gotIndex, tt.wantIndex)
			}
			gotArray := tt.list.ToArray()
			if !reflect.DeepEqual(gotArray, tt.wantArray) {
				t.Errorf("RemoveFirstOccurrence() gotArray = %v, wantArray %v", gotArray, tt.wantArray)
			}
			if tt.list.Size() != len(gotArray) {
				t.Errorf("RemoveFirstOccurrence() gotSize = %v, wantSize %v", len(gotArray), tt.list.Size())
			}
		})
	}
}

func TestConcurrentLinkedList_Remove(t *testing.T) {
	const (
		want1 = 1
		want2 = 2
		want3 = 3
		want4 = 4
	)
	list := NewConcurrentLinkedList[int]()
	list.AddLast(want1)
	list.AddLast(want2)
	list.AddLast(want3)
	list.AddLast(want4)
	got3, _ := list.Remove(2)
	assert.Equal(t, want3, got3)
	gotAr1 := list.ToArray()
	wantAr1 := []int{1, 2, 4}
	assert.Equal(t, wantAr1, gotAr1)

	got1, _ := list.Remove(0)
	assert.Equal(t, want1, got1)
	gotAr2 := list.ToArray()
	wantAr2 := []int{2, 4}
	assert.Equal(t, gotAr2, wantAr2)

	got4, _ := list.Remove(1)
	assert.Equal(t, want4, got4)
	gotAr3 := list.ToArray()
	wantAr3 := []int{2}
	assert.Equal(t, wantAr3, gotAr3)

	got2, _ := list.Remove(0)
	assert.Equal(t, want2, got2)
	gotAr4 := list.ToArray()
	assert.Equal(t, 0, len(gotAr4))
}

func TestConcurrentLinkedList_Remove_last(t *testing.T) {
	const expected1 = "value 1"
	const expected2 = "value 2"
	list := NewConcurrentLinkedList[string]()
	list.AddLast(expected1)
	list.AddLast(expected2)
	actual, err := list.Remove(list.Size() - 1)
	assert.Nil(t, err)
	assert.Equal(t, expected2, actual)
	assert.Equal(t, 1, list.Size())
	first, _ := list.GetFirst()
	assert.Equal(t, expected1, first)
	last, _ := list.GetLast()
	assert.Equal(t, expected1, last)
	assert.Same(t, list.first, list.last)
}

func TestConcurrentLinkedList_Remove_single(t *testing.T) {
	const expected = "single value"
	list := NewConcurrentLinkedList[string]()
	list.AddLast(expected)
	actual, err := list.Remove(0)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	assert.Equal(t, 0, list.Size())
	assert.Nil(t, list.first, "the first item should be nil")
	assert.Nil(t, list.last, "the last item should be nil")
}

func TestConcurrentLinkedList_Remove_fail(t *testing.T) {
	list := NewConcurrentLinkedList[string]()
	actual, err := list.Remove(0)
	assert.ErrorIs(t, err, ErrIndexOutOfRange, "expected an 'index is out of range' error")
	assert.Equal(t, "", actual)
	list.AddLast("value")
	actual, err = list.Remove(1)
	assert.ErrorIs(t, err, ErrIndexOutOfRange, "expected an 'index is out of range' error")
	assert.Equal(t, "", actual)
}

func TestConcurrentLinkedList_RemoveFirst(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(1)
	list.AddLast(2)
	list.AddLast(3)
	expectedSize := list.Size()
	for i := 0; i < 3; i++ {
		actual, ok := list.RemoveFirst()
		assert.True(t, ok, "the first element must exist")
		expectedValue := i + 1
		assert.Equal(t, expectedValue, actual)
		expectedSize--
		assert.Equal(t, expectedSize, list.Size())
	}
	actual, ok := list.RemoveFirst()
	assert.False(t, ok, "the list must be empty")
	assert.Equal(t, 0, actual)
}

func TestConcurrentLinkedList_RemoveFirst_before_last(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	list.AddLast(2)
	assert.Equal(t, 2, list.Size())
	actual, ok := list.RemoveFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, actual)
	assert.Equal(t, 1, list.Size())
	assert.NotNil(t, list.first, "the first element must exist")
	assert.NotNil(t, list.last, "the last element must exist")
	assert.Nil(t, list.first.prev, "the 'prev' value of the first element must be nil")
	assert.Nil(t, list.first.next, "the 'next' value of the first element must be nil")
	assert.Nil(t, list.last.prev, "the 'prev' value of the last element must be nil")
	assert.Nil(t, list.last.next, "the 'next' value of the last element must be nil")
	assert.Same(t, list.first, list.last, "values 'first' and 'last' must be the same")
	last, _ := list.GetLast()
	assert.Equal(t, 2, last)
}

func TestConcurrentLinkedList_RemoveFirst_single(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	assert.Equal(t, 1, list.Size())
	actual, ok := list.RemoveFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, actual)
	assert.Equal(t, 0, list.Size())
	assert.Nil(t, list.first, "the first item should be nil")
	assert.Nil(t, list.last, "the last item should be nil")
}

func TestConcurrentLinkedList_RemoveFirst_empty(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	actual, ok := list.RemoveFirst()
	assert.False(t, ok)
	assert.Equal(t, 0, actual, "0 is expected")
}

func TestConcurrentLinkedList_RemoveLast(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	list.AddFirst(2)
	list.AddFirst(3)
	expectedSize := list.Size()
	for i := 0; i < 3; i++ {
		actual, ok := list.RemoveLast()
		assert.True(t, ok, "the last element must exist")
		expectedValue := i + 1
		assert.Equal(t, expectedValue, actual)
		expectedSize--
		assert.Equal(t, expectedSize, list.Size())
	}
	actual, ok := list.RemoveLast()
	assert.False(t, ok, "the list should be empty")
	assert.Equal(t, 0, actual)
}

func TestConcurrentLinkedList_RemoveLast_before_last(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(2)
	list.AddFirst(1)
	assert.Equal(t, 2, list.Size(), "incorrect list size")
	actual, ok := list.RemoveLast()
	assert.True(t, ok)
	assert.Equal(t, 2, actual)
	assert.Equal(t, 1, list.Size())
	assert.NotNil(t, list.first, "the first element must exist")
	assert.NotNil(t, list.last, "the last element must exist")
	assert.Nil(t, list.first.prev, "'prev' value of the first element must be nil")
	assert.Nil(t, list.first.next, "'next' value of the first element must be nil")
	assert.Nil(t, list.last.prev, "'prev' value of the last element must be nil")
	assert.Nil(t, list.last.next, "'next' value of the last element must be nil")
	assert.Same(t, list.first, list.last, "values 'first' and 'last' must be the same")
	first, _ := list.GetFirst()
	assert.Equal(t, 1, first)
}

func TestConcurrentLinkedList_RemoveLast_single(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	assert.Equal(t, 1, list.Size(), "incorrect list size")
	actual, ok := list.RemoveLast()
	assert.True(t, ok)
	assert.Equal(t, 1, actual)
	assert.Equal(t, 0, list.Size())
	assert.Nil(t, list.first, "the first value should be nil")
	assert.Nil(t, list.last, "the last value should be nil")
}
func TestConcurrentLinkedList_RemoveLast_empty(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	actual, ok := list.RemoveLast()
	assert.Equal(t, 0, list.Size(), "incorrect list size")
	assert.False(t, ok)
	assert.Equal(t, 0, actual, "unexpected value")
}
func TestConcurrentLinkedList_Get(t *testing.T) {
	crt := func(num int) string {
		return fmt.Sprint("list item ", num)
	}
	list := NewConcurrentLinkedList[string]()
	list.AddFirst(crt(3))
	list.AddFirst(crt(2))
	list.AddFirst(crt(1))
	list.AddLast(crt(4))
	list.AddLast(crt(5))

	assert.Equal(t, 5, list.Size(), "incorrect list size")

	for i := 0; i < list.Size(); i++ {
		actual, err := list.Get(i)
		assert.Nil(t, err, "unexpected error:", err)
		expected := crt(i + 1)
		assert.Equal(t, expected, actual, "index:", i)
	}
}
func TestConcurrentLinkedList_Get_fail(t *testing.T) {
	list := NewConcurrentLinkedList[string]()
	val, err := list.Get(-1)
	assert.ErrorIs(t, err, ErrIndexOutOfRange, "unexpected error")
	assert.Equal(t, "", val, "incorrect default value")
}
func TestConcurrentLinkedList_ToArray_empty(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	actual := list.ToArray()
	assert.Equal(t, 0, len(actual), "an empty array is expected")
}
func TestConcurrentLinkedList_ToArray(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(3)
	list.AddFirst(2)
	list.AddFirst(1)
	list.AddLast(4)
	list.AddLast(5)
	assert.Equal(t, 5, list.Size(), "incorrect list size")
	actual := list.ToArray()
	assert.Equal(t, list.Size(), len(actual), "incorrect array size")
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, actual, "incorrect array")
}
func TestConcurrentLinkedList_AddLast(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(1)
	list.AddLast(2)
	list.AddLast(3)

	assert.Equal(t, 3, list.Size(), "incorrect list size")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 3, last, "incorrect last value")

	first, ok := list.GetFirst()
	assert.True(t, ok, "first value does not exist")
	assert.Equal(t, 1, first, "incorrect first value")
}
func TestConcurrentLinkedList_AddLast_first(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(1)
	assert.Equal(t, 1, list.Size(), "incorrect list size")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 1, last, "incorrect last value")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "first value doesn't exist")
	assert.Equal(t, 1, actual, "incorrect first value")

	assert.Equal(t, last, actual, "the last and first values aren't the same")
}
func TestConcurrentLinkedList_AddFirst(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	list.AddFirst(2)
	list.AddFirst(3)
	assert.Equal(t, 3, list.Size(), "incorrect list size")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "first value does not exist")
	assert.Equal(t, 3, actual, "incorrect first value")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 1, last, "incorrect last value")
}
func TestConcurrentLinkedList_AddFirst_first(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	assert.Equal(t, 1, list.Size(), "incorrect list size")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "the value has not been added")
	assert.Equal(t, 1, actual, "incorrect first value")
	last, lok := list.GetLast()
	assert.True(t, lok, "the last value does not exist")
	assert.Equal(t, 1, last, "incorrect last value")
	assert.Same(t, list.first, list.last, "the last and first values are not the same")
}
func TestConcurrentLinkedList_GetFirst(t *testing.T) {
	tests := []listTestStruct{{name: "struct1", value: 1}, {name: "struct2", value: 2}, {name: "struct3", value: 3}}
	list := NewConcurrentLinkedListItems[listTestStruct](tests...)
	actual, ok := list.GetFirst()
	assert.True(t, ok, "the item does not exist")
	assert.Equal(t, tests[0], actual, "unexpected item")
}
func TestConcurrentLinkedList_GetLast(t *testing.T) {
	tests := []listTestStruct{{name: "struct1", value: 1}, {name: "struct2", value: 2}, {name: "struct3", value: 3}}
	list := NewConcurrentLinkedListItems[listTestStruct](tests...)
	actual, ok := list.GetLast()
	assert.True(t, ok, "the item does not exist")
	assert.Equal(t, tests[2], actual, "unexpected item")
}

func TestConcurrentLinkedList_GetLast_empty_list(t *testing.T) {
	list := NewConcurrentLinkedList[*listTestStruct]()
	actual, ok := list.GetLast()
	assert.Falsef(t, ok, "the item exists")
	assert.Nil(t, actual, "nil value is expected")
}

func TestConcurrentLinkedList_GetLast_empty_list_not_nil(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	actual, ok := list.GetLast()
	assert.False(t, ok, "the item exists")
	assert.Equal(t, 0, actual, "0 value is expected")
}

func TestConcurrentLinkedList_GetFirst_empty_list(t *testing.T) {
	list := NewConcurrentLinkedList[*listTestStruct]()
	actual, ok := list.GetFirst()
	assert.Falsef(t, ok, "the item exists")
	assert.Nil(t, actual, "nil value is expected")
}
func TestConcurrentLinkedList_GetFirst_empty_list_not_nil(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	actual, ok := list.GetFirst()
	assert.False(t, ok, "the item exists")
	assert.Equal(t, 0, actual, "0 value is expected")
}

func TestNewConcurrentLinkedListItems(t *testing.T) {
	list := NewConcurrentLinkedListItems[string]("string 1", "string 2", "string 3")
	assert.Equal(t, 3, list.Size(), "incorrect list size")
	actual1, _ := list.Get(0)
	assert.Equal(t, "string 1", actual1)
	actual2, _ := list.Get(1)
	assert.Equal(t, "string 2", actual2)
	actual3, _ := list.Get(2)
	assert.Equal(t, "string 3", actual3)
}

func TestNewConcurrentLinkedList(t *testing.T) {
	list := NewConcurrentLinkedList[string]()
	assert.Nil(t, list.first, "the first doesn't equal nil")
	assert.Nil(t, list.last, "the last doesn't equal nil")
	assert.Equal(t, 0, list.size, "the list size doesn't equal 0")
}

type listTestStruct struct {
	name  string
	value int
}
