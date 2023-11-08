[![Go](https://github.com/PavloVM7/go-concurrency/actions/workflows/go.yml/badge.svg)](https://github.com/PavloVM7/go-concurrency/actions/workflows/go.yml)
# go-concurrency

This module contains some thread safe entities and collections

To install this package use the command:

```
go get -u github.com/PavloVM7/go-concurrency
```

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
defer ticker.Stop()
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

## ConcurrentSet

`ConcurrentSet` is a thread safe set.

### How to use

```go
func main() {
	const (
		count   = 100_000
		threads = 100
	)
	adds := make([]int, threads)
	set := NewConcurrentSetCapacity[int](count)
	chStart := make(chan struct{})
	chEnd := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			<-chStart
			for j := 1; j <= count; j++ {
				if !set.Contains(j) && set.Add(j) {
					adds[num]++
				}
			}
			<-chEnd
			wg.Done()
		}(i)
	}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			if set.Contains(count) {
				close(chEnd)
				return
			}
		}
	}()
	close(chStart)
	wg.Wait()
	sum := 0
	for _, v := range adds {
		sum += v
	}
	fmt.Println("sum=", sum) // prints 'sum= 100000'
}
```

## ⌨️ Author
[@PavloVM7](https://github.com/PavloVM7) - Idea & Initial work
