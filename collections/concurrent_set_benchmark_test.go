package collections

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkConcurrentSet_Add(b *testing.B) {
	const count = 100_000
	set := NewConcurrentSetCapacity[int](count)
	benchmarks := []struct {
		threads int
	}{
		{threads: 4},
		{threads: 100},
		{threads: 1000},
	}
	addFnc := func(threads int) {
		var run int32
		putF := func() {
			//revive:disable:empty-block
			for atomic.LoadInt32(&run) == 0 {
				// waiting for a start
			}
			//revive:enable:empty-block
			for i := 0; i < count; i++ {
				set.Add(i)
			}
		}
		var wg sync.WaitGroup
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				putF()
			}()
		}
		atomic.StoreInt32(&run, 1)
		wg.Wait()
	}
	for _, bm := range benchmarks {
		bmi := bm
		b.Run(fmt.Sprintf("Add() cnt%d thr%d", count, bmi.threads), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				set.Clear()
				b.StartTimer()
				addFnc(bmi.threads)
				b.StopTimer()
				if set.Size() != count {
					b.Fatal("wrong map size", "expected:", count, "actual:", set.Size())
				}
				b.StartTimer()
			}
		})
	}
}
