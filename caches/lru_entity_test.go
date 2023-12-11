package caches

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_lruEntity_removeYourself_prev(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	entity2.insertBefore(entity1)
	entity2.insertAfter(entity3)

	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity2, entity3.prev)
	assert.Nil(t, entity1.prev)
	assert.Nil(t, entity3.next)

	entity1.removeYourself()

	assert.Nil(t, entity2.prev)
}
func Test_lruEntity_removeYourself_last(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	entity2.insertBefore(entity1)
	entity2.insertAfter(entity3)

	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity2, entity3.prev)
	assert.Nil(t, entity1.prev)
	assert.Nil(t, entity3.next)

	entity3.removeYourself()

	assert.Nil(t, entity2.next)
}
func Test_lruEntity_removeYourself(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	entity2.insertBefore(entity1)
	entity2.insertAfter(entity3)

	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity2, entity3.prev)
	assert.Nil(t, entity1.prev)
	assert.Nil(t, entity3.next)

	entity2.removeYourself()
	assert.Same(t, entity3, entity1.next)
	assert.Same(t, entity1, entity3.prev)
}
func Test_lruEntity_insertAfter(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	entity1.insertAfter(entity3)
	entity1.insertAfter(entity2)

	assert.Nil(t, entity1.prev)
	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity1, entity2.prev)

	assert.Same(t, entity3, entity2.next)
	assert.Same(t, entity2, entity3.prev)
	assert.Nil(t, entity3.next)
}
func Test_lruEntity_insertAfter_no_next(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)

	entity1.insertAfter(entity2)
	assert.Nil(t, entity1.prev)
	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity1, entity2.prev)
}
func Test_lruEntity_insertBefore_no_prev(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)

	entity2.insertBefore(entity1)

	assert.Nil(t, entity1.prev)
	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity1, entity2.prev)
}
func Test_lruEntity_insertBefore(t *testing.T) {
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	entity3.insertBefore(entity1)
	entity3.insertBefore(entity2)

	assert.Nil(t, entity1.prev)
	assert.Same(t, entity2, entity1.next)
	assert.Same(t, entity1, entity2.prev)

	assert.Same(t, entity3, entity2.next)
	assert.Same(t, entity2, entity3.prev)
	assert.Nil(t, entity3.next)
}

func createTestEntity(num int) *lruEntity[int, string] {
	return &lruEntity[int, string]{key: num, value: fmt.Sprintf("value%d", num)}
}
