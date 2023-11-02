// Copyright Ⓒ 2023 Pavlo Moisieienko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collections

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkConcurrentMap_Put(b *testing.B) {
	cm := NewConcurrentMap[int, int]()
	const count = 100_000
	benchmarks := []struct {
		name    string
		threads int
		count   int
		fnc     func(k int, v int) (bool, int)
	}{
		{
			name:    "PutIfNotExists",
			fnc:     cm.PutIfNotExists,
			threads: 4,
			count:   count,
		},
		{
			name:    "PutIfNotExistsDoubleCheck",
			fnc:     cm.PutIfNotExistsDoubleCheck,
			threads: 4,
			count:   count,
		},
		{
			name:    "PutIfNotExists",
			fnc:     cm.PutIfNotExists,
			threads: 100,
			count:   count,
		},
		{
			name:    "PutIfNotExistsDoubleCheck",
			fnc:     cm.PutIfNotExistsDoubleCheck,
			threads: 100,
			count:   count,
		},
		{
			name:    "PutIfNotExists",
			fnc:     cm.PutIfNotExists,
			threads: 1000,
			count:   count,
		},
		{
			name:    "PutIfNotExistsDoubleCheck",
			fnc:     cm.PutIfNotExistsDoubleCheck,
			threads: 1000,
			count:   count,
		},
	}
	putFunc := func(threads int, count int, fnc func(k int, v int) (bool, int)) {
		var run int32
		putF := func() {
			//revive:disable:empty-block
			for atomic.LoadInt32(&run) == 0 {
				// waiting for a start
			}
			//revive:enable:empty-block
			for i := 0; i < count; i++ {
				fnc(i, i)
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
		bmv := bm
		b.Run(fmt.Sprintf("cnt%d th%d %s", bm.count, bm.threads, bm.name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				cm.Clear()
				b.StartTimer()
				putFunc(bmv.threads, bmv.count, bmv.fnc)
				b.StopTimer()
				if cm.Size() != bmv.count {
					b.Fatal("wrong map size", "expected:", bmv.count, "actual:", cm.Size())
				}
				b.StartTimer()
			}
		})
	}
}
