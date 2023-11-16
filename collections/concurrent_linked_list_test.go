package collections

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, list.first, list.last, "values 'first' and 'last' must be the same")
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
	assert.Equal(t, 2, list.Size(), "wrong list size")
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
	assert.Equal(t, list.first, list.last, "values 'first' and 'last' must be the same")
	first, _ := list.GetFirst()
	assert.Equal(t, 1, first)
}

func TestConcurrentLinkedList_RemoveLast_single(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	assert.Equal(t, 1, list.Size(), "wrong list size")
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
	assert.Equal(t, 0, list.Size(), "wrong list size")
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

	assert.Equal(t, 5, list.Size(), "wrong list size")

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
	assert.Equal(t, "", val, "wrong default value")
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
	assert.Equal(t, 5, list.Size(), "wrong list size")
	actual := list.ToArray()
	assert.Equal(t, list.Size(), len(actual), "wrong array size")
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, actual, "wrong array")
}
func TestConcurrentLinkedList_AddLast(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(1)
	list.AddLast(2)
	list.AddLast(3)

	assert.Equal(t, 3, list.Size(), "wrong list size")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 3, last, "wrong last value")

	first, ok := list.GetFirst()
	assert.True(t, ok, "first value does not exist")
	assert.Equal(t, 1, first, "wrong first value")
}
func TestConcurrentLinkedList_AddLast_first(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddLast(1)
	assert.Equal(t, 1, list.Size(), "wrong list size")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 1, last, "wrong last value")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "first value doesn't exist")
	assert.Equal(t, 1, actual, "wrong first value")

	assert.Equal(t, last, actual, "the last and first values aren't the same")
}
func TestConcurrentLinkedList_AddFirst(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	list.AddFirst(2)
	list.AddFirst(3)
	assert.Equal(t, 3, list.Size(), "wrong list size")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "first value does not exist")
	assert.Equal(t, 3, actual, "wrong first value")

	last, lok := list.GetLast()
	assert.True(t, lok, "last value doesn't exist")
	assert.Equal(t, 1, last, "wrong last value")
}
func TestConcurrentLinkedList_AddFirst_first(t *testing.T) {
	list := NewConcurrentLinkedList[int]()
	list.AddFirst(1)
	assert.Equal(t, 1, list.Size(), "wrong list size")
	actual, ok := list.GetFirst()
	assert.True(t, ok, "the value has not been added")
	assert.Equal(t, 1, actual, "wrong first value")
	last, lok := list.GetLast()
	assert.True(t, lok, "the last value does not exist")
	assert.Equal(t, 1, last, "wrong last value")
	assert.Equal(t, list.first, list.last, "the last and first values are not the same")
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
	assert.Equal(t, 3, list.Size(), "wrong list size")
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
