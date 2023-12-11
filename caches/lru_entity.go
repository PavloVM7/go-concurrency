// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package caches

import "fmt"

type lruEntity[K any, V any] struct {
	key   K
	value V
	prev  *lruEntity[K, V]
	next  *lruEntity[K, V]
}

func (e *lruEntity[K, V]) insertBefore(entity *lruEntity[K, V]) {
	entity.prev = e.prev
	entity.next = e
	e.prev = entity
	if entity.prev != nil {
		entity.prev.next = entity
	}
}
func (e *lruEntity[K, V]) insertAfter(entity *lruEntity[K, V]) {
	entity.next = e.next
	entity.prev = e
	e.next = entity
	if entity.next != nil {
		entity.next.prev = entity
	}
}
func (e *lruEntity[K, V]) removeYourself() {
	if e.prev != nil {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
}
func (e *lruEntity[K, V]) String() string {
	return fmt.Sprintf("('%v':'%v')", e.key, e.value)
}
