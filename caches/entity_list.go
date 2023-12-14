// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package caches

type entityList[K any, V any] struct {
	head *lruEntity[K, V]
	tail *lruEntity[K, V]
}

func (el *entityList[K, V]) setHead(entity *lruEntity[K, V]) {
	entity.prev = nil
	if el.head != nil {
		el.head.insertBefore(entity)
	} else {
		el.tail = entity
	}
	el.head = entity
}
func (el *entityList[K, V]) moveToHead(entity *lruEntity[K, V]) {
	if el.head == entity {
		return
	}
	el.removeEntity(entity)
	el.setHead(entity)
}
func (el *entityList[K, V]) removeEntity(entity *lruEntity[K, V]) {
	entity.removeYourself()
	if el.head == entity {
		el.head = entity.next
	}
	if el.tail == entity {
		el.tail = entity.prev
	}
}
func (el *entityList[K, V]) clear() {
	el.head = nil
	el.tail = nil
}
