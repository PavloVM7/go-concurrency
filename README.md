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
package main

import (
	"fmt"
	"github.com/PavloVM7/go-concurrency/collections"
	"runtime"
	"sync"
	"time"
)

func main() {
	println("👉 Example of using ConcurrentSet")
	using := func(funcs string) {
		fmt.Println("=== using ", funcs)
	}
	set := collections.NewConcurrentSetCapacity[int](10)
	showSet := func() {
		fmt.Printf(">>> ConcurrentSet size: %d, elements: %v\n", set.Size(), set.ToSlice())
	}
	isSetEmpty := func() {
		fmt.Println("~~~ is set empty? -", set.IsEmpty())
	}
	isSetEmpty()

	using("AddAll()")
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if set.AddAll(values...) {
		showSet()
		isSetEmpty()
	}
	subvalues := values[1:8]
	if !set.AddAll(subvalues...) {
		fmt.Printf("- the set already contains values %v\n", subvalues)
	}
	using("Add()")
	val := 11
	if set.Add(val) {
		fmt.Printf("value %d was added to the set\n", val)
		showSet()
	}
	if !set.Add(val) {
		fmt.Println("- the set already contains the value", val)
	}
	using("Contains()")
	showSet()
	if set.Contains(3) {
		fmt.Println("+ the set contains the value 3")
	}
	if set.Contains(4) {
		fmt.Println("+ the set contains the value 4")
	}
	if !set.Contains(123) {
		fmt.Println("- there is no value 123 in the set")
	}

	using("Remove()")
	if set.Remove(3) {
		fmt.Printf("+ the value %d was removed from the set\n", 3)
	}
	if set.Remove(4) {
		fmt.Printf("+ the value %d was removed from the set\n", 4)
	}
	if !set.Remove(123) {
		fmt.Printf("- the value %d was not removed from the set because the set did not contain it\n", 123)
	}
	showSet()

	using("Clear()")
	set.Clear()
	showSet()
	isSetEmpty()

	using("TrimToSize()")
	const amount = 1_000_000
	fillSet(set, amount, 2)
	fmt.Println(">>> set size =", set.Size())

	getMemStats := func() runtime.MemStats {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem
	}

	memToString := func(mem runtime.MemStats) string { return fmt.Sprintf("%d Kb", mem.Alloc/1024) }

	runtime.GC()

	fmt.Printf(">>> set size: %d, memory usage: %s\n", set.Size(), memToString(getMemStats()))

	removeValues(set, 21, amount, 3)

	runtime.GC()

	fmt.Printf("after removing memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()

	set.TrimToSize()

	runtime.GC()

	fmt.Printf("after TrimToSize() memory usage: %s, set size: %d\n", memToString(getMemStats()), set.Size())
	showSet()
}

func fillSet(set *collections.ConcurrentSet[int], amount, threads int) {
	fmt.Printf("* filling set, amount: %d, threads: %d\n", amount, threads)
	start := time.Now()
	chStart := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			<-chStart
			n := 0
			for !set.Contains(amount) {
				n++
				if set.Add(n) {
					adds[num]++
				}
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> the set was filled, duration: %v, amount: %d, threads: %d, each thread added: %v\n",
		time.Since(start), set.Size(), threads, adds)
}
func removeValues(set *collections.ConcurrentSet[int], start, end, threads int) {
	fmt.Printf("* remove values from set, from %d to %d , threads: %d\n", start, end, threads)
	st := time.Now()
	chStart := make(chan struct{})
	var wg sync.WaitGroup
	adds := make([]int, threads)
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			<-chStart
			for val := start; val <= end; val++ {
				if set.Remove(val) {
					adds[num]++
				}
			}
		}(i)
	}
	close(chStart)
	wg.Wait()
	fmt.Printf(">>> values were removed, duration: %v, set size: %d, threads: %d, each thread removed: %v\n",
		time.Since(st), set.Size(), threads, adds)
}
```

outputs like this:

```
👉 Example of using ConcurrentSet
~~~ is set empty? - true
=== using  AddAll()
>>> ConcurrentSet size: 10, elements: [1 3 6 7 2 4 5 8 9 10]
~~~ is set empty? - false
- the set already contains values [2 3 4 5 6 7 8]
=== using  Add()
value 11 was added to the set
>>> ConcurrentSet size: 11, elements: [10 11 2 4 5 8 9 1 3 6 7]
- the set already contains the value 11
=== using  Contains()
>>> ConcurrentSet size: 11, elements: [3 6 7 1 4 5 8 9 10 11 2]
+ the set contains the value 3
+ the set contains the value 4
- there is no value 123 in the set
=== using  Remove()
+ the value 3 was removed from the set
+ the value 4 was removed from the set
- the value 123 was not removed from the set because the set did not contain it
>>> ConcurrentSet size: 9, elements: [2 5 8 9 10 11 1 6 7]
=== using  Clear()
>>> ConcurrentSet size: 0, elements: []
~~~ is set empty? - true
=== using  TrimToSize()
* filling set, amount: 1000000, threads: 2
>>> the set was filled, duration: 296.042ms, amount: 1000000, threads: 2, each thread added: [570956 429044]
>>> set size = 1000000
>>> set size: 1000000, memory usage: 21898 Kb
* remove values from set, from 21 to 1000000 , threads: 3
>>> values were removed, duration: 467.082917ms, set size: 20, threads: 3, each thread removed: [346462 192433 461085]
after removing memory usage: 21900 Kb, set size: 20
>>> ConcurrentSet size: 20, elements: [15 8 13 16 18 1 6 17 12 3 2 19 14 11 9 4 20 5 7 10]
after TrimToSize() memory usage: 100 Kb, set size: 20
>>> ConcurrentSet size: 20, elements: [9 4 5 10 2 14 11 12 3 15 16 1 6 17 8 13 18 19 20 7]

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
