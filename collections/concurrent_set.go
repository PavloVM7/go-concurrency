package collections

import "sync"

type ConcurrentSet[T comparable] struct {
	sync.RWMutex
	mp       map[T]interface{}
	capacity int
}

// Add adds a specified value to the set.
// Returns true if the value did not exist and was added to the set, otherwise returns false.
func (cset *ConcurrentSet[T]) Add(value T) bool {
	cset.Lock()
	defer cset.Unlock()
	if _, ok := cset.mp[value]; !ok {
		cset.mp[value] = nil
		return true
	}
	return false
}
func (cset *ConcurrentSet[T]) Contains(value T) bool {
	cset.RLock()
	_, res := cset.mp[value]
	cset.RUnlock()
	return res
}
func (cset *ConcurrentSet[T]) Size() int {
	cset.RLock()
	defer cset.RUnlock()
	return len(cset.mp)
}
func (cset *ConcurrentSet[T]) ToSlice() []T {
	cset.RLock()
	result := make([]T, 0, len(cset.mp))
	for k := range cset.mp {
		result = append(result, k)
	}
	cset.RUnlock()
	return result
}

// NewConcurrentSet returns a new empty ConcurrentSet instance
//   - T - value type
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]interface{})}
}

// NewConcurrentSetCapacity returns a new empty ConcurrentSet instance with an initial space size (capacity)
//   - T - value type
//   - capacity - initial space size
func NewConcurrentSetCapacity[T comparable](capacity int) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]interface{}, capacity), capacity: capacity}
}
