package caches

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testLruLimit = 3

func TestLRU_Evict(t *testing.T) {
	keys := []int{1, 2, 3}
	values := []string{"value1", "value2", "value3"}
	lru := createTestLru()
	for i := 0; i < len(keys); i++ {
		lru.Put(keys[i], values[i])
	}
	assert.Equal(t, len(keys), lru.Size())
	for i := 0; i < len(keys); i++ {
		ok, val := lru.Evict(keys[i])
		assert.True(t, ok)
		assert.Equal(t, values[i], val)
	}
	ok, val := lru.Evict(123)
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestLRU_Get_evicted(t *testing.T) {
	keys := []int{1, 2, 3, 4, 5}
	values := []string{"value1", "value2", "value3", "value4", "value5"}

	lru := createTestLru()

	for i := 0; i < len(keys); i++ {
		lru.PutIfNotExists(keys[i], values[i])
	}

	assert.Equal(t, testLruLimit, lru.Size())

	n := len(keys) - testLruLimit
	exists := keys[n:]
	existsValues := values[n:]
	n = testLruLimit - 1
	notExists := keys[:n]
	//notExistsValues := values[:n]

	for i := 0; i < len(exists); i++ {
		ok, actual := lru.Get(exists[i])
		assert.True(t, ok)
		assert.Equal(t, actual, existsValues[i])
	}
	for i := 0; i < len(notExists); i++ {
		ok, actual := lru.Get(notExists[i])
		assert.False(t, ok)
		assert.Equal(t, "", actual)
	}

}
func TestLRU_Get(t *testing.T) {
	keys := []int{1, 2, 3}
	values := []string{"value1", "value2", "value3"}

	lru := createTestLru()
	for i := 0; i < len(keys); i++ {
		lru.Put(keys[i], values[i])
	}

	assert.Equal(t, testLruLimit, lru.Size())

	for i := 0; i < len(keys); i++ {
		ok, actual := lru.Get(keys[i])
		assert.True(t, ok)
		assert.Equal(t, values[i], actual)
	}
}

func TestLRU_PutIfNotExists_evict(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	value2 := "value2"
	value3 := "value3"

	lru.PutIfNotExists(1, value1)
	lru.PutIfNotExists(2, value2)
	lru.PutIfNotExists(3, value3)

	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())

	value4 := "value4"
	lru.PutIfNotExists(4, value4)

	assert.Equal(t, value4, lru.entities.head.value)
	assert.Equal(t, value2, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())
}
func TestLRU_PutIfNotExists_no_override(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	lru.PutIfNotExists(1, value1)
	value2 := "value2"
	lru.PutIfNotExists(2, value2)
	value3 := "value3"
	lru.PutIfNotExists(3, value3)
	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())

	ok, val := lru.PutIfNotExists(1, "other value for key 1")

	assert.False(t, ok)
	assert.Equal(t, value1, val)
	assert.Equal(t, testLruLimit, lru.Size())

	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)
}
func TestLRU_PutIfNotExists(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	ok, val := lru.PutIfNotExists(1, value1)
	assert.True(t, ok)
	assert.Equal(t, value1, val)
	assert.Equal(t, value1, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)
	value2 := "value2"
	ok, val = lru.PutIfNotExists(2, value2)
	assert.True(t, ok)
	assert.Equal(t, value2, val)
	assert.Equal(t, value2, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)
	value3 := "value3"
	ok, val = lru.PutIfNotExists(3, value3)
	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())
}

func TestLRU_Put_evict(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	value2 := "value2"
	value3 := "value3"

	lru.Put(1, value1)
	lru.Put(2, value2)
	lru.Put(3, value3)

	assert.Equal(t, testLruLimit, lru.Size())

	value4 := "value4"
	lru.Put(4, value4)

	assert.Equal(t, value4, lru.entities.head.value)
	assert.Equal(t, value2, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())
}
func TestLRU_Put_override(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	lru.Put(1, value1)
	value2 := "value2"
	lru.Put(2, value2)
	value3 := "value3"
	lru.Put(3, value3)
	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)

	value11 := "other value for key 1"

	lru.Put(1, value11)

	assert.Equal(t, testLruLimit, lru.Size())
	assert.Equal(t, value11, lru.entities.head.value)
	assert.Equal(t, value2, lru.entities.tail.value)
}
func TestLRU_Put(t *testing.T) {
	lru := createTestLru()
	value1 := "value1"
	lru.Put(1, value1)
	assert.Equal(t, value1, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)
	value2 := "value2"
	lru.Put(2, value2)
	assert.Equal(t, value2, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)
	value3 := "value3"
	lru.Put(3, value3)
	assert.Equal(t, value3, lru.entities.head.value)
	assert.Equal(t, value1, lru.entities.tail.value)

	assert.Equal(t, testLruLimit, lru.Size())
}

func createTestLru() *LRU[int, string] {
	return NewLRU[int, string](testLruLimit)
}
