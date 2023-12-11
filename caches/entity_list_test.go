package caches

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_entityList_moveToHead(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	list.setHead(entity1)
	list.setHead(entity2)
	list.setHead(entity3)

	assert.Same(t, entity3, list.head)
	assert.Same(t, entity1, list.tail)
	assert.Same(t, entity2, list.head.next)
	assert.Same(t, entity2, list.tail.prev)
	assert.Nil(t, list.head.prev)
	assert.Nil(t, list.tail.next)

	list.moveToHead(entity2)

	assert.Same(t, entity2, list.head)
	assert.Same(t, entity1, list.tail)

	list.moveToHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)
	assert.Same(t, entity2, list.head.next)
	assert.Same(t, entity2, list.tail.prev)
	assert.Nil(t, list.head.prev)
	assert.Nil(t, list.tail.next)
}
func Test_entityList_removeEntity_sole(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)

	list.setHead(entity1)
	assert.Same(t, entity1, list.head)
	assert.Same(t, entity1, list.tail)

	list.removeEntity(entity1)

	assert.Nil(t, list.head)
	assert.Nil(t, list.tail)
}
func Test_entityList_removeEntity_first(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	list.setHead(entity3)
	list.setHead(entity2)
	list.setHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)

	list.removeEntity(entity1)

	assert.Same(t, entity2, list.head)
	assert.Same(t, entity3, list.tail)

	assert.Same(t, entity3, list.head.next)
	assert.Same(t, entity2, list.tail.prev)
}
func Test_entityList_removeEntity_tail(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	list.setHead(entity3)
	list.setHead(entity2)
	list.setHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)

	list.removeEntity(list.tail)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity2, list.tail)
	assert.Same(t, entity2, list.head.next)
	assert.Same(t, entity1, list.tail.prev)
}
func Test_entityList_removeEntity_last(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	list.setHead(entity3)
	list.setHead(entity2)
	list.setHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)

	list.removeEntity(entity3)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity2, list.tail)
	assert.Same(t, entity2, list.head.next)
	assert.Same(t, entity1, list.tail.prev)
}

func Test_entityList_removeEntity(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	list.setHead(entity3)
	list.setHead(entity2)
	list.setHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)

	list.removeEntity(entity2)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity3, list.tail)
	assert.Same(t, entity1, entity3.prev)
	assert.Same(t, entity3, entity1.next)
}
func Test_entityList_setHead(t *testing.T) {
	list := createTestList()
	entity1 := createTestEntity(1)
	entity2 := createTestEntity(2)
	entity3 := createTestEntity(3)

	assert.Nil(t, list.head)
	assert.Nil(t, list.tail)

	list.setHead(entity1)

	assert.Same(t, entity1, list.head)
	assert.Same(t, entity1, list.tail)

	list.setHead(entity2)

	assert.Same(t, entity2, list.head)
	assert.Same(t, entity1, list.tail)

	list.setHead(entity3)

	assert.Same(t, entity3, list.head)
	assert.Same(t, entity1, list.tail)
	assert.Same(t, entity2, list.head.next)
	assert.Same(t, entity2, list.tail.prev)
	assert.Nil(t, list.head.prev)
	assert.Nil(t, list.tail.next)
}

func createTestList() *entityList[int, string] {
	return &entityList[int, string]{}
}
