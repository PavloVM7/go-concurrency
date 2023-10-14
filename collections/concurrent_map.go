package collections

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	sync.RWMutex
	mp map[K]V
}

func (cmap *ConcurrentMap[K, V]) PutIfNotExists(key K, value V) (bool, V) {
	cmap.Lock()
	defer cmap.Unlock()
	old, ok := cmap.mp[key]
	if ok {
		return false, old
	}
	cmap.mp[key] = value
	return true, value
}
func (cmap *ConcurrentMap[K, V]) PutIfNotExistsDoubleCheck(key K, value V) (bool, V) {
	old, ok := cmap.Get(key)
	if ok {
		return false, old
	}
	return cmap.PutIfNotExists(key, value)
}

func (cmap *ConcurrentMap[K, V]) RemoveIfExists(key K) (bool, V) {
	old, ok := cmap.Get(key)
	if !ok {
		return false, old
	}
	cmap.Lock()
	defer cmap.Unlock()
	old, ok = cmap.mp[key]
	if !ok {
		return false, old
	}
	delete(cmap.mp, key)
	return true, old
}
func (cmap *ConcurrentMap[K, V]) Remove(key K) {
	cmap.Lock()
	delete(cmap.mp, key)
	cmap.Unlock()
}
func (cmap *ConcurrentMap[K, V]) Put(key K, value V) {
	cmap.Lock()
	cmap.mp[key] = value
	cmap.Unlock()
}
func (cmap *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	cmap.RLock()
	val, ok := cmap.mp[key]
	cmap.RUnlock()
	return val, ok
}
func (cmap *ConcurrentMap[K, V]) Size() int {
	cmap.RLock()
	defer cmap.RUnlock()
	return len(cmap.mp)
}

func (cmap *ConcurrentMap[K, V]) Copy() map[K]V {
	cmap.RLock()
	result := make(map[K]V, len(cmap.mp))
	for key, value := range cmap.mp {
		result[key] = value
	}
	cmap.RUnlock()
	return result
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{mp: make(map[K]V)}
}
func NewConcurrentMapCapacity[K comparable, V any](capacity int) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{mp: make(map[K]V, capacity)}
}
