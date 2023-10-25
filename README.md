[![Go](https://github.com/PavloVM7/go-concurrency/actions/workflows/go.yml/badge.svg)](https://github.com/PavloVM7/go-concurrency/actions/workflows/go.yml)
# go-concurrency

This module contains some thread safe entities and collections

## ConcurrentMap

`ConcurrentMap` is a thread-safe map implementation

### How to use

``` go
cm := NewConcurrentMap[int, int]() // or NewConcurrentMapCapacity[int, int](128) with initial capacity 128

go func() {
    for ... {
        if ok, _ := cm.PutIfNotExistsDoubleCheck(i, num); ok {
            // do something
        }
    }
}()

go func() {
    for ... {
        if ok, old := cm.PutIfNotExistsDoubleCheck(i, num); !ok {
            // do something with old value
        }
    }
}()

ticker := time.NewTicker(1 * time.Second)
go func() {
    for range ticker.C {
        if cm.Size() > 100_000 {
            cm.ForEachRead(func(key int, value int) { 
                // process each (key, value) pair 
            })
            cm.Clear()
        }
    }
}()
```
