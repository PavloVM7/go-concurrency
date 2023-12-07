// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"errors"
	"sync"
)

var (
	// ErrIndexOutOfRange error: 'index is out of range'
	ErrIndexOutOfRange = errors.New("index is out of range")
)

// ConcurrentLinkedList is a thread safe implementation of a double-linked list
type ConcurrentLinkedList[T any] struct {
	mu    sync.RWMutex
	first *listItem[T]
	last  *listItem[T]
	size  int
}

// RemoveFirst removes the first item from this list and returns its value and true if it exists.
// If the list is empty, a default value (zero value) of type T and false is returned.
func (clist *ConcurrentLinkedList[T]) RemoveFirst() (T, bool) {
	var res T
	if clist.first != nil {
		res = clist.removeItem(clist.first)
		return res, true
	}
	return res, false
}

// RemoveLast removes the last item from this list and returns its value and true if it exists.
// If the list is empty, a default value of type T (zero value) and false is returned.
func (clist *ConcurrentLinkedList[T]) RemoveLast() (T, bool) {
	var res T
	clist.mu.Lock()
	defer clist.mu.Unlock()
	if clist.last != nil {
		res = clist.removeItem(clist.last)
		return res, true
	}
	return res, false
}

// Remove removes the element at the specified position in this list and returns its value
// or a default value (zero value) of type T and an error if the index is out of range.
func (clist *ConcurrentLinkedList[T]) Remove(index int) (T, error) {
	clist.mu.Lock()
	item, err := clist.getByIndex(index)
	var res T
	if err == nil {
		res = clist.removeItem(item)
	}
	clist.mu.Unlock()
	return res, err
}
func (clist *ConcurrentLinkedList[T]) removeItem(item *listItem[T]) T {
	res := item.value
	item.removeYourself()
	if clist.first == item {
		clist.first = item.next
	}
	if clist.last == item {
		clist.last = item.prev
	}
	clist.size--
	return res
}

// RemoveLastOccurrence removes from the list the last occurrence of an element that satisfies the condition
// specified by the needToRemove function (when traversing the list from tail to head).
// Returns the value and index of the removed element, or the zero value of type T and -1 if no element was removed.
//   - needToRemove - a function that is applied to each element to determine if it should be deleted
func (clist *ConcurrentLinkedList[T]) RemoveLastOccurrence(needToRemove func(value T) bool) (T, int) {
	clist.mu.Lock()
	defer clist.mu.Unlock()
	index := clist.size
	item := clist.last
	for item != nil {
		index--
		if needToRemove(item.value) {
			return clist.removeItem(item), index
		}
		item = item.prev
	}
	var res T
	return res, -1
}

// RemoveFirstOccurrence removes from the list the first occurrence of an element that satisfies the condition
// specified by the function (when traversing the list from head to tail).
// Returns the value and index of the removed element, or the zero value of type T and -1 if no element was removed.
//   - needToRemove - a function that is applied to each element to determine if it should be deleted
func (clist *ConcurrentLinkedList[T]) RemoveFirstOccurrence(needToRemove func(value T) bool) (T, int) {
	index := -1
	clist.mu.Lock()
	defer clist.mu.Unlock()
	item := clist.first
	for item != nil {
		index++
		if needToRemove(item.value) {
			return clist.removeItem(item), index
		}
		item = item.next
	}
	var res T
	return res, -1
}

// RemoveAll removes from the list all elements that satisfy the condition specified by the needToRemove function.
// Returns the number of elements removed
//   - needToRemove - a function that is applied to each element to determine if it should be deleted
func (clist *ConcurrentLinkedList[T]) RemoveAll(needRemove func(value T) bool) int {
	result := 0
	clist.mu.Lock()
	item := clist.first
	for item != nil {
		if needRemove(item.value) {
			clist.removeItem(item)
			result++
		}
		item = item.next
	}
	clist.mu.Unlock()
	return result
}

// AddFirst inserts specified element to the beginning this list.
//   - value - the value to be inserted
func (clist *ConcurrentLinkedList[T]) AddFirst(value T) {
	item := &listItem[T]{value: value}
	clist.mu.Lock()
	if clist.first != nil {
		clist.first.insert(item)
	} else {
		clist.last = item
	}
	clist.first = item
	clist.size++
	clist.mu.Unlock()
}

// AddLast appends specified element to the end of this list.
//   - value - the value to be appended
func (clist *ConcurrentLinkedList[T]) AddLast(value T) {
	item := &listItem[T]{value: value}
	clist.mu.Lock()
	clist.addLastInner(item)
	clist.mu.Unlock()
}
func (clist *ConcurrentLinkedList[T]) addLastInner(item *listItem[T]) {
	if clist.last != nil {
		clist.last.append(item)
	} else {
		clist.first = item
	}
	clist.last = item
	clist.size++
}

// GetFirst returns the first element of this list and true if it exists.
// If the list is empty, this method returns the zero value of type T and false
func (clist *ConcurrentLinkedList[T]) GetFirst() (T, bool) {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	if clist.first != nil {
		return clist.first.value, true
	}
	var res T
	return res, false
}

// GetLast returns the last element of this list and true if it exists.
// If the lis is empty, this method returns the zero value of type T and false.
func (clist *ConcurrentLinkedList[T]) GetLast() (T, bool) {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	if clist.last != nil {
		return clist.last.value, true
	}
	var res T
	return res, false
}

// Get returns an item at the specified position in this list
// or the zero value of type T and an error if the index is out of range.
func (clist *ConcurrentLinkedList[T]) Get(index int) (T, error) {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	item, err := clist.getByIndex(index)
	if err != nil {
		var res T
		return res, err
	}
	return item.value, nil
}
func (clist *ConcurrentLinkedList[T]) getByIndex(index int) (*listItem[T], error) {
	if index >= 0 && index < clist.size {
		for i, item := 0, clist.first; item != nil; i, item = i+1, item.next {
			if i == index {
				return item, nil
			}
		}
	}
	return nil, ErrIndexOutOfRange
}

// ToArray returns an array containing all elements of this list in the proper sequence
// (from the first to the last element).
func (clist *ConcurrentLinkedList[T]) ToArray() []T {
	clist.mu.RLock()
	result := make([]T, clist.size)
	for i, item := 0, clist.first; item != nil; i, item = i+1, item.next {
		result[i] = item.value
	}
	clist.mu.RUnlock()
	return result
}

// Clear clears this list
func (clist *ConcurrentLinkedList[T]) Clear() {
	clist.mu.Lock()
	clist.first = nil
	clist.last = nil
	clist.size = 0
	clist.mu.Unlock()
}

// Size returns the number of elements in this list
func (clist *ConcurrentLinkedList[T]) Size() int {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	return clist.size
}
func NewConcurrentLinkedList[T any]() *ConcurrentLinkedList[T] {
	return &ConcurrentLinkedList[T]{}
}
func NewConcurrentLinkedListItems[T any](values ...T) *ConcurrentLinkedList[T] {
	result := NewConcurrentLinkedList[T]()
	result.mu.Lock()
	for _, val := range values {
		result.addLastInner(&listItem[T]{value: val})
	}
	result.mu.Unlock()
	return result
}
