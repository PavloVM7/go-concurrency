// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import "sync"

// ConcurrentSet is a thread safe set.
// ConcurrentSet is safe for concurrent use by multiple goroutines.
//   - T - value type
type ConcurrentSet[T comparable] struct {
	mu       sync.RWMutex
	mp       map[T]struct{}
	capacity int
}

// ForEach performs a given action for each value of the ConcurrentSet
//   - f - the function, that will be called for each value in ConcurrentSet
//
// It should not be used to modify values if the value type (T) is a reference type,
// because a read lock is used under the hood.
func (cset *ConcurrentSet[T]) ForEach(f func(value T)) {
	cset.mu.RLock()
	for k := range cset.mp {
		f(k)
	}
	cset.mu.RUnlock()
}

// AddAll adds all the specified values to the ConcurrentSet.
// Returns true if this ConcurrentSet changed as result of the call.
func (cset *ConcurrentSet[T]) AddAll(values ...T) bool {
	changed := false
	cset.mu.Lock()
	for _, value := range values {
		if _, ok := cset.mp[value]; !ok {
			cset.mp[value] = struct{}{}
			changed = true
		}
	}
	cset.mu.Unlock()
	return changed
}

// Add adds a specified value to the set.
// Returns true if the value did not exist and was added to the set, otherwise returns false.
func (cset *ConcurrentSet[T]) Add(value T) bool {
	cset.mu.Lock()
	defer cset.mu.Unlock()
	if _, ok := cset.mp[value]; !ok {
		cset.mp[value] = struct{}{}
		return true
	}
	return false
}

// Contains returns true if the set contains the value
func (cset *ConcurrentSet[T]) Contains(value T) bool {
	cset.mu.RLock()
	_, res := cset.mp[value]
	cset.mu.RUnlock()
	return res
}

// Clear clears the set
func (cset *ConcurrentSet[T]) Clear() {
	cset.mu.Lock()
	if cset.capacity > 0 {
		cset.mp = make(map[T]struct{}, cset.capacity)
	} else {
		cset.mp = make(map[T]struct{})
	}
	cset.mu.Unlock()
}

// Size returns the current size of the ConcurrentSet
func (cset *ConcurrentSet[T]) Size() int {
	cset.mu.RLock()
	defer cset.mu.RUnlock()
	return len(cset.mp)
}

// IsEmpty returns true if the ConcurrentSet does not contain any values
func (cset *ConcurrentSet[T]) IsEmpty() bool {
	cset.mu.RLock()
	defer cset.mu.RUnlock()
	return len(cset.mp) == 0
}

// ToSlice returns a slice of ConcurrentSet elements
func (cset *ConcurrentSet[T]) ToSlice() []T {
	cset.mu.RLock()
	result := make([]T, 0, len(cset.mp))
	for k := range cset.mp {
		result = append(result, k)
	}
	cset.mu.RUnlock()
	return result
}

// NewConcurrentSet returns a new empty ConcurrentSet instance
//   - T - value type
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]struct{})}
}

// NewConcurrentSetCapacity returns a new empty ConcurrentSet instance with an initial space size (capacity)
//   - T - value type
//   - capacity - initial space size
func NewConcurrentSetCapacity[T comparable](capacity int) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]struct{}, capacity), capacity: capacity}
}

// NewConcurrentSetWithValues returns a new instance of ConcurrentSet containing specified values
//   - values ...T - values that the ConcurrentSet will contain
func NewConcurrentSetWithValues[T comparable](values ...T) *ConcurrentSet[T] {
	result := NewConcurrentSetCapacity[T](len(values))
	result.AddAll(values...)
	return result
}
