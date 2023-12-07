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

type ConcurrentLinkedList[T any] struct {
	mu    sync.RWMutex
	first *listItem[T]
	last  *listItem[T]
	size  int
}

func (clist *ConcurrentLinkedList[T]) RemoveFirst() (T, bool) {
	var res T
	if clist.first != nil {
		res = clist.removeItem(clist.first)
		return res, true
	}
	return res, false
}
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

func (clist *ConcurrentLinkedList[T]) GetFirst() (T, bool) {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	if clist.first != nil {
		return clist.first.value, true
	}
	var res T
	return res, false
}
func (clist *ConcurrentLinkedList[T]) GetLast() (T, bool) {
	clist.mu.RLock()
	defer clist.mu.RUnlock()
	if clist.last != nil {
		return clist.last.value, true
	}
	var res T
	return res, false
}
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
func (clist *ConcurrentLinkedList[T]) ToArray() []T {
	clist.mu.RLock()
	result := make([]T, clist.size)
	for i, item := 0, clist.first; item != nil; i, item = i+1, item.next {
		result[i] = item.value
	}
	clist.mu.RUnlock()
	return result
}
func (clist *ConcurrentLinkedList[T]) Clear() {
	clist.mu.Lock()
	clist.first = nil
	clist.last = nil
	clist.size = 0
	clist.mu.Unlock()
}
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
