package caches

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkLRU_PutIfNotExists(b *testing.B) {
	lru := NewLRU[int, string](10)
	b.ResetTimer()
	var (
		val string
		ok  bool
	)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		lru.Clear()
		b.StartTimer()
		ok, val = lru.PutIfNotExists(1, "value")
	}
	b.StopTimer()
	assert.True(b, ok)
	assert.Equal(b, "value", val)
}

func BenchmarkLRU_Put(b *testing.B) {
	const limit = 10
	lru := NewLRU[int, string](limit)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		lru.Clear()
		b.StartTimer()
		lru.Put(1, "value")
	}
	b.StopTimer()
	assert.Equal(b, 1, lru.Size())
}
func createTestValues(count int) []string {
	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, fmt.Sprintf("value%d", i))
	}
	return result
}
