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
	sync.RWMutex
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
	clist.Lock()
	defer clist.Unlock()
	if clist.last != nil {
		res = clist.removeItem(clist.last)
		return res, true
	}
	return res, false
}
func (clist *ConcurrentLinkedList[T]) Remove(index int) (T, error) {
	clist.Lock()
	item, err := clist.getByIndex(index)
	var res T
	if err == nil {
		res = clist.removeItem(item)
	}
	clist.Unlock()
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
func (clist *ConcurrentLinkedList[T]) AddFirst(value T) {
	item := &listItem[T]{value: value}
	clist.Lock()
	if clist.first != nil {
		clist.first.insert(item)
	} else {
		clist.last = item
	}
	clist.first = item
	clist.size++
	clist.Unlock()
}
func (clist *ConcurrentLinkedList[T]) AddLast(value T) {
	item := &listItem[T]{value: value}
	clist.Lock()
	clist.addLastInner(item)
	clist.Unlock()
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
	clist.RLock()
	defer clist.RUnlock()
	if clist.first != nil {
		return clist.first.value, true
	}
	var res T
	return res, false
}
func (clist *ConcurrentLinkedList[T]) GetLast() (T, bool) {
	clist.RLock()
	defer clist.RUnlock()
	if clist.last != nil {
		return clist.last.value, true
	}
	var res T
	return res, false
}
func (clist *ConcurrentLinkedList[T]) Get(index int) (T, error) {
	clist.RLock()
	defer clist.RUnlock()
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
	clist.RLock()
	result := make([]T, clist.size)
	for i, item := 0, clist.first; item != nil; i, item = i+1, item.next {
		result[i] = item.value
	}
	clist.RUnlock()
	return result
}
func (clist *ConcurrentLinkedList[T]) Clear() {
	clist.Lock()
	clist.first = nil
	clist.last = nil
	clist.size = 0
	clist.Unlock()
}
func (clist *ConcurrentLinkedList[T]) Size() int {
	clist.RLock()
	defer clist.RUnlock()
	return clist.size
}
func NewConcurrentLinkedList[T any]() *ConcurrentLinkedList[T] {
	return &ConcurrentLinkedList[T]{}
}
func NewConcurrentLinkedListItems[T any](values ...T) *ConcurrentLinkedList[T] {
	result := NewConcurrentLinkedList[T]()
	result.Lock()
	for _, val := range values {
		result.addLastInner(&listItem[T]{value: val})
	}
	result.Unlock()
	return result
}
