// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package caches contains some thread safe collections.
package caches

import "sync"

type LRU[K comparable, V any] struct {
	mu       sync.RWMutex
	mp       map[K]*lruEntity[K, V]
	entities *entityList[K, V]
	limit    int
}

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
func (lru *LRU[K, V]) PutIfNotExists(key K, value V) (bool, V) {
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
func (lru *LRU[K, V]) Clear() {
	lru.mu.Lock()
	lru.mp = make(map[K]*lruEntity[K, V], lru.limit)
	lru.entities.clear()
	lru.mu.Unlock()
}
func (lru *LRU[K, V]) Size() int {
	lru.mu.RLock()
	defer lru.mu.RUnlock()
	return len(lru.mp)
}
func NewLRU[K comparable, V any](limit int) *LRU[K, V] {
	return &LRU[K, V]{mp: make(map[K]*lruEntity[K, V], limit), entities: &entityList[K, V]{}, limit: limit}
}
