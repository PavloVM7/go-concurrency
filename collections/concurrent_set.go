package collections

import "sync"

type ConcurrentSet[T comparable] struct {
	sync.RWMutex
	mp       map[T]interface{}
	capacity int
}

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
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]interface{})}
}
func NewConcurrentSetCapacity[T comparable](capacity int) *ConcurrentSet[T] {
	return &ConcurrentSet[T]{mp: make(map[T]interface{}, capacity), capacity: capacity}
}
