// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package caches contains some thread safe collections.
package caches

import (
	"fmt"
	"sync"
)

// LRU (least recently used) is a cache that deletes the least-recently-used items.
// The LRU is safe for concurrent use by multiple goroutines.
// - K - comparable key type
// - V - value type
type LRU[K comparable, V any] struct {
	mu       sync.RWMutex
	mp       map[K]*lruEntity[K, V]
	entities *entityList[K, V]
	limit    int
}

// Put maps the specified key to the specified value
//   - key - the key with which a specified value is to be assigned
//   - value - the value to be associated with the specified key
func (lru *LRU[K, V]) Put(key K, value V) {
	lru.mu.Lock()
	entity, ok := lru.mp[key]
	if !ok {
		entity = &lruEntity[K, V]{key: key, value: value}
		lru.putEntity(entity)
	} else {
		entity.value = value
		lru.entities.moveToHead(entity)
	}
	lru.mu.Unlock()
}
func (lru *LRU[K, V]) putEntity(entity *lruEntity[K, V]) {
	lru.mp[entity.key] = entity
	lru.entities.setHead(entity)
	if len(lru.mp) > lru.limit {
		lru.evictEntity(lru.entities.tail)
	}
}

// PutIfAbsent maps the specified key to the specified value
// if the key doesn't exist returns true and a new value.
// If the key exists, the new value will not be mapped to it, the method returns false and the previous key value.
//   - key - the key with which a specified value is to be assigned
//   - value - the value to be associated with the specified key
func (lru *LRU[K, V]) PutIfAbsent(key K, value V) (bool, V) {
	lru.mu.Lock()
	entity, ok := lru.mp[key]
	if !ok {
		entity = &lruEntity[K, V]{key: key, value: value}
		lru.putEntity(entity)
	}
	lru.mu.Unlock()
	return !ok, entity.value
}

func (lru *LRU[K, V]) evictEntity(entity *lruEntity[K, V]) {
	lru.entities.removeEntity(entity)
	entity.prev = nil
	entity.next = nil
	delete(lru.mp, entity.key)
}

// Get returns the value to which the specified key is mapped and the sign of existence of this value.
// If a value for the key exists, its value is returned and true,
// otherwise the default value for the value type is returned and false.
//   - key - the key whose value will be returned
func (lru *LRU[K, V]) Get(key K) (bool, V) {
	var res V
	lru.mu.Lock()
	entity, ok := lru.mp[key]
	if ok {
		res = entity.value
		lru.entities.moveToHead(entity)
	}
	lru.mu.Unlock()
	return ok, res
}

// Evict evicts the value to which the specified key is mapped.
//   - key - the key that needs to be removed
func (lru *LRU[K, V]) Evict(key K) (bool, V) {
	var res V
	lru.mu.Lock()
	entity, ok := lru.mp[key]
	if ok {
		res = entity.value
		lru.evictEntity(entity)
	}
	lru.mu.Unlock()
	return ok, res
}

// Copy returns a shallow copy of this LRU cache instance: the keys and the values themselves are not copies.
func (lru *LRU[K, V]) Copy() map[K]V {
	lru.mu.RLock()
	result := make(map[K]V, len(lru.mp))
	for k, e := range lru.mp {
		result[k] = e.value
	}
	lru.mu.RUnlock()
	return result
}

// Clear clears the cache.
//
//revive:disable:confusing-naming
func (lru *LRU[K, V]) Clear() {
	lru.mu.Lock()
	lru.mp = make(map[K]*lruEntity[K, V], lru.limit)
	lru.entities.clear()
	lru.mu.Unlock()
} //revive:enable:confusing-naming

// Size returns the number of key-value mappings in this cache.
func (lru *LRU[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return len(lru.mp)
}

// String prints the LRU cache limit value and the number of key-value mappings in this cache
func (lru *LRU[K, V]) String() string {
	lru.mu.RLock()
	lmt := lru.limit
	sz := len(lru.mp)
	lru.mu.RUnlock()
	return fmt.Sprintf("LRU{limit: %d; size: %d}", lmt, sz)
}

// NewLRU creates and returns a new LRU cache.
// - limit - specifies the max number of key-value pairs that we want to keep.
// - K - comparable key type
// - V - value type
func NewLRU[K comparable, V any](limit int) *LRU[K, V] {
	return &LRU[K, V]{mp: make(map[K]*lruEntity[K, V], limit), entities: &entityList[K, V]{}, limit: limit}
}
