// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collections contains some thread safe collections.
package collections

import "sync"

// ConcurrentMap is a thread safe map.
// A ConcurrentMap is safe for concurrent use by multiple goroutines.
//   - K - comparable key type;
//   - V - value type.
type ConcurrentMap[K comparable, V any] struct {
	mu       sync.RWMutex
	mp       map[K]V
	capacity int
}

// ForEachRead performs a given action for each (key, value)
//   - f - the function, that will be called for each (key, value) pair in ConcurrentMap
//
// It should not be used to modify values if the value type (V) is a reference type,
// because a read lock is used under the hood.
// Note! ConcurrentMap methods, such as Get and Size can be used inside the 'f' function.
// However, you should not use methods that modify ConcurrentMap, as this will cause a deadlock.
func (cmap *ConcurrentMap[K, V]) ForEachRead(f func(key K, value V)) {
	cmap.mu.RLock()
	for k, v := range cmap.mp {
		f(k, v)
	}
	cmap.mu.RUnlock()
}

// ForEach performs a given action for each (key, value)
//   - f - the function, that will be called for each (key, value) pair in ConcurrentMap
//
// If the value type (V) is a reference type, this method can be used to modify values
// Note! Do NOT USE ConcurrentMap methods inside the 'f' function, as this will cause a deadlock.
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) ForEach(f func(key K, value V)) {
	cmap.mu.Lock()
	for k, v := range cmap.mp {
		f(k, v)
	}
	cmap.mu.Unlock()
} //revive:enable:confusing-naming

// PutIfNotExists maps the specified key (key) to the specified value (value)
// if the key doesn't exist returns true and a new value (value).
// If the key exists, the new value will not be mapped to it, the method returns false and the previous key (key) value.
//   - key - the key with which a specified value is to be assigned
//   - value - the value to be associated with the specified key
func (cmap *ConcurrentMap[K, V]) PutIfNotExists(key K, value V) (bool, V) {
	cmap.mu.Lock()
	defer cmap.mu.Unlock()
	if old, ok := cmap.mp[key]; ok {
		return false, old
	}
	cmap.mp[key] = value
	return true, value
}

// PutIfNotExistsDoubleCheck does the same thing as PutIfNotExists, but before doing so,
// it checks the existence of the key (key) using the Get method.
//   - key - the key with which a specified value is to be assigned
//   - value - the value to be associated with the specified key
func (cmap *ConcurrentMap[K, V]) PutIfNotExistsDoubleCheck(key K, value V) (bool, V) {
	old, ok := cmap.Get(key)
	if ok {
		return false, old
	}
	return cmap.PutIfNotExists(key, value)
}

// RemoveIfExistsDoubleCheck removes the key and its corresponding value,
// before this method checks the existence of the key using the Get method.
//   - key - the key that needs to be removed
func (cmap *ConcurrentMap[K, V]) RemoveIfExistsDoubleCheck(key K) (bool, V) {
	old, ok := cmap.Get(key)
	if !ok {
		return false, old
	}
	return cmap.RemoveIfExists(key)
}

// RemoveIfExists removes the key and its corresponding value.
// If the key exists, the method returns true and the value corresponding to that key,
// otherwise it returns false and the default value for the value type.
//   - key - the key that needs to be removed
func (cmap *ConcurrentMap[K, V]) RemoveIfExists(key K) (bool, V) {
	cmap.mu.Lock()
	defer cmap.mu.Unlock()
	old, ok := cmap.mp[key]
	if !ok {
		return false, old
	}
	delete(cmap.mp, key)
	return true, old
}

// Remove removes the key and its corresponding value from the ConcurrentMap.
//   - key - the key that needs to be removed
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) Remove(key K) {
	cmap.mu.Lock()
	delete(cmap.mp, key)
	cmap.mu.Unlock()
} //revive:enable:confusing-naming

// Put maps the specified key (key) to the specified value (value).
// The value can be retrieved by calling the Get method with a key that is equal to the original key.
//   - key - the key with which a specified value is to be assigned
//   - value - the value to be associated with the specified key
func (cmap *ConcurrentMap[K, V]) Put(key K, value V) {
	cmap.mu.Lock()
	cmap.mp[key] = value
	cmap.mu.Unlock()
}

// Get returns the value to which the specified key is mapped and the sign of existence of this value.
//   - key - the key whose value will be returned
//
// If a value for the key exists, its value is returned and true,
// otherwise the default value for the value type is returned and false.
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	cmap.mu.RLock()
	val, ok := cmap.mp[key]
	cmap.mu.RUnlock()
	return val, ok
} //revive:enable:confusing-naming

// Keys returns a slice of the keys contained in this map
func (cmap *ConcurrentMap[K, V]) Keys() []K {
	cmap.mu.RLock()
	result := make([]K, 0, len(cmap.mp))
	for k := range cmap.mp {
		result = append(result, k)
	}
	cmap.mu.RUnlock()
	return result
}

// Size returns the number of key-value mappings in this map.
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) Size() int {
	cmap.mu.RLock()
	defer cmap.mu.RUnlock()
	return len(cmap.mp)
} //revive:enable:confusing-naming

// IsEmpty returns true if the ConcurrentMap does not contain any (key, value) pairs
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) IsEmpty() bool {
	cmap.mu.RLock()
	defer cmap.mu.RUnlock()
	return len(cmap.mp) == 0
} //revive:enable:confusing-naming

// Copy returns a shallow copy of this ConcurrentMap instance: the keys and the values themselves are not copies.
func (cmap *ConcurrentMap[K, V]) Copy() map[K]V {
	cmap.mu.RLock()
	result := make(map[K]V, len(cmap.mp))
	for key, value := range cmap.mp {
		result[key] = value
	}
	cmap.mu.RUnlock()
	return result
}

// TrimToSize trims the capacity of this ConcurrentMap instance to be the map's current size.
// An application can use this operation to minimize the storage of a ConcurrentMap instance.
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) TrimToSize() {
	cmap.mu.Lock()
	tmp := make(map[K]V, len(cmap.mp))
	for k, v := range cmap.mp {
		tmp[k] = v
	}
	cmap.mp = tmp
	cmap.mu.Unlock()
} //revive:enable:confusing-naming

// Clear clears the map
//
//revive:disable:confusing-naming
func (cmap *ConcurrentMap[K, V]) Clear() {
	cmap.mu.Lock()
	if cmap.capacity > 0 {
		cmap.mp = make(map[K]V, cmap.capacity)
	} else {
		cmap.mp = make(map[K]V)
	}
	cmap.mu.Unlock()
} //revive:enable:confusing-naming

// NewConcurrentMap creates and returns a new empty ConcurrentMap instance.
//   - K - comparable key type;
//   - V - value type.
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{mp: make(map[K]V)}
}

// NewConcurrentMapCapacity creates and returns a new empty ConcurrentMap with an initial space size (capacity).
//   - K - comparable key type;
//   - V - value type;
//   - capacity - initial space size.
func NewConcurrentMapCapacity[K comparable, V any](capacity int) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{mp: make(map[K]V, capacity), capacity: capacity}
}
