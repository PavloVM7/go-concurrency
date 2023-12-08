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

## ConcurrentLinkedList

`ConcurrentLinkedList` is a thread safe linked list realisation

```go
package main

import (
	"fmt"
	"sync"

	"github.com/PavloVM7/go-concurrency/collections"
)

func main() {
	list := collections.NewConcurrentLinkedList[int]()
	using := func(funcs string) {
		fmt.Printf("=== using %s\n", funcs)
	}
	var wg sync.WaitGroup
	using("AddLast() and AddFirst()")
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 11; i <= 20; i++ {
			list.AddLast(i) // adds items to the end of the list
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 10; i > 0; i-- {
			list.AddFirst(i) // adds items to the head of the list
		}
	}()
	wg.Wait()
	showList := func() {
		fmt.Printf(">>> list size: %d, items: %v\n", list.Size(), list.ToArray())
	}
	showList()

	using("Get() and Remove()")
	item10, err := list.Get(10)
	fmt.Printf("before remove 10th item = %d, err = %v\n", item10, err)
	item10, err = list.Remove(10) // removes 10th item
	fmt.Printf("removed item10 = %d, err = %v\n", item10, err)
	item10, err = list.Get(10)
	fmt.Printf("after remove 10th item = %d, err = %v\n", item10, err)
	showList()

	using("GetFirst() and RemoveFirst()")
	first, firstOk := list.GetFirst()
	fmt.Printf("before remove first element: %d, exists: %t\n", first, firstOk)
	first, firstOk = list.RemoveFirst()
	fmt.Printf("first element: %d, removed: %t\n", first, firstOk)
	first, firstOk = list.GetFirst()
	fmt.Printf("current first element: %d, exists: %t\n", first, firstOk)
	showList()

	using("GetLast() and RemoveLast()")
	last, lastOk := list.GetLast()
	fmt.Printf("before remove last element: %d, exists: %t\n", last, lastOk)
	last, lastOk = list.RemoveLast()
	fmt.Printf("last element: %d, removed: %t\n", last, lastOk)
	last, lastOk = list.GetLast()
	fmt.Printf("current last element: %d, exists: %t\n", last, lastOk)
	showList()

	using("RemoveFirstOccurrence()")
	rFirst, fIndex := list.RemoveFirstOccurrence(func(value int) bool {
		return value%2 != 0
	})
	fmt.Printf("removed first odd value: %d, index: %d\n", rFirst, fIndex)
	showList()

	using("RemoveLastOccurrence()")
	rLast, lIndex := list.RemoveLastOccurrence(func(value int) bool {
		return value%2 == 0
	})
	fmt.Printf("removed last even value: %d, index: %d\n", rLast, lIndex)
	showList()

	using("RemoveAll()")
	count := list.RemoveAll(func(value int) bool {
		return value%3 == 0
	})
	fmt.Printf("%d elements that are dividable by 3 have been removed\n", count)
	showList()

	using("Clear()")
	list.Clear()
	showList()
}
```

output:

```
=== using AddLast() and AddFirst()
>>> list size: 20, items: [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20]
=== using Get() and Remove()
before remove 10th item = 11, err = <nil>
removed item10 = 11, err = <nil>
after remove 10th item = 12, err = <nil>
>>> list size: 19, items: [1 2 3 4 5 6 7 8 9 10 12 13 14 15 16 17 18 19 20]
=== using GetFirst() and RemoveFirst()
before remove first element: 1, exists: true
first element: 1, removed: true
current first element: 2, exists: true
>>> list size: 18, items: [2 3 4 5 6 7 8 9 10 12 13 14 15 16 17 18 19 20]
=== using GetLast() and RemoveLast()
before remove last element: 20, exists: true
last element: 20, removed: true
current last element: 19, exists: true
>>> list size: 17, items: [2 3 4 5 6 7 8 9 10 12 13 14 15 16 17 18 19]
=== using RemoveFirstOccurrence()
removed first odd value: 3, index: 1
>>> list size: 16, items: [2 4 5 6 7 8 9 10 12 13 14 15 16 17 18 19]
=== using RemoveLastOccurrence()
removed last even value: 18, index: 14
>>> list size: 15, items: [2 4 5 6 7 8 9 10 12 13 14 15 16 17 19]
=== using RemoveAll()
4 elements that are dividable by 3 have been removed
>>> list size: 11, items: [2 4 5 7 8 10 13 14 16 17 19]
=== using Clear()
>>> list size: 0, items: []
```
## ⌨️ Author
[@PavloVM7](https://github.com/PavloVM7) - Idea & Initial work
